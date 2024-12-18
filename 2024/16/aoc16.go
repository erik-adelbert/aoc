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
	"sync"
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
	data  [][]byte
	start Cell
	goal  Cell
	best  int
}

func main() {
	var start, goal Cell

	data := make([][]byte, 0, MAXDIM)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		if i := strings.Index(line, "S"); i >= 0 {
			start = Cell{r: len(data), c: i}
		}
		if i := strings.Index(line, "E"); i >= 0 {
			goal = Cell{r: len(data), c: i}
		}
		data = append(data, []byte(line))
	}

	maze := &Maze{data: data, start: start, goal: goal}
	score, tiles := maze.solve()

	fmt.Println(score, len(tiles)) // part 1 & 2
}

var DIRS = [4]Cell{{r: -1, c: 0}, {r: 0, c: 1}, {r: 1, c: 0}, {r: 0, c: -1}}

func key(x Cell) int {
	return x.r*MAXDIM*4 + x.c*4 + x.dir
}

func (m Maze) String() string {
	var sb strings.Builder
	for _, line := range m.data {
		sb.WriteString(string(line))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Maze) prune() {
	H, W := len(m.data), len(m.data[0])

	queue := make([]Cell, 0, 395)
	seen := make([][]bool, H)
	for i := range seen {
		seen[i] = make([]bool, W)
	}

	for r := 1; r < H-1; r++ {
		for c := 1; c < W-1; c++ {
			if m.data[r][c] == '.' {
				nopen := 0
				for _, δ := range DIRS {
					rr, cc := r+δ.r, c+δ.c
					if m.data[rr][cc] == '.' || m.data[rr][cc] == 'S' || m.data[rr][cc] == 'E' {
						nopen++
					}
				}
				if nopen <= 1 {
					queue = append(queue, Cell{r: r, c: c})
					seen[r][c] = true
				}
			}
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		m.data[cur.r][cur.c] = '#'

		for _, dir := range DIRS {
			rr, cc := cur.r+dir.r, cur.c+dir.c
			if m.data[rr][cc] == '.' && !seen[rr][cc] {
				nopen := 0
				for _, δ := range DIRS {
					nr, nc := rr+δ.r, cc+δ.c
					if m.data[nr][nc] == '.' || m.data[nr][nc] == 'S' || m.data[nr][nc] == 'E' {
						nopen++
					}
				}
				if nopen <= 1 {
					queue = append(queue, Cell{r: rr, c: cc})
					seen[rr][cc] = true
				}
			}
		}
	}
}

// shortest path from start to end
func (m *Maze) forward() {
	H, W := len(m.data), len(m.data[0])

	var pq Heap
	heap.Push(&pq, State{Cell{r: m.start.r, c: m.start.c, dir: 0}, 0})

	best := 0
	dist := make([]int, MAXDIM*MAXDIM*4)
	seen := make([]bool, MAXDIM*MAXDIM*4)

	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(State)

		if seen[key(cur.Cell)] {
			continue
		}

		seen[key(cur.Cell)] = true

		if dist[key(cur.Cell)] == 0 {
			dist[key(cur.Cell)] = cur.dist
		}

		if cur.r == m.goal.r && cur.c == m.goal.c && best == 0 {
			best = cur.dist
		}

		// forward
		δr, δc := DIRS[cur.dir].r, DIRS[cur.dir].c
		rr, cc := cur.r+δr, cur.c+δc
		if rr >= 0 && rr < H && cc >= 0 && cc < W && m.data[rr][cc] != '#' {
			heap.Push(&pq, State{Cell{r: rr, c: cc, dir: cur.dir}, cur.dist + 1})
		}

		// rotate clockwise
		heap.Push(&pq, State{Cell{r: cur.r, c: cur.c, dir: (cur.dir + 1) % 4}, cur.dist + 1000})

		// rotate counterclockwise
		heap.Push(&pq, State{Cell{r: cur.r, c: cur.c, dir: (cur.dir + 3) % 4}, cur.dist + 1000})
	}

	m.best = best
	m.dist1 = dist
}

// shortest path single source to all other cells
func (m *Maze) backward() {
	H, W := len(m.data), len(m.data[0])
	var pq Heap

	// push all directions from the end point
	for dir := range DIRS {
		heap.Push(&pq, State{Cell{r: m.goal.r, c: m.goal.c, dir: dir}, 0})
	}

	dist := make([]int, MAXDIM*MAXDIM*4)
	seen := make([]bool, MAXDIM*MAXDIM*4)

	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(State)

		if seen[key(cur.Cell)] {
			continue
		}

		seen[key(cur.Cell)] = true

		if dist[key(cur.Cell)] == 0 {
			dist[key(cur.Cell)] = cur.dist
		}

		// move backwards
		δr, δc := DIRS[(cur.dir+2)%4].r, DIRS[(cur.dir+2)%4].c
		rr, cc := cur.r+δr, cur.c+δc
		if rr >= 0 && rr < H && cc >= 0 && cc < W && m.data[rr][cc] != '#' {
			heap.Push(&pq, State{Cell{r: rr, c: cc, dir: cur.dir}, cur.dist + 1})
		}

		// clockwise
		heap.Push(&pq, State{Cell{r: cur.r, c: cur.c, dir: (cur.dir + 1) % 4}, cur.dist + 1000})

		// counterclockwise
		heap.Push(&pq, State{Cell{r: cur.r, c: cur.c, dir: (cur.dir + 3) % 4}, cur.dist + 1000})
	}

	m.dist2 = dist

	return
}

// find all optimal path tiles
func (m *Maze) solve() (int, []int) {
	var wg sync.WaitGroup

	m.prune()

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.forward()
	}()

	m.backward()
	wg.Wait()

	best, d1, d2 := m.best, m.dist1, m.dist2
	tiles := make([]int, 0, 602)

	for i := range d1 {
		if d1[i]+d2[i] == best {
			i := i >> 2
			tiles = append(tiles, i)
		}
	}

	slices.Sort(tiles)
	return best, slices.Compact(tiles) // unique tiles
}

type State struct {
	Cell
	dist int
}

// heap as priority queue
type Heap []State

// heap.Interface for MinHeap
func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(State))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
