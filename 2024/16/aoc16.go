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
	"slices"
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
	dist1 []int
	dist2 []int
	tiles map[int]struct{}
	data  []string
	start Cell
	goal  Cell
	best  int
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

	best, tiles := maze.search()

	p1, p2 := best, len(tiles)
	fmt.Println(p1, p2) // part 1 & 2
}

var DIRS = [4]Cell{{-1, 0, 0}, {0, 1, 0}, {1, 0, 0}, {0, -1, 0}}

func key(r, c, dir int) int {
	return r*MAXDIM*4 + c*4 + dir
}

// shortest path from start to end
func (m Maze) forward() Maze {
	H, W := len(m.data), len(m.data[0])

	var pq Heap
	heap.Push(&pq, Item{dist: 0, r: m.start.r, c: m.start.c, dir: 0})

	best := 0
	// dist := make(map[int]int, MAXLEN)
	dist := make([]int, MAXDIM*MAXDIM*4)
	seen := [MAXDIM][MAXDIM][4]bool{}

	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(Item)

		if seen[cur.r][cur.c][cur.dir] {
			continue
		}

		seen[cur.r][cur.c][cur.dir] = true

		if dist[key(cur.r, cur.c, cur.dir)] == 0 {
			dist[key(cur.r, cur.c, cur.dir)] = cur.dist
		}

		if cur.r == m.goal.r && cur.c == m.goal.c && best == 0 {
			best = cur.dist
		}

		// forward
		δr, δc := DIRS[cur.dir].r, DIRS[cur.dir].c
		rr, cc := cur.r+δr, cur.c+δc
		if rr >= 0 && rr < H && cc >= 0 && cc < W && m.data[rr][cc] != '#' {
			heap.Push(&pq, Item{dist: cur.dist + 1, r: rr, c: cc, dir: cur.dir})
		}

		// rotate clockwise
		heap.Push(&pq, Item{dist: cur.dist + 1000, r: cur.r, c: cur.c, dir: (cur.dir + 1) % 4})

		// rotate counterclockwise
		heap.Push(&pq, Item{dist: cur.dist + 1000, r: cur.r, c: cur.c, dir: (cur.dir + 3) % 4})
	}

	m.best = best
	m.dist1 = dist

	return m
}

// shortest path single source to all other cells
func (m Maze) backward() Maze {
	H, W := len(m.data), len(m.data[0])
	var pq Heap

	// push all directions from the end point
	for dir := range DIRS {
		heap.Push(&pq, Item{dist: 0, r: m.goal.r, c: m.goal.c, dir: dir})
	}

	dist := make([]int, MAXDIM*MAXDIM*4)
	seen := [MAXDIM][MAXDIM][4]bool{}

	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(Item)

		if seen[cur.r][cur.c][cur.dir] {
			continue
		}

		seen[cur.r][cur.c][cur.dir] = true

		if dist[key(cur.r, cur.c, cur.dir)] == 0 {
			dist[key(cur.r, cur.c, cur.dir)] = cur.dist
		}

		// move backwards (opposite direction)
		δr, δc := DIRS[(cur.dir+2)%4].r, DIRS[(cur.dir+2)%4].c
		rr, cc := cur.r+δr, cur.c+δc
		if rr >= 0 && rr < H && cc >= 0 && cc < W && m.data[rr][cc] != '#' {

			heap.Push(&pq, Item{dist: cur.dist + 1, r: rr, c: cc, dir: cur.dir})
		}

		// clockwise
		heap.Push(&pq, Item{dist: cur.dist + 1000, r: cur.r, c: cur.c, dir: (cur.dir + 1) % 4})

		// counterclockwise
		heap.Push(&pq, Item{dist: cur.dist + 1000, r: cur.r, c: cur.c, dir: (cur.dir + 3) % 4})
	}

	m.dist2 = dist

	return m
}

// find all optimal path tiles
func (m Maze) search() (int, []int) {
	m = m.forward()
	m = m.backward()

	best, d1, d2 := m.best, m.dist1, m.dist2
	tiles := make([]int, 0, 602)

	for i := range d1 {
		if d1[i]+d2[i] == best {
			i := i >> 2
			tiles = append(tiles, i)
		}
	}

	slices.Sort(tiles)
	return m.best, slices.Compact(tiles)
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
