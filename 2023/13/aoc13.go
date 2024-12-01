package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {

	fmt.Println('a')

	sum1, sum2 := 0, 0

	// yabadabadoo
	//    y   a   b   a   d   a   b   a   d   o   o
	//   24   0   1   0   3   0   1   0   3  14  14     ASCII codes - 97/'a'
	// [0 1 0 1 0 3 0 1 0 7 0 1 0 5 0 1 0 1 0 1 2 1 0]  all palindromes in Θ(n)!!!
	// fmt.Println(fastLongestPalindrome([]uint32{
	// 	24, 0, 1, 0, 3, 0, 1, 0, 3, 14, 14,
	// }))

	terrain := newArea()
	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		line := input.Text()
		switch len(line) {
		case 0:
			n1, n2 := solve(terrain)
			sum1 += n1
			sum2 += n2
			j, terrain = -1, newArea()
		default:
			terrain.H = max(terrain.H, j+1)
			terrain.W = max(terrain.W, len(line))
			for i := range line {
				if line[i] == '#' {
					terrain.setbit(j, i)
				}
			}
		}
	}
	// last input
	n1, n2 := solve(terrain)
	sum1 += n1
	sum2 += n2

	fmt.Println(sum1, sum2)
}

const (
	RC = false // row major
	CR = !RC   // col major
)

type area struct {
	ord  bool // default to ROWMAJ
	data bitarray32
	H, W int
}

func newArea() (a *area) {
	a = new(area)
	return
}

func (a *area) setbit(j, i int) {
	a.data.set(j, i)
}

func (a *area) flipbit(j, i int) {
	a.data.flip(j, i)
}

func (a *area) order(ord bool) *area {
	if a.ord != ord {
		switch {
		case a.ord == CR:
			a.data.trans32()
			a.ord, a.H, a.W = RC, a.W, a.H
		case a.ord == RC:
			a.data.trans32()
			a.ord, a.H, a.W = CR, a.W, a.H
		}
	}
	return a
}

func (a *area) offsets() (j0, i0 int) {
	if a.ord == CR {
		j0, i0 = 32-a.H, 32-a.W
	}
	return
}

func (a *area) String() string {
	var sb strings.Builder

	r := strings.NewReplacer(string(byte(0)), ".")

	j0, i0 := a.offsets()
	for _, n := range a.data[j0 : j0+a.H] {
		if a.ord == CR {
			n >>= i0
		}

		for i := 0; i < a.W; i++ {
			sb.WriteByte('#' * byte(n&1))
			n >>= 1
		}
		sb.WriteByte('\n')
	}
	fmt.Fprintf(&sb, "H: %d W: %d CR:%v\n", a.H, a.W, a.ord)

	return r.Replace(sb.String())
}

func (a *area) scan(s []int) ([]int, bool) {
	iseven := func(n int) bool {
		return n&1 == 0
	}

	matches := make([]int, 0)
	for i := range s {
		if iseven(i) && s[i] > 1 { // palindrome pivot between rows (i/2-1) and (i/2)
			half, l, r := (s[i]/2)-1, (i-2)/2, i/2

			if l-half == 0 || r+half == a.H-1 {
				switch a.ord {
				case RC:
					r *= 100
				case CR:
					r = a.H - r
				}
				matches = append(matches, r)
			}
		}
	}

	if len(matches) > 0 {
		return matches, true
	}

	return []int{-1}, false
}

func (a *area) score1D(order bool) (score []int, ok bool) {
	a.order(order)
	j0, _ := a.offsets()
	data := flp(a.data[j0 : j0+a.H])
	score, ok = a.scan(data)
	return
}

func (a *area) score() []int {
	rscore, _ := a.score1D(RC)
	cscore, ok := a.score1D(CR)

	if !ok {
		return rscore
	}
	return cscore
}

func solve(a *area) (int, int) {
	var smudge []int
	var clean int

	smudge = a.score()
	clean = a.clean(smudge)

	return smudge[0], clean
}

func (a *area) clean(base []int) int {
	var clean int

	for _, o := range []bool{a.ord, !a.ord} {
		a.order(o)
		j0, _ := a.offsets()

		d := a.data[j0 : j0+a.H]

		for i := 0; i < a.H-1; i++ {
			for ii := i + 1; ii < a.H; ii++ {
				if n := d[i] ^ d[ii]; popcnt(n) == 1 {
					d[i] ^= n
					if cleans, ok := a.score1D(o); ok {
						for _, clean = range cleans {
							if clean != base[0] {
								return clean
							}
						}
					}
					d[i] ^= n
				}
			}
		}
	}

	panic("unreachable")
}

var flp = fastLongestPalindrome

// https://www.akalin.com/longest-palindrome-linear-time
func fastLongestPalindrome(seq []uint32) []int {
	l := make([]int, 0, len(seq))

	i, pallen := 0, 0
	// Loop invariant: seq[(i - palLen):i] is a palindrome.
	// Loop invariant: len(l) >= 2 * i - palLen. The code path that
	// increments palLen skips the l-filling inner-loop.
	// Loop invariant: len(l) < 2 * i + 1. Any code path that
	// increments i past seqLen - 1 exits the loop early and so skips
	// the l-filling inner loop.
SCAN:
	for i < len(seq) {
		// First, see if we can extend the current palindrome.  Note
		// that the center of the palindrome remains fixed.
		if i > pallen && seq[i-pallen-1] == seq[i] {
			pallen += 2
			i += 1
			continue
		}

		l = append(l, pallen)

		// Now to make further progress, we look for a smaller
		// palindrome sharing the right edge with the current
		// palindrome.  If we find one, we can try to expand it and see
		// where that takes us.  At the same time, we can fill the
		// values for l that we neglected during the loop above. We
		// make use of our knowledge of the length of the previous
		// palindrome (palLen) and the fact that the values of l for
		// positions on the right half of the palindrome are closely
		// related to the values of the corresponding positions on the
		// left half of the palindrome.

		// Traverse backwards starting from the second-to-last index up
		// to the edge of the last palindrome.
		s := len(l) - 2
		e := s - pallen
		for j := s; j > e; j-- {
			// d is the value l[j] must have in order for the
			// palindrome centered there to share the left edge with
			// the last palindrome.  (Drawing it out is helpful to
			// understanding why the - 1 is there.)
			d := j - e - 1

			// We check to see if the palindrome at l[j] shares a left
			// edge with the last palindrome.  If so, the corresponding
			// palindrome on the right half must share the right edge
			// with the last palindrome, and so we have a new value for
			// palLen.
			//
			// An exercise for the reader: in this place in the code you
			// might think that you can replace the == with >= to improve
			// performance.  This does not change the correctness of the
			// algorithm but it does hurt performance, contrary to
			// expectations.  Why?
			if l[j] == d {
				pallen = d
				continue SCAN
			}

			// Otherwise, we just copy the value over to the right
			// side.  We have to bound l[i] because palindromes on the
			// left side could extend past the left edge of the last
			// palindrome, whereas their counterparts won't extend past
			// the right edge.
			l = append(l, min(d, l[j]))
		}

		// This code is executed in two cases: when the for loop
		// isn't taken at all (palLen == 0) or the inner loop was
		// unable to find a palindrome sharing the left edge with
		// the last palindrome.  In either case, we're free to
		// consider the palindrome centered at seq[i].
		pallen = 1
		i++
	}
	// We know from the loop invariant that len(l) < 2 * seqLen + 1, so
	// we must fill in the remaining values of l.

	// Obviously, the last palindrome we're looking at can't grow any
	// more.
	l = append(l, pallen)

	// Traverse backwards starting from the second-to-last index up
	// until we get l to size 2 * seqLen + 1. We can deduce from the
	// loop invariants we have enough elements.
	s := len(l) - 2
	e := s - (2*len(seq) + 1 - len(l))
	for i := s; i > e; i-- {
		// The d here uses the same formula as the d in the inner loop
		// above.  (Computes distance to left edge of the last
		// palindrome.)
		d := i - e - 1
		// # We bound l[i] with min for the same reason as in the inner
		// # loop above.
		l = append(l, min(d, l[i]))
	}

	return l
}

/*
 * bitarray32 type
 ********************/

type bitarray32 [32]uint32

// hacker's delight H.S. Warren, Jr. ISBN10: 0-201-91465-4 p. 113
func (BA *bitarray32) trans32() *bitarray32 {
	j, m := 16, uint32(0x0000FFFF)
	for j != 0 {
		for k := 0; k < 32; k = (k + j + 1) & ^j {
			t := (BA[k] ^ (BA[k+j] >> j)) & m
			BA[k] = BA[k] ^ t
			BA[k+j] = BA[k+j] ^ (t << j)
		}
		j >>= 1
		m ^= (m << j)
	}
	return BA
}

func (BA *bitarray32) set(j, i int) *bitarray32 {
	BA[j] |= 1 << i
	return BA
}

func (BA *bitarray32) get(j, i int) int {
	return int((BA[j] >> i) & 1)
}

func (BA *bitarray32) clear(j, i int) *bitarray32 {
	BA[j] &= ^(1 << i)
	return BA
}

func (BA *bitarray32) flip(j, i int) *bitarray32 {
	BA[j] ^= 1 << i
	return BA
}

func (BA bitarray32) String() string {
	var sb strings.Builder

	for j := range BA {
		fmt.Fprintf(&sb, "%032b\n", BA[j])
	}
	return sb.String()
}

var popcnt, trail0 = bits.OnesCount32, bits.TrailingZeros32

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
