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
	nums := NewCounter()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for _, n := range strings.Fields(input.Text()) {
			nums.Add(atoi(n), 1)
		}
	}

	for i := 0; i < 25; i++ {
		nums = nums.MemBlink()
	}
	sum1 := nums.Sum()

	for i := 0; i < 50; i++ {
		nums = nums.MemBlink()
	}
	sum2 := nums.Sum()

	fmt.Println(sum1, sum2)

}

type Counter struct {
	data map[int]int
}

func NewCounter() *Counter {
	return &Counter{
		data: make(map[int]int, MAXN),
	}
}

func blink(stone int) []int {
	word := strconv.Itoa(stone)
	switch {
	case stone == 0:
		return []int{1}
	case len(word) > 1 && len(word)%2 == 0:
		m := len(word) / 2
		return []int{atoi(word[:m]), atoi(word[m:])}
	default:
		return []int{2024 * stone}
	}
}

func (c *Counter) MemBlink() *Counter {
	next := NewCounter()
	for n, count := range c.data {
		for _, r := range blink(n) {
			next.Add(r, count)
		}
	}
	return next
}

// Add increments the value for the given key.
func (c *Counter) Add(key int, value int) {
	c.data[key] = c.data[key] + value
}

func (c *Counter) Sum() int {
	sum := 0
	for _, n := range c.data {
		sum += n
	}
	return sum
}

func strip(s string) string {
	stripped := strings.TrimLeft(s, "0")
	if stripped == "" {
		return "0"
	}
	return stripped
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
