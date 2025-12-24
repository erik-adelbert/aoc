// aoc1.go --
// advent of code 2025 day 1
//
// https://adventofcode.com/2025/day/1
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const M = 100

func main() {
	t0 := time.Now()

	var acc1, acc2 int // passwords for parts 1 and 2

	// process input lines
	input := bufio.NewScanner(os.Stdin)

	p := 50 // initial dial position
	for input.Scan() {
		buf := input.Bytes()
		dir, n := buf[0], atoi(buf[1:]) // parse direction and number

		// part 2: full wraps
		acc2 += n / M

		// part2: remaining steps
		r := n % M

		s := 1
		if dir == Left {
			s = -1
		}

		// i₀ is the first click index at which the dial lands on 0
		// if i₀ == 0, that corresponds to “already at 0” and must be ignored
		i0 := mod(-p*s, M) // pos + i·s ≡ 0 (mod 100) ⇒ i ≡ -pos·s (mod 100)

		// otherwise, the dial crosses 0 once in the remainder iff i₀ ≤ r
		if i0 != 0 && i0 <= r {
			acc2++
		}

		// part 1: final position
		if p = mod(p+s*n, M); p == 0 {
			acc1++
		}
	}

	fmt.Println(acc1, acc2, time.Since(t0)) // output passwords
}

func mod(a, b int) int {
	if a %= b; a < 0 {
		a += b
	}
	return a
}

const (
	Left  = 'L'
	Right = 'R'
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}

	return
}
