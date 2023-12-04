// aoc4.go --
// advent of code 2023 day 4
//
// https://adventofcode.com/2023/day/4
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2023-12-4: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	deck := make([]int, 256) // deck[i] is the count of card #i+1
	score, ncard := 0, 0     // part 1 & 2 results

	input := bufio.NewScanner(os.Stdin)
	for i := 0; input.Scan(); i++ {
		input := input.Text()
		// input is: ^Game\s(\s|\d)\d:\s(\d+\s)+|\s(\d+\s)+$
		// ditch '^Game \d+:\s' prefix, split winning and cards numbers
		raw := Split(input[Index(input, ":")+1:], " | ")
		w, card := Fields(raw[0]), Fields(raw[1])

		// map winning numbers into a set
		wins := make([]bool, 100) // fast adhoc set
		for i := range w {
			wins[atoi(w[i])] = true
		}

		// match card numbers against winning ones
		nmatch := 0
		for i := range card {
			if wins[atoi(card[i])] {
				nmatch++
			}
		}

		// compute part1
		// 2^(nmatch-1) | 0 if nmatch == 0
		score += 1 << nmatch >> 1

		// update deck and fwd duplicate cards
		deck[i] += 1
		for ii := i + 1; ii < (i+1)+nmatch; ii++ {
			deck[ii] += deck[i]
		}

		// compute part2
		ncard += deck[i]
	}
	fmt.Println(score, ncard) // parts 1 & 2
}

// package strings wrappers/sugars
var Fields, Index, Split = strings.Fields, strings.Index, strings.Split

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
