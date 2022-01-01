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
		outs := strings.Fields(strings.TrimSpace(args[1]))
		for _, t := range outs {
			counts[len(t)]++
		}
	}

	sum := 0
	for _, n := range []int{1, 4, 7, 8} {
		sum += counts[digits[n]]
	}
	fmt.Println(sum)
}
