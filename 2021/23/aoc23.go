package main

import (
	"bufio"
	"bytes"
	hp "container/heap"
	"fmt"
	"os"
	"strings"
)

const (
	// from challenge
	RLEN = 7
	BLEN = 14
)

type (
	// a burrow is a flattened, $ delimited, RLENxBLEN
	// row-major byte matrix: j is row, i is column
	// https://www.ce.jhu.edu/dalrymple/classes/602/Class12.pdf
	//
	// ex:
	//
	//	RLEN = 7, BLEN = 14, $ = \n, ε = \0, _ = 0x20 (ASCII space)
	//
	//	       part#2                part#1
	//	j\i|0123456789abcd|   j\i|0123456789abcd|
	//	0  |#############$|   0  |#############$| ↑
	//	1  |#...........#$|   1  |#...........#$| |
	//	2  |###A#B#C#D###$|   2  |###A#B#C#D###$| |
	//	3  |__#A#B#C#D#__$|   3  |__#A#B#C#D#__$| RLEN
	//	4  |__#A#B#C#D#__$|   4  |__#########__$| |
	//	5  |__#A#B#C#D#__$|   5  |εεεεεεεεεεεεε$| |
	//	6  |__#########__$|   6  |εεεεεεεεεεεεε$| ↓
	//	    <--- BLEN --->
	//
	//	buro memory layout [RLEN * BLEN]byte:
	//
	//	      RLEN * BLEN = 0x62
	//	index|<------------------- RLEN * BLEN -/~/--------------------->| limit
	//	  j: |0 <- BLEN -> |1            |2     /~/       |6             | RLEN
	//	raw: [#############$#...........#$###A#B/~/#C#D#  $  #########  $]
	//	  i: |0123456789abcd0123456789abcd012345/~/6789abcd0123456789abcd| BLEN
	//	 ii: |0123456789abcdef...               /~/                  ...ω| RLEN * BLEN
	//	                                                     ω = 0x62 - 1
	//
	// 2D burrow from/to buro:  ii = j*BLEN + i  <=> j = ii/BLEN, i = ii%BLEN
	buro [RLEN * BLEN]byte

	cost int

	// a move represent a game state and its cost
	//   - it is designed for A* operations: prio(), setprio()
	//   - it boasts a classical interface: move(), moves()
	//   - A* is attached to it: solve()
	move struct {
		b    *buro
		c, S cost
	}
)

// main entry point
func main() {
	input := bufio.NewScanner(os.Stdin)

	// muxed part 1, 2
	parts := mkburos(input)
	for p := range parts {
		fmt.Println(newMove(&parts[p], 0).solve())
	}
}

// uncomment and fix for runtime basic metrics
// var (
// 	nallocs int
// 	maxheap int
// )

// pawn weight scale for cost calculation
var weights = []cost{'A': 1, 'B': 10, 'C': 100, 'D': 1000}

// buro routines

func (b *buro) String() string {
	i := 1 + bytes.LastIndex(b[:], []byte{'#'})
	return string(b[:i])
}

func (b *buro) get(j, i int) byte {
	return b[j*BLEN+i]
}

func (b *buro) set(j, i int, c byte) {
	b[j*BLEN+i] = c
}

// home index of home(a) in i space
func (b *buro) home(a byte) int {
	switch {
	case ispawn(a):
		return int(2*(a-'A') + 3)
	default:
		return 0
	}
}

// heavy lift peek()/pop()
func (b *buro) popx(i int, pop bool) (byte, cost) {
	var j int
	for j = 1; j < RLEN; j++ {
		x := b.get(j, i)
		switch {
		case ispawn(x):
			if pop {
				b.set(j, i, '.')
			}
			return x, cost(j - 1)
		case x == '#': // buro bottom row
			return '.', 0
		}
	}

	panic("unreachable")
}

// peek at buro.room(i) top element
func (b *buro) peek(i int) (x byte) {
	x, _ = b.popx(i, false)
	return
}

// pop buro.room(i) top element
func (b *buro) pop(i int) (byte, cost) {
	return b.popx(i, true)
}

// obvious ispawn
func ispawn(a byte) bool {
	return ('A' <= a && a <= 'D')
}

// push pawn a to buro.room(i)
func (b *buro) push(i int, a byte) cost {
	var j int

	for j = 1; j < RLEN; j++ {
		if x := b.get(j, i); ispawn(x) || x == '#' {
			b.set(j-1, i, a)
			return cost(j - 2)
		}
	}

	panic("unreachable")
}

// obvious setrow
func (b *buro) setrow(j int, s string) {
	const ROWINIT = "             \n" // 13 spaces + '\n'
	safe := func(raw string) []byte {
		buf := []byte(ROWINIT)
		copy(buf, raw)
		return buf
	}

	low := j * BLEN
	max := low + BLEN
	copy(b[low:max:max], safe(s))
}

// obvious getrow
func (b *buro) getrow(j int) string {
	low := j * BLEN
	max := low + BLEN
	return string(b[low:max:max])
}

// buro maker routine for part 1&2
func mkburos(input *bufio.Scanner) []buro {
	var buros [2]buro // part 1&2

	for j := 0; input.Scan(); j++ {
		buros[0].setrow(j, input.Text())
	}

	buros[1] = buros[0]
	buros[1].setrow(3, "  #D#C#B#A#") // /!\ mandatory 2 spaces prefix
	buros[1].setrow(4, "  #D#B#A#C#")
	buros[1].setrow(5, buros[0].getrow(3))
	buros[1].setrow(6, "  #########")

	return buros[:]
}

// move helpers

// obvious isfull buro.room(i)
func (b *buro) isfull(i int) bool {
	j := 1
	if ishome(i) {
		j++
	}
	return b.get(j, i) != '.'
}

// isclear checks if hallway cells between s and t are free
func (b *buro) isclear(t, s int) bool {
	for i := min(t, s); i <= max(t, s); i++ {
		if i != s && ishall(i) && b.peek(i) != '.' {
			return false
		}
	}
	return true
}

// iscosy checks game rules
//   - room is cozy for pawn `a` only if it is home to it
//     and either empty or populated (even crowded) by homies
//   - an empty hallway is always cozy
func (b *buro) iscosy(i int, a byte) bool {
	if ishall(i) || i == b.home(a) {
		var j int
	VSCAN: // vertical scan room cells
		for j = 1; j < RLEN; j++ {
			x := b.get(j, i)
			switch {
			case x == '#': // buro bottom row
				return true
			case x != a && x != '.':
				break VSCAN
			}
		}
	}
	return false
}

// https://tinyurl.com/ycy4jwfm
//   - dead1 detects a deadlock in the middle section of buro
//   - dead2 detects a deadlock at either edge of buro
//   - either way is fatal
func (b *buro) isdead() bool {
	dead1 := func() bool {
		// eqz is true if x == y or x == 0 or y == 0
		eqz := func(x, y byte) bool {
			x -= '.'
			y -= '.'
			return (x-y)*x*y == 0
		}

		// #...D...B...#  D, B deadlock in the hallway
		// ###.#A#C#.###
		for i := 4; i < BLEN-5; i += 2 {
			x := b.peek(i)
			for ii := i + 2; ii < BLEN-5; ii += 2 {
				y := b.peek(ii)
				if !eqz(x, y) && b.home(x) > ii && b.home(y) < i {
					return true
				}
			}
		}
		return false
	}

	dead2 := func() bool {
		type edge struct {
			x   byte
			off int
		}

		// #C..A...D.B.#  C BC A deadlock on the left
		// ###B#.#.#A###  D AD B deadlock on the right
		//   #C#.#.#D#
		//   #########
		edges := []edge{
			{'A', -1}, // left
			{'D', +1}, // right
		}

		for _, e := range edges {
			x, off := e.x, e.off
			hx := b.home(x)

			if b.peek(hx-off) == x {
				nspace, nalien := 0, 0

			ESCAN: // scan halleway edge rooms
				for i := 1; i < 3; i++ {
					x := b.peek(hx + i*off)
					switch {
					case x == '.':
						nspace++
					case ispawn(x):
						break ESCAN
					}
				}

			VSCAN: // vertical scan home cells
				for j := 1; j < RLEN; j++ {
					x := b.get(j, hx)
					switch {
					case x != '.' && x != 'x':
						nalien++
					case x == '#':
						break VSCAN
					}
				}

				if nalien > nspace {
					return true
				}
			}
		}
		return false
	}

	return dead1() || dead2()
}

// A* heuristic cost is entropy S as a disorder value:
// for a buro, it is the sum  of every misplaced pawn
// (weighted) distance to home without accounting for
// collisions.
//
// It has features we can profit for A*:
//  1. it is *admissible* (never overestimates goal cost)
//  2. it is *consistent* (never overestimates move cost)
//  3. it is zero (highest piority) for goals by design
//
// see properties:
// https://en.wikipedia.org/wiki/A*_search_algorithm
func (b *buro) hcost() cost {
	var S cost             // entropy (disorder)
	popcnts := [BLEN]int{} // home population counts

	for ii := range b[BLEN:] {
		if !ispawn(b[ii]) {
			continue
		}

		j, s := (ii / BLEN), ii%BLEN // depth, source room
		t := b.home(b[ii])           // target is home
		popcnts[t]++                 // account for homecoming

		if s == t { // already home, no cost
			continue
		}

		var manh int = abs(t - s) // home dist
		if ishome(s) {            // hallway dist
			manh += j - 1
		}
		manh += popcnts[t] // home cell dist

		S += cost(manh) * weights[b[ii]] // add weighted total dist
	}
	return S
}

// move type

// newMove is a move constructor
//
// it takes a buro (board) and an initial cost and returns a move object
func newMove(b *buro, c cost) *move {
	return &move{
		b: b, c: c,
		S: 0, // lazy S, computed only if selected by A*
	}
}

func (m move) String() string {
	var sb *strings.Builder
	sb.WriteString(m.b.String())
	sb.WriteString(fmt.Sprintf("   @%p c: %d, S: %d", &m, m.c, m.S))
	return sb.String()
}

func ishome(i int) bool { return i&1 == 1 && 2 < i && i < BLEN-4 }
func ishall(i int) bool { return !ishome(i) }

// move a pawn from t to s
//   - it returns a move and an ok bool
//   - if inplace, this move is m modified
//   - otherwise it is a modified clone of m (allocation)
//   - on success ok is true
//   - on failure it returns m unmodified and ok is false
func (m *move) move(t, s int, inplace bool) (*move, bool) {
	b := m.b

	islegit := func(t, s int) bool {
		// game rules:
		//   - x is a pawn, moving from one place to another if the path is clear
		//   - iff x is not cosy at home and is moving to a spacious and cosy place
		type rule func(byte) bool
		rules := []rule{
			// extra rule: forbid hallway to hallway moves
			// func(x byte) bool {
			// 	return !(ishall(t) && ishall(s))
			// },
			func(x byte) bool {
				return ispawn(x) && t != s && b.isclear(t, s)
			},
			func(x byte) bool {
				return !(s == b.home(x) && b.iscosy(s, x)) && (!b.isfull(t) && b.iscosy(t, x))
			},
		}

		match := func(x byte) bool {
			// rules ander
			for _, r := range rules {
				if !r(x) {
					return false // no match for rule r
				}
			}
			return true
		}

		return match(b.peek(s))
	}

	if islegit(t, s) {
		// move!

		nxt := b
		if !inplace {
			// nallocs++ // uncomment for basic metrics
			buf := *b // clone
			nxt = &buf
		}

		x, cs := nxt.pop(s)
		ct := nxt.push(t, x)
		manh := (cs + cost(abs(t-s)) + ct) * weights[x]

		return newMove(nxt, m.c+manh), true
	}
	return m, false
}

var MBUF [32]*move // static move buffer

// moves generates all legal moves from the current move
func (m *move) moves() []*move {
	if m.b.isdead() { // deadlocked board
		return []*move{} // no move
	}

	// step1 - go back home
	// always an abolute move, make all such moves at once!
	var cur *move
	nxt, done := m, false
MSCAN: // scan move
	for !done { // find them all
		cur = nxt
		done = true

		for s := 1; s < BLEN-2; s++ {
			x := m.b.peek(s)
			if h := m.b.home(x); h != 0 {
				ok := false
				if nxt, ok = cur.move(h, s, true); ok { // no alloc
					// new homecoming, restart scan!!
					done = false
					continue MSCAN
				}
			}
		}
	}

	if cur.S == 0 { // entropy is 0, cur is goal!
		return []*move{cur} // return this absolute winning move
	}

	// step2 - move out from other's home to hallway
	moves, ok := MBUF[:0], false     // reset
	for s := 3; s < BLEN-4; s += 2 { // home index
		for t := 1; t < BLEN-2; t++ { // room index
			if ishall(t) { // filter hallway, move() is slower to reject
				if nxt, ok = cur.move(t, s, false); ok {
					moves = append(moves, nxt)
				}
			}
		}
	}

	return moves
}

// set heuristic cost as piority component used by A*
// see func (*buro).hcost()
func (m *move) setprio() *move {
	m.S = m.b.hcost() // compute S when selected by A*
	return m
}

// A* move priority is the sum of the move cost and
// the resulting board entropy see func (*buro).hcost()
func (m *move) prio() cost {
	return m.c + m.S
}

// canonical A* algorithm
// https://en.wikipedia.org/wiki/A*_search_algorithm
func (m *move) solve() cost {
	const (
		// tuned hints, primes have no impact whatsoever
		// on container/heap performance
		MAXALLOC = 47_981 // tune here
		MAXHEAP  = 7_993  // tune here
	)

	costs := make(map[buro]cost, MAXALLOC)
	// costs[*m.b] = m.c  // correct but useless

	// uncomment for winning game moves:
	// start := m
	// from := make(map[buro]move)

	heap := make(heap, 0, MAXHEAP)
	hp.Init(&heap)
	hp.Push(&heap, m.setprio()) // from m as start...
	for heap.Len() > 0 {
		// get prioritized move
		m := hp.Pop(&heap).(*move) // shadow m!

		if m.S == 0 { // entropy is zero, goal!

			// uncomment and fix to print:
			// basic metrics:
			// fmt.Println("ncosts =", len(costs))
			// fmt.Println("nallocs =", nallocs)
			// fmt.Println("maxheap =", maxheap)
			// nallocs, maxheap = 0, 0

			// winning game moves:
			// for x := m; x.b != start.b; x = from[*x.b] {
			// 	fmt.Println(x)
			// }

			return m.c // cost is minimal by design
		}

		// gen moves
		for _, x := range m.moves() {
			if known, seen := costs[*x.b]; !seen || known > x.c {
				// from[*x.b] = m
				costs[*x.b] = x.c
				hp.Push(&heap, x.setprio()) // prioritize next move
			}
		}
	}

	panic("unreachable")
}

// prority queue concrete type and interface
//
// https://en.wikipedia.org/wiki/Priority_queue
// Insert is (heap).Push(), Pull is (heap).Pop()
//
// https://pkg.go.dev/container/heap#Interface
type heap []*move

// sort interface
func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool {
	return h[i].prio() < h[j].prio()
}

func (h heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// heap interface
func (h *heap) Pop() interface{} {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], nil
	return pop
}

func (h *heap) Push(x interface{}) {
	// maxheap = max(maxheap, len(*h)) // uncomment for basic metrics
	*h = append(*h, x.(*move))
}

// helpers
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// var DEBUG = false
//
// func debug(a ...any) {
// 	if DEBUG {
// 		fmt.Println(a...)
// 	}
// }
//
// func debugf(format string, a ...any) {
// 	if DEBUG {
// 		fmt.Printf(format, a...)
// 	}
// }

// goodies:
//
//	you can see the bench of this non-crypto hashing function, it is efficient,
//	faster than stdlib and neat but somehow the overall result is slower when
//	using it
//
// https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
// FNV-1
func (b *buro) hash() (h uint64) {
	const (
		o = 0xcbf29ce484222325 // fnv_offset_basis
		p = 0x100000001b3      // fnv_prime
	)

	h = o
	for ii := range b {
		h *= p
		h ^= uint64(b[ii])
	}
	return
}
