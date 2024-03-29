// aoc13.go --
// advent of code 2021 day 13
//
// https://adventofcode.com/2021/day/13
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-13: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MaxInt an MinInt are defined in the idiomatic way
const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -MaxInt - 1
)

type vec struct { // dot or axis
	x, y int
}

func main() {
	dots := make([]vec, 0, 1024)

	fold := func(a byte, n int) { // fold along axis (x|y) = n
		for i, d := range dots {
			if d.x > n && a == 'x' {
				d.x = 2*n - d.x
			}
			if d.y > n && a == 'y' {
				d.y = 2*n - d.y
			}
			dots[i] = vec{d.x, d.y}
		}
	}

	first := true
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if args := strings.Split(line, ","); len(args) > 1 { // dots (coded)
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			dots = append(dots, vec{x, y})
		}
		if args := strings.Split(line, "="); len(args) > 1 { // fold (decode)
			a := args[0][len(args[0])-1] // a = (x|y)
			n, _ := strconv.Atoi(args[1])
			fold(a, n)

			if first {
				first = false
				_, frame := frame(dots)
				fmt.Println(len(frame)) // part1
			}
		}
	}

	display(dots) // part2
}

func frame(dots []vec) (bbox, map[vec]int) {
	b, frame := BBox(), make(map[vec]int)
	for _, d := range dots {
		frame[d]++ // could be color
		b.add(d)   // bounding box
	}
	return b, frame
}

func display(dots []vec) {
	const (
		Black = ' '
		White = '\uFFFD' // undefined is very bright
	)

	b, frame := frame(dots)

	fb := make([][]rune, b.h()) // frame buffer
	for y := range fb {
		fb[y] = make([]rune, b.w())
		for x := range fb[y] {
			fb[y][x] = Black // init
		}
	}

	for d := range frame { // rasterise
		d = b.rebase(d)      // this translation convert vectors (dots) to actual pixels
		fb[d.y][d.x] = White // light up! could use color LUT
	}

	for _, r := range fb { // scan display
		fmt.Println(string(r))
	}
}

type bbox struct { // aabb
	a, b vec
}

// BBox constructs an aabb object
func BBox() bbox {
	return bbox{
		vec{MaxInt, MaxInt},
		vec{MinInt, MinInt},
	}
}

func (b *bbox) add(c vec) {
	b.a = min(b.a, c)
	b.b = max(b.b, c)
	return
}

func (b bbox) rebase(c vec) vec {
	c.y -= b.ymin()
	c.x -= b.xmin()
	return c
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

func min(a, b vec) vec {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	return vec{min(a.x, b.x), min(a.y, b.y)}
}

func max(a, b vec) vec {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	return vec{max(a.x, b.x), max(a.y, b.y)}
}
