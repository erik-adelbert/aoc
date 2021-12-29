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
	xmin, xmax := MaxInt, MinInt // bounding box
	ymin, ymax := MaxInt, MinInt
	frame := make(map[coo]int)
	for _, d := range dots {
		frame[d]++ // could be intensity/color here
		if d.x < xmin {
			xmin = d.x
		}
		if d.x > xmax {
			xmax = d.x
		}
		if d.y < ymin {
			ymin = d.y
		}
		if d.y > ymax {
			ymax = d.y
		}
	}

	fb := make([][]byte, (ymax-ymin)+1) // frame buffer
	for y := range fb {
		fb[y] = make([]byte, (xmax-xmin)+1)
		for x := range fb[y] {
			fb[y][x] = ' ' // init to black
		}
	}

	for d := range frame { // rasterize
		fb[d.y-ymin][d.x-xmin] = '@' // light up, could use intensity/color here
	}

	for _, r := range fb { // scan display
		fmt.Println(string(r))
	}
}
