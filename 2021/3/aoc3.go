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

const (
	O2  = true
	CO2 = !O2
)

func rate(nums []string, gas bool) (int64, error) {
	n := append(nums[:0:0], nums...) // clone

	for i := 0; i < width && len(n) > 1; i++ {
		popcnts := popcounts(n)

		bits := map[bool]string{ // most popular bit filters ordered by gas
			O2:  "01",
			CO2: "10",
		}

		j := 0
		for _, s := range n {
			switch {
			case s[i] == bits[gas][0] && len(n) > 2*popcnts[i]:
				n[j], j = s, j+1
			case s[i] == bits[gas][1] && len(n) <= 2*popcnts[i]:
				n[j], j = s, j+1
			}
		}
		n = n[:j]
	}
	return strconv.ParseInt(n[0], 2, 32)
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
