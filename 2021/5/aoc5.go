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
	w = 1000
	h = 1000
)

type canvas [w * h]int

var can1, can2 canvas

// plot draws a line on the canvas according to part 1&2 constraints
func plot(a, b point) {
	Δx, Δy := b.x-a.x, b.y-a.y

	switch {
	case Δx == 0:
		x := a.x
		ymin, ymax := minax(a.y, b.y)
		for y := ymin; y <= ymax; y++ {
			can1[x*w+y]++
			can2[x*w+y]++
		}
	case Δy == 0:
		y := a.y
		xmin, xmax := minax(a.x, b.x)
		for x := xmin; x <= xmax; x++ {
			can1[x*w+y]++
			can2[x*w+y]++
		}
	default:
		m := Δy / Δx
		c := a.y - a.x*m
		xmin, xmax := minax(a.x, b.x)
		for x := xmin; x <= xmax; x++ {
			y := m*x + c
			can2[x*w+y]++
		}
	}
}

func popcounts() (int, int) {
	pop1, pop2 := 0, 0
	for i := 0; i < len(can1); i++ {
		if can1[i] > 1 {
			pop1++
		}
		if can2[i] > 1 {
			pop2++
		}
	}
	return pop1, pop2
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := strings.Replace(input.Text(), "->", ",", 1)
		args := strings.Split(line, ",")

		a, b := Point(args[0], args[1]), Point(args[2], args[3])
		plot(a, b)
	}

	fmt.Println(popcounts()) // part1 & part2
}

func minax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}
