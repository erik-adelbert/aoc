package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	width   = 12 // bits
	bitmask = (1 << width) - 1
)

func main() {
	bitpops := make([]int, width)

	n, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for i, c := range input.Text() {
			if c == '1' {
				bitpops[i]++
			}
		}
		n++
	}

	var ε int
	for _, bpop := range bitpops {
		ε <<= 1
		if 2*bpop <= n {
			ε |= 1
		}
	}
	γ := ^ε & bitmask
	fmt.Println(ε * γ)
}
