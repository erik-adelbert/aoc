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
	var (
		wlen int
		seen [128]int
	)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		first := true

		line := input.Bytes()
		// slide over input:
		for i, c := range line {
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

			// display and loop (part1) or terminate (part2)
			switch {
			case first && wlen == 4:
				fmt.Println(i + 1)
				first = false
			case wlen == 14:
				fmt.Println(i + 1)
				break // done!
			}
		}

		wlen = 0
		seen = [len(seen)]int{} // zero
	}
}
