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
)

type Cell struct {
	r, c, dir int
}

type Problem struct {
	start, goal Cell
}

var DIRS = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func main() {
	maze := make([]string, 0, MAXDIM)
	problem := Problem{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		if i := strings.Index(line, "S"); i >= 0 {
			problem.start = Cell{len(maze), i, 0}
		}
		if i := strings.Index(line, "E"); i >= 0 {
			problem.goal = Cell{len(maze), i, 0}
		}
		maze = append(maze, line)
	}

	best, dist1 := forward(maze, problem)
	dist2 := backward(maze, problem)

	tiles := search(dist1, dist2, best, len(maze), len(maze[0]))

	p1, p2 := best, len(tiles)
	fmt.Println(p1, p2) // part 1 & 2
}

// shortest path from start to end
func forward(maze []string, p Problem) (int, map[Cell]int) {
	var pq Heap
	heap.Push(&pq, struct {
		dist, r, c, dir int
	}{dist: 0, r: p.start.r, c: p.start.c, dir: 1})

	best := 0
	dist := make(map[Cell]int)
	seen := make(map[Cell]bool)

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(struct {
			dist, r, c, dir int
		})

		if seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] {
			continue
		}

		seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = true

		if _, exists := dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}]; !exists {
			dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = curr.dist
		}

		if curr.r == p.goal.r && curr.c == p.goal.c && best == 0 {
			best = curr.dist
		}

		// forward
		δr, δc := DIRS[curr.dir][0], DIRS[curr.dir][1]
		rr, cc := curr.r+δr, curr.c+δc
		if rr >= 0 && rr < len(maze) && cc >= 0 && cc < len(maze[0]) && maze[rr][cc] != '#' {
			heap.Push(&pq, struct {
				dist, r, c, dir int
			}{dist: curr.dist + 1, r: rr, c: cc, dir: curr.dir})
		}

		// rotate clockwise
		heap.Push(&pq, struct {
			dist, r, c, dir int
		}{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 1) % 4})

		// rotate counterclockwise
		heap.Push(&pq, struct {
			dist, r, c, dir int
		}{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 3) % 4})
	}

	return best, dist
}

// shortest path single source to all other cells
func backward(maze []string, p Problem) map[Cell]int {
	var pq Heap

	// push all directions from the end point
	for dir := 0; dir < 4; dir++ {
		heap.Push(&pq, struct {
			dist, r, c, dir int
		}{dist: 0, r: p.goal.r, c: p.goal.c, dir: dir})
	}

	// distance map to store the shortest distance for (r, c, dir)
	dist := make(map[Cell]int)

	// set to store visited states
	seen := make(map[Cell]bool)

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(struct {
			dist, r, c, dir int
		})

		if seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] {
			continue
		}

		seen[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = true

		if _, ok := dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}]; !ok {
			dist[Cell{r: curr.r, c: curr.c, dir: curr.dir}] = curr.dist
		}

		// move backwards (opposite direction)
		δr, δc := DIRS[(curr.dir+2)%4][0], DIRS[(curr.dir+2)%4][1]
		rr, cc := curr.r+δr, curr.c+δc
		if rr >= 0 && rr < len(maze) && cc >= 0 && cc < len(maze[0]) && maze[rr][cc] != '#' {
			heap.Push(&pq, struct {
				dist, r, c, dir int
			}{dist: curr.dist + 1, r: rr, c: cc, dir: curr.dir})
		}

		// clockwise
		heap.Push(&pq, struct {
			dist, r, c, dir int
		}{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 1) % 4})

		// counterclockwise
		heap.Push(&pq, struct {
			dist, r, c, dir int
		}{dist: curr.dist + 1000, r: curr.r, c: curr.c, dir: (curr.dir + 3) % 4})
	}

	return dist
}

// find all optimal path tiles
func search(d1, d2 map[Cell]int, best, H, W int) map[Cell]struct{} {
	bests := make(map[Cell]struct{})

	for r := 0; r < H; r++ {
		for c := 0; c < W; c++ {
			for dir := 0; dir < 4; dir++ {
				_, ok1 := d1[Cell{r: r, c: c, dir: dir}]
				_, ok2 := d2[Cell{r: r, c: c, dir: dir}]

				if ok1 && ok2 {
					if d1[Cell{r: r, c: c, dir: dir}]+d2[Cell{r: r, c: c, dir: dir}] == best {
						bests[Cell{r: r, c: c}] = struct{}{} // optimal!
					}
				}
			}
		}
	}

	return bests
}

// Heap for the priority queue
type Heap []struct {
	dist, r, c, dir int
}

// Implementing heap.Interface for MinHeap
func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(struct{ dist, r, c, dir int }))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
