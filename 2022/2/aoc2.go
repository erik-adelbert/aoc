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

// see day2 notes
func main() {
	scale := [][]int{
		{1, 2, 0, 0},
		{0, 1, 2, 0},
		{2, 0, 1, 0},
		{0, 0, 0, 0},
	}

	scores := make([]int, 4)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Bytes()

		r := int(line[0] - 'A') // opponent move
		c := int(line[2] - 'X') // our move or goal

		scores[Part1] += 1 + c + 3*scale[r][c]
		scores[Part2] += 1 + scale[2-c][r] + 3*c // apply symmetry
	}

	fmt.Println(scores[:2])
}
