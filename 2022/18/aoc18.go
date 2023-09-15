package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	w := make(world, 4096)
	b := new(AABB)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := strings.Split(input.Text(), ",")
		p := XYZ{
			atoi(line[0]), atoi(line[1]), atoi(line[2]),
		}
		b.add(p)
		w.add(p)
	}

	// part1
	fmt.Println(w.area())

	// part2
	p0 := b[Min].sub(XYZ{-1, -1, -1})
	fmt.Println(w.flood(p0, b))
}

type world map[XYZ]struct{}

func (w world) add(p XYZ) {
	w[p] = struct{}{}
}

func (w world) area() int {
	area := 0
	for x := range w {
		for _, δ := range Δ {
			if _, ok := w[x.add(δ)]; !ok {
				area++
			}
		}
	}
	return area
}

func (w world) flood(p XYZ, b *AABB) int {
	queue := make([]XYZ, 0, 2048)
	push := func(p XYZ) {
		queue = append(queue, p)
	}
	popf := func() XYZ {
		var p XYZ
		p, queue[0] = queue[0], XYZ{}
		queue = queue[1:]
		return p
	}
	empty := func() bool {
		return len(queue) == 0
	}

	flood := 0
	seen := make(world, 2048)
	push(p)
	for !empty() {
		x := popf()
		for _, δ := range Δ {
			δx := x.add(δ)
			if _, ok := w[δx]; ok {
				flood++
				continue
			}

			if _, ok := seen[δx]; !ok && b.contains(δx) {
				push(δx)
				seen[δx] = struct{}{}
			}
		}
	}

	return flood
}

// Δ is neighbor offsets
var Δ = []XYZ{
	{+0, +1, +0}, // up
	{+1, +0, +0}, // right
	{+0, -1, +0}, // down
	{-1, +0, +0}, // left
	{+0, +0, +1}, // fwd
	{+0, +0, -1}, // back
}

// Min, Max indices for AABB
const (
	Min = iota
	Max
)

// AABB is axis aligned bounding box
type AABB [2]XYZ

// resize AABB to contain p
func (b *AABB) add(p XYZ) {
	b[Min][X] = min(b[Min][X], p[X])
	b[Min][Y] = min(b[Min][Y], p[Y])
	b[Min][Z] = min(b[Min][Z], p[Z])

	b[Max][X] = max(b[Max][X], p[X])
	b[Max][Y] = max(b[Max][Y], p[Y])
	b[Max][Z] = max(b[Max][Z], p[Z])
}

func (b *AABB) contains(p XYZ) bool {
	m := b[Min].sub(XYZ{1, 1, 1})
	M := b[Max].add(XYZ{1, 1, 1})
	return m.lte(p) && M.gte(p)
}

// axis indices
const (
	X = iota
	Y
	Z
)

// XYZ is 3D point
type XYZ [3]int

// add points
func (a XYZ) add(b XYZ) XYZ {
	return XYZ{
		a[X] + b[X], a[Y] + b[Y], a[Z] + b[Z],
	}
}

// subtract points
func (a XYZ) sub(b XYZ) XYZ {
	return XYZ{
		a[X] - b[X], a[Y] - b[Y], a[Z] - b[Z],
	}
}

func (a XYZ) lte(b XYZ) bool {
	return a[X] <= b[X] && a[Y] <= b[Y] && a[Z] <= b[Z]
}

func (a XYZ) gte(b XYZ) bool {
	return a[X] >= b[X] && a[Y] >= b[Y] && a[Z] >= b[Z]
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return n
}

// minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
