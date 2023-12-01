// aoc23.go --
// advent of code 2022 day 23
//
// https://adventofcode.com/2022/day/23
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-23: initial commit
// 2023-11-23: improve readability

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	clock, gol := 1, newGame(bufio.NewScanner(os.Stdin))
	//fmt.Println(gol)

	// part 1
	for clock < 11 {
		clock++
		gol.tick()
	}
	area, popcnt := gol.poparea()
	fmt.Println(area - popcnt)

	// part 2
	for gol.tick() {
		clock++
	}
	fmt.Println(clock)
}

type golife struct {
	cells, n, s, w, e []uint256
	head              int // heading
}

func newGame(input *bufio.Scanner) (g *golife) {
	const (
		off    = 74
		heigth = 222
	)

	n := make([]uint256, heigth)
	s := make([]uint256, heigth)
	w := make([]uint256, heigth)
	e := make([]uint256, heigth)

	cells := make([]uint256, heigth)
	for j := 0; input.Scan(); j++ {
		for i, c := range input.Text() {
			if c == '#' {
				cells[j+off].setbit(i + off) // cells[j][i] is alive
			}
		}
	}

	return &golife{cells, n, s, w, e, 0}
}

func (g *golife) poparea() (area int, popcnt int) {
	b, popcnt := g.bbox()
	return b.area(), popcnt
}

type AABB struct {
	ymin, ymax, xmin, xmax int
}

func (a AABB) area() int {
	return (a.ymax - a.ymin) * (a.xmax - a.xmin)
}

func (g *golife) bbox() (box AABB, popcnt int) {
	var mask uint256

	cells := g.cells
	box.ymin = len(cells)

	for j := range cells {
		if n := cells[j].popcnt(); n > 0 {
			box.ymin = min(box.ymin, j)
			box.ymax = max(box.ymax, j)
			mask = mask.or(cells[j])
			popcnt += n
		}
	}
	box.ymax++

	box.xmin, box.xmax = mask.trail0(), uint256size-mask.lead0()
	return
}

func (g *golife) tick() (alive bool) {
	// sugars
	cells, head := g.cells, g.head
	n, s, w, e := g.n, g.s, g.w, g.e

	b, _ := g.bbox()
	min, max := b.ymin-1, b.ymax+2
	alive = false

	var old, cur, nxt uint256

	cur, nxt = not(cur), not(nxt)
	for j := min; j < max; j++ {
		// roll 3 lines window
		old, cur = cur, nxt
		nxt = not(cells[j+1].or(cells[j+1].lsh(1), cells[j+1].rsh(1)))

		// plan moves
		v := not(zero256.or(cells[j-1 : j+2]...)) // vertical

		u, d := old, nxt                         // up, down
		l, r := v.lsh(1), v.rsh(1)               // left, right
		still := cells[j].andnot(u.and(d, l, r)) // mask not moving baseline

		heading := []*uint256{&u, &d, &l, &r}
		for i := range "NSWE" {
			x := heading[(head+i)&0x3] // left rotate N,S,W,E and select u,d,l,r accordingly
			*x = x.and(still)
			still = still.andnot(*x)
		}

		// store planned moves
		n[j-1] = u
		s[j+1] = d
		w[j] = l.rsh(1)
		e[j] = r.lsh(1)
	}

	// cancel moves ending up in the same cell
	for j := min; j < max; j++ {
		// alias nswe to up, down, left, right
		u, d, l, r := n[j], s[j], w[j], e[j]

		// cancel vertical moves
		n[j] = n[j].andnot(d)
		s[j] = s[j].andnot(u)

		// cancel horizontal moves
		w[j] = w[j].andnot(r)
		e[j] = e[j].andnot(l)
	}

	// make all moves
	for j := min; j < max; j++ {
		// not moving
		still := cells[j].andnot(
			n[j-1].or(s[j+1], w[j].lsh(1), e[j].rsh(1)),
		)
		// moving
		moved := n[j].or(s[j], w[j], e[j])

		// move!
		cells[j] = still.or(moved)
		alive = alive || !moved.isZero()
	}
	g.head++

	return alive
}

func (g *golife) String() string {
	var sb strings.Builder

	cells := g.cells
	b, popcnt := g.bbox()

	fmt.Fprintf(
		&sb, "head: %d, box: %v, pop:%d\n", g.head, b, popcnt,
	)
	for j := range cells[b.ymin:b.ymax] {
		fmt.Fprintf(&sb, "%03d: ", j)

		x := cells[b.ymin+j].rsh(b.xmin)
		for i := b.xmin; i < b.xmax; i++ {
			sb.WriteByte(".#"[x.w0&1])
			x = x.rsh(1)
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

const uint256size = 256

var zero256 uint256

type uint256 struct {
	w0, w1, w2, w3 uint64
}

// setbit sets bit n-th n = 0 is LSB.
// n must be <= 255.
func (u *uint256) setbit(n int) {
	switch n >> 6 {
	case 3:
		u.w3 |= (1 << (n & 0x3f))
	case 2:
		u.w2 |= (1 << (n & 0x3f))
	case 1:
		u.w1 |= (1 << (n & 0x3f))
	case 0:
		u.w0 |= (1 << (n & 0x3f))
	}
}

func (u *uint256) popcnt() (n int) {
	n += bits.OnesCount64(u.w3)
	n += bits.OnesCount64(u.w2)
	n += bits.OnesCount64(u.w1)
	n += bits.OnesCount64(u.w0)
	return
}

func (u *uint256) lead0() (n int) {
	if n = bits.LeadingZeros64(u.w3); n != 64 {
		return
	}
	if n += bits.LeadingZeros64(u.w2); n != 128 {
		return
	}
	if n += bits.LeadingZeros64(u.w1); n != 192 {
		return
	}
	n += bits.LeadingZeros64(u.w0)
	return
}

func (u *uint256) trail0() (n int) {
	if n = bits.TrailingZeros64(u.w0); n != 64 {
		return
	}
	if n += bits.TrailingZeros64(u.w1); n != 128 {
		return
	}
	if n += bits.TrailingZeros64(u.w2); n != 192 {
		return
	}
	return n + bits.TrailingZeros64(u.w3)
}

// isZero returns true if u == 0
func (u uint256) isZero() bool {
	return (u.w3 | u.w2 | u.w1 | u.w0) == 0
}

func (u uint256) lsh(n int) uint256 {
	var a, b uint64

	switch {
	case n > 256:
		return uint256{}
	case n > 192:
		u.w3, u.w2, u.w1, u.w0 = u.w0, 0, 0, 0
		n -= 192
		goto sh192
	case n > 128:
		u.w3, u.w2, u.w1, u.w0 = u.w1, u.w0, 0, 0
		n -= 128
		goto sh128
	case n > 64:
		u.w3, u.w2, u.w1, u.w0 = u.w2, u.w1, u.w0, 0
		n -= 64
		goto sh64
	}

	// remaining shifts
	a = u.w0 >> (64 - n)
	u.w0 = u.w0 << n

sh64:
	b = u.w1 >> (64 - n)
	u.w1 = (u.w1 << n) | a

sh128:
	a = u.w2 >> (64 - n)
	u.w2 = (u.w2 << n) | b

sh192:
	u.w3 = (u.w3 << n) | a

	return u
}

func (u uint256) rsh(n int) uint256 {
	var a, b uint64

	switch {
	case n > 256:
		return uint256{}
	case n > 192:
		u.w3, u.w2, u.w1, u.w0 = 0, 0, 0, u.w3
		n -= 192
		goto sh192
	case n > 128:
		u.w3, u.w2, u.w1, u.w0 = 0, 0, u.w3, u.w2
		n -= 128
		goto sh128
	case n > 64:
		u.w3, u.w2, u.w1, u.w0 = 0, u.w3, u.w2, u.w1
		n -= 64
		goto sh64
	}

	// remaining shifts
	a = u.w3 << (64 - n)
	u.w3 = u.w3 >> n

sh64:
	b = u.w2 << (64 - n)
	u.w2 = (u.w2 >> n) | a

sh128:
	a = u.w1 << (64 - n)
	u.w1 = (u.w1 >> n) | b

sh192:
	u.w0 = (u.w0 >> n) | a

	return u
}

func (u uint256) and(m ...uint256) uint256 {
	for i := range m {
		u.w3 &= m[i].w3
		u.w2 &= m[i].w2
		u.w1 &= m[i].w1
		u.w0 &= m[i].w0
	}
	return u
}

func (u uint256) andnot(m uint256) uint256 {
	u.w3 &= ^m.w3
	u.w2 &= ^m.w2
	u.w1 &= ^m.w1
	u.w0 &= ^m.w0
	return u
}

func not(u uint256) uint256 {
	u.w3 = ^u.w3
	u.w2 = ^u.w2
	u.w1 = ^u.w1
	u.w0 = ^u.w0
	return u
}

func (u uint256) or(m ...uint256) uint256 {
	for i := range m {
		u.w3 |= m[i].w3
		u.w2 |= m[i].w2
		u.w1 |= m[i].w1
		u.w0 |= m[i].w0
	}
	return u
}

func (u uint256) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%016x%016x%016x%016x", u.w3, u.w2, u.w1, u.w0)
	return sb.String()
}

func rot(a []uint256, n int) {
	if n = n % len(a); n != 0 {
		copy(a, append(a[len(a)-n:], a[:len(a)-n]...))
	}
}

const DEBUG = true

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}
