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
	engine := newSchema()

	input := bufio.NewScanner(os.Stdin)
	for j := 1; input.Scan(); j++ {
		engine.setrow(j, input.Bytes())
	}
	fmt.Println(engine.inventory())
}

type gear struct {
	count, ratio int
}

type schema struct {
	φ     func(j, i int) int
	W, H  int
	schem [MAXN * MAXN]byte // original schematic
	gears [MAXN * MAXN]gear // static gear map

	// number, symbol and gear bitmaps
	nums [MAXN]uint192
	syms [MAXN]uint192
	cogs [MAXN]uint192
}

func newSchema() (sc *schema) {
	sc = new(schema)
	sc.φ = func(j, i int) int {
		return j*(sc.W+2) + i // 1-based surrounded by empty cells
	}

	for i := range sc.gears {
		sc.gears[i].ratio = 1
	}

	return
}

func (sc *schema) setrow(j int, row []byte) {
	sc.H, sc.W = max(sc.H, j), len(row)
	// sugars
	φ, nums, syms, cogs := sc.φ, &sc.nums, &sc.syms, &sc.cogs

	copy(sc.schem[φ(j, 1):], row) // 1-based

	pre, cur, nxt := j-1, j, j+1
	for i, c := range row {
		i++ // 1-based

		base := one192.lsh(i)
		mask := base.or(base.lsh(1), base.rsh(1))

		// bitmap
		switch {
		case isdigit(c):
			nums[j].setbit(i)
		case c == '*':
			cogs[pre] = cogs[pre].or(mask)
			cogs[nxt] = cogs[nxt].or(mask)
			cogs[cur] = cogs[cur].or(mask.xor(base)) // do not count "*" itself
		case issymbol(c):
			syms[pre] = syms[pre].or(mask)
			syms[nxt] = syms[nxt].or(mask)
			syms[cur] = syms[cur].or(mask)
		}
	}
}

func (sc *schema) inventory() (sum, ratio int) {
	// sugars
	schem, φ, H, W := &sc.schem, sc.φ, sc.H, sc.W
	nums, syms, cogs := &sc.nums, &sc.syms, &sc.cogs

	buf := make([]byte, 0, 4)

	// vscan
	for j := 1; j <= H; j++ {
		part, gear := false, 0
		parts := nums[j].and(syms[j].or(cogs[j])) // gear is also a symbol
		gears := nums[j].and(cogs[j])

		getpart := func() {
			n := atoi(buf)

			sum += n     // part1
			part = false // unset flag

			if gear > 0 {
				// also a gear
				sc.gears[gear].count++
				sc.gears[gear].ratio *= n
				gear = 0
			}
		}

		// hscan
		for i := 1; i <= W; i++ {
			c := schem[φ(j, i)]
			switch {
			case isdigit(c): // candidate part number digit
				buf = append(buf, c)

				part = part || parts.getbit(i) > 0 // immutable flag once set

				if gear == 0 && gears.getbit(i) > 0 { // set once per part number
					gear = φ(j, i)
				}
			case c == '*': // candidate gear
				base := one192.lsh(i)
				mask := base.or(base.lsh(1), base.rsh(1))

				// no boundary check because boundary is neutral to ops by design
				pre := nums[j-1].and(cogs[j-1], mask)
				nxt := nums[j+1].and(cogs[j+1], mask)
				// do not include '*' itself
				cur := gears.and(mask.xor(base))

				// compute surrounding part count
				pop := func(row uint192) int {
					pop := row.popcnt()
					// fuse adjacent digits
					if row.lead0()+row.trail0() == uint192size-pop {
						// adjacent!
						return 1 // fuse!
					}
					return pop
				}

				// sum surrounding part counts (3x3 window)
				sc.gears[φ(j, i)].count += pop(pre)
				sc.gears[φ(j, i)].count += pop(cur)
				sc.gears[φ(j, i)].count += pop(nxt)

				if !part {
					continue // hscan
				}
				fallthrough // part is next to '*'
			case part:
				getpart()
				fallthrough
			default:
				buf = buf[:0]
			}
		}
		// part ends on schematic row boundary, get it now
		if part {
			getpart()
			buf = buf[:0]
		}
	}

	// sum gears ratios
	gears := sc.gears
	for i := range gears[0 : (H+2)*(W+2)] {
		if g := gears[i]; g.count == 2 {

			// 3x3 window is garanteed to get only two numbers
			// grid is surrounded by empty cells: no boundary check
			n := 1
			for jj := -1; jj < 2; jj++ {
				for ii := i - 1; ii < i+2; ii++ {
					n *= gears[φ(jj, ii)].ratio
				}
			}
			ratio += n // part2
		}
	}
	return
}

func (sc *schema) String() string {
	φ := sc.φ
	var sb strings.Builder

	schem, H := sc.schem, sc.H
	for j := 0; j < H; j++ {
		fmt.Fprintln(&sb, string(schem[φ(j, 0):φ(j+1, 0)]))
	}
	return sb.String()
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
