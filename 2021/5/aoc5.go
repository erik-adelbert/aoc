package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

type field [w * h]int

var part1, part2 field

func draw(a, b point) {
	Δx, Δy := b.x-a.x, b.y-a.y

	switch {
	case Δx == 0:
		x := a.x
		ymin, ymax := min(a.y, b.y), max(a.y, b.y)
		for y := ymin; y <= ymax; y++ {
			part1[x*w+y]++
			part2[x*w+y]++
		}
	case Δy == 0:
		y := a.y
		xmin, xmax := min(a.x, b.x), max(a.x, b.x)
		for x := xmin; x <= xmax; x++ {
			part1[x*w+y]++
			part2[x*w+y]++
		}
	default:
		m := Δy / Δx
		c := a.y - a.x*m
		xmin, xmax := min(a.x, b.x), max(a.x, b.x)
		for x := xmin; x <= xmax; x++ {
			y := m*x + c
			part2[x*w+y]++
		}
	}
}

func popcounts() (int, int) {
	p1, p2 := 0, 0
	for i := 0; i < len(part1); i++ {
		if part1[i] > 1 {
			p1++
		}
		if part2[i] > 1 {
			p2++
		}
	}
	return p1, p2
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := strings.Replace(input.Text(), "->", ",", 1)
		args := strings.Split(line, ",")

		a, b := Point(args[0], args[1]), Point(args[2], args[3])
		draw(a, b)
	}

	fmt.Println(popcounts()) // part1 & part2
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
