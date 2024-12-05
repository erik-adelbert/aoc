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
	"slices"
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
	slices.Sort(left)
	slices.Sort(right)

	sum, sim := 0, 0
	for i := range left {
		sum += abs(left[i] - right[i])          // part 1
		sim += left[i] * popcnt(right, left[i]) // part 2
	}
	fmt.Println(sum, sim) // part 1 & 2
}

func popcnt(slice []int, n int) (count int) {
	// find the first occurrence of n using binary search
	start := sort.Search(len(slice), func(i int) bool {
		return slice[i] >= n
	})

	// count the number of occurrences of n
	for _, x := range slice[start:] {
		if x != n {
			return
		}
		count++
	}

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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
