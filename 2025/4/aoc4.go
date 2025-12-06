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

	// scan for roll removal using single buffer + queue0 approach

	// preallocate double buffer queues
	queue0 := make([][2]int, 0, MaxGridSize*MaxGridSize)
	queue1 := make([][2]int, 0, MaxGridSize*MaxGridSize)

	// preallocate double buffer presence maps
	seen := make([]bool, MaxGridSize*MaxGridSize)
	// seen1 := make([]bool, MaxGridSize*MaxGridSize)

	// initially, queue all roll positions
	for r := range grid.size {
		for c := range grid.size {
			if grid.data[r*grid.size+c] == Roll {
				pos := [2]int{r, c} // next candidate position

				queue0 = append(queue0, pos)
			}
		}
	}

	updates := make([]int, 0, RemoveSizeHint) // roll delete list

	for {
		// process current queue - collect removals without modifying grid
		for _, pos := range queue0 {
			r, c := pos[0], pos[1]
			i := r*grid.size + c

			if grid.data[i] != Roll {
				continue // skip if not a roll
			}

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
						pos := [2]int{nr, nc} // next candidate position

						i := nr*grid.size + nc
						if !seen[i] {
							queue1 = append(queue1, pos)
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

		// prepare for next iteration
		if nremove == 0 {
			break // no more removals
		}

		clear(seen)           // reset presence map
		updates = updates[:0] // reset updates

		queue0 = queue1     // swap queues
		queue1 = queue1[:0] // reset length, keep capacity

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
