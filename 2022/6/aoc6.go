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
			if first && wlen == 4 { // part1
				fmt.Println(i + 1)
				first = false
			}

			if wlen == 14 { // part2
				fmt.Println(i + 1)
				break // all done!
			}
		}

		wlen = 0
		seen = [len(seen)]int{} // zero
	}
}
