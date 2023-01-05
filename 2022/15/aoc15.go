package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var sensors []sensor = make([]sensor, 0, 64)

// YMAX is world max depth
var YMAX int = 4_000_000

func main() {
	// part2
	A := make(map[int]any, 64)
	B := make(map[int]any, 64)
	C := make(map[int]any, 64)
	D := make(map[int]any, 64)

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
		A[s.O[X]-s.O[Y]+s.R+1] = any(nil)
		B[s.O[X]-s.O[Y]-s.R-1] = any(nil)
		C[s.O[X]+s.O[Y]+s.R+1] = any(nil)
		D[s.O[X]+s.O[Y]-s.R-1] = any(nil)
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][X] < ranges[j][X]
	})

	lims, ngaps := ranges[0], 0
	for _, R := range ranges {
		ngaps, lims[1] = max(0, R[0]-lims[1]-1), max(lims[1], R[1])
	}

	// part1
	fmt.Println(lims[1] - lims[0] - ngaps)

	inter(A, B)
	inter(C, D)

	a := pop(A)
	b := pop(C)

	// part2
	fmt.Println((a+b)*YMAX/2 + (b-a)/2)
}

func inter(A, B map[int]any) {
	for i := range A {
		if _, ok := B[i]; !ok {
			delete(A, i)
		}
	}
}

func pop(A map[int]any) int {
	for i := range A {
		return i
	}
	panic("unreachable")
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

// strconv.Atoi modified core loop
// s is ^-?\d+.*$
func atoi(s string) int {
	var n int
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}
	for _, c := range s {
		if c < '0' || '9' < c {
			break
		}
		n = 10*n + int(c-'0')
	}
	return n * neg
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
