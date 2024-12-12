// aoc12.go --
// advent of code 2024 day 12
//
// https://adventofcode.com/2024/day/12
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-12: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cell struct {
	r, c int
}

func main() {
	grid := make([][]rune, 0, 140)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		grid = append(grid, []rune(line))
	}

	sum1, sum2 := 0, 0
	regions := decompose(grid)
	for _, region := range regions {
		sum1 += region.area * region.perim
		sum2 += region.sides * region.area
	}

	fmt.Println(sum1, sum2)
	// fmt.Println(regions)
}

type Region struct {
	char  rune
	area  int
	perim int
	sides int
	cells []Cell
}

// decompose the matrix into regions by dfs flood filling
func decompose(matrix [][]rune) []Region {
	H, W := len(matrix), len(matrix[0])

	var cells []Cell

	seen := make([][]bool, len(matrix))
	for i := range seen {
		seen[i] = make([]bool, len(matrix[0]))
	}

	var research func(int, int, rune) (int, int) // recursive dfs
	research = func(r, c int, char rune) (int, int) {
		if r < 0 || r >= H || c < 0 || c >= W || matrix[r][c] != char {
			return 0, 1 // out of bounds or different character contributes to perimeter
		}

		if seen[r][c] {
			return 0, 0
		}
		seen[r][c] = true

		area, perim := 1, 0
		cells = append(cells, Cell{r, c})

		for _, dir := range dirs {
			rr, rc := r+dir.r, c+dir.c
			suba, subp := research(rr, rc, char)
			area += suba
			perim += subp
		}

		return area, perim
	}

	var regions []Region
	for r := range matrix {
		for c := range matrix[r] {
			if !seen[r][c] {
				cells = make([]Cell, 0, len(matrix)*len(matrix[0]))
				area, perimeter := research(r, c, matrix[r][c])
				regions = append(regions, Region{
					char:  matrix[r][c],
					area:  area,
					perim: perimeter,
					sides: sidecount(cells),
					cells: cells,
				})
			}
		}
	}

	return regions
}

// neighboring cells (up, down, left, right)
var dirs = []Cell{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // UDLR
}

func sidecount(region []Cell) (count int) {
	cells := make(map[Cell]bool)
	for _, cell := range region {
		cells[cell] = true
	}

	seen := make(map[[4]int]bool)

	for _, cell := range region {
		r, c := cell.r, cell.c

		for _, δ := range dirs {
			δr, δc := δ.r, δ.c

			// check if the neighboring cell is in the group
			if _, found := cells[Cell{r + δr, c + δc}]; found {
				continue
			}

			// find the corner side
			rr, rc := r, c
			for {
				// check if the neighboring cell in the direction is in the group
				if cells[Cell{rr + δc, rc + δr}] {
					if !cells[Cell{rr + δr, rc + δc}] {
						rr += δc
						rc += δr
						continue
					}
				}
				break
			}

			if !seen[[4]int{rr, rc, δr, δc}] {
				seen[[4]int{rr, rc, δr, δc}] = true
				count++
			}
		}
	}

	return count
}
