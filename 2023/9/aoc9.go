package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const NELEM = 21

func main() {
	sumL, sumR := 0, 0
	history := make([]int, NELEM)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fields := fields(input.Text())
		for i := range fields {
			history[i] = atoi(fields[i])
		}
		history = history[:len(fields)]
		L, R := next(history)

		sumL += L // part2
		sumR += R // part1
	}

	fmt.Println(sumR, sumL)
}

func isZero(a []int) bool {
	sum := 0
	for i := range a {
		sum |= a[i]
	}
	return sum == 0
}

func next(a []int) (int, int) {
	difs := make([]int, len(a)-1)

	for i := range a[:len(difs)] {
		difs[i] = a[i+1] - a[i]
	}

	if isZero(difs) {
		return a[0], a[0]
	}

	L, R := next(difs)
	return a[0] - L, a[len(a)-1] + R
}

var fields = strings.Fields

// strconv.Atoi simplified core loop
// s is ^-?\d+$
func atoi(s string) (n int) {
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}

	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return neg * n
}
