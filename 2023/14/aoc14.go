package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	p1, p2 := 0, 0
	w := newWorld()

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		w.readline(j, input.Text())
	}

	seen := make(map[uint64]int)

	// first cycle
	for _, d := range []int{South, East, North, West} {
		w.tilt(d)

		if d == South {
			p1 = w.load()
		}
	}
	seen[w.hash()] = 1

	i, len := 0, 0
CYCLE:
	for i = 2; ; i++ {
		w.cycle(0, 1)

		if x, ok := seen[w.hash()]; ok {
			len = i - x
			break CYCLE
		}

		seen[w.hash()] = i
	}

	end := (1_000_000_000 - i) % len
	p2 = w.cycle(0, end)

	fmt.Println(p1, p2)
}

/*
 * world type --
 ********************/
type world struct {
	walls   bitarray128
	state   bitarray128
	D, H, W int
}

const (
	East = iota
	South
	West
	North
)

func newWorld() (w *world) {
	w = new(world)
	return
}

func (w *world) cycle(start, end int) int {
	for i := start; i < end; i++ {
		for _, d := range []int{South, East, North, West} {
			w.tilt(d)
		}
	}
	return w.load()
}

func (w *world) face(dir int) {
	for w.D != dir {
		w.state.rotCW()
		w.walls.rotCW()
		w.H, w.W = w.W, w.H
		w.D = (w.D + 1) % 4
	}
}

// https://stackoverflow.com/a/12996028
func (w *world) hash() uint64 {
	j0, i0 := w.locate()

	seed := uint64(w.H)
	for j := j0; j < w.H+j0; j++ {
		x := w.state[j].rsh(i0).w0
		x = ((x >> 32) ^ x) * 0xD2E23944245D9F3B
		x = ((x >> 32) ^ x) * 0xD2E23944245D9F3B
		x = (x >> 32) ^ x
		seed = x + 0x3BBCD6C79E3779B9 + (seed << 6) + (seed >> 2)
	}
	return seed
}

func (w *world) load() int {
	old := w.D
	w.face(East)

	sum := 0
	j0, _ := w.locate()
	for j := j0; j < j0+w.H; j++ {
		sum += (int(w.H) + j0 - j) * w.state[j].popcnt()
	}

	w.face(old)
	return sum
}

func (w *world) locate() (int, int) {
	switch w.D {
	case East:
		return 0, 0
	case South:
		return uint128size - w.H, 0
	case West:
		return uint128size - w.H, uint128size - w.H
	case North:
		return 0, uint128size - w.H
	}

	panic("unreachable")
}

func (w *world) readline(j int, s string) {
	w.H = max(w.H, int(j+1))
	w.W = max(w.W, int(len(s)))
	for i := range s {
		switch s[i] {
		case '#':
			w.walls.set(j, i)
		case 'O':
			w.state.set(j, i)
		}
	}
}

func (w *world) tilt(dir int) {
	w.face(dir)

	j0, i0 := w.locate()
	for j := j0; j < j0+w.H; j++ {
		walls := w.walls[j]

		old, cur, done := i0-1, 0, zero128
		for !not(done).isZero() {
			var balls uint128
			if balls = w.state[j]; balls.isZero() { // no ball
				break
			}

			cur = walls.trail0()                                  // current obstacle
			mask := not(zero128).rsh(uint128size - cur).xor(done) // between old and cur obstacle

			balls = w.state[j].and(mask) // get balls in between old and cur
			n := balls.popcnt()

			w.state[j] = w.state[j].and(not(balls.or(mask))) // remove balls

			for i := 0; i < n; i++ {
				w.state[j] = w.state[j].set(old + 1 + i) // group balls on top of old
			}

			walls = walls.clear(cur)       // remove obstacle
			done, old = done.or(mask), cur // expand done mask over cur
		}

	}

	return
}

func (w *world) String() string {
	old := w.D
	w.face(East)

	var sb strings.Builder
	j0, i0 := w.locate()
	for j := j0; j < j0+w.H; j++ {
		for i := i0; i < i0+w.W; i++ {
			switch {
			case w.state[j].get(i) > 0:
				sb.WriteByte('O')
			case w.walls[j].get(i) > 0:
				sb.WriteByte('#')
			default:
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}
	fmt.Fprintf(&sb, "H: %d W: %d L:%d, H:%0x\n", w.H, w.W, w.load(), w.hash())

	w.face(old)
	return sb.String()
}

/*
 * bitarray128 type --
 ********************/

type bitarray128 [128]uint128

func (BA *bitarray128) isZero() bool {
	return BA[0].or(BA[1:]...).isZero()
}

func (BA *bitarray128) set(j, i int) *bitarray128 {
	BA[j] = BA[j].set(i)
	return BA
}

func (BA *bitarray128) clear(j, i int) *bitarray128 {
	BA[j] = BA[j].clear(i)
	return BA
}

func (BA *bitarray128) get(j, i int) int {
	return BA[j].get(i)
}

func (BA *bitarray128) trans() *bitarray128 {
	var L0, L1, H0, H1 bitarray64

	// split into 4 bitarray64
	// A bitarray128 = {
	// 		L0, H0,
	// 		L1, H1,
	// }
	for i := range BA {
		switch {
		case i < 64:
			L0[i], H0[i] = BA[i].w0, BA[i].w1
		default:
			L1[i-64], H1[i-64] = BA[i].w0, BA[i].w1
		}
	}

	// transpose all quarters
	for _, A := range []*bitarray64{&L0, &L1, &H0, &H1} {
		A.trans64()
	}

	// rebuild transposed BA
	// A bitarray128 = {
	// 		L0, L1,
	// 		H0, H1,
	// }
	for i := range BA {
		if i < 64 {
			BA[i] = uint128{L0[i], L1[i]}
		} else {
			BA[i] = uint128{H0[i-64], H1[i-64]}
		}
	}

	return BA
}

func (BA *bitarray128) rotCW() *bitarray128 {
	//transpose
	BA.trans()

	// mirror horizontally
	for i := range BA[:63] {
		BA[i], BA[127-i] = BA[127-i], BA[i]
	}

	return BA
}

func (BA *bitarray128) rotCCW() *bitarray128 {
	//transpose
	BA.trans()

	// mirror
	for i := range BA[:63] {
		BA[i], BA[127-i] = BA[127-i].reverse(), BA[i].reverse()
	}

	return BA
}

func (BA bitarray128) String() string {
	var sb strings.Builder
	for j := range BA {
		fmt.Fprintf(&sb, "%064b%064b\n", BA[j].w1, BA[j].w0)
	}
	return sb.String()
}

/*
 * uint128 type --
 ********************/

const uint128size = 128

var (
	zero128 uint128
	one128  = uint128{1, 0}
)

type uint128 struct {
	w0, w1 uint64
}

func (u uint128) isZero() bool {
	return (u.w1 | u.w0) == 0
}

// set sets bit n-th n = 0 is LSB.
// n must be <= 128.
func (u uint128) set(n int) uint128 {
	switch n >> 6 {
	case 1:
		u.w1 |= (1 << (n & 0x3f))
	case 0:
		u.w0 |= (1 << (n & 0x3f))
	}
	return u
}

func (u uint128) get(n int) int {
	x := u.rsh(n)
	return int(x.w0 & 1)
}

func (u uint128) clear(n int) uint128 {
	switch n >> 6 {
	case 1:
		u.w1 &= ^(1 << (n & 0x3f))
	case 0:
		u.w0 &= ^(1 << (n & 0x3f))
	}
	return u
}

func (u uint128) hi64() uint64 {
	return u.w1
}

func (u uint128) lo64() uint64 {
	return u.w0
}

func (u uint128) popcnt() (n int) {
	n += bits.OnesCount64(u.w1)
	n += bits.OnesCount64(u.w0)
	return
}

func (u uint128) lead0() (n int) {
	if n += bits.LeadingZeros64(u.w1); n != 64 {
		return
	}
	n += bits.LeadingZeros64(u.w0)
	return
}

func (u uint128) trail0() (n int) {
	if n = bits.TrailingZeros64(u.w0); n != 64 {
		return
	}
	return n + bits.TrailingZeros64(u.w1)
}

func (u uint128) reverse() uint128 {
	return uint128{bits.Reverse64(u.w1), bits.Reverse64(u.w0)}
}

func (u uint128) lsh(n int) uint128 {
	var a uint64

	switch {
	case n > 128:
		return uint128{}
	case n > 64:
		u.w1, u.w0 = u.w0, 0
		n -= 64
		goto sh64
	}

	// remaining shifts
	a = u.w0 >> (64 - n)
	u.w0 = u.w0 << n

sh64:
	u.w1 = (u.w1 << n) | a

	return u
}

func (u uint128) rsh(n int) uint128 {
	var a uint64

	switch {
	case n > 128:
		return uint128{}
	case n > 64:
		u.w1, u.w0 = 0, u.w1
		n -= 64
		goto sh64
	}

	// remaining shifts
	a = u.w1 << (64 - n)
	u.w1 = u.w1 >> n

sh64:
	u.w0 = (u.w0 >> n) | a

	return u
}

func not(u uint128) uint128 {
	u.w1 = ^u.w1
	u.w0 = ^u.w0
	return u
}

func (u uint128) and(m ...uint128) uint128 {
	for i := range m {
		u.w1 &= m[i].w1
		u.w0 &= m[i].w0
	}
	return u
}

func (u uint128) or(m ...uint128) uint128 {
	for i := range m {
		u.w1 |= m[i].w1
		u.w0 |= m[i].w0
	}
	return u
}

func (u uint128) xor(m ...uint128) uint128 {
	for i := range m {
		u.w1 ^= m[i].w1
		u.w0 ^= m[i].w0
	}
	return u
}

func (u uint128) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%064b%064b", u.w1, u.w0)
	return sb.String()
}

/*
 * bitarray64 type
 ********************/

type bitarray64 [64]uint64

func (BA *bitarray64) trans64() *bitarray64 {
	var mask = [12]uint64{
		0x5555555555555555, 0xAAAAAAAAAAAAAAAA,
		0x3333333333333333, 0xCCCCCCCCCCCCCCCC,
		0x0F0F0F0F0F0F0F0F, 0xF0F0F0F0F0F0F0F0,
		0x00FF00FF00FF00FF, 0xFF00FF00FF00FF00,
		0x0000FFFF0000FFFF, 0xFFFF0000FFFF0000,
		0x00000000FFFFFFFF, 0xFFFFFFFF00000000,
	}

	for j := 5; j >= 0; j-- {
		s := 1 << j
		for p := 0; p < 32/s; p++ {
			for i := 0; i < s; i++ {
				i0 := (p*s)<<1 + i
				i1 := i0 + s

				t0 := (BA[i0] & mask[j<<1]) | ((BA[i1] & mask[j<<1]) << s)
				t1 := ((BA[i0] & mask[j<<1|1]) >> s) | (BA[i1] & mask[j<<1|1])
				BA[i0] = t0
				BA[i1] = t1
			}
		}
	}
	return BA
}

func (BA *bitarray64) hsym() *bitarray64 {
	for l, r := 0, 63; l < r; l, r = l+1, r-1 {
		BA[l], BA[r] = BA[r], BA[l]
	}
	return BA
}

func (BA *bitarray64) rotCW() *bitarray64 {
	BA.trans64()
	BA.hsym()
	return BA
}

func (BA *bitarray64) set(j, i int) *bitarray64 {
	BA[j] |= 1 << i
	return BA
}

func (BA *bitarray64) get(j, i int) int {
	return int((BA[j] >> i) & 1)
}

func (BA bitarray64) String() string {
	var sb strings.Builder
	for j := range BA {
		fmt.Fprintf(&sb, "%064b\n", BA[j])
	}
	return sb.String()
}
