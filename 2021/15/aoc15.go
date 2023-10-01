package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"os"
	"strings"
)

type cell struct {
	y, x, v int
}

type heap []*cell

func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool { return h[i].v < h[j].v }

func (h heap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *heap) Push(x interface{}) {
	c := x.(cell)
	*h = append(*h, &c)
}

func (h *heap) Pop() interface{} {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], nil
	return pop
}

type grid struct {
	d    [][]int
	h, w int
}

func newGrid() *grid {
	N, g := 128, new(grid)
	g.d = make([][]int, N)
	for i := range g.d {
		g.d[i] = make([]int, N)
	}
	return g
}

func (g *grid) get(y, x int) int {
	j, y := y/g.h, y%g.h
	i, x := x/g.w, x%g.w

	v, inc := g.d[y][x], j+i
	if v+inc > 9 {
		return v + inc - 9
	}
	return v + inc
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
}

func (g grid) String() string {
	var sb *strings.Builder
	for i := 0; i < g.h; i++ {
		fmt.Fprintln(sb, g.d[i][:g.w]) // ignore errors
	}
	return sb.String()
}

// safest computes the shortest distance between the upper-left and the lower-right of the grid.
// It uses a shortest-path tree from (0,0) built by a canonical Dijkstra. The grid is virtual:
// it can be arbitrarily scaled up using a factor.
func safest(g *grid, factor int) int {
	const MaxInt = int(^uint(0) >> 1)

	h, w := factor*g.h, factor*g.w

	dist := make([][]int, h)
	for j := range dist {
		dist[j] = make([]int, w)
		for i := range dist[j] {
			dist[j][i] = MaxInt
		}
	}

	δy := []int{+0, 1, 0, -1}
	δx := []int{-1, 0, 1, +0}

	heap := make(heap, 0, 1024)
	hp.Init(&heap)

	hp.Push(&heap, cell{})
	dist[0][0] = g.get(0, 0)

	for heap.Len() > 0 {
		v := hp.Pop(&heap).(*cell)
		dist[v.y][v.x] = v.v

		if v.y == h-1 && v.x == w-1 {
			break
		}

		for i := range δy {
			u := cell{y: v.y + δy[i], x: v.x + δx[i]}
			if !(u.y >= 0 && u.y < h && u.x >= 0 && u.x < w) {
				continue // discard out of bounds
			}

			if dist[u.y][u.x] > dist[v.y][v.x]+g.get(u.y, u.x) {
				dist[u.y][u.x] = dist[v.y][v.x] + g.get(u.y, u.x)
				hp.Push(&heap, cell{u.y, u.x, dist[u.y][u.x]})
			}
		}
	}

	return dist[h-1][w-1] - dist[0][0]
}

func main() {
	cave := newGrid()

	j, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for i, b := range input.Bytes() {
			cave.d[j][i] = int(b - '0')
		}
		j++
	}
	cave.redim(j, j)

	fmt.Println(safest(cave, 1)) // part1
	fmt.Println(safest(cave, 5)) // part2
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
