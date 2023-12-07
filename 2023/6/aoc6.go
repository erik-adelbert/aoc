// aoc6.go --
// advent of code 2023 day 6
//
// https://adventofcode.com/2023/day/6
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2023-12-6: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := bufio.NewScanner(os.Stdin)

	// parse multiple (part1) and single (part2) race data
	times, T := parse(input) // first line
	dists, D := parse(input) // second line

	// solve V.x = d with V = (t - x) <=> (x - t) * x - d = 0
	// this quadratic formula leads to
	// Δ = √(t² + 4d) / 2
	// x₀₁ = t ± √(t² + 4d) / 2
	// solution r = | ⌈x₀⌉ - ⌊x₁⌋ | + 1
	solve := func(t, d int) int {
		Δ := isqrt(t*t - 4*d)

		x0 := (t - Δ) / 2
		x1 := t - x0

		if x0*(t-x0) <= d {
			x0++ // ceil
		}

		if x1*(t-x1) <= d {
			x1-- // floor
		}

		return x1 - x0 + 1
	}

	Π := 1
	for i := range times {
		Π *= solve(times[i], dists[i])
	}

	fmt.Println(Π, solve(T, D))
}

func parse(input *bufio.Scanner) ([]int, int) {

	input.Scan() // advance input reading

	// ditch header, split fields
	line := input.Text()
	fields := Fields(line[Index(line, ":")+1:])

	a := make([]int, 0, len(fields)) // part1
	var A strings.Builder            // part2
	for _, s := range fields {
		a = append(a, atoi(s)) // convert/collect for part1
		A.WriteString(s)       // concatenate for part2
	}

	return a, atoi(A.String())
}

// Go package strings wrapper/sugar
var Index, Fields = strings.Index, strings.Fields

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

// isqrt
var tab64 = [64]uint64{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5,
}

func log2(n uint64) uint64 {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return tab64[((n-(n>>1))*0x07EDD5E59A4E28C2)>>58]
}

func isqrt(x int) int {
	x64 := uint64(x)
	var b, r uint64
	for p := uint64(1 << ((uint(log2(x64)) >> 1) << 1)); p != 0; p >>= 2 {
		b = r | p
		r >>= 1
		if x64 >= b {
			x64 -= b
			r |= p
		}
	}
	return int(r)
}

func debug(a ...any) {
	if false {
		fmt.Println(a...)
	}
}
