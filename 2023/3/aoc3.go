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
	eng := newSchema()

	input := bufio.NewScanner(os.Stdin)
	for j := 1; input.Scan(); j++ {
		input := input.Bytes()
		eng.setrow(j, input)

	}
	fmt.Println(eng.analyze())
}

type schema struct {
	data  [MAXN * MAXN]byte
	gears [MAXN * MAXN]gear
	W, H  int

	// number, symbol and star bitmaps
	nums [MAXN]uint192
	syms [MAXN]uint192
	cogs [MAXN]uint192
}

func newSchema() (sc *schema) {
	sc = new(schema)

	for i := range sc.gears {
		sc.gears[i].ratio = 1
	}

	return
}

func (sc *schema) idx(j, i int) int {
	return j*(sc.W+2) + i // 1-based surrounded by empty cells
}

func (sc *schema) setrow(j int, row []byte) {
	sc.H, sc.W = j, len(row)
	idx := sc.idx

	copy(sc.data[idx(j, 1):], row)

	pre, cur, nxt := j-1, j, j+1
	for i, c := range row {
		i++ // 1-based

		mask := one192.lsh(i)
		mask = mask.or(mask.lsh(1), mask.rsh(1))

		// make bitmap
		switch {
		case isdigit(c):
			sc.nums[j].setbit(i)
		case c == '*':
			sc.cogs[pre] = sc.cogs[pre].or(mask)
			sc.cogs[nxt] = sc.cogs[nxt].or(mask)
			// do not count "*" itself
			sc.cogs[cur] = sc.cogs[cur].or(mask.xor(one192.lsh(i)))
			fallthrough // star is also a symbol
		case issymbol(c):
			sc.syms[pre] = sc.syms[pre].or(mask)
			sc.syms[nxt] = sc.syms[nxt].or(mask)
			sc.syms[cur] = sc.syms[cur].or(mask)
		}
	}
}

func (sc *schema) analyze() (sum, ratio int) {
	data, nums, syms, cogs, idx, H, W := sc.data, sc.nums, sc.syms, sc.cogs, sc.idx, sc.H, sc.W
	buf := make([]byte, 0, 4)

	for j := 1; j <= H; j++ {
		read, gear := false, 0
		parts := nums[j].and(syms[j])
		gears := nums[j].and(cogs[j])

		getpart := func() {
			n := atoi(buf)

			sum += n // part1
			read = false

			if gear > 0 {
				// also a gear
				sc.gears[gear].count++
				sc.gears[gear].ratio *= n
				gear = 0
			}
		}

		for i := 1; i <= W; i++ {
			c := data[idx(j, i)]
			switch {
			case isdigit(c):
				buf = append(buf, c)

				read = read || parts.getbit(i) > 0 // permanent flag once set

				if gear == 0 && gears.getbit(i) > 0 { // set once per number
					gear = idx(j, i)
				}
			case c == '*':
				mask := one192.lsh(i)
				mask = mask.or(mask.lsh(1), mask.rsh(1))

				pre := nums[j-1].and(cogs[j-1], mask)
				nxt := nums[j+1].and(cogs[j+1], mask)
				// do not count "*" itself
				cur := gears.and(mask.xor(one192.lsh(i)))

				// compute row part count
				pop := func(row uint192) int {
					// merge adjacent digits
					if row.popcnt()+row.lead0()+row.trail0() == uint192size {
						// adjacent!
						return 1
					}
					return row.popcnt()
				}

				// sum neigboring part counts (3x3 window)
				for _, u := range []uint192{pre, cur, nxt} {
					sc.gears[idx(j, i)].count += pop(u)
				}

				if !read {
					continue
				}
				fallthrough // part is next to '*'
			case read:
				getpart()
				fallthrough
			default:
				// consume buffer
				buf = buf[:0]
			}
		}
		// number ends on row boundary
		if len(buf) > 0 && read {
			getpart()
			buf = buf[:0]
		}
	}

	gears := sc.gears
	for i := range gears[0 : (H+2)*(W+2)] {
		if g := gears[i]; g.count == 2 {
			n := 1

			// 3x3 window is garanteed to get only two numbers
			// grid is surrounded by empty cells: no boundary check
			for jj := -1; jj < 2; jj++ {
				for ii := i - 1; ii < i+2; ii++ {
					n *= gears[idx(jj, ii)].ratio
				}
			}

			ratio += n // part2
		}
	}
	return
}

func (sc *schema) String() string {
	idx := sc.idx
	var sb strings.Builder

	data, H := sc.data, sc.H
	for j := 0; j < H; j++ {
		fmt.Fprintln(&sb, string(data[idx(j, 0):idx(j+1, 0)]))
	}
	return sb.String()
}

type gear struct {
	count, ratio int
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
// n must be <= 191.
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
