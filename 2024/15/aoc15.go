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

func main() {
	moves := make([]string, 0, MAXMOVE)

	rob1 := Cell{0, 0}
	mat1 := make(Grid, 0, MAXDIM)

	state := MATRIX
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		switch {
		case len(line) == 0:
			state = MOVES
		case state == MATRIX:
			mat1 = append(mat1, []rune(line))
			if i := strings.Index(line, "@"); i != -1 {
				rob1.r, rob1.c = len(mat1)-1, i
				mat1[rob1.r][rob1.c] = '.'
			}
		case state == MOVES:
			moves = append(moves, line)
		}
	}

	rob2 := Cell{rob1.r, 2 * rob1.c}
	mat2 := mat1.expand()

	for _, dirs := range moves {
		for _, dir := range dirs {
			rob1 = mat1.move(rob1, dir)
			rob2 = mat2.move(rob2, dir)
		}
	}
	sum1 := mat1.score()
	sum2 := mat2.score()

	fmt.Println(sum1, sum2) // part 1 & 2
}

type Grid [][]rune

func (g Grid) move(x Cell, dir rune) Cell {
	cur, nxt := x, x.move(dir)

	car := g[nxt.r][nxt.c]
	switch car {
	case '#':
		nxt = cur
	case 'O':
		if !g.push(nxt, dir) {
			nxt = cur
		}
	case '[':
		if !g.push(nxt, dir) {
			nxt = cur
		} else {
			g.clear(nxt)
		}
	case ']':
		if !g.push(nxt, dir) {
			nxt = cur
		} else {
			g.clear(nxt)
		}
	}
	return nxt
}

func (g Grid) push(x Cell, dir rune) (ok bool) {

	type up struct { // update
		Cell
		val rune
	}

	updates := make([]up, 0, 32)

	todo := func(x ...up) {
		updates = append(updates, x...)
	}

	var repush func(Cell, rune) bool
	repush = func(nxt Cell, dir rune) bool {
		var cur Cell

		cur, nxt = nxt, nxt.move(dir)

		car := g[cur.r][cur.c]
		switch car {
		case 'O':
			if repush(nxt, dir) {
				todo(up{nxt, car})
				todo(up{cur, '.'})
				return true
			}
		case '[', ']':
			// update context
			lcar, rcar := g[cur.r][cur.c], g[cur.r][cur.c+1]
			lcur, rcur := cur, Cell{cur.r, cur.c + 1}
			lnxt, rnxt := nxt, rcur.move(dir)
			if car == ']' {
				lcar, rcar = g[cur.r][cur.c-1], g[cur.r][cur.c]
				lcur, rcur = Cell{cur.r, cur.c - 1}, cur
				lnxt, rnxt = lcur.move(dir), nxt
			}

			switch dir {
			case 'v':
				if repush(lnxt, dir) && repush(rnxt, dir) {
					todo(up{lnxt, lcar}, up{rnxt, rcar})
					todo(up{lcur, '.'}, up{rcur, '.'})
					return true
				}
			case '^':
				if repush(rnxt, dir) && repush(lnxt, dir) {
					todo(up{rcur, '.'}, up{lcur, '.'})
					todo(up{rnxt, rcar}, up{lnxt, lcar})
					return true
				}
			case '<', '>':
				if repush(nxt, dir) {
					todo(up{nxt, car})
					return true
				}
			}
		case '.':
			return true
		}

		return false
	}

	ok = repush(x, dir)
	if ok && len(updates) > 0 {
		seen := make(map[Cell]bool, len(updates))

		var u up

		for _, u = range updates {
			if seen[u.Cell] && u.val == '.' {
				continue
			}
			seen[u.Cell] = true

			to, v := u.Cell, u.val
			g[to.r][to.c] = v
		}
	}

	return
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

func (g Grid) clear(x Cell) {
	g[x.r][x.c] = '.'
}

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

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
