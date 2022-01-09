package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MaxInt and MinInt are defined in the idiomatic way
const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -MaxInt - 1
)

type vec struct {
	x, y int
}

func main() {
	dots := make([]vec, 0, 1024)

	vfold := func(x int) { // fold along vertical axis
		for i, d := range dots {
			if d.x > x {
				d.x = 2*x - d.x
			}
			dots[i] = vec{d.x, d.y}
		}
	}

	hfold := func(y int) { // fold along horizontal axis
		for i, d := range dots {
			if d.y > y {
				d.y = 2*y - d.y
			}
			dots[i] = vec{d.x, d.y}
		}
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if args := strings.Split(line, ","); len(args) > 1 { // dots
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			dots = append(dots, vec{x, y})
		}
		if args := strings.Split(line, "="); len(args) > 1 { // folding
			n, _ := strconv.Atoi(args[1])
			i := len(args[0]) - 1
			if args[0][i] == 'x' { // last car of arg0
				vfold(n)
			} else {
				hfold(n)
			}
			break // on first fold
		}
	}

	frame := make(map[vec]int)
	for _, d := range dots {
		frame[d]++
	}

	fmt.Println(len(frame))
}
