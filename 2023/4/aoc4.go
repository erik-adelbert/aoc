// aoc4.go --
// advent of code 2023 day 4
//
// https://adventofcode.com/2023/day/4
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2023-12-4: initial commit
// 2023-12-4: implement u/masklinn ideas

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	const MAXMATCH = 16
	score, ncard := 0, 0 // part 1 & 2 results

	deck := make([]int, MAXMATCH) // ring buffer, deck[i%MAXMATCH] stores card#(i+1) count
	θ := func(i int) int {        // deck circular index
		return i & (MAXMATCH - 1)
	}

	input := bufio.NewScanner(os.Stdin)
	for i := 0; input.Scan(); i++ {

		input := input.Text()
		// input is: ^Game\s(\s|\d)\d:\s((\s|\d)\d+\s)+|\s((\s|\d)\d\s)+(\s|\d)\d$

		// ditch '^Game \d+:\s' prefix, split winning and cards numbers
		raw := split(input[index(input, ":")+1:], " | ")
		w, card := fields(raw[0]), fields(raw[1])

		// map winning numbers into a set
		wins := nullset // fast adhoc int set
		for i := range w {
			wins.set(atoi(w[i]))
		}

		// match card numbers against winning ones
		nmatch := 0
		for i := range card {
			if wins.get(atoi(card[i])) > 0 {
				nmatch++
			}
		}

		// compute part1
		// 2^(nmatch-1) | 0 if nmatch == 0
		score += 1 << nmatch >> 1

		// update deck and fwd duplicate cards
		deck[θ(i)] += 1
		for ii := i + 1; ii < (i+1)+nmatch; ii++ {
			deck[θ(ii)] += deck[θ(i)]
		}

		// compute part2
		ncard += deck[θ(i)]
		deck[θ(i)] = 0 // consume deck
	}
	fmt.Println(score, ncard) // parts 1 & 2
}

// package strings wrappers/sugars
var fields, index, split = strings.Fields, strings.Index, strings.Split

const uint128size = 128

type uint128 struct {
	w0, w1 uint64
}

var (
	zero128 uint128
	nullset = zero128 // sugar
)

// setbit sets bit n-th n = 0 is LSB.
// n must be <= 127.
func (u *uint128) set(n int) {
	switch n >> 6 {
	case 1:
		u.w1 |= (1 << (n & 0x3f))
	case 0:
		u.w0 |= (1 << (n & 0x3f))
	}
}

func (u *uint128) get(n int) uint64 {
	x := u.rsh(n)
	return x.w0 & 1
}

func (u uint128) rsh(n int) uint128 {
	var a uint64

	switch {
	case n > 128:
		return zero128
	case n > 64:
		u.w1, u.w0 = 0, u.w1
		n -= 64
		goto sh64
	}

	// remaining shifts
	a = u.w1 << (64 - n)
	u.w1 = u.w1 >> n

sh64:
	u.w0 = (u.w0 >> n) | a

	return u
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
