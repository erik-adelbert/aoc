package main

import (
	hp "container/heap"
	"fmt"
	"strings"
	"sync"
)

type board [11]string // hhRhRhRhRhh h(allway), R(oom)

func cost(p rune) int {
	costs := [...]int{1, 10, 100, 1000}
	return costs[rtoi(p)]
}

func goal(p rune) int {
	return 2 + 2*int(rtoi(p)) // 'A': 2, 'B': 4, 'C':6 'D': 8
}

func room(r int) bool {
	return 1 < r && r < 9 && r&1 == 0 // true for 2, 4, 6, 8
}

var (
	halls = []int{0, 1, 3, 5, 7, 9, 10} // hallway cells
	rooms = []int{2, 4, 6, 8}
)

// free checks if hallway cells between s,t are free
func (b board) free(s, t int) bool {
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
func (b board) granted(r int, p rune) bool { // room, pawn
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

func (b board) pawn(r int) rune { // room
	for _, c := range b[r] {
		if c != '.' {
			return c
		}
	}
	return 0
}

func (b board) rem(r int, p rune) (string, int) { // room, pawn -> board, cost
	if i := strings.IndexRune(b[r], p); i > -1 {
		cell := []rune(b[r])
		cell[i] = '.'
		return string(cell), i + 1
	}
	return b[r], 0
}

func (b board) add(r int, p rune) (string, int) { // room, pawn -> board, cost
	if i := strings.Count(b[r], "."); i != 0 { // room has free cells
		cell := []rune(b[r])
		cell[i-1] = p // take the deeper one
		return string(cell), i
	}
	return b[r], 0
}

// https://github.com/pemoreau/advent-of-code-2021/blob/main/go/23/day23.go#L147-L204
// dead1 detects an interlock in the middle section of the board
func (b board) dead1() bool {
	for i := range []int{3, 5, 7} {
		for j := i + 2; j < 8; j += 2 {
			x, y := b.pawn(i), b.pawn(j)
			if x*y*(y-x) != 0 && goal(x) >= j && goal(y) <= i {
				return true
			}
		}
	}
	return false
}

// dead2 detects an interlock at either edge of the board
func (b board) dead2() bool {
	edges := []struct {
		r   rune
		off int
	}{
		{'D', +1},
		{'A', -1},
	}

	for _, e := range edges {
		g := goal(e.r)
		if b.pawn(g-e.off) == e.r {
			nspace := 0
			if b.pawn(g+e.off) == 0 {
				nspace++
			}
			if b.pawn(g+2*e.off) == 0 {
				nspace++
			}
			nalien := 0
			if !b.granted(g, e.r) {
				for _, r := range b[g] {
					if r != e.r {
						nalien++
					}
				}
			}
			if nalien > nspace {
				return true
			}
		}
	}
	return false
}

type cboard struct {
	b *board
	c int
}

func (b board) moves() []cboard {
	if b.dead1() || b.dead2() { // prune deadlocked board
		return []cboard(nil) // no move
	}

	// step1 - go back home
	nxt := cboard{&b, 0}
	done := false
	for !done { // always a good move, make all such moves at once!
		done = true
		for i := range nxt.b {
			b := nxt.b
			p := b.pawn(i)
			g := goal(p)
			if p == 0 || (i == g && b.granted(g, p)) { // skip empty & @home
				continue
			}
			if b.free(i, g) && b.granted(g, p) {
				new, cost := b.move(i, g)
				nxt.b, nxt.c = new, nxt.c+cost
				done = false
			}
		}
	}

	if nxt.c > 0 { // send back for prioritization
		return []cboard{nxt}
	}

	// step2 - move out (later)
	moves := make([]cboard, 0, 28)
	for _, s := range rooms {
		b := nxt.b
		p := b.pawn(s)
		if p == 0 || b.granted(s, p) { // skip empty & @home
			continue
		}
		for _, t := range halls {
			if b.free(s, t) {
				new, cost := b.move(s, t)
				moves = append(moves, cboard{new, nxt.c + cost})
			}
		}
	}

	return moves
}

// move moves the top pawn of s to t, it returns the resulting board and the move cost
func (b board) move(s, t int) (*board, int) {
	nxt := b // array copy
	p := b.pawn(s)

	n, dist := 0, 0
	nxt[s] = "."
	if room(s) {
		nxt[s], n = b.rem(s, p)
		dist += n
	}
	nxt[t] = string(p)
	if room(t) {
		nxt[t], n = b.add(t, p)
		dist += n
	}
	dist += abs(s - t)
	return &nxt, dist * cost(p)
}

// hcost is a heuristic function from:
// https://github.com/pemoreau/advent-of-code-2021/blob/main/go/23/day23.go#L377-L401
// see:
// https://www.reddit.com/r/adventofcode/comments/rzvsjq/comment/hswxkbr/?utm_source=share&utm_medium=web2x&context=3
func hcost(b *board) int {
	popcnts := make([]int, len(b))
	entropy := 0

SCAN:
	for s := range b { // for all cells as sources
		var pawns strings.Builder
		pawns.Grow(len(b[s]))
		for _, p := range b[s] {
			if p != '.' {
				pawns.WriteRune(p)
			}
		}

		for j, p := range pawns.String() { // for all pawns in a source cell
			t := goal(p) // target

			if s == t && b.granted(t, p) { // pawns already @home
				continue SCAN // discard
			}

			dist := abs(s - t)             // walk p back home
			dist += len(b[t]) - popcnts[t] // put it in there
			popcnts[t]++

			if room(s) { // get it out first
				dist += len(b[s]) - len(pawns.String())
				dist += 1 + j
			}

			entropy += dist * cost(p)
		}
	}
	return entropy
}

type heap []cboard

func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool {
	return h[i].c < h[j].c
}

func (h heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heap) Push(x interface{}) {
	b := x.(cboard)
	*h = append(*h, b)
}

func (h *heap) Pop() interface{} {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], cboard{}
	return pop
}

func (b board) solve(goal *board, costs map[string]int) int {
	concat := func(b *board) string {
		var sb strings.Builder
		sb.Grow(24)
		for _, s := range *b {
			sb.WriteString(s)
		}
		return sb.String()
	}

	heap := make(heap, 0, 5920)
	hp.Init(&heap)

	hp.Push(&heap, cboard{&b, 0}) // from the start...
	for heap.Len() > 0 {          // ...play all possible games
		cur := hp.Pop(&heap).(cboard).b // pop a (sub)game board
		if *cur == *goal {              // this cut works because heap is kinda sorted by costs
			return costs[concat(goal)]
		}
		for _, move := range cur.moves() {
			nxt, cost := move.b, move.c
			cost += costs[concat(cur)]

			nkey := concat(nxt)
			if known, seen := costs[nkey]; !seen || known > cost {
				costs[nkey] = cost                // if it's the best move so far...
				prio := cost + hcost(nxt)         // prioritize by hypercosts (heuristic/entropy)
				hp.Push(&heap, cboard{nxt, prio}) // ...send subgame to resolution
			}
		}
	}
	panic("solve() unreachable")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		costs := make(map[string]int, 4800)
		goal := board{
			".", ".", "AA", ".", "BB", ".", "CC", ".", "DD", ".", ".",
		}
		part1 := board{
			".", ".", "AB", ".", "DC", ".", "BA", ".", "DC", ".", ".",
		}
		fmt.Println(part1.solve(&goal, costs)) // part1
		// fmt.Println(len(costs))
		wg.Done()
	}()

	go func() {
		costs := make(map[string]int, 42200)
		goal := board{
			".", ".", "AAAA", ".", "BBBB", ".", "CCCC", ".", "DDDD", ".", ".",
		}
		part2 := board{
			".", ".", "ADDB", ".", "DCBC", ".", "BBAA", ".", "DACC", ".", ".",
		}
		fmt.Println(part2.solve(&goal, costs)) // part2
		// fmt.Println(len(costs))
		wg.Done()
	}()

	wg.Wait()
}

func rtoi(r rune) rune {
	return r - 'A'
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
