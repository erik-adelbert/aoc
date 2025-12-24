// aoc2.go --
// advent of code 2025 day 2
//
// https://adventofcode.com/2025/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-2: initial commit
// 2025-12-4: Tim Visée's approach
// 2025-12-5: improved approach discussed with hm - sub ms runtime

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"iter"
	"os"
	"time"
)

func main() {
	t0 := time.Now() // start timer

	var acc1, acc2 int // parts 1 and 2 accumulators

	// read input spans
	input := bufio.NewReader(os.Stdin)
	buf, _, _ := input.ReadLine()

	for a, b := range allSpans(buf) { // split into aligned subranges
		switch {
		case a >= 1e9:
			const seed1, seed2, lcm = 100_001, 101_010_101, 1_111_111_111

			// sum all multiples of seed1 over [a, b]
			acc1 += sm(a, b, seed1)

			// sum multiples of seed2, subtracting common multiples of seed1&2 already counted in sub1
			acc2 += sm(a, b, seed2) - sm(a, b, lcm)
		case a >= 1e8:
			const seed2a, seed2b = 1_001_001, 111_111_111

			acc2 += sm(a, b, seed2a)
			acc2 += sm(a, b, seed2b)
		case a >= 1e7:
			const seed1, seed2 = 10_001, 11_111_111

			acc1 += sm(a, b, seed1)
			acc2 += sm(a, b, seed2)
		case a >= 1e6:
			const seed2 = 1_111_111

			acc2 += sm(a, b, seed2)
		case a >= 1e5:
			const seed1, seed2, lcm = 1_001, 10_101, 111_111

			acc1 += sm(a, b, seed1)
			acc2 += sm(a, b, seed2) - sm(a, b, lcm)
		case a >= 1e4:
			const seed2 = 11_111

			acc2 += sm(a, b, seed2)
		case a >= 1e3:
			const seed1 = 101

			acc1 += sm(a, b, seed1)
		case a >= 1e2:
			const seed2 = 111

			acc2 += sm(a, b, seed2)
		case a >= 1e1:
			const seed1 = 11

			acc1 += sm(a, b, seed1)
		}
	}

	acc2 += acc1 // part 2 includes part 1
	fmt.Println(acc1, acc2, time.Since(t0))
}

// allSpans iterates over subranges [a, b] split at ten powers boundaries
// e.g., [95, 105] -> [95, 99], [100, 105]
func allSpans(buf []byte) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		spans := bytes.SplitSeq(buf, []byte(","))

		for span := range spans {
			lhs, rhs, _ := bytes.Cut(span, []byte("-")) // parse range

			start, end := atoi(lhs), atoi(rhs)

			for x := 10; x <= 1e9; x *= 10 {
				if x > start && x <= end {
					if !yield(start, x-1) { // [start, x-1]
						return
					}

					start = x
				}
			}

			// yield the final range [last_split, b] and return anyway
			if start <= end {
				if !yield(start, end) {
					return
				}
			}
		}
	}
}

// sm computes the sum of all multiples of x in the range [l, r]
func sm(l, r, x int) int {
	// first multiple of x in [l, r]: ⎡l/x⎤ * x
	// x when no multiple in [l, r]
	α := ((l + x - 1) / x) * x

	// last multiple of x in [l, r]: ⎣r/x⎦ * x
	// 0 when no multiple in [l, r]
	ω := (r / x) * x

	// count of multiples of x in [α, ω]
	// 0 if none (ω-α == -x)
	n := (ω-α)/x + 1

	// sum all multiples using arithmetic series
	// 0 if none (n == 0)
	return n * (α + ω) / 2
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}

	return
}
