// aoc7.go --
// advent of code 2025 day 7
//
// https://adventofcode.com/2025/day/7
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-7: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

func main() {
	t0 := time.Now() // start timer

	var acc1, acc2 int // parts 1 and 2 accumulators

	grid := make([][]byte, 0, MaxWidth)

	paths := make([]int, MaxWidth) // count of paths at each position

	input := bufio.NewScanner(os.Stdin)

	// read all input lines, filtering out empty dot lines because they are not contributing
	for input.Scan() {
		buf := input.Bytes()

		if bytes.ContainsAny(buf, "^S") {
			grid = append(grid, bytes.Clone(buf))
		}
	}

	i := bytes.Index(grid[0], []byte("S")) // find starting position 'S' in first row
	paths[i] = 1                           // start with 1 path at S

	// process remaining rows
	for _, row := range grid[1:] {
		for i := range paths {
			if row[i] == Prism {
				if paths[i] > 0 {
					acc1 += 1 // part1: count splits
				}

				// split paths to left and right
				paths[i-1] += paths[i] // add to left
				paths[i+1] += paths[i] // add to right
				paths[i] = 0           // clear paths at split position
			}
		}
	}

	// part 2: sum all path counts
	for i := range paths {
		acc2 += paths[i]
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

const (
	MaxWidth = 141 // maximum grid width from prior runs
	Prism    = '^'
)
