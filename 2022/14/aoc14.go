package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	mark any
	XY   [2]int
)

func (a XY) add(b XY) XY {
	a[0] += b[0]
	a[1] += b[1]
	return a
}

func (a XY) cmp(b XY) XY {
	return XY{cmp(b[0], a[0]), cmp(b[1], a[1])}
}

func (a XY) eq(b XY) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func (a XY) idx() int {
	return a[1]*512 + (a[0] - 300)
}

var (
	walls  [256][512]byte
	grains [256][512]byte
)

func mkwalls(s string) int {
	ymax := -1
	shape := make([]XY, 0, 256)
	for _, segs := range strings.Split(s, "->") {
		var seg XY
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
			walls[p[1]][p[0]-300] = 1
		}
		walls[b[1]][b[0]-300] = 1
	}
	return ymax
}

func drop(floor int) XY {
	free := func(p XY) bool {
		r := walls[p[1]][p[0]]
		g := grains[p[1]][p[0]]
		return r == 0 && g == 0 && p[1] < floor
	}

	cur, nxt := XY{-1, -1}, XY{200, 0}
	for !cur.eq(nxt) {
		cur = nxt
		for _, δ := range []XY{{0, 1}, {-1, 1}, {1, 1}} {
			if free(cur.add(δ)) {
				nxt = cur.add(δ)
				break
			}
		}
	}
	return cur
}

func fill(floor int, ymax int) int {
	pop := 0
	dst := drop(floor)
	for dst[1] != ymax {
		grains[dst[1]][dst[0]] = 1
		pop++
		dst = drop(floor)
	}
	return pop
}

func main() {
	depth := 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		y := mkwalls(input.Text())
		if depth < y {
			depth = y
		}
	}

	part1 := fill(depth+1, depth)
	// no reset of part1 grains
	part2 := part1 + 1 + fill(depth+2, 0)

	fmt.Println(part1, part2)
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
