package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	Part1 = iota
	Part2
)

type (
	dot struct{}
	pos [2]int
)

func (a *pos) add(b pos) {
	a[0] += b[0]
	a[1] += b[1]
}

func (a *pos) sub(b pos) {
	a[0] -= b[0]
	a[1] -= b[1]
}

func (a pos) dir() pos {
	dir := pos{0, 0}

	for i := 0; i < 2; i++ {
		switch {
		case a[i] < 0:
			dir[i] = -1
		case a[i] > 0:
			dir[i] = 1
		}
	}

	return dir
}

func (a pos) len2() int {
	return a[0]*a[0] + a[1]*a[1]
}

func main() {
	visits := [2]map[pos]dot{
		make(map[pos]dot),
		make(map[pos]dot),
	}

	visits[Part1][pos{0, 0}] = dot{}
	visits[Part2][pos{0, 0}] = dot{}

	off := map[byte]pos{
		'U': {+0, +1}, 'L': {-1, +0},
		'D': {+0, -1}, 'R': {+1, +0},
	}

	knots := [16]pos{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// input text: ^([U,D,L,R]) (\d)$
		line := input.Text()
		θ, n := line[0], atoi(line[2:])

		for ; n > 0; n-- {
			// move head
			knots[0].add(off[θ])

			for i, seg := range knots[:len(knots)-1] { // head
				seg.sub(knots[i+1]) // head - tail

				if seg.len2() >= 4 {
					// update tail
					knots[i+1].add(seg.dir())
				}
			}

			visits[Part1][knots[1]] = dot{}
			visits[Part2][knots[9]] = dot{}
		}
	}

	fmt.Println(len(visits[Part1]), len(visits[Part2]))
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int
	for _, c := range []byte(s) {
		n = 10*n + int(c-'0')
	}
	return n
}
