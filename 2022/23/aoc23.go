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
	g := new(grid)
	g.init(bufio.NewScanner(os.Stdin))

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
	g.crop()
	fmt.Println(i)
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
		sb.WriteString(g.cell[j].String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

const (
	NROWS = 256
)

func (g *grid) init(input *bufio.Scanner) {
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
}

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

func (a uint128) iszero() bool {
	return a == uint128{}
}

// func (a uint128) add(b uint128) uint128 {
// 	c, lo := bits.Add64(a.lo, b.lo, 0)
// 	return uint128{a.hi + b.hi + c, lo}
// }

// func (a uint128) sub(b uint128) uint128 {
// 	c, hi := bits.Sub64(a.hi, b.hi, 0)
// 	return uint128{hi, a.lo - b.lo - c}
// }

func (a uint128) popcnt() int {
	return bits.OnesCount64(a.hi) + bits.OnesCount64(a.lo)
}

func (a uint128) lead0() int {
	n := bits.LeadingZeros64(a.hi)
	if n < 64 {
		return n
	}
	return n + bits.LeadingZeros64(a.lo)
}

func (a uint128) trail0() int {
	n := bits.TrailingZeros64(a.lo)
	if n < 64 {
		return n
	}
	return n + bits.TrailingZeros64(a.hi)
}

func (a uint128) lsh(i int) uint128 {
	if i >= 64 {
		return uint128{
			a.lo << (i - 64),
			0,
		}
	}
	return uint128{
		a.hi<<i | a.lo>>(64-i),
		a.lo << i,
	}
}

func (a uint128) rsh(i int) uint128 {
	if i >= 64 {
		return uint128{
			0,
			a.hi >> (i - 64),
		}
	}
	return uint128{
		a.hi >> i,
		a.lo>>i | a.hi<<(64-i),
	}
}

func (a uint128) not() uint128 {
	return uint128{
		^a.hi,
		^a.lo,
	}
}

func (a uint128) and(b uint128) uint128 {
	return uint128{
		a.hi & b.hi,
		a.lo & b.lo,
	}
}

func (a uint128) or(b uint128) uint128 {
	return uint128{
		a.hi | b.hi,
		a.lo | b.lo,
	}
}

// func (a uint128) xor(b uint128) uint128 {
// 	return uint128{
// 		a.hi ^ b.hi,
// 		a.lo ^ b.lo,
// 	}
// }

func (a uint128) get(i int) bool {
	if i >= 64 {
		x := uint64(1 << (i - 64))
		return a.hi&x == x
	}
	x := uint64(1 << i)
	return a.lo&x == x
}

func (a uint128) set(i int) uint128 {
	if i >= 64 {
		a.hi |= 1 << (i - 64)
		return a
	}
	a.lo |= (1 << i)
	return a
}

// func (a uint128) clr(i int) uint128 {
// 	if i >= 64 {
// 		a.hi &= ^(1 << (i - 64))
// 		return a
// 	}
// 	a.lo &= ^(1 << i)
// 	return a
// }

// func (a uint128) flp(i int) uint128 {
// 	if i >= 64 {
// 		a.hi ^= (1 << (i - 64))
// 		return a
// 	}
// 	a.lo ^= 1 << i
// 	return a
// }

var u128print func(uint128) string = func(a uint128) string {
	if a.hi != 0 {
		return fmt.Sprintf("%x%016x", a.hi, a.lo)
	}
	return fmt.Sprintf("%x", a.lo)
}

func fmtu128(v string) {
	switch v {
	case "x":
		u128print = func(a uint128) string {
			if a.hi != 0 {
				return fmt.Sprintf("%x%016x", a.hi, a.lo)
			}
			return fmt.Sprintf("%x", a.lo)
		}
	case "0x":
		u128print = func(a uint128) string {
			return fmt.Sprintf("%016x%016x", a.hi, a.lo)
		}
	case "b":
		u128print = func(a uint128) string {
			if a.hi != 0 {
				return fmt.Sprintf("%b%064b", a.hi, a.lo)
			}
			return fmt.Sprintf("%064b", a.lo)
		}
	case "0b":
		u128print = func(a uint128) string {
			return fmt.Sprintf("%064b%064b", a.hi, a.lo)
		}
	}
}

func (a uint128) String() string {
	return u128print(a)
}

const uint256size = 256

type uint256 struct {
	hi, lo uint128
}

var (
	zero256 = uint256{}
	one256  = uint256{zero128, one128}
)

func (a uint256) iszero() bool {
	return a == uint256{}
}

func (a uint256) get(i int) bool {
	if i >= 128 {
		return a.hi.get(i - 128)
	}
	return a.lo.get(i)
}

func (a uint256) set(i int) uint256 {
	if i >= 128 {
		a.hi = a.hi.set(i - 128)
		return a
	}
	a.lo = a.lo.set(i)
	return a
}

// func (a uint256) clr(i int) uint256 {
// 	if i >= 128 {
// 		a.hi = a.hi.clr(i - 128)
// 		return a
// 	}
// 	a.lo = a.lo.clr(i)
// 	return a
// }

// func (a uint256) flp(i int) uint256 {
// 	if i >= 128 {
// 		a.hi = a.hi.flp(i - 128)
// 		return a
// 	}
// 	a.lo = a.lo.flp(i)
// 	return a
// }

func (a uint256) popcnt() int {
	return a.hi.popcnt() + a.lo.popcnt()
}

func (a uint256) lead0() int {
	n := a.hi.lead0()
	if n < 128 {
		return n
	}
	return n + a.lo.lead0()
}

func (a uint256) trail0() int {
	n := a.lo.trail0()
	if n < 128 {
		return n
	}
	return n + a.hi.trail0()
}

func (a uint256) lsh(i int) uint256 {
	if i >= 128 {
		return uint256{
			a.lo.lsh(i - 128),
			uint128{},
		}
	}
	return uint256{
		a.hi.lsh(i).or(a.lo.rsh(128 - i)),
		a.lo.lsh(i),
	}
}

func (a uint256) rsh(i int) uint256 {
	if i >= 128 {
		return uint256{
			uint128{},
			a.hi.rsh(i - 128),
		}
	}
	return uint256{
		a.hi.rsh(i),
		a.lo.rsh(i).or(a.hi.lsh(128 - i)),
	}
}

func (a uint256) not() uint256 {
	return uint256{
		a.hi.not(),
		a.lo.not(),
	}
}

func (a uint256) and(b uint256) uint256 {
	return uint256{
		a.hi.and(b.hi),
		a.lo.and(b.lo),
	}
}

func (a uint256) or(b uint256) uint256 {
	return uint256{
		a.hi.or(b.hi),
		a.lo.or(b.lo),
	}
}

func (a uint256) String() string {
	if a.hi.iszero() {
		return u128print(a.lo)
	}
	return u128print(a.hi) + u128print(a.lo)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
