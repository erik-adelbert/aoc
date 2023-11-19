// aoc3.go --
// advent of code 2021 day 3
//
// https://adventofcode.com/2021/day/3
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-3: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	width   = 12 // bits
	bitmask = (1 << width) - 1
)

func popcounts(nums []string) []int {
	popcnts := make([]int, width)
	for _, n := range nums {
		for i, c := range n {
			if c == '1' {
				popcnts[i]++
			}
		}
	}
	return popcnts
}

type gas int

const (
	o2 gas = iota
	co2
)

func rate(nums []string, g gas) (int64, error) {
	bs := append(nums[:0:0], nums...) // clone binary strings

	bits := [...]string{o2: "01", co2: "10"}[g] // most popular bit filters by gas

	for i := 0; i < width && len(bs) > 1; i++ {
		popcnts := popcounts(bs)

		j := 0
		for _, s := range bs {
			switch {
			case s[i] == bits[0] && len(bs) > 2*popcnts[i]:
				bs[j], j = s, j+1
			case s[i] == bits[1] && len(bs) <= 2*popcnts[i]:
				bs[j], j = s, j+1
			}
		}
		bs = bs[:j]
	}
	return strconv.ParseInt(bs[0], 2, 32)
}

func main() {
	nums := make([]string, 0, 1024)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		nums = append(nums, line)
	}

	rates := make(chan int64)
	defer close(rates)

	go func() {
		n, _ := rate(nums, o2)
		rates <- n
	}()

	go func() {
		n, _ := rate(nums, co2)
		rates <- n
	}()

	var ε int
	for _, popcnt := range popcounts(nums) {
		ε <<= 1
		if 2*popcnt <= len(nums) {
			ε |= 1
		}
	}
	γ := ^ε & bitmask
	fmt.Println(ε * γ) // part1

	n := int64(1)
	for i := 0; i < 2; i++ {
		n *= <-rates
	}
	fmt.Println(n) // part2
}
