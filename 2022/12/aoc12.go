// aoc12.go --
// advent of code 2022 day 12
//
// https://adventofcode.com/2022/day/12
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-12: initial commit

package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"os"
)

// elevation map
type grid struct {
	d    [][]int
	h, w int
}

func main() {
	area := newGrid()
	var one, all []*cell
	var end *cell

	h, w, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Bytes()
		for i, b := range input.Bytes() {
			area.d[h][i] = int(b)

			switch b {
			case 'S':
				// part1
				area.d[h][i] = int('a')
				one = append(one, &cell{h, i, 0})
			case 'a':
				// part2
				area.d[h][i] = int('a')
				all = append(all, &cell{h, i, 0})
			case 'E':
				area.d[h][i] = int('z')
				end = &cell{h, i, 0}
			}
		}
		w = len(line)
		h++
	}
	area.redim(h, w)

	all = append(one, all...)

	fmt.Println(area.shortest(all, end)) // part 1&2
}

// solve computes the shortest distance between e and all
// the cells of the grid. It uses a shortest-path tree
// built by a canonical Dijkstra.
func (g *grid) shortest(start []*cell, end *cell) (int, int) {
	const MaxInt = int(^uint(0) >> 1)

	// paths := []int{}

	dist := make([][]int, g.h)
	for j := range dist {
		dist[j] = make([]int, g.w)
		for i := range dist[j] {
			dist[j][i] = MaxInt
		}
	}

	heap := make(heap, 0, 1024)

	δy := []int{+0, 1, 0, -1}
	δx := []int{-1, 0, 1, +0}

	heap = heap[:0]
	hp.Init(&heap)

	// push end as single source
	hp.Push(&heap, end)
	dist[end.y][end.x] = 0

	// dijkstra
	for heap.Len() > 0 {
		a := hp.Pop(&heap).(*cell)
		dist[a.y][a.x] = a.d

		for i := range δy {
			b := cell{y: a.y + δy[i], x: a.x + δx[i]}

			outside := func() bool {
				return !(b.y >= 0 && b.y < g.h && b.x >= 0 && b.x < g.w)
			}

			// constraint from day12
			low := func() bool {
				return g.get(a.y, a.x)-g.get(b.y, b.x) > 1
			}

			switch {
			case outside() || low():
				// discard
				continue
			case dist[b.y][b.x] > dist[a.y][a.x]+1:
				dist[b.y][b.x] = dist[a.y][a.x] + 1
				hp.Push(&heap, &cell{b.y, b.x, dist[b.y][b.x]})
			}
		}
	}

	min := start[0]
	for _, s := range start[1:] {
		if dist[s.y][s.x] < dist[min.y][min.x] {
			min = s
		}
	}

	return dist[start[0].y][start[0].x], dist[min.y][min.x]
}

func newGrid() *grid {
	H, W, g := 64, 128, new(grid)
	g.d = make([][]int, H)
	for i := range g.d {
		g.d[i] = make([]int, W)
	}
	return g
}

func (g grid) get(y, x int) int {
	return g.d[y][x]
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
}

// types and interface used with hp
type cell struct {
	y, x int
	d    int
}

type heap []*cell

func (h heap) Len() int           { return len(h) }
func (h heap) Less(i, j int) bool { return h[i].d < h[j].d }
func (h heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *heap) Push(x any) {
	c := x.(*cell)
	*h = append(*h, c)
}

func (h *heap) Pop() any {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], nil
	return pop
}
