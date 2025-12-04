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

const (
	Empty byte = '.'
	Roll  byte = '@'
)

func main() {
	var part1, part2 int

	// prepare grids
	grids := []*grid{
		newGrid(MaxGridSize),
		newGrid(MaxGridSize),
	}

	grid := grids[0]

	// read input grid
	input := bufio.NewScanner(os.Stdin)

	for i := 0; input.Scan(); i++ {
		buf := input.Bytes()

		grid.size = len(buf)

		copy(grid.data[i*grid.size:], buf)
	}

	// setup double buffer
	cur := grids[0]
	nxt := grids[1]

	nxt.size = cur.size

	// scan for roll removal
	done := false

	for !done {
		done = true // scan until no more removals

		nremove := 0
		for r := range cur.size {
			for c := range cur.size {
				i := r*cur.size + c // linear index

				nxt.data[i] = cur.data[i] // default copy

				if cur.data[i] != Roll {
					continue // skip non-rolls entirely
				}

				// define neighbor bounds
				rmin := max(0, r-1)
				rmax := min(cur.size-1, r+1)
				cmin := max(0, c-1)
				cmax := min(cur.size-1, c+1)

				// scan neighbors -- including center roll
				nrolls := 0

				for nr := rmin; nr <= rmax; nr++ {
					for nc := cmin; nc <= cmax; nc++ {
						if cur.data[nr*cur.size+nc] == Roll {
							nrolls++
						}
					}
				}

				// decide removal
				if nrolls <= MinRolls { // include center roll
					done = false
					nxt.data[i] = Empty
					nremove++
				}
			}
		}

		// update counts
		if part1 == 0 {
			part1 = nremove // first removal count
		}
		part2 += nremove

		cur, nxt = nxt, cur // swap buffers
	}

	fmt.Println(part1, part2)
	// fmt.Println(cur) // Uncomment to see final grid
}

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
