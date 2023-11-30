// aoc1.go --
// advent of code 2022 day 1
//
// https://adventofcode.com/2022/day/1
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var m1, m2, m3, sum int

	max3 := func() {
		switch {
		case sum > m1:
			m1, m2, m3 = sum, m1, m2
		case sum > m2:
			m2, m3 = sum, m2
		case sum > m3:
			m3 = sum
		}
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if line := input.Text(); len(line) > 0 {
			sum += atoi(line)
			continue
		}

		max3()
		sum = 0
	}
	max3()

	fmt.Println(m1, m1+m2+m3) // part 1 & 2
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
