package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"os"
	"strings"
)

// types and interface used with hp
type cell struct {
	y, x int
	h    int
}

type heap []*cell

func (h heap) Len() int           { return len(h) }
func (h heap) Less(i, j int) bool { return h[i].h < h[j].h }
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

// elevation map
type grid struct {
	d    [][]int
	h, w int
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

func (g grid) String() string {
	var sb strings.Builder
	for i := 0; i < g.h; i++ {
		fmt.Fprintln(&sb, g.d[i][:g.w])
	}
	return sb.String()
}

// solve computes the shortest distance between s and e,
// two cells of the grid. It uses a shortest-path tree
// built by a canonical Dijkstra.
func (g *grid) shortest(start []*cell, end *cell) int {
	const MaxInt = int(^uint(0) >> 1)

	paths := []int{}

	dist := make([][]int, g.h)
	for j := range dist {
		dist[j] = make([]int, g.w)
	}

	heap := make(heap, 0, 1024)

	for _, s := range start {

		for j := range dist {
			for i := range dist[j] {
				dist[j][i] = MaxInt
			}
		}

		δy := []int{+0, 1, 0, -1}
		δx := []int{-1, 0, 1, +0}

		heap = heap[:0]
		hp.Init(&heap)

		// push start
		hp.Push(&heap, s)
		dist[s.y][s.x] = 0

		for heap.Len() > 0 {
			v := hp.Pop(&heap).(*cell)
			dist[v.y][v.x] = v.h

			if v.y == end.y && v.x == end.x {
				break
			}

			for i := range δy {
				u := cell{y: v.y + δy[i], x: v.x + δx[i]}

				outside := func() bool {
					return !(u.y >= 0 && u.y < g.h && u.x >= 0 && u.x < g.w)
				}

				high := func() bool {
					return g.get(u.y, u.x)-g.get(v.y, v.x) > 1
				}

				switch {
				// discard out of bounds & too high
				case outside() || high():
					continue
				case dist[u.y][u.x] > dist[v.y][v.x]+1:
					dist[u.y][u.x] = dist[v.y][v.x] + 1
					hp.Push(&heap, &cell{u.y, u.x, dist[u.y][u.x]})
				}
			}
		}
		paths = append(paths, dist[end.y][end.x]-dist[s.y][s.x])
	}

	return min(paths)
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
				area.d[h][i] = int('a')
				one = append(one, &cell{h, i, 0})
			case 'a':
				area.d[h][i] = int('a')
				all = append(all, &cell{h, i, 0})
			case 'E':
				area.d[h][i] = int('z')
				end = &cell{h, i, int('z')}
			}
		}
		w = len(line)
		h++
	}
	area.redim(h, w)

	fmt.Println(area.shortest(one, end)) // part1
	fmt.Println(area.shortest(all, end)) // part2
}

func min(A []int) int {
	min := A[0]
	for _, v := range A {
		if v < min {
			min = v
		}
	}
	return min
}
