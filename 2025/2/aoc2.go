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
)

func main() {
	var acc1, acc2 int // parts 1 and 2 accumulators

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	line := bytes.TrimSpace(input.Bytes()) // single line input

	for span := range bytes.SplitSeq(line, []byte(",")) {
		bufA, bufB, _ := bytes.Cut(span, []byte("-")) // parse range

		a, b := atoi(bufA), atoi(bufB)

		sub1, sub2 := 0, 0                      // partial sums for this span
		for span := range allSplitSpans(a, b) { // split into subranges
			a, b := span[0], span[1]

			switch {
			case a >= 1_000_000_000:
				const seed1, seed2 = 100_001, 101_010_101

				// sum multiples of seed1
				sub1 += sumMultiples(a, b, seed1)

				// sum multiples of seed2, excluding common multiples already counted in sub1
				sub2 += sumMultiples(a, b, seed2) - sumMultiples(a, b, lcm(seed2, seed1))
			case a >= 100_000_000:
				const seed2a, seed2b = 1_001_001, 111_111_111

				sub2 += sumMultiples(a, b, seed2a)
				sub2 += sumMultiples(a, b, seed2b)
			case a >= 10_000_000:
				const seed1, seed2 = 10_001, 11_111_111

				sub1 += sumMultiples(a, b, seed1)

				sub2 += sumMultiples(a, b, seed2) - sumMultiples(a, b, lcm(seed2, seed1))
			case a >= 1_000_000:
				const seed2 = 1_111_111

				sub2 += sumMultiples(a, b, seed2)
			case a >= 100_000:
				const seed1, seed2 = 1_001, 10_101

				sub1 += sumMultiples(a, b, seed1)

				sub2 += sumMultiples(a, b, seed2) - sumMultiples(a, b, lcm(seed2, seed1))
			case a >= 10_000:
				const seed2 = 11_111

				sub2 += sumMultiples(a, b, seed2)
			case a >= 1_000:
				const seed1 = 101

				sub1 += sumMultiples(a, b, seed1)
			case a >= 100:
				const seed2 = 111

				sub2 += sumMultiples(a, b, seed2)
			case a >= 10:
				const seed1 = 11

				sub1 += sumMultiples(a, b, seed1)
			}
		}

		acc1 += sub1
		acc2 += sub2
	}
	acc2 += acc1 // part 2 includes part 1

	fmt.Println(acc1, acc2)
}

// allSplitSpans iterates over subranges of [a, b] split at ten powers boundaries
func allSplitSpans(a, b int) iter.Seq[[2]int] {
	return func(yield func([2]int) bool) {
		var splitPoints = [...]int{
			10, 100, 1_000, 10_000, 100_000, 1_000_000,
			10_000_000, 100_000_000, 1_000_000_000,
		}

		start := a

		for _, x := range splitPoints {
			if x > start && x <= b {
				if !yield([2]int{start, x - 1}) {
					return
				}
				start = x
			}
		}

		// yield the final range [last_split, b] and return anyway
		if start <= b {
			yield([2]int{start, b})
		}
	}
}

// sumMultiples computes the sum of all multiples of x in the range [a, b]
func sumMultiples(a, b, x int) int {
	var first, last int // first and last multiples of x in [a, b]

	// first multiple: ceiling(a/x) * x
	if first = ((a + x - 1) / x) * x; first > b {
		return 0 // no multiples in range
	}

	// last multiple: floor(b/x) * x
	last = (b / x) * x

	// count of multiples
	count := (last-first)/x + 1

	// sum using arithmetic series
	return count * (first + last) / 2
}

// gcd computes the greatest common divisor of a and b
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm computes the least common multiple of a and b
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
