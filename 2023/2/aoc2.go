// aoc2.go --
// advent of code 2023 day 2
//
// https://adventofcode.com/2023/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2023-12-2: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	idsum, pwsum := 0, 0

	input := bufio.NewScanner(os.Stdin)
	for i := 1; input.Scan(); i++ {
		valid, power := true, [3]int{}

		// input is ^Game \d+:\s(.*;\s)+(.*)$
		input := input.Text()
		draws := Split(input[Index(input, ": ")+2:], "; ") // ditch "^Game X: " prefix, split tail
		for j := range draws {
			draws := Split(draws[j], ", ") // split game drawings
			for i := range draws {
				datas := Fields(draws[i])           // split RGB component and count
				color := Index("bgr", datas[1][:1]) // single char 'r', 'g' or 'b' -> R, G, B
				count := atoi(datas[0])

				// check for validity, part1
				valid = valid && count <= 14-color

				// record max power RGB, part2
				power[color] = max(power[color], count)
			}
		}
		if valid {
			idsum += i // part1
		}
		pwsum += power[B] * power[G] * power[R] // part2
	}
	fmt.Println(idsum, pwsum)
}

// package strings wrappers/sugars
var Index, Fields, Split = strings.Index, strings.Fields, strings.Split

const (
	B = iota
	G
	R
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

var DEBUG = true

func debug(a ...any) (int, error) {
	if DEBUG {
		return fmt.Println(a...)
	}
	return 0, nil
}
