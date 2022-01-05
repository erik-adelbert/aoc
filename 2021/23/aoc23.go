package main

import (
	hp "container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

type board [11]string // hhRhRhRhRhh h(allway), R(oom)

func cost(p rune) int {
	costs := [...]int{1, 10, 100, 1000}
	return costs[p-'A']
}

func goal(p rune) int {
	return 2 + 2*int(p-'A') // 'A': 2, 'B': 4, 'C':6 'D': 8
}

func room(r int) bool {
	return 1 < r && r < 9 && r&1 == 0 // true for 2, 4, 6, 8
}

// free checks if hallway cells between s,t are free
func (b board) free(s, t int) bool { // s(rc), t(arget)
	l, r := min(s, t), max(s, t)
	for i := l; i <= r; i++ {
		if i != s && !room(i) && b[i] != "." {
			return false
		}
	}
	return true
}

// granted checks if a room is either empty or populated with
// homies only
func (b board) granted(r int, p rune) bool { // r(oom), p(awn)
	if r != goal(p) {
		return false
	}
	for _, c := range b[r] {
		if c != '.' && c != p {
			return false
		}
	}
	return true
}

func (b board) pawn(r int) rune { // r(oom)
	for _, c := range b[r] {
		if c != '.' {
			return c
		}
	}
	return 0
}

func (b board) get(r int) (string, int) { // r(oom)
	cell := []byte(b[r])
	for i, c := range b[r] {
		if c != '.' {
			cell[i] = '.'
			return string(cell), i + 1
		}
	}
	return string(cell), 0
}

func (b board) put(r int, p rune) (string, int) { // r(oom), p(awn)
	cell := []byte(b[r])
	if i := strings.Count(b[r], "."); i != 0 { // room has free cells
		cell[i-1] = byte(p) // take the deeper one
		return string(cell), i
	}
	return b[r], 0
}

func (b board) moves(r int) []int { // r(oom)

	p := b.pawn(r) // pawn to move

	if r == goal(p) && b.granted(r, p) { // pawn already at destination
		return []int(nil)
	}

	if !room(r) { // pawn in the hallway, moving to goal is the only move
		if b.free(r, goal(p)) && b.granted(goal(p), p) {
			return []int{goal(p)} // move if way is free and room is open
		}
		return []int(nil)
	}

	moves := make([]int, 0, 8)
	for i := 0; i < len(b); i++ {
		switch {
		case i == r: // skip starting room
		case i != goal(p) && room(i): // enter no room except for...
		case i == goal(p) && !b.granted(i, p): // ...the goal one, if not closed ...
		case b.free(r, i): // ... and free way
			moves = append(moves, i)
		}
	}
	return moves
}

func (b board) move(s, t int) (board, int) { // s(ource), t(arget) -> board, cost
	nxt, p := b, b.pawn(s)

	n, dist := 0, 0
	nxt[s] = "."
	if room(s) {
		nxt[s], n = b.get(s)
		dist += n
	}
	nxt[t] = string(p)
	if room(t) {
		nxt[t], n = b.put(t, p)
		dist += n
	}
	dist += abs(s - t)
	return nxt, dist * cost(p)
}

var costs map[board]int

func init() {
	costs = make(map[board]int, 16411)
}

type cboard struct {
	b board
	c int
}

type heap []*cboard

func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool {
	return h[i].c < h[j].c
}

func (h heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heap) Push(x interface{}) {
	b := x.(cboard)
	*h = append(*h, &b)
}

func (h *heap) Pop() interface{} {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], nil
	return pop
}

func (b board) solve(goal board) int {
	heap := make(heap, 0, 16411)
	hp.Init(&heap)

	hp.Push(&heap, cboard{b, 0}) // from start...
	for heap.Len() > 0 {         // ...play all possible games
		b := hp.Pop(&heap).(*cboard).b // pop a (sub)game
		if b == goal {                 // it works because heap is sorted by costs
			return costs[goal]
		}
		for i := range b { // for all cells
			if b.pawn(i) == 0 { // empty cell, nothing to do
				continue
			}
			for _, j := range b.moves(i) { // for all legal moves from cell i...
				sub, cost := b.move(i, j) // ...play one
				cost += costs[b]
				if costs[sub] == 0 || costs[sub] > cost {
					costs[sub] = cost                 // if it's the best move so far...
					hp.Push(&heap, cboard{sub, cost}) // ...send subgame to resolution
				}
			}
		}
	}
	return costs[goal]
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	goal := board{
		".", ".", "AA", ".", "BB", ".", "CC", ".", "DD", ".", ".",
	}
	part1 := board{
		".", ".", "AB", ".", "DC", ".", "BA", ".", "DC", ".", ".",
	}
	fmt.Println(part1.solve(goal))

	goal = board{
		".", ".", "AAAA", ".", "BBBB", ".", "CCCC", ".", "DDDD", ".", ".",
	}
	part2 := board{
		".", ".", "ADDB", ".", "DCBC", ".", "BBAA", ".", "DACC", ".", ".",
	}
	fmt.Println(part2.solve(goal))

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
