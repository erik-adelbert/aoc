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
	gol := newGame(bufio.NewScanner(os.Stdin))

	var clock int
	for gol.tick() {
		clock++

		// part 1
		if clock == 10 {
			fmt.Println(gol.crop() - gol.popcnt())
		}
	}
	// part 2
	fmt.Println(clock)
}

const NROWS = 150

type golife struct {
	cells []uint256
	heads []byte // round-robin heading order "NSWE"
	H     int    // height
}

func newGame(input *bufio.Scanner) (g *golife) {
	var (
		H, i int
		c    rune
	)

	cells := make([]uint256, NROWS)
	for H = 0; input.Scan(); H++ {
		for i, c = range input.Text() {
			if c == '#' {
				cells[H].setbit(i) // cells[H][i] is alive
			}
		}
	}
	return &golife{cells, []byte("NSWE"), H}
}

func (g *golife) crop() int {
	H, cells := g.H, g.cells

	var min int
	for min = range cells[:H] {
		if !cells[min].isZero() {
			break
		}
	}

	for H = range cells[min:] {
		if cells[H].isZero() {
			break
		}
	}
	g.H = H

	rot(cells, NROWS-min)

	mask := zero256.or(cells[:H]...)

	lead0, trail0 := mask.lead0(), mask.trail0()
	for j := range cells[:H] {
		cells[j] = cells[j].rsh(trail0)
	}
	W := uint256size - (trail0 + lead0)

	return W * H
}

func (g *golife) extend() {
	H, cells := g.H, g.cells

	// ensure empty first column
	for j := range cells[:H] {
		if cells[j].trail0() == 0 { // little-endian
			for j := range cells[:H] {
				cells[j] = cells[j].lsh(1)
			}
			break
		}
	}

	// ensure empty first row
	if !cells[0].isZero() {
		rot(cells, 1)
		H++
	}

	// ensure empty last two rows
	for !cells[H-2].or(cells[H-1]).isZero() {
		H++
	}

	g.H = H
	return
}

func (g *golife) tick() bool {
	g.extend()

	var old, cur, nxt plan

	H, cells := g.H, g.cells
	alive := false

	cur = g.plan(cells[0], cells[1], cells[2])

	const (
		Left  = -1
		None  = +0
		Right = +1
	)

	untie := func(tied uint256, planned, unmoved *uint256, shift int) {
		still := tied.and(*planned)
		*planned = planned.andnot(still)
		switch shift {
		case Left:
			still = still.lsh(1)
		case Right:
			still = still.rsh(1)
		}
		*unmoved = unmoved.or(still)
	}

	cells[0] = cur.n
	for i := range cells[:H-3] {
		// i:old, j:cur, k:nxt
		j, k := i+1, i+2

		nxt = g.plan(cells[k-1], cells[k], cells[k+1])
		tied := old.s.and(nxt.n)
		untie(tied, &old.s, &cells[i], None)
		untie(tied, &nxt.n, &nxt.o, None)

		tied = cur.w.and(cur.e)
		untie(tied, &cur.w, &cur.o, Left)
		untie(tied, &cur.e, &cur.o, Right)

		moved := old.s.or(cur.w, cur.e, nxt.n)
		if !moved.isZero() {
			alive = true
		}

		cells[j] = cur.o.or(moved)
		old.s, cur = cur.s, nxt
	}
	cells[H-2] = cells[H-2].or(old.s)
	rot(g.heads, 3)

	return alive
}

type plan struct {
	o, n, s, w, e uint256
}

func (g *golife) plan(north, cur, south uint256) plan {
	if cur.isZero() {
		return plan{}
	}

	var n, s, w, e uint256

	N := north.or(
		north.lsh(1), north.rsh(1),
	)
	S := south.or(south.lsh(1), south.rsh(1))
	W, E := cur.lsh(1), cur.rsh(1)

	ok := not(not(cur).or(N, S, W, E))
	nok := cur.andnot(ok)

	for _, d := range g.heads {
		switch d {
		case 'N':
			n = nok.andnot(N)
			nok = nok.andnot(n)
		case 'S':
			s = nok.andnot(S)
			nok = nok.andnot(s)
		case 'W':
			w = not(nok).or(
				W, north.lsh(1), south.lsh(1),
			)
			nok, w = nok.and(w), not(w).rsh(1)
		case 'E':
			e = not(nok).or(
				E, north.rsh(1), south.rsh(1),
			)
			nok, e = nok.and(e), not(e).lsh(1)
		}
	}
	return plan{ok.or(nok), n, s, w, e}
}

func (g *golife) popcnt() int {
	count, cells, H := 0, g.cells, g.H
	for j := range cells[:H] {
		count += cells[j].popcnt()
	}
	return count
}

func (g *golife) String() string {
	var sb strings.Builder
	for j := range g.cells[:g.H] {
		fmt.Fprintln(&sb, g.cells[j])
	}
	return sb.String()
}

const uint256size = 256

type uint256 struct {
	w0, w1, w2, w3 uint64
}

var zero256 = uint256{0, 0, 0, 0}

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

func (u uint256) popcnt() (n int) {
	n += bits.OnesCount64(u.w3)
	n += bits.OnesCount64(u.w2)
	n += bits.OnesCount64(u.w1)
	n += bits.OnesCount64(u.w0)
	return
}

func (u uint256) lead0() (n int) {
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

func (u uint256) and(m uint256) uint256 {
	u.w3 &= m.w3
	u.w2 &= m.w2
	u.w1 &= m.w1
	u.w0 &= m.w0
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

func rot[V uint256 | byte](a []V, n int) {
	if n = n % len(a); n != 0 {
		copy(a, append(a[len(a)-n:], a[:len(a)-n]...))
	}
}

const DEBUG = false

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}
