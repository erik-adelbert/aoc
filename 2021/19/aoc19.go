// aoc19.go --
// advent of code 2021 day 19
//
// https://adventofcode.com/2021/day/19
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-19: initial commit

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type safeVecMap struct {
	mu sync.Mutex
	m  map[vec]bool
}

var (
	reads []reading
	scans reading
	fixed = new(safeVecMap)
)

// niladic init
//
// https://go.dev/doc/effective_go#initialization
func init() {
	reads = make([]reading, 0, 32)
	scans = append(make([]vec, 0, 27), vec{0, 0, 0}) // origin
	fixed.m = make(map[vec]bool, 337)
}

func main() {
	points := make(reading, 0, 32)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch line := input.Text(); line != "" {
		default:
			args := strings.Split(line, ",") // trust input
			points = append(points, vec{
				X: atoi(args[0]),
				Y: atoi(args[1]),
				Z: atoi(args[2]),
			})

		case strings.HasPrefix(line, "---"):
			if len(points) > 0 {
				reads = append(reads, clone(points))
				points = points[:0] // reset
			}
		}
	}
	if points != nil {
		reads = append(reads, points) // last reading
	}

	// dice trick!
	// exhibit runtime ±22% relative to readings order
	rand.Shuffle(len(reads), func(i, j int) {
		reads[i], reads[j] = reads[j], reads[i]
	})

	for _, p := range reads[0] { // origin
		fixed.m[p] = true
	}
	reads = reads[1:] // shift

	// align all readings gradually
	for len(reads) > 0 {
		i := 0
		for _, r := range reads {
			if !ralign(r) { // populate scans[] and fixed[]
				reads[i] = r // no match yet, push back
				i++
			}
		}
		reads = reads[:i] // retry with pushed back
	}
	fmt.Println(len(fixed.m)) // part1

	diam := 0
	for i, v0 := range scans {
		for _, v1 := range scans[i+1:] {
			if dist := v1.sub(v0).manh(); dist > diam {
				diam = dist // max
			}
		}
	}
	fmt.Println(diam) // part2
}

// Axis
const (
	X = iota
	Y
	Z
)

type vec [3]int

func (v vec) sub(u vec) vec {
	return vec{
		v[X] - u[X], v[Y] - u[Y], v[Z] - u[Z],
	}
}

func (v vec) manh() int {
	return abs(v[X]) + abs(v[Y]) + abs(v[Z])
}

type reading []vec

// clone is a vector field abstraction
func clone(points []vec) reading {
	buf := make(reading, len(points))
	copy(buf, points) // clone
	return buf
}

func (r reading) π2rots() <-chan reading {
	out := make(chan reading, 24)
	defer close(out)

	roll := func(r reading) reading {
		for i := range r {
			r[i] = vec{r[i][X], r[i][Z], -r[i][Y]}
		}
		return r
	}

	turn := func(r reading) reading {
		for i := range r {
			r[i] = vec{-r[i][Y], r[i][X], r[i][Z]}
		}
		return r
	}

	//    _________
	//   /  top   /| t
	//  /_______ / |s  half: top, near, east
	// |        | a|   roll (R): near -> top
	// |  near  |e /   turn (T): near <- east
	// |        | /
	// |________|/
	//
	// RTTTRTTTRTTT -> first half 12 rotations
	// RTR          -> switch half
	// RTTTRTTTRTTT -> second half 12 rotations
	//
	// https://tinyurl.com/yckc8c5b (SO discussion)

	for half := 0; half < 2; half++ {
		for top := 0; top < 3; top++ {
			// 3xRTTT
			r = roll(r) // R
			out <- clone(r)
			for near := 0; near < 3; near++ {
				r = turn(r) // T
				out <- clone(r)
			}
		}
		r = roll(turn(roll(r))) // RTR
	}

	return out
}

func list(vmap *safeVecMap) reading {
	vmap.mu.Lock()
	defer vmap.mu.Unlock()

	list := make(reading, len(vmap.m))
	i := 0
	for v := range vmap.m {
		list[i] = v
		i++
	}
	return list
}

func difs(r reading) reading {
	difs := make([]vec, len(r)-1)
	for i, v1 := range r[1:] { // i, v1 = {0, r[1]}, {1, r[2]}, ...
		difs[i] = v1.sub(r[i]) // difs[0] = [r[1]-r[0], r[2]-r[1] ...]
	}
	return difs
}

func inter(a, b reading, first bool) reading {
	if len(a) > len(b) {
		a, b = b, a
	}

	m := make(map[vec]bool, len(a))
	for _, v := range a {
		m[v] = true
	}

	inter := make(reading, 0, len(m))
	for _, v := range b {
		if m[v] {
			inter = append(inter, v)
			if first {
				break
			}
		}
	}
	return inter
}

// in-place rebase
func rebase(r reading, o vec) reading {
	for i := range r {
		r[i] = r[i].sub(o)
	}
	return r
}

func merge(done <-chan struct{}, cs ...<-chan bool) <-chan bool {
	var wg sync.WaitGroup
	out := make(chan bool)

	// start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan bool) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func ralign(r reading) bool {
	done := make(chan struct{})
	defer close(done) // cancellation channel

	a1 := align(done, r.π2rots()) // rotate-align pipeline
	a2 := align(done, r.π2rots())

	for match := range merge(done, a1, a2) {
		if match {
			return true // and cancel
		}
	}
	return false // and cancel
}

func align(done <-chan struct{}, in <-chan reading) <-chan bool {
	out := make(chan bool)

	index := func(r reading, x vec) int { // Go 1.21
		return slices.Index(r, x)
	}

	go func() {
		defer close(out)
	ALIGN:
		for r := range in {
			const TRESH = 3 // from challenge

			known := list(fixed)
			for a := range []int{X, Y, Z} {
				sort.Slice(r, func(i, j int) bool {
					return r[i][a] < r[j][a]
				})

				sort.Slice(known, func(i, j int) bool {
					return known[i][a] < known[j][a]
				})

				rdifs, kdifs := difs(r), difs(known)

				const (
					ALL   = false
					FIRST = !ALL
				)

				if matches := inter(rdifs, kdifs, FIRST); len(matches) == 1 {
					pivot := matches[0]
					i := index(rdifs, pivot)
					j := index(kdifs, pivot)
					o := r[i].sub(known[j])
					rebased := rebase(r, o)

					if len(inter(known, rebased, ALL)) >= TRESH {
						fixed.mu.Lock()
						for _, v := range rebased {
							fixed.m[v] = true
						}
						fixed.mu.Unlock()
						scans = append(scans, o)
						select {
						case out <- true:
						case <-done:
							return
						}
						continue ALIGN
					}
				}
			}
			select {
			case out <- false:
			case <-done:
				return
			}
		}
	}()
	return out
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s) // trust s
	return n
}
