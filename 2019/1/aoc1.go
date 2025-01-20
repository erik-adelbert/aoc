// aoc1.go --
// advent of code 2019 day 1
//
// https://adventofcode.com/2019/day/1
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Σ1, Σ2 := 0, 0

	f := func(x int) int { return x/3 - 2 }

	// fofo...f(x) reapeatedly apply f to x while x > 0
	fofo := func(x int) (Σ int) {
		for x > 0 {
			Σ, x = Σ+x, f(x)
		}
		return
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		x := atoi(input.Text())
		Σ1 += f(x)
		Σ2 += fofo(f(x))
	}
	fmt.Println(Σ1, Σ2) // part 1 & 2
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
