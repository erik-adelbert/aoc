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

func less(a, b *cell) bool {
	if a.v != b.v {
		return a.v < b.v
	}
	if a.x != b.x {
		return a.x < b.x
	}
	return a.y < b.y
}

type heap []*cell

func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool {
	return less(h[i], h[j])
}

func (h heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

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

func (g grid) String() string {
	rows := make([]string, g.h)
	for j := 0; j < g.h; j++ {
		rows[j] = fmt.Sprintln(g.d[j][:g.w])
	}
	return strings.Join(rows, "")
}

func newGrid() *grid {
	var g grid
	g.d = make([][]int, 128)
	for i := 0; i < 128; i++ {
		g.d[i] = make([]int, 128)
	}
	return &g
}

func (g grid) get(y, x int) int {
	next := func(v, inc int) int {
		if v+inc > 9 {
			return v + inc - 9
		}
		return v + inc
	}

	j, y := y/g.h, y%g.h
	i, x := x/g.w, x%g.w

	return next(g.d[y][x], j+i)
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
}

func safest(g *grid, factor int) int {
	const MaxInt = int(^uint(0) >> 1)

	h, w := factor*g.h, factor*g.w

	dist := make([][]int, h)
	for j := range dist {
		dist[j] = make([]int, w)
		for i := 0; i < w; i++ {
			dist[j][i] = MaxInt
		}
	}

	δy := []int{+0, 1, 0, -1}
	δx := []int{-1, 0, 1, +0}

	valid := func(y, x int) bool {
		return !(y < 0 || y >= h || x < 0 || x >= w)
	}

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

		for i := 0; i < 4; i++ {
			u := cell{y: v.y + δy[i], x: v.x + δx[i]}

			if !valid(u.y, u.x) {
				continue
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
