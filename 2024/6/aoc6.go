// aoc6.go --
// advent of code 2024 day 6
//
// https://adventofcode.com/2024/day/6
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-6: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

func main() {
	var origin Point
	matrix := make([][]rune, 0, 130)

	j := 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		matrix = append(matrix, []rune(line))

		if i := strings.Index(line, "^"); i != -1 {
			origin = Point{j, i}
			matrix[j][i] = '.'
		}

		j++
	}

	maze := Maze{scan(matrix), origin}
	_, path := maze.run()
	_ = path

	// cells := path.cells()
	// count1 := len(cells)

	count1, count2 := 0, 0
	// for _, p := range cells {
	// 	maze.add(p)
	// 	if cycle, _ := maze.run(); cycle {
	// 		count2++
	// 	}
	// 	maze.del(p)
	// }

	fmt.Println(count1, count2)
}

type Point struct {
	y, x int
}

type Matrix struct {
	rows   [][]int
	cols   [][]int
	blocks []Point
}

func scan(matrix [][]rune) Matrix {
	H, W := len(matrix), len(matrix[0])

	// Initialize rows and cols
	rows := make([][]int, H)
	for j := range rows {
		rows[j] = make([]int, 0, W/5)
	}
	cols := make([][]int, W)
	for i := range cols {
		cols[i] = make([]int, 0, H/5)
	}
	blocks := make([]Point, 0, H*W/20)

	// Iterate through the matrix to find '#' and populate rows and cols
	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			if matrix[j][i] == '#' {
				rows[j] = append(rows[j], i) // Append column index to the row
				cols[i] = append(cols[i], j) // Append row index to the column
				blocks = append(blocks, Point{j, i})
			}
		}

	}

	for i := 0; i < W; i++ {
		cols[i] = append([]int{-2}, cols[i]...)
		cols[i] = append(cols[i], H+1)
	}

	for j := 0; j < H; j++ {
		rows[j] = append([]int{-2}, rows[j]...)
		rows[j] = append(rows[j], W+1)
	}

	return Matrix{rows, cols, blocks}
}

type Maze struct {
	Matrix
	o Point
}

func (mat Matrix) width() int {
	return len(mat.cols)
}

func (mat Matrix) height() int {
	return len(mat.rows)
}

func (mat Matrix) inbounds(p Point) bool {
	return 0 <= p.x && p.x < mat.height() && 0 <= p.y && p.y < mat.width()
}

func (mat Matrix) next(p Point, heading int) (Point, int) {

	switch heading {
	case NORTH:
		j, _ := slices.BinarySearch(mat.cols[p.x], p.y)
		p.y = mat.cols[p.x][j-1] + 1
		heading = EAST
	case EAST:
		i, _ := slices.BinarySearch(mat.rows[p.y], p.x)
		p.x = mat.rows[p.y][i] - 1
		heading = SOUTH
	case SOUTH:
		j, _ := slices.BinarySearch(mat.cols[p.x], p.y)
		p.y = mat.cols[p.x][j] - 1
		heading = WEST
	case WEST:
		i, _ := slices.BinarySearch(mat.rows[p.y], p.x)
		p.x = mat.rows[p.y][i-1] + 1
		heading = NORTH
	}

	return p, heading
}

// func (m Maze) add(p Point) {
// 	m.rows[p.y] = insert(m.rows[p.y], p.x)
// 	m.cols[p.x] = insert(m.cols[p.x], p.y)
// }

// func (m Maze) del(p Point) {
// 	m.rows[p.y] = delete(m.rows[p.y], p.x)
// 	m.cols[p.x] = delete(m.cols[p.x], p.y)
// }

// // Binary insertion function to keep slices sorted
// func insert(slice []int, value int) []int {
// 	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
// 	return append(slice[:i], append([]int{value}, slice[i:]...)...)
// }

// // Binary deletion function to keep slices sorted
// func delete(slice []int, value int) []int {
// 	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
// 	if i < len(slice) && slice[i] == value {
// 		slice = append(slice[:i], slice[i+1:]...)
// 	}
// 	return slice
// }

func (m Maze) run() (bool, Path) {
	H, W := m.height(), m.width()

	// run the maze
	seen := make([][4]int, H*W)

	clock := func() func() int {
		start := 0
		return func() int {
			now := start
			start++
			return now
		}
	}()

	heading, to := NORTH, NORTH
	stamp, pre, cur := clock(), m.o, m.o

	mark := func() {
		xmin, xmax := max(min(pre.x, cur.x), 0), min(max(pre.x, cur.x)+1, W)
		ymin, ymax := max(min(pre.y, cur.y), 0), min(max(pre.y, cur.y)+1, H)

		switch to {
		case NORTH, SOUTH:
			fmt.Println("see y:", ymin, ymax, cur.x, to, ymax-ymin)
			for y := ymin; y < ymax; y++ {
				i := y*W + cur.x
				seen[i][to] = stamp
			}
		case EAST, WEST:
			fmt.Println("see x:", xmin, xmax, cur.y, to, xmax-xmin)
			for x := xmin; x < xmax; x++ {
				i := cur.y*W + x
				seen[i][to] = stamp
			}
		}
	}

	for m.inbounds(cur) {
		stamp = clock()
		i := cur.y*W + cur.x
		seen[i][heading] = stamp

		to, pre = heading, cur
		cur, heading = m.next(cur, heading)

		// update the seen matrix
		mark()
	}
	// update the exit segment
	mark()

	now := time.Now()
	count1, count2, count3 := 0, 0, 0
	for i := range seen {
		if slices.Max(seen[i][:]) > 0 {
			count1++

			// this looks for a sure cycle: one that branches in the past
			for _, h0 := range []int{NORTH, EAST, SOUTH, WEST} {
				h1, h2 := (h0+1)%4, (h0+2)%4 // (NORTH, EAST, SOUTH) -> (EAST, SOUTH, WEST) -> ...
				if seen[i][h0] > 0 {
					nxt, _ := m.next(Point{i / W, i % W}, h1)
					if m.inbounds(nxt) {
						snxt := seen[nxt.y*W+nxt.x]
						switch {
						case 0 < snxt[h2] && snxt[h2] < seen[i][h0]:
							fmt.Println("from:", h0, "to:", h1, "coo:", i/W, i%W)
							count2++
						case slices.Max(snxt[:]) == 0:
							count3++
						}
					}
				}
			}
		}
	}
	elapsed := time.Since(now)
	fmt.Println("counts:", count1, count2, count3, elapsed)

	return false, []Point{}
}

type Path []Point

func (p Path) cells() Path {
	cells := make(Path, 0, 5000)

	for i := range p[:len(p)-1] {
		start := p[i]
		end := p[i+1]

		// Horizontal segment
		if start.y == end.y {
			start.x, end.x = min(start.x, end.x), max(start.x, end.x)
			for x := start.x; x <= end.x; x++ {
				cells = append(cells, Point{start.y, x})
			}
		} else {
			start.y, end.y = min(start.y, end.y), max(start.y, end.y)
			for y := start.y; y <= end.y; y++ {
				cells = append(cells, Point{y, start.x})
			}
		}
	}

	slices.SortFunc(cells, func(a, b Point) int {
		if a.y == b.y {
			return a.x - b.x
		}
		return a.y - b.y
	})

	cells = slices.CompactFunc(cells, func(a, b Point) bool {
		return a == b
	})

	return cells
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
