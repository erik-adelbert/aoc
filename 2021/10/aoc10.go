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

	closing := map[byte]byte{
		'(': ')', '[': ']', '{': '}', '<': '>',
	}
	scale1 := map[byte]int{
		')': 3, ']': 57, '}': 1197, '>': 25137,
	}
	scale2 := map[byte]int64{
		')': 1, ']': 2, '}': 3, '>': 4,
	}

	scores := make([]int64, 0, 128)
	sum, input := 0, bufio.NewScanner(os.Stdin)
SCAN:
	for input.Scan() {
		stack = stack[:0] // reset
		for _, b := range input.Bytes() {
			switch b {
			case '(', '[', '{', '<':
				push(closing[b])
			case ')', ']', '}', '>':
				if a := pop(); a != b { // discard corrupted
					sum += scale1[b]
					continue SCAN
				}
			}
		}

		var n int64
		for !empty() {
			n = 5*n + scale2[pop()]
		}
		if n > 0 {
			scores = append(scores, n)
		}
	}

	fmt.Println(sum) // part1

	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	fmt.Println(scores[len(scores)/2]) // median, part2
}
