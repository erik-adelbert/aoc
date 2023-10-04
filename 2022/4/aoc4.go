// aoc4.go --
// advent of code 2022 day 4
//
// https://adventofcode.com/2022/day/4
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-4: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// part indices
const (
	Part1 = iota
	Part2
)

func main() {
	counts := [2]int{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// input text: ^(\d+)-(\d+),(\d+)-(\d+)$
		// fields:        0     1     2     3
		// varname:       l1    r1    l2    r2
		fields := bytes.FieldsFunc(
			input.Bytes(),
			func(r rune) bool {
				return r == '-' || r == ','
			},
		)

		l1 := atoi(fields[0])
		r1 := atoi(fields[1])
		l2 := atoi(fields[2])
		r2 := atoi(fields[3])

		// closed segments layout ex.
		//  l1        r1
		//  |----------|
		//          |-------|
		//          l2     r2
		switch {
		case (l1-l2)*(r1-r2) <= 0: // 1D contains
			counts[Part1]++
		case (l1-r2)*(r1-l2) <= 0: // 1D intersect
			counts[Part2]++
		}
	}

	// every contained segment is intersecting as well
	fmt.Println(counts[Part1], counts[Part1]+counts[Part2])
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) int {
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return n
}
