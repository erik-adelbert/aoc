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
			//   first time sym or
			//   outside current window?
			//   extend window!
			// or
			//   repeating inside?
			//   redim window!
			switch {
			case seen[line[i]] == 0 && line[i] != line[0]:
				fallthrough
			case i-seen[line[i]] > wlen:
				wlen++ // extend
			case i-seen[line[i]] < wlen:
				wlen = i - seen[line[i]] // redim
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
