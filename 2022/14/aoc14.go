package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// world map
var world [256][512]byte

// worl facts
const (
	// world is contained in 256x512
	// so 512 is ok for +inf
	INF = 512
	// world is uselessly translated too far east
	XOFF = 300
	// sand is poured from X=200, Y=0 in our translated world
	XORG = 500 - XOFF
	YORG = 0
)

func main() {
	box := AABB{{INF, 0}, {0, 0}}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		box.merge(mkworld(input.Text()))
	}

	// offset to have sand boundaries
	box[Min][X]--
	box[Max][X]++
	box[Max][Y]++

	depth := box[1][1]

	part1 := fill(depth+1, depth, box)
	part2 := part1 + 1 + fill(depth+2, 0, box)

	fmt.Println(part1, part2)

	// uncomment for visualization
	// worldmap()
}

func fill(floor int, depth int, box AABB) int {
	// free() slices world where the action is (ie. the rocks)
	// the world physics garantee that the final shape will
	// always be an isosceles right triangle centered at X=200
	//
	// we don't want to simulate aisles as we can compute them
	// easily
	free := func(p XY) bool {
		return world[p[1]][p[0]] == 0 &&
			p[1] < floor &&
			box.contains(p)
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
		return XY{XORG, YORG}
	}

	// dfs iterator
	next := func() XY {
		cur, nxt := XY{-1, -1}, pop() // backtrack
		for !cur.eq(nxt) {
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
	for dst := next(); dst[Y] != depth; dst = next() {
		world[dst[Y]][dst[X]] = '.'
		cnt++
		// on world boundaries, bring back sliced parts
		if dst[X] == box[Min][X] || dst[X] == box[Max][X] {
			cnt += box[Max][Y] - dst[Y]
		}
	}
	return cnt
}

func mkworld(s string) AABB {
	wall := make([]XY, 0, 128)
	box := AABB{{INF, INF}, {0, 0}}
	for _, segs := range strings.Split(s, "->") {
		var seg XY
		for i, s := range strings.Split(segs, ",") {
			seg[i] = atoi(s)
			if i == 0 { // translate x to fit
				seg[i] -= XOFF
			}
		}
		box.add(seg)
		wall = append(wall, seg)
	}
	for i := range wall[:len(wall)-1] {
		a, b := wall[i], wall[i+1]

		δ := a.cmp(b)
		for p := a; !p.eq(b); p = p.add(δ) {
			world[p[1]][p[0]] = '#'
		}
		world[b[1]][b[0]] = '#'
	}
	return box
}

// pretty print worldmap
// !!rise your term resolution!!
func worldmap() {
	var worldmap strings.Builder
	for _, row := range world {
		for _, b := range row {
			switch b {
			case 0:
				worldmap.WriteByte(' ')
			default:
				worldmap.WriteByte(b)
			}
		}
		worldmap.WriteByte('\n')
	}
	fmt.Println(worldmap.String())
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

func (a XY) eq(b XY) bool {
	return a[X] == b[X] && a[Y] == b[Y]
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
// s is ^\s+\d+.*
// front spaces are trimmed
func atoi(s string) int {
	var n int
	for _, c := range s {
		if c == ' ' {
			continue
		}
		n = 10*n + int(c-'0')
	}
	return n
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
