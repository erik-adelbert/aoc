package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Infinity is equal to MaxInt
const Infinity = int(^uint(0) >> 1)

type point struct {
	x, y int
}

// Point constructs a point object
func Point(a, b string) point {
	x, _ := strconv.Atoi(strings.TrimSpace(a))
	y, _ := strconv.Atoi(strings.TrimSpace(b))

	return point{x, y}
}

const (
	w = 1024
	h = 1024
)

var field [w * h]int

func draw(a, b point) {
	Δx, Δy := b.x-a.x, b.y-a.y

	switch {
	case Δx == 0:
		x := a.x
		ymin, ymax := min(a.y, b.y), max(a.y, b.y)
		for y := ymin; y <= ymax; y++ {
			field[x*w+y]++
		}
	case Δy == 0:
		y := a.y
		xmin, xmax := min(a.x, b.x), max(a.x, b.x)
		for x := xmin; x <= xmax; x++ {
			field[x*w+y]++
		}
	}
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := strings.Replace(input.Text(), "->", ",", 1)
		args := strings.Split(line, ",")

		a, b := Point(args[0], args[1]), Point(args[2], args[3])
		draw(a, b)
	}

	n := 0
	for _, v := range field {
		if v > 1 {
			n++
		}
	}
	fmt.Println(n)
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
