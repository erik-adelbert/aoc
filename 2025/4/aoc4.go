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
	"strings"
)

const (
	RemoveSizeHint = 1500 // from previous runs
	MaxGridSize    = 140  // maximum grid size
	MinRolls       = 4    // minimum neighbors to keep a cell
)

func main() {
	var acc1, acc2 int // parts 1 and 2 accumulators

	grid := newGrid(MaxGridSize)

	// read input grid
	input := bufio.NewScanner(os.Stdin)

	for i := 0; input.Scan(); i++ { // enumerate input rows
		buf := input.Bytes()

		grid.size = len(buf)

		copy(grid.data[i*grid.size:], buf)
	}

	// scan for roll removal using single buffer + double-buffered queue approach

	// preallocate double buffer queues
	queue0 := make([]int, 0, MaxGridSize*MaxGridSize)
	queue1 := make([]int, 0, MaxGridSize*MaxGridSize)

	// preallocate presence maps
	seen := make([]bool, MaxGridSize*MaxGridSize)

	// initially, queue all roll positions
	for r := range grid.size {
		for c := range grid.size {
			if grid.data[r*grid.size+c] == Roll {
				i := r*grid.size + c // linear index

				queue0 = append(queue0, i)
			}
		}
	}

	updates := make([]int, 0, RemoveSizeHint) // roll delete list

	for {
		// process current queue - collect removals without modifying grid
		for _, i := range queue0 {
			if grid.data[i] != Roll {
				continue // skip if not a roll
			}

			r, c := i/grid.size, i%grid.size

			// define neighbor bounds
			rmin := max(0, r-1)
			rmax := min(grid.size-1, r+1)
			cmin := max(0, c-1)
			cmax := min(grid.size-1, c+1)

			// scan neighbors -- including center roll
			nrolls := 0

			for nr := rmin; nr <= rmax; nr++ {
				for nc := cmin; nc <= cmax; nc++ {
					i := nr*grid.size + nc

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
		for _, i := range updates { // indirect addressing of data
			grid.data[i] = Empty
		}

		// queue neighbors of removed rolls for next iteration
		for _, i := range updates {
			r, c := i/grid.size, i%grid.size

			// branchless neighbor bounds
			rmin := max(0, r-1)
			rmax := min(grid.size-1, r+1)
			cmin := max(0, c-1)
			cmax := min(grid.size-1, c+1)

			for nr := rmin; nr <= rmax; nr++ {
				for nc := cmin; nc <= cmax; nc++ {
					if grid.data[nr*grid.size+nc] == Roll { // only queue remaining rolls
						ni := nr*grid.size + nc // linear index

						if !seen[ni] {
							queue1 = append(queue1, ni)
							seen[ni] = true
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

	fmt.Println(acc1, acc2)
	// fmt.Println(grid) // uncomment to see the final grid
}

// sugars
const (
	Empty byte = '.'
	Roll  byte = '@'
)

// grid represents a 2D grid of bytes in row-major order
type grid struct {
	data []byte // flat data
	size int
}

// newGrid creates a new grid of given size
func newGrid(size int) *grid {
	return &grid{
		data: make([]byte, MaxGridSize*MaxGridSize),
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
