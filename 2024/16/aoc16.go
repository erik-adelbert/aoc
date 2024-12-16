// aoc16.go --
// advent of code 2024 day 16
//
// https://adventofcode.com/2024/day/16
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-16: initial commit

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strings"
)

const (
	MAXDIM = 141
	MAXLEN = 40468
)

type Cell struct {
	r, c, dir int
}

type Maze struct {
	data        []string
	start, goal Cell
}

func main() {
	var start, goal Cell

	data := make([]string, 0, MAXDIM)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		if i := strings.Index(line, "S"); i >= 0 {
			start = Cell{len(data), i, 0}
		}
		if i := strings.Index(line, "E"); i >= 0 {
			goal = Cell{len(data), i, 0}
		}
		data = append(data, line)
	}

	maze := Maze{data: data, start: start, goal: goal}

	best, dist1 := forward(maze)
	dist2 := backward(maze)

	tiles := search(dist1, dist2, best, len(data), len(data[0]))

	p1, p2 := best, len(tiles)
	fmt.Println(p1, p2) // part 1 & 2

	fmt.Println(len(dist1), len(dist2), len(tiles))
}

var DIRS = [4]Cell{{-1, 0, 0}, {0, 1, 0}, {1, 0, 0}, {0, -1, 0}}

// shortest path from start to end
func forward(m Maze) (int, map[Cell]int) {
	H, W := len(m.data), len(m.data[0])

	var pq Heap
	heap.Push(&pq, Item{dist: 0, r: m.start.r, c: m.start.c, dir: 0})

	best := 0
	dist := make(map[Cell]int, MAXLEN)
	seen := make(map[Cell]bool, MAXLEN)

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(Item)

		if seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] {
			continue
		}

		seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = true

		if _, exists := dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}]; !exists {
			dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = curr.dist
		}

		if curr.r == m.goal.r && curr.c == m.goal.c && best == 0 {
			best = curr.dist
		}

		// forward
		δr, δc := DIRS[curr.dir].r, DIRS[curr.dir].c
		rr, cc := curr.r+δr, curr.c+δc
		if rr >= 0 && rr < H && cc >= 0 && cc < W && m.data[rr][cc] != '#' {
			heap.Push(&pq, Item{dist: curr.dist + 1, r: rr, c: cc, dir: curr.dir})
		}

		// rotate clockwise
		heap.Push(&pq, Item{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 1) % 4})

		// rotate counterclockwise
		heap.Push(&pq, Item{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 3) % 4})
	}

	return best, dist
}

// shortest path single source to all other cells
func backward(m Maze) map[Cell]int {
	H, W := len(m.data), len(m.data[0])
	var pq Heap

	// push all directions from the end point
	for dir := range DIRS {
		heap.Push(&pq, Item{dist: 0, r: m.goal.r, c: m.goal.c, dir: dir})
	}

	dist := make(map[Cell]int, MAXLEN)
	seen := make(map[Cell]bool, MAXLEN)

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(Item)

		if seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] {
			continue
		}

		seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = true

		if _, ok := dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}]; !ok {
			dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = curr.dist
		}

		// move backwards (opposite direction)
		δr, δc := DIRS[(curr.dir+2)%4].r, DIRS[(curr.dir+2)%4].c
		rr, cc := curr.r+δr, curr.c+δc
		if rr >= 0 && rr < H && cc >= 0 && cc < W && m.data[rr][cc] != '#' {
			heap.Push(&pq, Item{dist: curr.dist + 1, r: rr, c: cc, dir: curr.dir})
		}

		// clockwise
		heap.Push(&pq, Item{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 1) % 4})

		// counterclockwise
		heap.Push(&pq, Item{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 3) % 4})
	}

	return dist
}

// find all optimal path tiles
func search(d1, d2 map[Cell]int, best, H, W int) map[Cell]struct{} {
	bests := make(map[Cell]struct{}, 504)

	for k, d1 := range d1 {
		if d2, ok := d2[k]; ok && d1+d2 == best {
			k := Cell{r: k.r, c: k.c}
			bests[k] = struct{}{}
		}
	}

	return bests
}

type Item struct {
	dist, r, c, dir int
}

// Heap for the priority queue
type Heap []Item

// Implementing heap.Interface for MinHeap
func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
