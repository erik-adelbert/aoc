package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -MaxInt - 1
)

type coo struct { // dot or axis
	x, y int
}

func main() {
	dots := make([]coo, 0, 1024)

	vfold := func(x int) { // fold along vertical axis
		for i, d := range dots {
			if d.x > x {
				d.x = 2*x - d.x
			}
			dots[i] = coo{d.x, d.y}
		}
	}

	hfold := func(y int) { // fold along horizontal axis
		for i, d := range dots {
			if d.y > y {
				d.y = 2*y - d.y
			}
			dots[i] = coo{d.x, d.y}
		}
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if args := strings.Split(line, ","); len(args) > 1 { // dots (code)
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			dots = append(dots, coo{x, y})
		}
		if args := strings.Split(line, "="); len(args) > 1 { // folding (decode)
			n, _ := strconv.Atoi(args[1])
			if args[0][len(args[0])-1] == 'x' {
				vfold(n)
			} else {
				hfold(n)
			}
		}
	}

	// display
	const (
		Black = ' '
		White = '\uFFFD' // undefined is very bright
	)

	b, frame := BBox(), make(map[coo]int)
	for _, d := range dots {
		frame[d]++  // could be color
		b.update(d) // bounding box
	}

	fb := make([][]rune, b.h()) // frame buffer
	for y := range fb {
		fb[y] = make([]rune, b.w())
		for x := range fb[y] {
			fb[y][x] = Black // init
		}
	}

	for d := range frame { // rasterize
		fb[d.y-b.ymin()][d.x-b.xmin()] = White // light up! could use color LUT
	}

	for _, r := range fb { // scan display
		fmt.Println(string(r))
	}
}

type bbox struct {
	a, b coo
}

func BBox() bbox {
	return bbox{
		coo{MaxInt, MaxInt},
		coo{MinInt, MinInt},
	}
}

func (b *bbox) update(c coo) {
	b.a = min(b.a, c)
	b.b = max(b.b, c)
	return
}

func (b bbox) h() int {
	return b.b.y - b.a.y + 1
}

func (b bbox) w() int {
	return b.b.x - b.a.x + 1
}

func (b bbox) xmin() int {
	return b.a.x
}

func (b bbox) ymin() int {
	return b.a.y
}

func min(a, b coo) coo {
	λ := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	return coo{λ(a.x, b.x), λ(a.y, b.y)}
}

func max(a, b coo) coo {
	λ := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	return coo{λ(a.x, b.x), λ(a.y, b.y)}
}
