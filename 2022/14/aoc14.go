package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	mark  any
	point [2]int
)

func (a point) add(b point) point {
	a[0] += b[0]
	a[1] += b[1]
	return a
}

func (a point) cmp(b point) point {
	return point{cmp(b[0], a[0]), cmp(b[1], a[1])}
}

func (a point) eq(b point) bool {
	return a[0] == b[0] && a[1] == b[1]
}

var lines, grains map[[2]int]any

func mklines(s string) int {
	ymax := -1
	shape := make([]point, 0, 256)
	for _, segs := range strings.Split(s, "->") {
		var seg point
		for i, s := range strings.Split(segs, ",") {
			seg[i] = atoi(s)
		}
		shape = append(shape, seg)
	}
	for i := range shape[:len(shape)-1] {
		a, b := shape[i], shape[i+1]
		δ := a.cmp(b)
		for p := a; !p.eq(b); p = p.add(δ) {
			if ymax < p[1] {
				ymax = p[1]
			}
			lines[p] = mark(nil)
		}
		lines[b] = mark(nil)
	}
	return ymax
}

func drop(floor int) point {
	free := func(p point) bool {
		_, r := lines[p]
		_, g := grains[p]
		return !r && !g && p[1] < floor
	}

	cur, nxt := point{-1, -1}, point{500, 0}
	for !cur.eq(nxt) {
		cur = nxt
		for _, δ := range []point{{0, 1}, {-1, 1}, {1, 1}} {
			if free(cur.add(δ)) {
				nxt = cur.add(δ)
				break
			}
		}
	}
	return cur
}

func fill(floor int, ymax int) int {
	dst := drop(floor)
	for dst[1] != ymax {
		grains[dst] = mark(nil)
		dst = drop(floor)
	}
	return len(grains)
}

func main() {
	depth := 0
	lines = make(map[[2]int]any)
	grains = make(map[[2]int]any)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		y := mklines(input.Text())
		if depth < y {
			depth = y
		}
	}

	// fmt.Println(len(lines), len(grains))
	fmt.Println(fill(depth+1, depth))
	fmt.Println(1 + fill(depth+2, 0))
}

// strconv.Atoi modified core loop
// s is ^\s+\d+.*
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

