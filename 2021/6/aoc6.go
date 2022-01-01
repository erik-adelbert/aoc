package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type popcnt [9]uint64

func incube(cur popcnt) popcnt {
	var nxt popcnt = popcnt{6: cur[0]}
	for i, v := range cur {
		nxt[(i+8)%9] += v // rotate left
	}
	return nxt
}

func sum(p popcnt) uint64 {
	var sum uint64
	for _, n := range p {
		sum += n
	}
	return sum
}

func main() {
	var fishes popcnt

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), ",")
		for _, arg := range args {
			i, _ := strconv.Atoi(arg)
			fishes[i]++
		}
	}

	for i := 0; i < 256; i++ {
		if i == 80 {
			fmt.Println(sum(fishes)) // part1
		}
		fishes = incube(fishes)
	}
	fmt.Println(sum(fishes)) // part2
}
