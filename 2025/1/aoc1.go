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

func main() {
	t0 := time.Now()

	var acc1, acc2 int // passwords for parts 1 and 2

	old, cur := MaxDial/2, MaxDial/2 // dial starts at 50

	// process input lines
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		buf := input.Bytes()
		dir, n := buf[0], atoi(buf[1:]) // parse direction and number

		// handle large movements
		acc2 += n / MaxDial // count full wraps
		n %= MaxDial        // reduce to within one wrap

		// move dial: default to left turn
		if cur = old - n; dir == Right {
			cur = old + n // adjust for right turn
		}

		// handle circular dial (0-99)
		if cur %= MaxDial; cur < 0 {
			cur += MaxDial // adjust negative remainder
		}

		switch {
		case old == 0:
			// cannot reach or cross zero from zero in less than a wrap
			// count nothing
		case cur == 0:
			// part1: count turns landing on zero
			acc1++
		case (old < cur) == (dir == Left): // position increased/decreased when turning left/right
			// part2: count turns crossing zero
			acc2++
		}

		old = cur // update
	}
	acc2 += acc1 // part2 includes part1

	fmt.Println(acc1, acc2, time.Since(t0)) // output passwords
}

const MaxDial = 100

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
