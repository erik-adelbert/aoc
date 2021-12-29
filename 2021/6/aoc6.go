package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type fishes [9]uint64

func next(cur fishes) fishes {
	var nxt fishes
	nxt[6] = cur[0]
	for i, v := range cur {
		nxt[(i+8)%9] += v // rotate left
	}
	return nxt
}

func sum(popcnt fishes) uint64 {
	sum := uint64(0)
	for _, v := range popcnt {
		sum += v
	}
	return sum
}

func main() {
	var popcnt fishes

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), ",")
		for _, arg := range args {
			i, _ := strconv.Atoi(arg)
			popcnt[i]++
		}
	}

	for i := 0; i < 256; i++ {
		if i == 80 {
			fmt.Println(sum(popcnt)) // part1
		}
		popcnt = next(popcnt)
	}
	fmt.Println(sum(popcnt)) // part2
}
