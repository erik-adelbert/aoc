package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const width = 12 // bits

func bitpops(nums []string) []int {
	bpops := make([]int, width)
	for _, n := range nums {
		for i, c := range n {
			if c == '1' {
				bpops[i]++
			}
		}
	}
	return bpops
}

func rate(numbers []string, o2 bool) (int64, error) {
	for i := 0; i < width && len(numbers) > 1; i++ {
		bpops := bitpops(numbers)
		matched := make([]string, 0, len(numbers))
		for _, s := range numbers {
			if o2 {
				switch {
				case s[i] == '0' && len(numbers) > 2*bpops[i]:
					matched = append(matched, s)
				case s[i] == '1' && len(numbers) <= 2*bpops[i]:
					matched = append(matched, s)
				}
			} else { // co2
				switch {
				case s[i] == '0' && len(numbers) <= 2*bpops[i]:
					matched = append(matched, s)
				case s[i] == '1' && len(numbers) > 2*bpops[i]:
					matched = append(matched, s)
				}
			}
		}
		numbers = matched
	}
	return strconv.ParseInt(numbers[0], 2, 32)
}

func main() {
	nums := make([]string, 0, 1024)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		arg := input.Text()
		nums = append(nums, arg)
	}

	rates := make(chan int64)

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

	n := int64(1)
	for i := 0; i < 2; i++ {
		n *= <-rates
	}
	fmt.Println(n)
}
