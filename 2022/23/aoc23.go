// aoc23.go --
// advent of code 2022 day 23
//
// https://adventofcode.com/2022/day/23
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-23: initial commit

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

type grid struct {
	cell []uint256
	dseq []byte
	life bool
	h    int
	w    int
}

func main() {
	g := newGrid(bufio.NewScanner(os.Stdin))

	var i int
	for i = 0; g.life; i++ {
		// part 1
		if i == 10 {
			g.crop()
			fmt.Println(g.h*g.w - g.popcnt())
		}
		g.evolve()
	}

	// part 2
	fmt.Println(i)
}

func newGrid(input *bufio.Scanner) (g *grid) {
	g = new(grid)
	g.dseq = []byte{'N', 'S', 'W', 'E'}

	g.cell = make([]uint256, NROWS)
	for j := 0; input.Scan(); j++ {
		for i, c := range input.Bytes() {
			if c == '#' {
				g.cell[j] = g.cell[j].set(i)
				g.h, g.w = max(g.h, j+1), max(g.w, i+1)
			}
		}
	}

	g.life = true
	return
}

func (g *grid) evolve() {
	g.extend()
	g.move()
	return
}

func (g *grid) move() {
	g.life = false

	var old, cur, nxt plan
	cur = g.plan(g.cell[0], g.cell[1], g.cell[2])

	untie := func(tied uint256, planed, unmoved *uint256, mode int) {
		still := tied.and(*planed)
		*planed = planed.and(still.not())
		switch mode {
		case -1:
			still = still.lsh(1)
		case 1:
			still = still.rsh(1)
		}
		*unmoved = unmoved.or(still)
	}

	g.cell[0] = cur.n
	for i := 0; i < g.h-3; i++ {
		// i:old, j:cur, k:nxt
		j, k := i+1, i+2

		nxt = g.plan(g.cell[k-1], g.cell[k], g.cell[k+1])
		tied := old.s.and(nxt.n)
		untie(tied, &old.s, &g.cell[i], 0)
		untie(tied, &nxt.n, &nxt.o, 0)

		tied = cur.w.and(cur.e)
		untie(tied, &cur.w, &cur.o, -1)
		untie(tied, &cur.e, &cur.o, +1)

		moved := old.s.or(cur.w).or(cur.e).or(nxt.n)
		if !moved.iszero() {
			g.life = true
		}

		g.cell[j] = cur.o.or(moved)
		old.s = cur.s
		cur = nxt
	}
	g.cell[g.h-2] = g.cell[g.h-2].or(old.s)
	rot(g.dseq, 3)
}

type plan struct {
	o, n, s, w, e uint256
}

func (p plan) String() string {
	return fmt.Sprintf(
		"{o: %v, n: %v, s: %v, w: %v, e: %v}",
		p.o, p.n, p.s, p.w, p.e,
	)
}

func (g *grid) plan(n, cur, s uint256) plan {
	var p plan
	if cur.iszero() {
		return p
	}

	W, E := cur.lsh(1), cur.rsh(1)
	N := n.or(n.lsh(1)).or(n.rsh(1))
	S := s.or(s.lsh(1)).or(s.rsh(1))
	ok := cur.not().or(W).or(E).or(N).or(S).not()
	nok := cur.and(ok.not())
	for _, d := range g.dseq {
		switch d {
		case 'N':
			p.n = nok.and(N.not())
			nok = nok.and(p.n.not())
		case 'S':
			p.s = nok.and(S.not())
			nok = nok.and(p.s.not())
		case 'W':
			w := nok.not().or(W).or(n.lsh(1)).or(s.lsh(1)).not()
			p.w = w.rsh(1)
			nok = nok.and(w.not())
		case 'E':
			e := nok.not().or(E).or(n.rsh(1)).or(s.rsh(1)).not()
			p.e = e.lsh(1)
			nok = nok.and(e.not())
		}
	}
	p.o = ok.or(nok)
	return p
}

func (g *grid) crop() {
	var jmin, jmax int
	for jmin = 0; jmin < g.h; jmin++ {
		if !g.cell[jmin].iszero() {
			break
		}
	}

	for jmax = g.h - 1; jmax > 0; jmax-- {
		if !g.cell[jmax].iszero() {
			break
		}
	}

	rot(g.cell, NROWS-jmin)

	g.h = jmax - jmin + 1

	mask := zero256
	for j := 0; j < g.h; j++ {
		mask = mask.or(g.cell[j])
	}
	n := mask.trail0()
	for j := 0; j < g.h; j++ {
		g.cell[j] = g.cell[j].rsh(n)
	}
	g.w = uint256size - n - mask.lead0()
}

func (g *grid) extend() {
	// ensure empty first column
	for j := 0; j < g.h; j++ {
		if g.cell[j].trail0() == 0 {
			for j = 0; j < g.h; j++ {
				g.cell[j] = g.cell[j].lsh(1)
			}
			g.w++
			break
		}
	}

	// ensure empty first row
	if !g.cell[0].iszero() {
		rot(g.cell, 1)
		g.h++
	}

	// ensure empty last two rows
	for !g.cell[g.h-2].iszero() || !g.cell[g.h-1].iszero() {
		g.h++
	}
}

func (g *grid) String() string {
	var sb strings.Builder
	for j := 0; j < g.h; j++ {
		fmt.Fprintln(&sb, g.cell[j])
	}
	return sb.String()
}

const (
	NROWS = 256
)

func (g *grid) popcnt() int {
	pop := 0
	for j := 0; j < g.h; j++ {
		pop += g.cell[j].popcnt()
	}
	return pop
}

const uint128size = 128

type uint128 struct {
	hi, lo uint64
}

var (
	zero128 = uint128{}
	one128  = uint128{0, 1}
)

func (u uint128) iszero() bool {
	return u.hi|u.lo == 0
}

func (u uint128) popcnt() int {
	count := bits.OnesCount64
	return count(u.hi) + count(u.lo)
}

func (u uint128) lead0() int {
	lead0 := bits.LeadingZeros64
	n := lead0(u.hi)
	if n < 64 {
		return n
	}
	return n + lead0(u.lo)
}

func (u uint128) trail0() int {
	trail0 := bits.TrailingZeros64
	n := trail0(u.lo)
	if n < 64 {
		return n
	}
	return n + trail0(u.hi)
}

func (u uint128) lsh(n int) uint128 {
	if n >= 64 {
		return uint128{u.lo << (n - 64), 0}
	}
	return uint128{u.hi<<n | u.lo>>(64-n), u.lo << n}
}

func (u uint128) rsh(n int) uint128 {
	if n >= 64 {
		return uint128{0, u.hi >> (n - 64)}
	}
	return uint128{u.hi >> n, u.lo>>n | u.hi<<(64-n)}
}

func (u uint128) not() uint128 { return uint128{^u.hi, ^u.lo} }

func (u uint128) and(m uint128) uint128 {
	return uint128{u.hi & m.hi, u.lo & m.lo}
}

func (u uint128) or(m uint128) uint128 {
	return uint128{u.hi | m.hi, u.lo | m.lo}
}

func (u uint128) get(n int) bool {
	if n >= 64 {
		x := uint64(1 << (n - 64))
		return u.hi&x == x
	}
	x := uint64(1 << n)
	return u.lo&x == x
}

func (u uint128) set(n int) uint128 {
	if n >= 64 {
		u.hi |= 1 << (n - 64)
		return u
	}
	u.lo |= (1 << n)
	return u
}

func (u uint128) String() string {
	var sb strings.Builder
	if u.hi != 0 {
		fmt.Fprintf(&sb, "%x%016x", u.hi, u.lo)
	} else {
		fmt.Fprintf(&sb, "%x", u.lo)
	}
	return sb.String()
}

const uint256size = 256

type uint256 struct {
	hi, lo uint128
}

var (
	zero256 = uint256{}
	one256  = uint256{zero128, one128}
)

func (u uint256) iszero() bool {
	return u.hi.or(u.lo).iszero()
}

func (u uint256) get(n int) bool {
	if n >= 128 {
		return u.hi.get(n - 128)
	}
	return u.lo.get(n)
}

func (u uint256) set(n int) uint256 {
	if n >= 128 {
		u.hi = u.hi.set(n - 128)
		return u
	}
	u.lo = u.lo.set(n)
	return u
}

func (u uint256) popcnt() int {
	return u.hi.popcnt() + u.lo.popcnt()
}

func (u uint256) lead0() int {
	n := u.hi.lead0()
	if n < 128 {
		return n
	}
	return n + u.lo.lead0()
}

func (u uint256) trail0() int {
	n := u.lo.trail0()
	if n < 128 {
		return n
	}
	return n + u.hi.trail0()
}

func (u uint256) lsh(n int) uint256 {
	if n >= 128 {
		return uint256{u.lo.lsh(n - 128), zero128}
	}
	return uint256{u.hi.lsh(n).or(u.lo.rsh(128 - n)), u.lo.lsh(n)}
}

func (a uint256) rsh(n int) uint256 {
	if n >= 128 {
		return uint256{zero128, a.hi.rsh(n - 128)}
	}
	return uint256{a.hi.rsh(n), a.lo.rsh(n).or(a.hi.lsh(128 - n))}
}

func (u uint256) not() uint256 {
	return uint256{u.hi.not(), u.lo.not()}
}

func (u uint256) and(m uint256) uint256 {
	return uint256{u.hi.and(m.hi), u.lo.and(m.lo)}
}

func (u uint256) or(m uint256) uint256 {
	return uint256{u.hi.or(m.hi), u.lo.or(m.lo)}
}

func (a uint256) String() string {
	var sb strings.Builder
	if a.hi.iszero() {
		fmt.Fprint(&sb, a.lo)
	} else {
		fmt.Fprintf(&sb, "%v%v", a.hi, a.lo)
	}
	return sb.String()
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

func rot[V uint256 | byte](a []V, n int) {
	n = n % len(a)
	if n != 0 {
		copy(a, append(a[len(a)-n:], a[:len(a)-n]...))
	}
}

const DEBUG = true

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}
