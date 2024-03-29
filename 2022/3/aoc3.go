// aoc3.go --
// advent of code 2022 day 3
//
// https://adventofcode.com/2022/day/3
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-3: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

// part indices
const (
	Part1 = iota
	Part2
)

func main() {
	// stricly positive integers are required
	// see day3 notes
	const (
		Head = 1
		Tail = 2
	)

	prios := [2]int{}

	count := func(p int, c rune) {
		prios[p] += int(c - 'a' + 1)
		if c < 'a' {
			prios[p] += 'a' - 'A' + 26
		}
	}

	nline := 0
	chunk := [2][]int{}
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		seen := make([]int, 128)

		line := input.Text()

		// part1 scan
		head, tail := line[:len(line)/2], line[len(line)/2:]

		// turn head into a set
		for _, c := range head {
			seen[c] = Head
		}

		// intersect while adding tail to the set
		for _, c := range tail {
			if seen[c] == Head {
				// head and tail common item
				count(Part1, c)
			}
			seen[c] = Tail
		}

		// part2
		if nline%3 == 2 {
			// chunk scan every 3 input lines
			for _, c := range line {
				if chunk[0][c]*chunk[1][c] > 0 {
					// chunk common badge
					count(Part2, c)
					break // done scanning!
				}
			}
		}

		// store set according to line parity
		chunk[nline&1] = seen
		nline++ // input is 300 lines
	}

	fmt.Println(prios[Part1], prios[Part2])
}
