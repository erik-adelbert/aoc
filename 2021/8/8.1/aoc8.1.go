package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var digits = []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6} // number of segments for 0..9

func main() {
	counts := make([]int, 8)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), "|")
		tokens := strings.Fields(strings.TrimSpace(args[1]))
		for _, t := range tokens {
			counts[len(t)]++
		}
	}

	n := 0
	for _, d := range []int{1, 4, 7, 8} {
		n += counts[digits[d]]
	}
	fmt.Println(n)
}
