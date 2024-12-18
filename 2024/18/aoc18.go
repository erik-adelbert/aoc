// aoc18.go --
// advent of code 2024 day 18
//
// https://adventofcode.com/2024/day/18
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-18: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	T0     = 1024
	MAXDIM = 71
)

func main() {
	rows := make([]int, 0, MAXDIM)
	cols := make([]int, 0, MAXDIM)

	t := 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), ",")
		rows = append(rows, atoi(args[1]))
		cols = append(cols, atoi(args[0]))
		t++
	}

	g := newGrid(slices.Max(rows)+1, slices.Max(cols)+1, t)
	for i := range rows {
		x := Cell{rows[i], cols[i]}
		g.set(x, i+1)
	}

	p1 := g.shortest(T0)
	p2 := g.locate(g.failfast())

	fmt.Printf("%d  %d,%d\n", p1, p2.c, p2.r)
}

type Grid struct {
	data    []int
	H, W, T int
}

const MaxInt = int(^uint(0) >> 1)

func newGrid(h, w, t int) *Grid {
	g := &Grid{make([]int, h*w), h, w, t}
	for i := range g.data {
		g.data[i] = MaxInt
	}
	return g
}

func (g *Grid) at(x Cell, t int) bool {
	r, c := x.r, x.c
	if r < 0 || r >= g.H || c < 0 || c >= g.W {
		return false
	}
	return g.get(x) > t
}

func (g *Grid) inbounds(x Cell) bool {
	return x.r >= 0 && x.r < g.H && x.c >= 0 && x.c < g.W
}

func (g *Grid) index(x Cell) int {
	return x.r*g.W + x.c
}

func (g *Grid) get(x Cell) int {
	return g.data[g.index(x)]
}

func (g *Grid) set(x Cell, t int) {
	g.data[g.index(x)] = t
}

func (g *Grid) locate(t int) Cell {
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			x := Cell{r, c}
			if g.get(x) == t {
				return x
			}
		}
	}
	return Cell{-1, -1}
}

func (g *Grid) String() string {
	var sb strings.Builder
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			x := Cell{r, c}
			if g.get(x) == MaxInt {
				sb.WriteString(" . ")
				continue
			}
			fmt.Fprintf(&sb, "%2d ", g.get(x))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type Cell struct {
	r, c int
}

func (x Cell) add(y Cell) Cell {
	return Cell{x.r + y.r, x.c + y.c}
}

func (g *Grid) shortest(t int) int {
	dirs := []Cell{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	α, ω := Cell{0, 0}, Cell{g.H - 1, g.W - 1}

	idx := g.index

	Q := []Cell{α}
	seen := make([]bool, g.H*g.W)
	seen[idx(α)] = true

	steps := 0
	for len(Q) > 0 {
		nxtQ := []Cell{}
		for _, cur := range Q {
			if cur == ω {
				return steps
			}

			for _, δ := range dirs {
				nxt := cur.add(δ)
				if g.inbounds(nxt) && !seen[idx(nxt)] && g.at(nxt, t) {
					seen[idx(nxt)] = true
					nxtQ = append(nxtQ, nxt)
				}
			}
		}

		Q = nxtQ
		steps++
	}

	return -1
}

func (g *Grid) failfast() int {
	low, high := T0+1, g.T

	for low < high {
		mid := (low + high) / 2
		if g.shortest(mid) == -1 {
			high = mid
		} else {
			low = mid + 1
		}
	}

	return low
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
