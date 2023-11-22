// aoc2.go --
// advent of code 2022 day 2
//
// https://adventofcode.com/2022/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-2: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

// see day2 notes
func main() {
	scale := [][]int{
		{1, 2, 0},
		{0, 1, 2},
		{2, 0, 1},
	}

	scores := 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		r := int(line[0] - 'A') // opponent move
		c := int(line[2] - 'X') // our move or goal

		// pack parts 1&2 scores
		scores += pack(1+c+3*scale[r][c], 1+scale[2-c][r]+3*c) // apply symmetry
	}

	fmt.Println(unpack(scores))

}

const WIDTH = 16

func pack(a, b int) int {
	return a<<WIDTH + b
}

func unpack(n int) (int, int) {
	return n >> WIDTH, n & ((1 << WIDTH) - 1)
}
