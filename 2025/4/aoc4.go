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
	MaxGridSize = 140 // maximum grid size
	MinRolls    = 4   // minimum roll neighbors to keep roll
)

func main() {
	var acc1, acc2 int // parts 1 and 2 accumulators

	grid := newGrid(MaxGridSize)

	// read input grid
	input := bufio.NewScanner(os.Stdin)

	for i := 0; input.Scan(); i++ {
		buf := input.Bytes()

		grid.size = len(buf)

		copy(grid.data[i*grid.size:], buf)
	}

	// scan for roll removal using single buffer + queue0 approach
	queue0 := make([][2]int, 0, MaxGridSize*MaxGridSize) // preallocate queue
	queue1 := make([][2]int, 0, MaxGridSize*MaxGridSize) // preallocate next queue
	seen0 := make([]bool, MaxGridSize*MaxGridSize)       // prevent duplicate queue entries

	// initially, queue all roll positions
	for r := range grid.size {
		for c := range grid.size {
			if grid.data[r*grid.size+c] == Roll {
				pos := [2]int{r, c}
				i := r*grid.size + c
				if !seen0[i] {
					queue0 = append(queue0, pos)
					seen0[i] = true
				}
			}
		}
	}

	for len(queue0) > 0 {
		queue1 = queue1[:0] // reset length, keep capacity
		seen1 := make([]bool, MaxGridSize*MaxGridSize)
		var toRemove []int // collect positions to remove

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
					if grid.data[nr*grid.size+nc] == Roll {
						nrolls++
					}
				}
			}

			// decide removal
			if nrolls <= MinRolls { // include center roll
				toRemove = append(toRemove, i)
			}
		}

		nremove := len(toRemove)

		// apply all removals at once
		for _, i := range toRemove {
			grid.data[i] = Empty
		}

		// queue neighbors of removed rolls for next iteration
		for _, i := range toRemove {
			r, c := i/grid.size, i%grid.size

			rmin := max(0, r-1)
			rmax := min(grid.size-1, r+1)
			cmin := max(0, c-1)
			cmax := min(grid.size-1, c+1)

			for nr := rmin; nr <= rmax; nr++ {
				for nc := cmin; nc <= cmax; nc++ {
					if grid.data[nr*grid.size+nc] == Roll { // only queue remaining rolls
						nxtPos := [2]int{nr, nc}
						i := nr*grid.size + nc
						if !seen1[i] {
							queue1 = append(queue1, nxtPos)
							seen1[i] = true
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

		queue0, seen0 = queue1, seen1
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
