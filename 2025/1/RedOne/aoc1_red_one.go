// aoc1_red_one.go --
// advent of code 2025 day 1
//
// https://adventofcode.com/2025/day/1
// https://github.com/erik-adelbert/aoc
// https://www.reddit.com/r/adventofcode/comments/1pb3y8p/2025_day_1_solutions/https://old.reddit.com/r/adventofcode/comments/1pb3yje/advent_of_code_2025_reddit_one_submissions/?
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-10: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"
)

const MaxDial = 100

func main() {
	t0 := time.Now()

	var acc1, acc2 int // passwords for parts 1 and 2

	cmds := make(chan []byte)

	go readin(cmds) // start input reader

	// initial state
	in := make(chan int, 3)

	in <- MaxDial / 2 // position
	in <- 0           // part 1 accumulator value
	in <- 0           // part 2 accumulator value

	// chain dial goroutines
	for cmd := range cmds {
		out := make(chan int, 3)
		go dial(cmd, in, out)
		in = out
	}

	// collect final results
	_ = <-in           // discard final position
	acc1 = <-in        // part 1 password
	acc2 = acc1 + <-in // part 2 password

	fmt.Println(acc1, acc2, time.Since(t0)) // output passwords
}

// readin reads all input lines and sends them to the channel
func readin(ch chan<- []byte) {
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		ch <- slices.Clone(input.Bytes())
	}

	close(ch)
}

// dial processes a single dial command
func dial(cmd []byte, in <-chan int, out chan<- int) {
	dir, n := cmd[0], atoi(cmd[1:])     // parse direction and number
	old, acc1, acc2 := <-in, <-in, <-in // get current position and accumulators

	// handle large movements
	acc2 += n / MaxDial // count full wraps
	n %= MaxDial        // reduce to within one wrap

	var cur int

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

	// send updated state
	out <- cur
	out <- acc1
	out <- acc2
}

// direction constants
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
