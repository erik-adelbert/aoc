// aoc13.go --
// advent of code 2024 day 13
//
// https://adventofcode.com/2024/day/13
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-13: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var sum1, sum2 int

	system := make([]int, 0, 4)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		switch {
		case line == "":
			// reset system on new line
			system = system[:0]

		case line[0] == 'B':
			// read any button
			args := strings.Split(line[10:], ", ")
			system = append(system, atoi(args[0][2:]), atoi(args[1][2:]))

		case line[0] == 'P':
			// read target X, Y
			args := strings.Split(line[7:], ", ")
			system = append(system, atoi(args[0][2:]), atoi(args[1][2:]))

			// solve for A, B
			A, B := solve(system, 0)
			sum1 += A*3 + B // part 1

			A, B = solve(system, 10_000_000_000_000)
			sum2 += A*3 + B // part 2
		}
	}

	fmt.Println(sum1, sum2) // part 1 & 2
}

// solve for A, B
func solve(sys []int, off int) (A, B int) {
	var Δ int

	// unpack system
	ax, ay, bx, by, X, Y := sys[0], sys[1], sys[2], sys[3], sys[4], sys[5]

	X += off
	Y += off

	// compute determinant
	Δ = ax*by - ay*bx

	if Δ == 0 {
		// no solution
		return
	}

	A = (X*by - Y*bx) / Δ
	B = (ax*Y - X*ay) / Δ

	// check solution
	if A < 0 || B < 0 || A*ax != X-B*bx || A*ay != Y-B*by {
		return 0, 0 // invalid!
	}

	// all done
	return
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
