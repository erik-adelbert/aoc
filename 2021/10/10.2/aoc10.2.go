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
		if empty() {
			return 0
		}

		b := stack[len(stack)-1]
		stack, stack[len(stack)-1] = stack[:len(stack)-1], 0
		return b
	}

	scores := make([]int64, 0, 128)
	input := bufio.NewScanner(os.Stdin)
SCAN:
	for input.Scan() {
		stack = stack[:0] // reset
		for _, b := range input.Bytes() {
			switch b {
			case '(', '[', '{', '<':
				push(b)
			case ')', ']', '}', '>':
				pair := map[byte]byte{
					'(': ')', '[': ']', '{': '}', '<': '>',
				}
				if a := pop(); a == 0 || pair[a] != b {
					continue SCAN
				}
			}
		}

		scale := map[byte]int64{
			'(': 1, '[': 2, '{': 3, '<': 4,
		}

		var n int64
		for n = 0; !empty(); {
			n = 5*n + scale[pop()]
		}
		if n > 0 {
			scores = append(scores, n)
		}
	}

	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	fmt.Println(scores[len(scores)/2])
}
