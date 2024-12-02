// aoc1.go --
// advent of code 2024 day 1
//
// https://adventofcode.com/2024/day/1
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	input := bufio.NewScanner(os.Stdin)

	var left, right []int
	for input.Scan() {
		words := strings.Fields(input.Text())
		left = append(left, atoi(words[0]))
		right = append(right, atoi(words[1]))
	}

	// presort
	sort.Ints(left)
	sort.Ints(right)

	sum, sim := 0, 0
	for i := range left {
		sum += abs(left[i] - right[i])          // part 1
		sim += left[i] * popcnt(right, left[i]) // part 2

	}
	fmt.Println(sum, sim) // part 1 & 2
}

func popcnt(slice []int, value int) int {
	// Find the first occurrence of value using binary search
	start := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= value
	})

	// If the value isn't in the slice, return 0
	if start == len(slice) || slice[start] != value {
		return 0
	}

	// Count occurrences of the value
	count := 0
	for i := start; i < len(slice) && slice[i] == value; i++ {
		count++
	}

	return count
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
