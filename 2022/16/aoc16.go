// aoc16.go --
// advent of code 2022 day 16
//
// https://adventofcode.com/2022/day/16
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-16: initial commit
// 2023-11-22: adapt from:
//   https://github.com/maneatingape/advent-of-code-rust/blob/11750514bb00915bf23fcee22c00fcaeb6a64a5c/src/year2022/day16.rs

package main

import (
	"bufio"
	"cmp"
	"fmt"
	"math/bits"
	"os"
	"slices"
	"sort"
	"strings"
)

func main() {
	w := readWorld(bufio.NewScanner(os.Stdin))
	part1(w)
	part2(w)
}

func part1(w *world) {
	score := 0
	highscore := func(_, flow int) int {
		score = max(score, flow) // capture score
		return score
	}
	part1 := &state{todo: w.state, from: w.from, time: 30, flow: 0}
	w.bbsolve(part1, highscore) // update score via highscore() closure

	fmt.Println(score)
}

func part2(w *world) {

	step1 := &state{todo: w.state, from: w.from, time: 26, flow: 0}
	score1, closed := 0, 0
	highscore := func(todo, flow int) int {
		if flow > score1 {
			score1 = max(score1, flow) // capture score1
			closed = todo              // capture left
		}
		return score1
	}
	w.bbsolve(step1, highscore) // update score1 & left via highscore() closure

	step2 := &state{todo: closed, from: w.from, time: 26, flow: 0}
	score2 := 0
	highscore = func(_, flow int) int {
		score2 = max(score2, flow) // capture score
		return score2
	}
	w.bbsolve(step2, highscore) // update score2 via highscore() closure

	step3 := &state{todo: w.state, from: w.from, time: 26, flow: 0}
	scores3 := make([]int, w.state+1)
	highscore = func(todo, flow int) int {
		done := w.state ^ todo
		scores3[done] = max(scores3[done], flow) // capture scores3
		return score2                            // use score2 as heuristic baseline
	}
	w.bbsolve(step3, highscore) // update scores3 via highscore() closure

	// sanitize and prepare scores3
	scores := make([]struct{ i, v int }, 0, len(scores3)/2)
	for i, v := range scores3 {
		if v > 0 {
			scores = append(scores, struct{ i, v int }{i, v})
		}
	}
	slices.SortFunc(scores, func(a, b struct{ i, v int }) int {
		return -cmp.Compare(a.v, b.v)
	})

	// maxout best score
	best := score1 + score2
	for i := range scores[:len(scores)-1] {
		mask1, score1 := scores[i].i, scores[i].v

		if 2*score1 <= best {
			break
		}

		for j := range scores[i+1:] {
			mask2, score2 := scores[j].i, scores[j].v

			if mask1&mask2 == 0 {
				// score1&2 are for disjoint movesets
				best = max(best, score1+score2)
				break
			}
		}

	}

	fmt.Println(best)
}

const (
	MaxInt = int(^uint(0) >> 1)
	Inf    = MaxInt
)

type world struct {
	size  int
	from  int
	state int
	flows []int
	dists []int
	nears []int
}

type valve struct {
	name  string
	flow  int
	links []string
}

func readWorld(input *bufio.Scanner) *world {
	valves := make([]valve, 0, 60)

	var r = strings.NewReplacer(
		"=", " ",
		";", "",
		",", "",
	)

	for input.Scan() {
		// input line is:
		// ^\w+\s([A-Z]{2})(\s\w+){3}=(\d+);(\s\w+){4}\s([A-Z]{2})(,\s([A-Z]{2}))*$
		// ex.
		// Valve GO has flow rate=0; tunnels lead to valves HO, DO
		//
		// replace '=' with ' ' and remove extraneous cars ';' and ',':
		// ^\w+\s([A-Z]{2})(\s\w+){3}\s(\d+)(\s\w+){4}(\s([A-Z]{2}))+$
		// split on space:
		//  0    1                        5                        10:
		// [\w+, [A-Z]{2}, \w+, \w+, \w+, \d+, \w+, \w+, \w+, \w+, ([A-Z]{2}))+ ]
		//       name                     flow                     links...
		args := strings.Fields(r.Replace(input.Text()))
		valves = append(valves, valve{
			name: args[1], flow: atoi(args[5]), links: args[10:],
		})
	}

	sort.Sort(byDescendingFlow(valves))

	// size is non-zero flow valve count plus 1 for "AA"
	size := 1
	for i := range valves {
		if valves[i].flow == 0 {
			break
		}
		size++
	}

	// valve name to index map
	vids := make(map[string]int, len(valves))
	for i, v := range valves {
		vids[v.name] = i
	}

	// distance between 2 given valves flatten index
	idx := func(from, to int) int {
		return from*size + to
	}

	// flatten valve distance matrix
	dists := mkIntSlice(Inf, size*size)
	for from, valve := range valves[:size] {
		dists[idx(from, from)] = 0

		for _, link := range valve.links {
			pre, cur := valve.name, link
			to := vids[cur]
			dist := 1

			for to >= size {
				for _, nxt := range valves[to].links {
					if nxt != pre {
						pre, cur, to, dist = cur, nxt, vids[nxt], dist+1
						break
					}
				}
			}
			dists[idx(from, to)] = dist
		}
	}

	// find all-pairs shortest distances
	// symetric floyd-warshall flooding
	for k := 0; k < size; k++ {
		for i := 0; i < size; i++ {
			for j := 0; j < i; j++ {
				if v := dists[idx(i, k)] + dists[idx(k, j)]; v > 0 {
					dists[idx(i, j)] = min(dists[idx(i, j)], v)
					dists[idx(j, i)] = dists[idx(i, j)]
				}
			}
		}
	}
	for i := range dists {
		dists[i]++ // offset 1mn for valve opening
	}

	w := new(world)
	w.size = size
	w.from = vids["AA"]
	w.state = 1<<(size-1) - 1
	w.dists = dists

	w.flows = make([]int, size)
	w.nears = make([]int, size)
	for i := range valves[:size] {
		w.flows[i] = valves[i].flow

		min := struct{ d, i int }{Inf, 0}

		lo, hi := idx(i, 0), idx(i, size)
		for i, d := range dists[lo:hi:hi] {
			if 1 < d && d < min.d {
				min.i, min.d = i, d
			}
		}

		w.nears[i] = min.d
	}

	return w
}

func (w *world) dist(from, to int) int {
	return w.dists[from*w.size+to]
}

type state struct {
	todo int
	from int
	time int
	flow int
}

func (w *world) bbsolve(s *state, fscore func(int, int) int) {
	todo, from, time, flow := s.todo, s.from, s.time, s.flow
	score := fscore(todo, flow)

	valves := todo
	for valves > 0 {
		to := bits.TrailingZeros(uint(valves))
		mask := 1 << to
		valves ^= mask

		trip := w.dist(from, to)
		if trip >= time {
			continue
		}

		time := time - trip
		todo := todo ^ mask
		flow := flow + time*w.flows[to]

		heuristic := func() int {
			valves, time, flow := todo, time, flow

			for valves > 0 && time > 3 {
				to := bits.TrailingZeros(uint(valves))
				valves ^= 1 << to
				time -= w.nears[to]
				flow += time * w.flows[to]
			}

			return flow
		}

		if heuristic() > score {
			next := &state{todo, to, time, flow}
			w.bbsolve(next, fscore)
		}
	}
}

// sort interface by descending flow then ascending name
type byDescendingFlow []valve

func (a byDescendingFlow) Len() int { return len(a) }

// sort by descending flow then ascending name
func (a byDescendingFlow) Less(i, j int) bool {
	if a[i].flow == a[j].flow {
		return a[i].name < a[j].name
	}
	return a[i].flow > a[j].flow
}
func (a byDescendingFlow) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// make an initialized int slice
func mkIntSlice(value int, size int) []int {
	s := make([]int, size)
	for i := range s {
		s[i] = value
	}
	return s
}

// strconv.Atoi modified core loop
// s is ^\d+.*$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

var DEBUG = false

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}

func debugf(format string, a ...any) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}
