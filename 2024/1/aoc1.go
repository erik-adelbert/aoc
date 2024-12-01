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

	popcnt := mkpopcnt()

	sum, sim := 0, 0
	for i := range left {
		sum += abs(left[i] - right[i])          // part 1
		sim += left[i] * popcnt(right, left[i]) // part 2

	}
	fmt.Println(sum, sim) // part 1 & 2
}

// mkpopcnt returns a closure that counts the number of occurrences of a value in a sorted slice
// The closure maintains a base index to avoid counting the same value multiple times
func mkpopcnt() func([]int, int) int {
	base := 0

	popcnt := func(slice []int, value int) int {
		slice = slice[base:]
		// Find the first occurrence of value using binary search
		start := search(slice, value)

		// If the value isn't in the slice, return 0
		if start == len(slice) || slice[start] != value {
			return 0
		}

		// Count occurrences of the value
		count := 0
		for i := start; i < len(slice) && slice[i] == value; i++ {
			count++
		}

		base = start + count

		return count
	}

	return popcnt
}

func search(slice []int, value int) int {
	return sort.Search(len(slice), func(i int) bool {
		return slice[i] >= value
	})
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
