// aoc2.go --
// advent of code 2024 day 2
//
// https://adventofcode.com/2024/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-2: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var reports [][]int

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		words := strings.Fields(input.Text())
		report := make([]int, len(words))
		for i, word := range words {
			report[i] = atoi(word)
		}
		reports = append(reports, report)
	}

	var count1, count2 int // safe reports with no or a single error
	for _, r := range reports {
		switch {
		case safe(r, 0):
			// safe without any error
			count1++
		case safe(r, 1):
			// safe by removing a misplaced element
			count2++
		}
	}

	fmt.Println(count1, count1+count2) // part 1 & 2
}

func safe(report []int, maxerr int) bool {
	if len(report) < 2 {
		return true // a single element or empty list is trivially "safe"
	}

	// try to remove the first element and check if the report is safe
	if maxerr > 0 && safe(report[1:], maxerr-1) {
		return true
	}

	// determine the trend from the first two elements
	increasing := report[1] > report[0]
	decreasing := report[1] < report[0]

	// find the first misplaced element if any
	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]

		unsafe := abs(diff) < 1 || abs(diff) > 3 || (increasing && diff <= 0) || (decreasing && diff >= 0)
		if unsafe {
			if maxerr == 0 {
				return false
			}

			// attempt to remove either the previous or the current element in turn and check
			// if the resulting report is safe
			left, right := remove(clone(report), i-1), remove(report, i)
			return safe(left, maxerr-1) || safe(right, maxerr-1)
		}

	}

	return true
}

func clone(slice []int) []int {
	return append([]int(nil), slice...)
}

func remove(slice []int, i int) []int {
	return append(slice[:i], slice[i+1:]...)
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
