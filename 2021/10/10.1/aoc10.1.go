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
		i := len(stack) - 1

		pop := stack[i]
		stack, stack[i] = stack[:i], 0
		return pop
	}

	closing := map[byte]byte{
		'(': ')', '[': ']', '{': '}', '<': '>',
	}

	scale := map[byte]int{
		')': 3, ']': 57, '}': 1197, '>': 25137,
	}

	n, input := 0, bufio.NewScanner(os.Stdin)
SCAN:
	for input.Scan() {
		stack = stack[:0] // reset
		for _, b := range input.Bytes() {
			switch b {
			case '(', '[', '{', '<':
				push(closing[b])
			case ')', ']', '}', '>':
				if a := pop(); a != b {
					n += scale[b]
					continue SCAN
				}
			}
		}
	}

	fmt.Println(n)
}
