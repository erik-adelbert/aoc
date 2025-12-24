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

		// part 2: full wraps ⎣n/M⎦
		acc2 += n / M

		// part 2: check for crossing 0 in the remainder

		//remaining steps
		r := n % M

		// step sign
		s := 1
		if dir == Left {
			s = -1
		}

		// if i is a click landing on 0 from p with step s,
		// let i₀ be the first of those clicks in [0, M-1]
		// we have:
		//     p + i·s ≡ 0 (mod M) => i·s ≡ -p (mod M)
		// s = ±1 and M are coprime, so the modular inverse of s is s itself
		// and we can multiply both sides by it, thus:
		//     i ≡ -p·s (mod M)
		// and, finally:
		//     i₀ = -p·s (mod M)
		i0 := mod(-p*s, M)

		// if i₀ == 0, that corresponds to “already at 0” and must be ignored
		// otherwise, the dial crosses 0 once in the remainder iff i₀ ≤ r
		if i0 != 0 && i0 <= r {
			acc2++
		}

		// part 1: advance p, final position is checked for 0
		if p = mod(p+s*n, M); p == 0 {
			acc1++
		}
	}

	// output clear passwords on stdout because why not?
	fmt.Println(acc1, acc2, time.Since(t0))
}

// sugar
const Left = 'L'

// mod computes euclidean modulo
func mod(a, b int) int {
	if a %= b; a < 0 {
		a += b
	}
	return a
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}

	return
}
