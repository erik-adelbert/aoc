// aoc15.go --
// advent of code 2024 day 15
//
// https://adventofcode.com/2024/day/15
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-15: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	MAXDIM = 50
	MATRIX = iota
	MOVES
	MAXMOVE = 20000
)

var DIRS = []Cell{
	'^': {-1, 0},
	'v': {1, 0},
	'<': {0, -1},
	'>': {0, 1},
}

type Cell struct {
	r, c int
}

func (c Cell) move(d rune) Cell {
	return Cell{c.r + DIRS[d].r, c.c + DIRS[d].c}
}

type Grid [][]rune

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (g Grid) move(x Cell, dir rune) Cell {

	switch g[x.r][x.c] {
	case '[':
		g[x.r][x.c+1] = '.'
	case ']':
		g[x.r][x.c-1] = '.'
	}
	g[x.r][x.c] = '.'
	old := x

	x = x.move(dir)
	switch g[x.r][x.c] {
	case '#':
		x = old
	case 'O', '[', ']':
		if !g.push(x, dir) {
			x = old
		}
	}

	// fmt.Println("move", x, string(dir), "\n", g)
	return x
}

func (g Grid) push(x Cell, dir rune) bool {

	var repush func(Cell, rune) bool
	repush = func(x Cell, dir rune) bool {
		var old Cell
		// fmt.Println("repush", x, string(dir), "\n", g)

		old, x = x, x.move(dir)
		box := g[old.r][old.c]
		_ = box
		switch g[x.r][x.c] {
		case '#':
			x = old
		case 'O':
			if repush(x, dir) {
				return true
			}
		case '[', ']':
			xx := Cell{x.r, x.c - 1}
			if g[x.r][x.c] == '[' {
				xx.c += 2
			}
			if repush(x, dir) && repush(xx, dir) {
				return true
			}
		case '.':
			switch box {
			case 'O':
				g[x.r][x.c] = 'O'
			case '[':
				if g[x.r][x.c+1] != '.' {
					return false
				}
				g[x.r][x.c] = '['
				g[x.r][x.c+1] = ']'
			case ']':
				if g[x.r][x.c-1] != '.' {
					return false
				}
				g[x.r][x.c] = ']'
				g[x.r][x.c-1] = '['
			}
			return true
		}
		return false
	}

	return repush(x, dir)
}

func (g Grid) expand() Grid {
	new := make(Grid, 0, len(g))
	for j, row := range g {
		new = append(new, make([]rune, 0, 2*len(row)))
		for _, cell := range row {
			switch cell {
			case '#':
				new[j] = append(new[j], '#', '#')
			case '.':
				new[j] = append(new[j], '.', '.')
			case 'O':
				new[j] = append(new[j], '[', ']')
			}
		}
	}
	return new
}

func (g Grid) score() int {
	var sum int
	for j, row := range g {
		for i, cell := range row {
			if cell == 'O' || cell == '[' {
				sum += 100*j + i
			}
		}
	}
	return sum
}

func main() {
	moves := make([]string, 0, MAXMOVE)

	robot1 := Cell{0, 0}
	matrix1 := make(Grid, 0, MAXDIM)

	state := MATRIX
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		switch {
		case len(line) == 0:
			state = MOVES
		case state == MATRIX:
			matrix1 = append(matrix1, []rune(line))
			if i := strings.Index(line, "@"); i != -1 {
				robot1.r, robot1.c = len(matrix1)-1, i
				matrix1[robot1.r][robot1.c] = '.'
			}
		case state == MOVES:
			moves = append(moves, line)
		}
	}

	robot2 := robot1
	matrix2 := matrix1.expand()
	fmt.Println(matrix2)

	for _, dirs := range moves {
		for _, dir := range dirs {
			robot1 = matrix1.move(robot1, dir)
			robot2 = matrix2.move(robot2, dir)
		}
	}
	sum1 := matrix1.score()
	sum2 := matrix2.score()

	// fmt.Println(matrix2)

	fmt.Println(sum1, sum2) // part 1 & 2

	// fmt.Println(matrix)
}

// strconv.Atoi modified loop
// s is ^\d+.*$
func atoi(s string) (n int) {
	for i := 0; i < len(s); i++ {
		n = 10*n + int(s[i]-'0')
	}
	return
}
