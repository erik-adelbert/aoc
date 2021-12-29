package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	stack := make([]byte, 0, 128)

	push := func(b byte) {
		stack = append(stack, b)
	}

	pop := func() byte {
		if len(stack) == 0 {
			return 0
		}

		b := stack[len(stack)-1]
		stack, stack[len(stack)-1] = stack[:len(stack)-1], 0
		return b
	}

	n, input := 0, bufio.NewScanner(os.Stdin)
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
				scale := map[byte]int{
					')': 3, ']': 57, '}': 1197, '>': 25137,
				}
				if a := pop(); a == 0 || pair[a] != b {
					n += scale[b]
					continue SCAN
				}
			}
		}
	}

	fmt.Println(n)
}
