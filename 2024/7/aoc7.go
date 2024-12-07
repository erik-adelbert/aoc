// aoc7.go --
// advent of code 2024 day 7
//
// https://adventofcode.com/2024/day/7
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-7: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	CONS = true
)

func main() {
	sum1, sum2 := 0, 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		words := strings.Fields(strings.Replace(line, ":", " ", 1))
		nums := make([]int, len(words))
		for i, word := range words {
			nums[i] = atoi(word)
		}

		switch {
		case check(nums, !CONS):
			sum1 += nums[0]
		case check(nums, CONS):
			sum2 += nums[0]
		}
	}
	fmt.Println(sum1, sum1+sum2) // part 1 & 2
}

func check(nums []int, hascons bool) bool {
	acc, nums := nums[0], nums[1:]

	if !hascons { // incorrect but ok because of main calling order
		slices.Reverse(nums)
	}

	var recheck func(int, []int) bool

	recheck = func(acc int, nums []int) bool {
		switch {
		case acc < 0:
			return false
		case len(nums) == 1:
			return acc == nums[0]
		case hascons && acc%mask(nums[0]) == nums[0] && recheck(acc/mask(nums[0]), nums[1:]):
			return true
		case acc%nums[0] == 0 && recheck(acc/nums[0], nums[1:]):
			return true
		case recheck(acc-nums[0], nums[1:]):
			return true
		}
		return false
	}

	return recheck(acc, nums)
}

func mask(n int) int {
	return pow10(log10(n))
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
	i := 1
	for n > 0 {
		i *= 10
		n--
	}
	return i
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
