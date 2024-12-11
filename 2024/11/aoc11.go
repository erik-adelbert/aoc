// aoc11.go --
// advent of code 2024 day 11
//
// https://adventofcode.com/2024/day/11
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-11: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const MAXN = 3799 // arbitrary

func main() {
	stones := NewCounter()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for _, n := range strings.Fields(input.Text()) {
			stones[atoi(n)] = 1
		}
	}

	blink := func(n int) {
		for i := 0; i < n; i++ {
			stones = stones.Blink()
		}
	}

	blink(25)
	count1 := stones.Popcnt()

	blink(50)
	count2 := stones.Popcnt()

	fmt.Println(count1, count2)
}

type Counter map[int]int

func NewCounter() Counter {
	return make(map[int]int, MAXN)
}

func (c Counter) Blink() Counter {
	next := NewCounter()
	for n, count := range c {
		for _, m := range blink(n) {
			next[m] += count
		}
	}
	return next
}

func (c Counter) Popcnt() int {
	pop := 0
	for _, n := range c {
		pop += n
	}
	return pop
}

func blink(n int) []int {
	ndigit := log10(n)

	switch {
	case n == 0:
		return []int{1}
	case ndigit%2 == 0:
		// split stone in half
		d := pow10(ndigit / 2)
		return []int{n / d, n % d}
	default:
		return []int{2024 * n}
	}
}

func log10(n int) int {
	i := 0
	for n > 0 {
		n /= 10
		i++
	}
	return i
}

func pow10(n int) int {
	return []int{
		1, 10, 100, 1000, 10000, 100000, 1000000,
	}[n]
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
