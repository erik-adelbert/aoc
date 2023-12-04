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

	deck := make([]int, MAXMATCH) // ring buffer, deck[i%MAXMATCH] is the count of card #(i+1)
	x := func(j int) int {        // ring deck index
		return j % MAXMATCH
	}

	input := bufio.NewScanner(os.Stdin)
	for i := 0; input.Scan(); i++ {

		input := input.Text()
		// input is: ^Game\s(\s|\d)\d:\s(\d+\s)+|\s(\d+\s)+$
		// ditch '^Game \d+:\s' prefix, split winning and cards numbers
		raw := Split(input[Index(input, ":")+1:], " | ")
		w, card := Fields(raw[0]), Fields(raw[1])

		// map winning numbers into a set
		wins := zero128 // fast adhoc set
		for i := range w {
			wins.setbit(atoi(w[i]))
		}

		// match card numbers against winning ones
		nmatch := 0
		for i := range card {
			if wins.getbit(atoi(card[i])) > 0 {
				nmatch++
			}
		}

		// compute part1
		// 2^(nmatch-1) | 0 if nmatch == 0
		score += 1 << nmatch >> 1

		// update deck and fwd duplicate cards
		deck[x(i)] += 1
		for ii := i + 1; ii < (i+1)+nmatch; ii++ {
			deck[x(ii)] += deck[x(i)]
		}

		// compute part2
		ncard += deck[x(i)]
		deck[x(i)] = 0 // consume deck
	}
	fmt.Println(score, ncard) // parts 1 & 2
}

// package strings wrappers/sugars
var Fields, Index, Split = strings.Fields, strings.Index, strings.Split

const uint128size = 128

var (
	zero128 uint128
	one28   = uint128{1, 0}
)

type uint128 struct {
	w0, w1 uint64
}

// setbit sets bit n-th n = 0 is LSB.
// n must be <= 127.
func (u *uint128) setbit(n int) {
	switch n >> 6 {
	case 1:
		u.w1 |= (1 << (n & 0x3f))
	case 0:
		u.w0 |= (1 << (n & 0x3f))
	}
}

func (u *uint128) getbit(n int) uint64 {
	x := u.rsh(n)
	return x.w0 & 1
}

func (u uint128) rsh(n int) uint128 {
	var a uint64

	switch {
	case n > 128:
		return uint128{}
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
