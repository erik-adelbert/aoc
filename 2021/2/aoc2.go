// aoc2.go --
// advent of code 2021 day 2
//
// https://adventofcode.com/2021/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-2: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	x, y, aim := 0, 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		arg, _ := strconv.Atoi(strings.Fields(line)[1])
		switch line[0] {
		case 'f': // forward
			x += arg
			y += aim * arg
		case 'u': // up
			aim -= arg
		case 'd': // down
			aim += arg
		}
	}
	fmt.Println(x * aim) // part1
	fmt.Println(x * y)   // part2
}
