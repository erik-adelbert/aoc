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
)

func main() {
	input := bufio.NewScanner(os.Stdin)

	old, cur := 50, 50 // dial starts at 50
	pwd1, pwd2 := 0, 0 // passwords

	// process input lines
	for input.Scan() {
		buf := input.Bytes()
		dir, n := buf[0], atoi(buf[1:]) // parse direction and number

		// handle large movements
		pwd2 += n / 100 // count full wraps
		n = mod(n, 100) // reduce to within one wrap

		// move dial
		cur += n // default to right turn
		if dir == 'L' {
			cur -= 2 * n
		}

		// handle circular dial (0-99)
		cur = mod(cur, 100)

		switch {
		case old == 0:
			// cannot cross or reach zero from zero
			// count nothing
		case cur == 0:
			// count landings on zero
			pwd1++
		case cur > old && dir == 'L':
			// count left turns crossing zero
			pwd2++
		case cur < old && dir == 'R':
			// count right turns crossing zero
			pwd2++
		}

		old = cur // update
	}

	fmt.Println(pwd1, pwd1+pwd2) // output passwords
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

// mod returns a modulo b, always non-negative
func mod(a, b int) (r int) {
	r = a % b
	if r < 0 {
		r += b
	}
	return
}
