// aoc2.go --
// advent of code 2025 day 2
//
// https://adventofcode.com/2025/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-2: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const MaxDigits = 10 // maximum digit count for our inputs

func main() {
	var acc1, acc2 int // parts 1 and 2

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	line := bytes.TrimSpace(input.Bytes()) // single line input

	for span := range bytes.SplitSeq(line, []byte(",")) {
		// parse range
		bufA, bufB, _ := bytes.Cut(span, []byte("-"))

		a, b := atoi(bufA), atoi(bufB)

		// https://github.com/timvisee/advent-of-code-2025/blob/master/day02b/src/main.rs
		for i := a; i <= b; i++ {
			switch {
			case i >= 1_000_000_000:
				switch {
				case i%100_001 == 0:
					acc1 += i
					fallthrough
				case i%101_010_101 == 0 || i%1_111_111_111 == 0:
					acc2 += i
				}
			case i >= 100_000_000:
				if i%1_001_001 == 0 || i%111_111_111 == 0 {
					acc2 += i
				}
			case i >= 10_000_000:
				switch {
				case i%10_001 == 0:
					acc1 += i
					fallthrough
				case i%1_010_101 == 0 || i%11_111_111 == 0:
					acc2 += i
				}
			case i >= 10_000_000:
				if i%1_010_101 == 0 || i%11_111_111 == 0 {
					acc2 += i
				}
			case i >= 1_000_000:
				if i%1_111_111 == 0 {
					acc2 += i
				}
			case i >= 100_000:
				switch {
				case i%1_001 == 0:
					acc1 += i
					fallthrough
				case i%10_101 == 0 || i%111_111 == 0:
					acc2 += i
				}
			case i >= 10_000:
				if i%11_111 == 0 {
					acc2 += i
				}
			case i >= 1_000:
				switch {
				case i%101 == 0:
					acc1 += i
					fallthrough
				case i%1_111 == 0:
					acc2 += i
				}
			case i >= 100:
				if i%111 == 0 {
					acc2 += i
				}
			case i >= 10:
				if i%11 == 0 {
					acc1 += i
					acc2 += i
				}
			}
		}
	}

	fmt.Println(acc1, acc2)
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
