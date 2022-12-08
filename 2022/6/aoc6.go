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

		line := input.Text()
		// slide over input:
		for i := range line {
			//   outside current window?
			//   extend window!
			// or
			//   repeating inside?
			//   shrink window!
			switch {
			case i-seen[line[i]] > wlen:
				wlen++ // extend right
			case i-seen[line[i]] < wlen:
				wlen = i - seen[line[i]] // shrink left
			}
			seen[line[i]] = i

			// display and loop (part1) or terminate (part2)
			switch {
			case first && wlen == 4:
				fmt.Println(i + 1)
				first = false
			case wlen == 14:
				fmt.Println(i + 1)
				break
			}
		}

		wlen = 0
		seen = [len(seen)]int{} // zero
	}
}
