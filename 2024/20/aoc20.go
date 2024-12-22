// aoc20.go --
// advent of code 2024 day 20
//
// https://adventofcode.com/2024/day/20
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-20: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	MAXDIM    = 141
	MAXTRACK  = 9500
	MAXSHORTS = 1011500
	MAXLOCALS = 3000
	MAXBATCH  = 15
	MAXWORKER = 4
)

type Cell struct {
	r, c int
}

type Track []Cell

type Maze struct {
	data        [][]byte
	dist        []int
	tree        *KDTree
	track       Track
	start, goal Cell
}

func (c Cell) String() string {
	return fmt.Sprintf("(%d, %d)", c.r, c.c)
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

	var shorts1, shorts2 []Shortcut
	maze := newMaze(data, start, goal)

	shorts1 = maze.shortcut(2, 100)
	shorts2 = maze.shortcut(20, 100)

	fmt.Println(len(shorts1), len(shorts2))
}

func newMaze(data [][]byte, start, goal Cell) *Maze {
	m := &Maze{
		data:  data,
		start: start,
		goal:  goal,
	}

	m.track = m.mktrack()
	m.dist = m.mkdist(m.track)
	m.tree = &KDTree{root: mktree(m.track, 0)}

	return m
}

type Shortcut struct {
	start, end Cell
	time       int
}

func (m *Maze) shortcut(tmax, lim int) []Shortcut {
	shorts := make([]Shortcut, 0, MAXSHORTS)
	work := make(chan []Cell, len(m.track)/10) // Batch cells into ranges.
	chunks := make(chan []Shortcut, len(m.track)/10)

	dist := m.getdist

	// Worker function to process batches of points
	shorter := func() {
		locals := make([]Shortcut, 0, MAXLOCALS)
		for batch := range work {
			locals = locals[:0]
			for _, x := range batch {

				// Core logic
				dx := dist(x)
				for _, xx := range m.tree.query(x, tmax) {
					md := manh(x, xx)
					dxx := md + dist(xx)
					if short := dx - dxx; short >= lim {
						locals = append(locals, Shortcut{start: x, end: xx, time: short})
					}
				}

			}
			chunks <- locals
		}
	}

	for i := 0; i < MAXWORKER; i++ {
		go shorter()
	}

	// Send work to the channel in batches
	const bsize = MAXBATCH
	for i := 0; i < len(m.track); i += bsize {
		end := i + bsize
		if end > len(m.track) {
			end = len(m.track)
		}
		work <- m.track[i:end]
	}
	close(work)

	// Collect
	for i := 0; i < len(m.track)/bsize+1; i++ {
		shorts = append(shorts, <-chunks...)
	}
	close(chunks)

	return shorts
}

func (m *Maze) mktrack() Track {
	H, W := len(m.data), len(m.data[0])
	seen := make([][]bool, H)
	for r := range seen {
		seen[r] = make([]bool, W)
	}

	track := make([]Cell, 0, MAXTRACK)
	cur, dist := m.start, 0

	for {
		// Mark the cell as visited and add it to the path
		seen[cur.r][cur.c] = true
		track = append(track, cur)

		// If this is the goal, break
		if cur.r == m.goal.r && cur.c == m.goal.c {
			break
		}

		// Explore neighbors: Up, Down, Left, Right
		dirs := []Cell{
			{r: -1, c: 0}, {r: 1, c: 0}, {r: 0, c: -1}, {r: 0, c: 1},
		}

		for _, δ := range dirs {
			rr, cc := cur.r+δ.r, cur.c+δ.c

			// If within bounds, not a wall, and not visited, move to the next cell
			if rr >= 0 && rr < H && cc >= 0 && cc < W &&
				m.data[rr][cc] != '#' && !seen[rr][cc] {
				cur = Cell{r: rr, c: cc}
				dist++
				break
			}
		}
	}

	return track
}

func (m *Maze) getdist(x Cell) int {
	W := len(m.data[0])
	return m.dist[x.r*W+x.c]
}

func (m *Maze) mkdist(track Track) []int {
	H, W := len(m.data), len(m.data[0])
	dist := make([]int, H*W)
	for i := range dist {
		dist[i] = -1
	}

	// Populate distances for the path
	for d, x := range track {
		dist[x.r*W+x.c] = len(track) - d - 1
	}

	return dist
}

func manh(a, b Cell) int {
	return abs(a.r-b.r) + abs(a.c-b.c)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *Maze) String() string {
	var sb strings.Builder

	sb.WriteString("Maze Data:\n")
	for r := range m.data {
		for c := range m.data[r] {
			if m.data[r][c] == '#' {
				sb.WriteString(" # ")
			} else if m.start.r == r && m.start.c == c {
				sb.WriteString(" S ")
			} else if m.goal.r == r && m.goal.c == c {
				sb.WriteString(" E ")
			} else {
				sb.WriteString(" . ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

type KDNode struct {
	left  *KDNode
	right *KDNode
	cell  Cell
	axis  int
}

type KDTree struct {
	root *KDNode
}

// mktree recursively builds the k-d tree
func mktree(track []Cell, depth int) *KDNode {
	if len(track) == 0 {
		return nil
	}

	axis := depth % 2 // Alternate between r (0) and c (1)
	mid := len(track) / 2

	// Sort cells based on the axis
	// Sorting by r-axis or c-axis depending on the current depth/axis
	if axis == 0 {
		slices.SortFunc(track, func(a, b Cell) int {
			return a.r - b.r
		})
	} else {
		slices.SortFunc(track, func(a, b Cell) int {
			return a.c - b.c
		})
	}

	// Create a new node and recursively build the left and right subtrees
	node := &KDNode{
		cell: track[mid],
		axis: axis,
	}
	node.left = mktree(track[:mid], depth+1)
	node.right = mktree(track[mid+1:], depth+1)

	return node
}

func (tree *KDTree) query(target Cell, tmax int) []Cell {

	// Range requery function to find all points within a Manhattan distance
	var requery func(*KDNode, Cell, int, *[]Cell)
	requery = func(node *KDNode, target Cell, tmax int, results *[]Cell) {
		if node == nil {
			return
		}

		// Calculate the Manhattan distance to the current node's point
		dist := manh(target, node.cell)
		if dist <= tmax {
			*results = append(*results, node.cell)
		}

		// Recursively search in the appropriate subtrees
		if node.axis == 0 { // Splitting by x-axis
			if target.r-tmax <= node.cell.r {
				requery(node.left, target, tmax, results)
			}
			if target.r+tmax >= node.cell.r {
				requery(node.right, target, tmax, results)
			}
		} else { // Splitting by y-axis
			if target.c-tmax <= node.cell.c {
				requery(node.left, target, tmax, results)
			}
			if target.c+tmax >= node.cell.c {
				requery(node.right, target, tmax, results)
			}
		}
	}
	var neighs []Cell

	requery(tree.root, target, tmax, &neighs)
	return neighs
}
