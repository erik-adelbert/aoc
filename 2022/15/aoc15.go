// aoc15.go --
// advent of code 2022 day 15
//
// https://adventofcode.com/2022/day/15
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-15: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var sensors []sensor = make([]sensor, 0, 64)

// YMAX is the world max depth
var YMAX int = 4_000_000

func main() {
	// part2
	A := make(set, 64)
	B := make(set, 64)
	C := make(set, 64)
	D := make(set, 64)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), "=")
		SB := []XY{
			{atoi(args[1]), atoi(args[2])}, // sensor
			{atoi(args[3]), atoi(args[4])}, // beacon
		}
		sensors = append(sensors, mksensor(SB))
	}

	ranges := make([]XY, 0, 64)
	for _, s := range sensors {
		// part1
		if δ := s.R - abs(YMAX/2-s.O[Y]); δ >= 0 {
			ranges = append(ranges, XY{s.O[X] - δ, s.O[X] + δ})
		}

		// part2, in between the sensors
		A.add(s.O[X] - s.O[Y] + s.R + 1)
		B.add(s.O[X] - s.O[Y] - s.R - 1)
		C.add(s.O[X] + s.O[Y] - s.R - 1)
		D.add(s.O[X] + s.O[Y] + s.R + 1)
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][X] < ranges[j][X]
	})

	ngaps, lims := 0, ranges[0]
	for _, R := range ranges {
		ngaps, lims[1] = max(0, R[0]-lims[1]-1), max(lims[1], R[1])
	}

	// part1
	fmt.Println(lims[1] - lims[0] - ngaps)

	A.inter(B)
	C.inter(D)

	a := A.pop()
	c := C.pop()

	// part2
	fmt.Println((a+c)*YMAX/2 + (c-a)/2)
}

type set map[int]struct{}

func (A set) add(i int) {
	A[i] = struct{}{}
}

func (A set) pop() int {
	for i := range A {
		return i
	}
	panic("unreachable")
}

func (A set) inter(B set) {
	for i := range A {
		if _, ok := B[i]; !ok {
			delete(A, i)
		}
	}
}

type sensor CIRC

func mksensor(SB []XY) sensor {
	S, B := SB[0], SB[1] // sensor, beacon
	R := S.manh(B)       // radius

	return sensor{S, R}
}

// CIRC is a manhattan circle
type CIRC struct {
	O XY
	R int
}

// indices for XY
const (
	X = iota
	Y
)

// XY is a 2D point
type XY [2]int

func (a XY) manh(b XY) int {
	return abs(a[X]-b[X]) + abs(a[Y]-b[Y])
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// strconv.Atoi modified loop
// s is ^-?\d+.*$
// the suffix part .* is ditched
func atoi(s string) (n int) {
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}

	for i := 0; i < len(s) && s[i]-'0' < 10; i++ {
		n = 10*n + int(s[i]-'0')
	}
	return n * neg
}
