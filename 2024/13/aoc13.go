// aoc13.go --
// advent of code 2024 day 13
//
// https://adventofcode.com/2024/day/13
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-13: initial commit

// 78751208820885
// 79352015273424

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	systems := make([][]int, 0, 320)

	input := bufio.NewScanner(os.Stdin)
	system := make([]int, 0, 4)
	for input.Scan() {
		line := input.Text()
		switch {
		case line == "":
			systems = append(systems, system)
			system = make([]int, 0, 4)
		case line[0] == 'B':
			args := strings.Split(line[10:], ", ")
			system = append(system, atoi(args[0][2:]), atoi(args[1][2:]))
		case line[0] == 'P':
			args := strings.Split(line[7:], ", ")
			system = append(system, atoi(args[0][2:]), atoi(args[1][2:]))
		}
	}
	systems = append(systems, system) // last system

	sum1, sum2 := 0, 0
	for _, sys := range systems {
		a, b := solve(sys, 0)
		sum1 += a*3 + b

		a, b = solve(sys, 10_000_000_000_000)
		sum2 += a*3 + b
	}

	fmt.Println(sum1, sum2)
}

func solve(system []int, offset int) (int, int) {
	ax, ay, bx, by, X, Y := system[0], system[1], system[2], system[3], system[4], system[5]

	X += offset
	Y += offset

	Δ := ax*by - ay*bx

	if Δ == 0 {
		return 0, 0
	}

	A := (X*by - Y*bx) / Δ
	B := (ax*Y - X*ay) / Δ

	if A < 0 || B < 0 || A*ax != X-B*bx || A*ay != Y-B*by {
		return 0, 0
	}

	return A, B
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
