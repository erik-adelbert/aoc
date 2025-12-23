// aoc4.go --
// advent of code 2025 day 4
//
// https://adventofcode.com/2025/day/4
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-4: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"
)

func main() {
	t0 := time.Now() // start timer

	var (
		acc1, acc2 int  // parts 1 and 2 accumulators
		grid       grid // input grid
	)

	// read input grid
	input := bufio.NewScanner(os.Stdin)

	sz := 0
	for i := 0; input.Scan(); i++ { // enumerate input rows
		buf := input.Bytes()

		if i == 0 {
			sz = len(buf)
		}

		copy(grid[i*sz:], buf) // flat copy into grid
	}

	// scan for roll removal using single buffer + double-buffered queue

	// preallocate double buffer queues
	queue0 := make([]int, 0, 7*sq(sz)/10) // read queue
	queue1 := make([]int, 0, 7*sq(sz)/10) // write queue

	// preallocate presence maps
	seen := make([]uint8, sq(sz))

	// preallocate roll delete list
	updates := make([]int, 0, RemoveSizeHint)

	// initially, queue all roll positions
	for i := range sq(sz) {
		if grid[i] == Roll {
			queue0 = append(queue0, i)
		}
	}

	// removal loop
	for gid := uint8(1); ; gid++ { // removal generation
		// collect removals
		for i := range slices.Values(queue0) {
			if grid[i] != Roll {
				continue // skip if not a roll
			}

			r, c := i/sz, i%sz
			// branchless neighbor bounds
			rα := max(0, r-1)
			rω := min(sz-1, r+1)
			cα := max(0, c-1)
			cω := min(sz-1, c+1)

			// scan neighbors -- including center roll
			nrolls := 0

			for r = rα; r <= rω; r++ {
				for c = cα; c <= cω; c++ {
					i := r*sz + c

					if grid[i] == Roll {
						nrolls++
					}
				}
			}

			// decide removal
			if nrolls <= MinRolls { // include center roll
				updates = append(updates, i)
			}
		}

		nremove := len(updates)
		if nremove == 0 {
			break // no more removals, all done!
		}

		// update counts
		if acc1 == 0 {
			acc1 = nremove // first removal count
		}
		acc2 += nremove

		// apply all removals atomically
		for i := range slices.Values(updates) { // indirect addressing of updates
			grid[i] = Empty
		}

		// queue neighbors of removed rolls for next iteration
		for i := range slices.Values(updates) { // indirect addressing of updates
			r, c := i/sz, i%sz

			// neighbor bounds
			rα := max(0, r-1)
			rω := min(sz-1, r+1)
			cα := max(0, c-1)
			cω := min(sz-1, c+1)

			for r = rα; r <= rω; r++ {
				for c = cα; c <= cω; c++ {
					i := r*sz + c // linear index

					if grid[i] != Roll { // only queue remaining rolls
						continue
					}

					if seen[i] != gid {
						queue1 = append(queue1, i)
						seen[i] = gid
					}
				}
			}
		}

		// prepare for next iteration
		updates = updates[:0]               // reset updates
		queue0, queue1 = queue1, queue0[:0] // swap queues, reset write queue
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

const (
	RemoveSizeHint = 1500 // from previous runs
	MaxGridSize    = 140  // maximum grid size
	MinRolls       = 4    // minimum neighbors to keep a cell
)

// sugars
const (
	Empty byte = '.'
	Roll  byte = '@'
)

// grid represents a 2D grid of bytes in row-major order
type grid = [MaxGridSize * MaxGridSize]byte // flat data

func sq(n int) int {
	return n * n
}
