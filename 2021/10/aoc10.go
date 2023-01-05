package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	stack := make([]byte, 0, 128)

	empty := func() bool {
		return len(stack) == 0
	}

	push := func(b byte) {
		stack = append(stack, b)
	}

	pop := func() byte {
		if i := len(stack) - 1; i >= 0 {
			pop := stack[i]
			stack, stack[i] = stack[:i], 0
			return pop
		}
		return 0
	}

	closing := [128]byte{
		'(': ')', '[': ']', '{': '}', '<': '>',
	}
	scales := [...][128]int64{
		1: {')': 3, ']': 57, '}': 1197, '>': 25137}, // for part1
		2: {')': 1, ']': 2, '}': 3, '>': 4},         // for part2
	}

	scores := make([]int64, 0, 128)
	sum, input := int64(0), bufio.NewScanner(os.Stdin)
SCAN:
	for input.Scan() {
		stack = stack[:0] // reset
		for _, b := range input.Bytes() {
			if v := closing[b]; v > 0 {
				push(v)
			} else if a := pop(); a != b { // discard corrupted
				sum += scales[1][b]
				continue SCAN
			}
		}

		var n int64
		for !empty() {
			n = 5*n + scales[2][pop()]
		}
		if n > 0 {
			scores = append(scores, n)
		}
	}

	fmt.Println(sum) // part1

	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	fmt.Println(scores[len(scores)/2]) // median, part2
}
