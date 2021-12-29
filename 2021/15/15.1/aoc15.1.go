package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"os"
	"sort"
	"strings"
)

type cell struct {
	x, y, v int
}

func (a *cell) smaller(b *cell) bool {
	if a.v == b.v {
		if a.x != b.x {
			return a.x < b.x
		}
		return a.y < b.y
	}
	return a.v < b.v
}

type heap []*cell

func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool {
	return h[i].smaller(h[j])
}

func (h heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heap) Push(x interface{}) {
	c := x.(cell)
	*h = append(*h, &c)
}

func (h *heap) Pop() interface{} {
	q, n := *h, len(*h)-1
	c := q[n] // last
	q[n], *h = nil, q[:n]
	return c
}

func (h *heap) remove(c cell) {
	i := sort.Search(len(*h), func(i int) bool { return (&c).smaller((*h)[i]) })
	if i < len(*h) {
		*h = append((*h)[:i], (*h)[i+1:]...)
		hp.Init(h)
	}
}

type grid struct {
	d    [][]int
	w, h int
}

func (g *grid) String() string {
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

func (g grid) get(x, y int) int {
	return g.d[y][x]
}

func (g *grid) redim(w, h int) {
	g.w, g.h = w, h
}

func shortest(g *grid) int {
	const MaxInt = int(^uint(0) >> 1)

	w, h := g.w, g.h

	dist := make([][]int, h)
	for j := range dist {
		dist[j] = make([]int, w)
		for i := 0; i < w; i++ {
			dist[j][i] = MaxInt
		}
	}

	δx := []int{-1, 0, 1, 0}
	δy := []int{0, 1, 0, -1}

	valid := func(x, y int) bool {
		return !(x < 0 || x >= g.w || y < 0 || y >= g.h)
	}

	heap := make(heap, 0, 16364)
	hp.Init(&heap)

	hp.Push(&heap, cell{0, 0, 0})
	dist[0][0] = g.get(0, 0)

	for heap.Len() > 0 {
		v := hp.Pop(&heap).(*cell)
		dist[v.y][v.x] = v.v

		if v.y == h-1 && v.x == w-1 {
			break
		}

		for i := 0; i < 4; i++ {
			u := cell{}
			u.x = v.x + δx[i]
			u.y = v.y + δy[i]

			if !valid(u.x, u.y) {
				continue
			}

			if dist[u.y][u.x] > dist[v.y][v.x]+g.get(u.x, u.y) {
				dist[u.y][u.x] = dist[v.y][v.x] + g.get(u.x, u.y)
				hp.Push(&heap, cell{u.x, u.y, dist[u.y][u.x]})
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

	fmt.Println(shortest(cave))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
