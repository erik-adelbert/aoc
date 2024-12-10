// aoc10.go --
// advent of code 2024 day 10
//
// https://adventofcode.com/2024/day/10
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-10: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXDIM = 60

type Cell struct {
	r, c int
}

func main() {
	grid := make([][]int, 0, MAXDIM)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Bytes()
		row := make([]int, 0, len(line))
		for _, c := range line {
			row = append(row, btoi(c))
		}
		grid = append(grid, row)
	}

	fmt.Println(solve(grid)) // part 1 & 2
}

var neighbors = []Cell{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func solve(grid [][]int) (int, int) {
	H, W := len(grid), len(grid[0])

	seen := make([][]bool, H)
	for i := range seen {
		seen[i] = make([]bool, W)
	}

	// recursive dfs search for paths
	var research func(p Cell, target int, goals map[Cell]bool) int
	research = func(p Cell, target int, goals map[Cell]bool) int {
		if p.r < 0 || p.r >= H || p.c < 0 || p.c >= W || seen[p.r][p.c] || grid[p.r][p.c] != target {
			return 0
		}

		// if we reach 9, we found a valid path
		if target == 9 {
			goals[p] = true // remember the goal
			return 1        // count the path
		}

		// mark the cell
		seen[p.r][p.c] = true

		// count the paths from the neighbors
		count := 0
		for _, x := range neighbors {
			count += research(add(p, x), target+1, goals)
		}

		// unmark the cell
		seen[p.r][p.c] = false
		return count
	}

	// path score
	scores := make(map[Cell]int)

	// find all starting points
	count2 := 0
	for r := 0; r < H; r++ {
		for c := 0; c < W; c++ {
			if grid[r][c] == 0 {
				goals := make(map[Cell]bool)
				count2 += research(Cell{r, c}, 0, goals)
				scores[Cell{r, c}] = len(goals)
			}
		}
	}

	count1 := 0
	for _, v := range scores {
		count1 += v
	}

	return count1, count2
}

func add(a, b Cell) Cell {
	return Cell{a.r + b.r, a.c + b.c}
}

func btoi(b byte) int {
	return int(b - '0')
}
