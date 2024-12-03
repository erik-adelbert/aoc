package main

import (
	"fmt"
	"os"
	"regexp"
)

const (
	IN = iota
	OUT
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`mul\((\d+),(\d+)\)|do\(\)|don\'t\(\)`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	sum1, sum2 := 0, 0
	state := IN
	for _, match := range matches {
		switch match[0] {
		case "do()":
			state = IN
		case "don't()":
			state = OUT
		default:
			n := atoi(match[1]) * atoi(match[2])
			sum1 += n
			if state == IN {
				sum2 += n
			}
		}
	}
	fmt.Println(sum1, sum2)
}

func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
