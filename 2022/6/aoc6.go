// aoc6.go --
// advent of code 2022 day 6
//
// https://adventofcode.com/2022/day/6
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-6: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// state handlers for parts 1&2
	// https://go.dev/talks/2011/lex.slide
	type state func(int) (next state, match bool)

	var (
		wlen         int
		seen         [128]int
		part1, part2 state
	)

	// https://go.dev/talks/2011/lex.slide#19
	part1 = func(wlen int) (next state, match bool) {
		if wlen == 4 {
			return part2, true // part1 done, transition to part2
		}
		return part1, false // part1 not done, no match
	}

	part2 = func(wlen int) (next state, match bool) {
		if wlen == 14 {
			return nil, true // part2 done, transition to end
		}
		return part2, false // part2 not done, no match
	}

	// initial state is solving for part1
	check := part1

	// slide over single line input
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	for i, c := range input.Bytes() {
		//   outside current window?
		//   extend window!
		// or
		//   repeating inside?
		//   shrink window!
		switch {
		case i-seen[c] > wlen:
			wlen++ // extend right
		case i-seen[c] < wlen:
			wlen = i - seen[c] // shrink left
		}
		seen[c] = i

		// dynamic state machine
		// states are self transitioning check functions
		var match bool
		if check, match = check(wlen); match {
			fmt.Println(i + 1)
			if check == nil { // terminal state
				return
			}
		}

		// terminal state check could be here (canonical) but
		// in this case it would only be an extraneous test as
		// we already know that transitions occurs on matches
		// see https://go.dev/talks/2011/lex.slide#20
		//
		// if check == nil { // terminal state
		// 	return
		// }

	}
}
