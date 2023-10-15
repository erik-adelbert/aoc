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
)

var (
	reads []reading
	scans reading
	fixed map[vec]bool
)

// niladic init
//
// https://go.dev/doc/effective_go#initialization
func init() {
	reads = make([]reading, 0, 32)
	scans = append(make([]vec, 0, 27), vec{0, 0, 0}) // origin
	fixed = make(map[vec]bool, 337)
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
		fixed[p] = true
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
	fmt.Println(len(fixed)) // part1

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

type rotator func() reading // rotations iterator

func (r reading) π2rots() rotator {
	rots := []struct {
		s, a vec // sign, axis order
	}{
		{vec{-1, -1, +1}, vec{X, Y, Z}},
		{vec{-1, +1, -1}, vec{X, Y, Z}},
		{vec{+1, -1, -1}, vec{X, Y, Z}},
		{vec{+1, +1, +1}, vec{X, Y, Z}},

		{vec{-1, -1, -1}, vec{X, Z, Y}},
		{vec{-1, +1, +1}, vec{X, Z, Y}},
		{vec{+1, -1, +1}, vec{X, Z, Y}},
		{vec{+1, +1, -1}, vec{X, Z, Y}},

		{vec{-1, -1, -1}, vec{Y, X, Z}},
		{vec{-1, +1, +1}, vec{Y, X, Z}},
		{vec{+1, -1, +1}, vec{Y, X, Z}},
		{vec{+1, +1, -1}, vec{Y, X, Z}},

		{vec{-1, -1, +1}, vec{Y, Z, X}},
		{vec{-1, +1, -1}, vec{Y, Z, X}},
		{vec{+1, -1, -1}, vec{Y, Z, X}},
		{vec{+1, +1, +1}, vec{Y, Z, X}},

		{vec{-1, -1, +1}, vec{Z, X, Y}},
		{vec{-1, +1, -1}, vec{Z, X, Y}},
		{vec{+1, -1, -1}, vec{Z, X, Y}},
		{vec{+1, +1, +1}, vec{Z, X, Y}},

		{vec{-1, -1, -1}, vec{Z, Y, X}},
		{vec{-1, +1, +1}, vec{Z, Y, X}},
		{vec{+1, -1, +1}, vec{Z, Y, X}},
		{vec{+1, +1, -1}, vec{Z, Y, X}},
	}

	i, turned := 0, make(reading, len(r))

	return func() reading { // rotator
		if i >= len(rots) { // capture i
			return nil
		}

		rot := rots[i]
		for j, v := range r {
			// π/2 rotation
			// shuffle X, Y, Z according to sign and axis order
			s, a := &rot.s, &rot.a
			turned[j] = vec{ // capture turned
				s[0] * v[a[0]], s[1] * v[a[1]], s[2] * v[a[2]],
			}
		}
		i++
		return turned // leak turned
	}
}

func list(m map[vec]bool) reading {
	list := make(reading, len(m))

	i := 0
	for v := range m {
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
}

func ralign(r reading) bool {
	next := r.π2rots()

	rot := next()
	for rot != nil {
		if align(rot) {
			return true
		}
		rot = next()
	}
	return false
}

func align(r reading) bool {

	// index := func(r reading, x vec) int {  // older Go
	// 	for i, v := range r {
	// 		if x == v {
	// 			return i
	// 		}
	// 	}
	// 	return len(r)
	// }

	index := func(r reading, x vec) int { // Go 1.21
		return slices.Index(r, x)
	}

	const TRESH = 12 // from challenge

	known := list(fixed)
	for a := X; a <= Z; a++ { // X, Y, Z
		sort.Slice(r, func(i, j int) bool {
			return r[i][a] <= r[j][a]
		})

		sort.Slice(known, func(i, j int) bool {
			return known[i][a] <= known[j][a]
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
				for _, v := range rebased {
					fixed[v] = true
				}
				scans = append(scans, o)
				return true
			}
		}
	}
	return false
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
