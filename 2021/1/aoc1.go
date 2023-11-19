// aoc1.go --
// advent of code 2021 day 1
//
// https://adventofcode.com/2021/day/1
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// MaxInt is defined in the idiomatic way
const MaxInt = int(^uint(0) >> 1)

func main() {
	old1, old2, old3 := MaxInt, MaxInt, MaxInt // 3 last depths window

	n1, n2 := 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cur, _ := strconv.Atoi(input.Text())
		if old1 < cur { // increase!
			n1++
		}
		if old3 < cur { // increase!
			n2++
		}
		old1, old2, old3 = cur, old1, old2 // shift/update window
	}
	fmt.Println(n1) // part1
	fmt.Println(n2) // part2
}

func atoi(s string) (n int) {
	n, _ = strconv.Atoi(s)
	return
}
