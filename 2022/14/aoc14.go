// aoc14.go --
// advent of code 2022 day 14
//
// https://adventofcode.com/2022/day/14
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-14: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var world worldmap

// world facts
const (
	// world is contained in 256x512, 512 is ok for +inf
	INF = 512
	// world is uselessly translated too far east
	XOFF = 300
	// sand is poured from X0=200, Y0=0 in our translated world
	X0 = 500 - XOFF
	Y0 = 0
)

func main() {
	box := AABB{{INF, 0}, {0, 0}}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		box.merge(mkworld(input.Text()))
	}

	// offset to sand boundaries
	box[Min][X]--
	box[Max][X]++
	box[Max][Y]++

	depth := box[Max][Y]

	part1 := fill(depth+1, depth, box)
	part2 := part1 + 1 + fill(depth+2, 0, box)

	fmt.Println(part1, part2)

	// !!rise your term resolution!!
	// uncomment next line for visualization
	// fmt.Println(world)
}

func fill(floor int, depth int, box AABB) int {
	// free() slices world where the action is (ie. the rocks)
	// the world physics garantee that the final shape will
	// always be an isosceles right triangle centered at X=200
	//
	// we don't want to simulate aisles as we can compute them
	// easily
	free := func(p XY) bool {
		return world.get(p) == 0 && p[Y] < floor && box.contains(p)
	}

	stack := make([]XY, 0, 256)

	push := func(p XY) {
		stack = append(stack, p)
	}

	pop := func() XY {
		if len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			return p
		}
		return XY{X0, Y0}
	}

	// dfs iterator
	next := func() XY {
		cur, nxt := XY{INF, INF}, pop() // backtrack
		for cur != nxt {
			cur = nxt
			// try in turn if stuck: south, sw, se
			for _, δ := range []XY{{0, 1}, {-1, 1}, {1, 1}} {
				if free(cur.add(δ)) {
					push(cur) // backup last moving location
					nxt = cur.add(δ)
					break
				}
			}
		}
		return cur
	}

	// pour and count sand grains
	cnt := 0
	dst := next()
	for dst[Y] != depth {
		world.set(dst, '.')
		cnt++
		// on world boundaries, bring back sliced parts
		if dst[X] == box[Min][X] || dst[X] == box[Max][X] {
			cnt += box[Max][Y] - dst[Y]
		}
		dst = next()
	}
	return cnt
}

func mkworld(s string) AABB {
	wall := make([]XY, 0, 128)
	box := AABB{{INF, INF}, {0, 0}}
	for _, segs := range strings.Split(s, " -> ") {
		var seg XY
		for i, s := range strings.Split(segs, ",") {
			seg[i] = atoi(s)
		}
		seg[0] -= XOFF // translate x to fit

		box.add(seg)
		wall = append(wall, seg)
	}
	for i := range wall[:len(wall)-1] {
		a, b := wall[i], wall[i+1]

		δ := a.cmp(b)
		for x := a; x != b; x = x.add(δ) {
			world.set(x, '#')
		}
		world.set(b, '#')
	}
	return box
}

type worldmap [256][512]byte

func (w worldmap) String() string {
	var sb strings.Builder
	for _, row := range w {
		for _, b := range row {
			if b == 0 {
				b = ' '
			}
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (w *worldmap) set(p XY, v byte) {
	w[p[Y]][p[X]] = v
}

func (w *worldmap) get(p XY) byte {
	return w[p[Y]][p[X]]
}

// indices for XY
const (
	X = iota
	Y
)

// XY is a 2D point
type XY [2]int

func (a XY) add(b XY) XY {
	a[X] += b[X]
	a[Y] += b[Y]
	return a
}

func (a XY) cmp(b XY) XY {
	return XY{cmp(b[X], a[X]), cmp(b[Y], a[Y])}
}

// indices for AABB
const (
	Min = iota
	Max
)

// AABB is axis aligned bounding box
type AABB [2]XY

func (a AABB) contains(p XY) bool {
	// a[Min][Y] is 0
	return a[Min][X] <= p[X] &&
		a[Max][X] >= p[X] &&
		a[Max][Y] >= p[Y]
}

func (a *AABB) add(p XY) {
	a[Min][X] = min(a[Min][X], p[X])
	a[Min][Y] = min(a[Min][Y], p[Y])
	a[Max][X] = max(a[Max][X], p[X])
	a[Max][Y] = max(a[Max][Y], p[Y])
}

func (a *AABB) merge(b AABB) {
	a[Min][X] = min(a[Min][X], b[Min][X])
	a[Min][Y] = min(a[Min][Y], b[Min][Y])
	a[Max][X] = max(a[Max][X], b[Max][X])
	a[Max][Y] = max(a[Max][Y], b[Max][Y])
}

// strconv.Atoi modified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

func cmp(a, b int) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
