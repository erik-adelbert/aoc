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

	paths := make([]int, MaxWidth) // count of paths at each position

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// get started with first row
	i := bytes.Index(input.Bytes(), []byte("S")) // find starting position 'S'
	paths[i] = 1                                 // start with 1 path at S

	// process all remaining rows
	for input.Scan() {
		buf := input.Bytes()

		for i := range buf {
			if buf[i] == Prism {
				if paths[i] > 0 {
					acc1 += 1 // part1: count splits

					// adding paths count to acc2 here is miai with #L55-57 below
					// acc2 += paths[i]
				}

				// split paths to left and right
				// always splitting is faster than including it inside #L41 test
				// because of the speculative execution
				paths[i-1] += paths[i] // add to left
				paths[i+1] += paths[i] // add to right
				paths[i] = 0           // clear paths at split position
			}
		}
	}

	// part 2: sum all active path counts
	// batch adding here is faster than adding one-by-one at #L43
	// because it is optimally sequential memory access
	for i := range paths {
		acc2 += paths[i]
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

const (
	MaxWidth = 141 // maximum grid width from prior runs
	Prism    = '^'
)
