// aoc9.go --
// advent of code 2022 day 9
//
// https://adventofcode.com/2022/day/9
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-9: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

// part indices
const (
	Part1 = iota
	Part2
)

// XY is 2D point
type XY [2]int

func main() {
	visits := [2]set{}

	visits[Part1].add(XY{0, 0})
	visits[Part2].add(XY{0, 0})

	off := [128]XY{ // 0-hashing map
		'U': {+0, +1}, 'L': {-1, +0},
		'D': {+0, -1}, 'R': {+1, +0},
	}

	knots := [10]XY{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// input text: ^([U,D,L,R]) (\d)$
		line := input.Text()

		θ, n := line[0], atoi(line[2:]) // heading, steps
		for n > 0 {
			// move head
			knots[0].add(off[θ])

			// scan rope vectors/knots
			for i, vec := range knots[:len(knots)-1] {
				// vec is k[i] - k[i+1] ie. head - tail
				vec.sub(knots[i+1])

				// len(vec)^2 >= 4 => abs(len(vec)) >= 2
				if vec.len2() >= 4 {
					// move tail knot
					knots[i+1].add(vec.dir())
				}
			}

			visits[Part1].add(knots[1])
			visits[Part2].add(knots[9])

			n--
		}
	}

	fmt.Println(visits[Part1].len(), visits[Part2].len())
}

func (a *XY) add(b XY) {
	a[0] += b[0]
	a[1] += b[1]
}

func (a *XY) sub(b XY) {
	a[0] -= b[0]
	a[1] -= b[1]
}

func (a XY) dir() XY {
	dir := XY{0, 0}

	for i := 0; i < len(dir); i++ {
		switch {
		case a[i] < 0:
			dir[i] = -1
		case a[i] > 0:
			dir[i] = 1
		}
	}

	return dir
}

// square of length
func (a XY) len2() int {
	return a[0]*a[0] + a[1]*a[1]
}

// faster adhoc set than map[XY]struct{}
const MAXN = 512

type set struct {
	data  [MAXN * MAXN / 2]bool
	count int
}

func (s set) len() int {
	return s.count
}

func (s *set) add(x XY) {
	hash := func() int {
		// offset XY by XY{MAXN/2, MAXN/2} to make it positive
		// and map to a flatten (MAXN/2) x MAXN 2D grid
		// x in [0..MAXN/2] and y in [0..MAXN]
		m := (MAXN >> 1)
		return m*(x[0]+m+1) + x[1]
	}
	if i := hash(); !s.data[i] {
		s.data[i] = true
		s.count++
	}
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
