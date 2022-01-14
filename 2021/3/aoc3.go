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

func rate(nums []string, o2 bool) (int64, error) {
	for i := 0; i < width && len(nums) > 1; i++ {
		popcnts := popcounts(nums)
		matched := make([]string, 0, len(nums))
		for _, s := range nums {
			if o2 {
				switch {
				case s[i] == '0' && len(nums) > 2*popcnts[i]:
					matched = append(matched, s)
				case s[i] == '1' && len(nums) <= 2*popcnts[i]:
					matched = append(matched, s)
				}
			} else { // co2
				switch {
				case s[i] == '0' && len(nums) <= 2*popcnts[i]:
					matched = append(matched, s)
				case s[i] == '1' && len(nums) > 2*popcnts[i]:
					matched = append(matched, s)
				}
			}
		}
		nums = matched
	}
	return strconv.ParseInt(nums[0], 2, 32)
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

	const (
		O2  = true
		CO2 = false
	)

	go func() {
		o2, _ := rate(nums, O2)
		rates <- o2
	}()

	go func() {
		co2, _ := rate(nums, CO2)
		rates <- co2
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
