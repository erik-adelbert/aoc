// aoc3.go --
// advent of code 2023 day 3
//
// https://adventofcode.com/2023/day/3
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2023-12-3: initial commit

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

const MAXN = 142

func main() {
	var g grid

	input := bufio.NewScanner(os.Stdin)
	for j := 1; input.Scan(); j++ {
		input := input.Bytes()
		g.H, g.W = j, len(input)
		g.setrow(j, input)
	}
	fmt.Println(g.decode())
}

type grid struct {
	data  [MAXN * MAXN]byte
	gears [MAXN * MAXN]gear
	W, H  int

	nums [MAXN]uint192
	syms [MAXN]uint192
	cogs [MAXN]uint192
}

func (g *grid) idx(j, i int) int {
	return j*(g.W+2) + i // 1-based surrounded by empty cells
}

func (g *grid) setrow(j int, row []byte) {
	idx := g.idx

	copy(g.data[idx(j, 1):], row)

	pre, cur, nxt := j-1, j, j+1
	for i, c := range row {
		i++ // 1-based
		g.gears[idx(j, i)].ratio = 1

		mask := one192.lsh(i)
		mask = mask.or(mask.lsh(1), mask.rsh(1))

		switch {
		case isdigit(c):
			g.nums[j].setbit(i)
		case c == '*':
			g.cogs[pre] = g.cogs[pre].or(mask)
			g.cogs[nxt] = g.cogs[nxt].or(mask)

			mask := mask.xor(one192.lsh(i))
			g.cogs[cur] = g.cogs[cur].or(mask)
			fallthrough
		case issymbol(c):
			g.syms[pre] = g.syms[pre].or(mask)
			g.syms[cur] = g.syms[cur].or(mask)
			g.syms[nxt] = g.syms[nxt].or(mask)
		}
	}
}

func (g *grid) decode() (sum, ratio int) {
	data, nums, syms, cogs, idx, H, W := g.data, g.nums, g.syms, g.cogs, g.idx, g.H, g.W
	buf := make([]byte, 0, 4)

	for j := 1; j <= H; j++ {
		read, gear := false, 0
		parts := nums[j].and(syms[j])
		gears := nums[j].and(cogs[j])

		for i := 1; i <= W; i++ {
			c := data[idx(j, i)]
			switch {
			case isdigit(c):
				buf = append(buf, c)

				read = read || parts.getbit(i) > 0

				if gear == 0 && gears.getbit(i) > 0 {
					gear = idx(j, i)
				}
			case c == '*':
				mask := one192.lsh(i)
				mask = mask.or(mask.lsh(1), mask.rsh(1))

				pre := nums[j-1].and(cogs[j-1], mask)
				nxt := nums[j+1].and(cogs[j+1], mask)

				mask.xor(one192.lsh(i))
				cur := gears.and(mask)

				// compute neighboring population
				pop := func(u uint192) int {
					// merge adjacent digits
					if u.popcnt()+u.lead0()+u.trail0() == uint192size {
						// adjacent!
						return 1
					}
					return u.popcnt()
				}

				for _, u := range []uint192{pre, cur, nxt} {
					g.gears[idx(j, i)].count += pop(u)
				}

				if !read {
					continue
				}
				fallthrough
			case read:
				n := atoi(buf)

				sum += n
				read = false

				if gear > 0 {
					g.gears[gear].count++
					g.gears[gear].ratio *= n
					gear = 0
				}
				fallthrough
			default:
				buf = buf[:0]
			}
		}
		if len(buf) > 0 && read {
			n := atoi(buf)

			sum += n
			read = false

			if gear > 0 {
				g.gears[gear].count++
				g.gears[gear].ratio *= n
				gear = 0
			}
			buf = buf[:0]
		}
	}

	gears := g.gears
	for i := range gears[0 : (H+2)*(W+2)] {
		if g := gears[i]; g.count == 2 {
			n := 1

			for voff := -(W + 2); voff <= (W + 2); voff += (W + 2) {
				for hoff := -1; hoff <= +1; hoff++ {
					if r := gears[i-voff+hoff].ratio; r > 0 {
						n *= r
					}
				}
			}

			ratio += n
		}
	}
	return
}

func (g *grid) String() string {
	idx := g.idx
	var sb strings.Builder

	data, H := g.data, g.H
	for j := 0; j < H; j++ {
		fmt.Fprintln(&sb, string(data[idx(j, 0):idx(j+1, 0)]))
	}
	return sb.String()
}

type gear struct {
	count int
	ratio int
}

func isdigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func issymbol(c byte) bool {
	return c > 0 && c != '.' && !isdigit(c)
}

const uint192size = 192

var (
	zero192 uint192
	one192  = uint192{1, 0, 0}
)

type uint192 struct {
	w0, w1, w2 uint64
}

// setbit sets bit n-th n = 0 is LSB.
// n must be <= 255.
func (u *uint192) setbit(n int) {
	switch n >> 6 {
	case 2:
		u.w2 |= (1 << (n & 0x3f))
	case 1:
		u.w1 |= (1 << (n & 0x3f))
	case 0:
		u.w0 |= (1 << (n & 0x3f))
	}
}

func (u *uint192) getbit(n int) uint64 {
	x := u.rsh(n)
	return x.w0 & 1
}

func (u *uint192) popcnt() (n int) {
	n += bits.OnesCount64(u.w2)
	n += bits.OnesCount64(u.w1)
	n += bits.OnesCount64(u.w0)
	return
}

func (u *uint192) lead0() (n int) {
	if n += bits.LeadingZeros64(u.w2); n != 64 {
		return
	}
	if n += bits.LeadingZeros64(u.w1); n != 128 {
		return
	}
	n += bits.LeadingZeros64(u.w0)
	return
}

func (u *uint192) trail0() (n int) {
	if n = bits.TrailingZeros64(u.w0); n != 64 {
		return
	}
	if n += bits.TrailingZeros64(u.w1); n != 128 {
		return
	}
	return n + bits.TrailingZeros64(u.w2)
}

func (u uint192) lsh(n int) uint192 {
	var a, b uint64

	switch {
	case n > 192:
		return uint192{}
	case n > 128:
		u.w2, u.w1, u.w0 = u.w0, 0, 0
		n -= 128
		goto sh128
	case n > 64:
		u.w2, u.w1, u.w0 = u.w1, u.w0, 0
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
	u.w2 = (u.w2 << n) | b

	return u
}

func (u uint192) rsh(n int) uint192 {
	var a, b uint64

	switch {
	case n > 192:
		return uint192{}
	case n > 128:
		u.w2, u.w1, u.w0 = 0, 0, u.w2
		n -= 128
		goto sh128
	case n > 64:
		u.w2, u.w1, u.w0 = 0, u.w2, u.w1
		n -= 64
		goto sh64
	}

	// remaining shifts
	a = u.w2 << (64 - n)
	u.w2 = u.w2 >> n

sh64:
	b = u.w1 << (64 - n)
	u.w1 = (u.w1 >> n) | a

sh128:
	u.w0 = (u.w0 >> n) | b

	return u
}

func (u uint192) and(m ...uint192) uint192 {
	for i := range m {
		u.w2 &= m[i].w2
		u.w1 &= m[i].w1
		u.w0 &= m[i].w0
	}
	return u
}

func (u uint192) or(m ...uint192) uint192 {
	for i := range m {
		u.w2 |= m[i].w2
		u.w1 |= m[i].w1
		u.w0 |= m[i].w0
	}
	return u
}

func (u uint192) xor(m ...uint192) uint192 {
	for i := range m {
		u.w2 ^= m[i].w2
		u.w1 ^= m[i].w1
		u.w0 ^= m[i].w0
	}
	return u
}

func (u uint192) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%016x%016x%016x", u.w2, u.w1, u.w0)
	return sb.String()
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
