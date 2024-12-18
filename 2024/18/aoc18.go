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

	mem := newGrid(slices.Max(rows)+1, slices.Max(cols)+1, t)
	for i := range rows {
		x := Cell{rows[i], cols[i]}
		mem.set(x, i+1)
	}

	p1 := mem.shortest(T0)
	p2 := mem.locate(mem.failfast(T0 + 1)) // we know the closing time is at least T0+1

	fmt.Printf("%d  %d,%d\n", p1, p2.c, p2.r)
}

func (g *Grid) failfast(t0 int) int {
	low, high := t0, g.T

	for low < high {
		mid := (low + high) / 2
		if g.shortest(mid) == -1 {
			high = mid
		} else {
			low = mid + 1
		}
	}

	// g.dumpat(low)
	return low
}

func (g *Grid) shortest(t int) int {
	dirs := []Cell{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	α, ω := Cell{0, 0}, Cell{g.H - 1, g.W - 1}

	idx := g.index

	// BFS
	const NALLOC = 80 // arbitrary but informed preallocation

	// double-buffered queue
	Q1, Q2 := make([]Cell, 0, NALLOC), make([]Cell, 0, NALLOC)

	// visited cells
	seen := make([]bool, g.H*g.W)

	seen[idx(α)] = true
	Q1 = append(Q1, α)
	nstep := 0
	for len(Q1) > 0 {
		for _, cur := range Q1 {
			if cur == ω {
				return nstep
			}

			for _, δ := range dirs {
				nxt := cur.add(δ)
				if g.inbounds(nxt) && !seen[idx(nxt)] && g.at(nxt, t) {
					seen[idx(nxt)] = true
					Q2 = append(Q2, nxt)
				}
			}
		}
		Q1, Q2 = Q2, Q1[:0] // swap queues + reset Q2
		nstep++
	}

	return -1
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
			if x := (Cell{r, c}); g.get(x) == t {
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
				sb.WriteString(" ∞ ")
				continue
			}
			fmt.Fprintf(&sb, "%2d ", g.get(x))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g *Grid) dumpat(t int) {
	var sb strings.Builder
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			if x := (Cell{r, c}); g.get(x) <= t {
				sb.WriteString(" # ")
			} else {
				sb.WriteString(" . ")
			}
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}

type Cell struct {
	r, c int
}

func (x Cell) add(y Cell) Cell {
	return Cell{x.r + y.r, x.c + y.c}
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
