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
	"strings"
	"time"
)

func main() {
	t0 := time.Now() // start timer

	var acc1, acc2 int // parts 1 and 2 accumulators

	grid := newGrid(MaxGridSize)

	// read input grid
	input := bufio.NewScanner(os.Stdin)

	for i := 0; input.Scan(); i++ { // enumerate input rows
		buf := input.Bytes()

		if i == 0 {
			grid.size = len(buf)
		}

		copy(grid.data[i*grid.size:], buf) // flat copy into grid
	}

	// scan for roll removal using single buffer + double-buffered queue approach

	// preallocate double buffer queues
	queue0 := make([]int, 0, 4*sq(grid.size)/5) // read queue
	queue1 := make([]int, 0, 4*sq(grid.size)/5) // write queue

	// preallocate presence maps
	seen := make([]bool, sq(grid.size))

	// initially, queue all roll positions
	for i := range sq(grid.size) {
		if grid.data[i] == Roll {
			queue0 = append(queue0, i)
		}
	}

	// preallocate roll delete list
	updates := make([]int, 0, RemoveSizeHint)

	for {
		// process current queue - collect removals without modifying grid
		for i := range slices.Values(queue0) {
			if grid.data[i] != Roll {
				continue // skip if not a roll
			}

			r, c := i/grid.size, i%grid.size

			// branchless neighbor bounds
			rα := max(0, r-1)
			rω := min(grid.size-1, r+1)
			cα := max(0, c-1)
			cω := min(grid.size-1, c+1)

			// scan neighbors -- including center roll
			nrolls := 0

			for r = rα; r <= rω; r++ {
				for c = cα; c <= cω; c++ {
					i := r*grid.size + c

					if grid.data[i] == Roll {
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

		// apply all removals at once
		for i := range slices.Values(updates) { // indirect addressing of updates
			grid.data[i] = Empty
		}

		// queue neighbors of removed rolls for next iteration
		for i := range slices.Values(updates) { // indirect addressing of updates
			r, c := i/grid.size, i%grid.size

			// neighbor bounds
			rα := max(0, r-1)
			rω := min(grid.size-1, r+1)
			cα := max(0, c-1)
			cω := min(grid.size-1, c+1)

			for r = rα; r <= rω; r++ {
				for c = cα; c <= cω; c++ {
					i := r*grid.size + c // linear index

					if grid.data[i] == Roll { // only queue remaining rolls
						if !seen[i] {
							queue1 = append(queue1, i)
							seen[i] = true
						}
					}
				}
			}
		}

		// update counts
		if acc1 == 0 {
			acc1 = nremove // first removal count
		}
		acc2 += nremove

		if nremove == 0 {
			break // no more removals
		}

		// prepare for next iteration
		clear(seen)           // reset presence map
		updates = updates[:0] // reset updates

		queue0 = queue1     // swap queues
		queue1 = queue1[:0] // reset queue1
	}

	fmt.Println(acc1, acc2, time.Since(t0))
	// fmt.Println(grid) // uncomment to see the final grid
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
type grid struct {
	data [MaxGridSize * MaxGridSize]byte // flat data
	size int
}

// newGrid creates a new grid of given size
func newGrid(size int) *grid {
	return &grid{
		size: size,
	}
}

func (g *grid) String() string {
	var sb strings.Builder

	for r := range g.size {
		sb.Write(g.data[r*g.size : (r+1)*g.size])
		sb.WriteByte('\n')
	}

	return sb.String()
}

func sq(n int) int {
	return n * n
}
