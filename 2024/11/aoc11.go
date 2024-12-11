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
	"strconv"
	"strings"
)

const MAXN = 3799

func main() {
	stones := NewCounter()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for _, n := range strings.Fields(input.Text()) {
			stones.Add(atoi(n), 1)
		}
	}

	for i := 0; i < 25; i++ {
		stones = stones.MemBlink()
	}
	sum1 := stones.Sum()

	for i := 0; i < 50; i++ {
		stones = stones.MemBlink()
	}
	sum2 := stones.Sum()

	fmt.Println(sum1, sum2)
}

type Counter struct {
	data map[int]int
}

func NewCounter() Counter {
	return Counter{
		data: make(map[int]int, MAXN),
	}
}

func blink(n int) []int {
	stone := strconv.Itoa(n)
	switch {
	case n == 0:
		return []int{1}
	case len(stone) > 1 && len(stone)%2 == 0:
		// split stone in half
		m := len(stone) / 2
		return []int{atoi(stone[:m]), atoi(stone[m:])}
	default:
		return []int{2024 * n}
	}
}

func (c Counter) MemBlink() Counter {
	next := NewCounter()
	for n, count := range c.data {
		for _, m := range blink(n) {
			next.Add(m, count)
		}
	}
	return next
}

// Add increments the value for the given key.
func (c Counter) Add(stone int, count int) {
	c.data[stone] = c.data[stone] + count
}

func (c Counter) Sum() int {
	sum := 0
	for _, n := range c.data {
		sum += n
	}
	return sum
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
