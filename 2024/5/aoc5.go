// aoc2.go --
// advent of code 2024 day 5
//
// https://adventofcode.com/2024/day/5
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-5: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	RULE = iota
	UPDATE
)

func main() {
	sum1, sum2 := 0, 0
	rules := [100][]int{}

	state := RULE
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "" {
			state = UPDATE
			continue
		}

		switch state {
		case RULE:
			words := strings.Split(line, "|")
			index, next := atoi(words[0]), atoi(words[1])
			rules[index] = append(rules[index], next)
		case UPDATE:
			words := strings.Split(line, ",")
			indices := make([]int, len(words))
			for i, w := range words {
				indices[i] = atoi(w)
			}

			if safe(indices, rules) {
				sum1 += median(indices)
			} else {
				sum2 += median(sort(indices, rules))
			}
		}
	}

	fmt.Println(sum1, sum2) // part 1 & 2
}

func safe(indices []int, rules [100][]int) bool {
	pre := indices[0]
	for _, cur := range indices[1:] {
		if !slices.Contains(rules[pre], cur) {
			return false
		}
		pre = cur
	}

	return true
}

func sort(indices []int, rules [100][]int) []int {
	return slices.SortedFunc(slices.Values(indices), func(a, b int) int {
		if slices.Contains(rules[b], a) {
			return -1
		}
		return 0
	})
}

func median(indices []int) int {
	return indices[len(indices)/2]
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
