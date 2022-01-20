package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type popcnts [9]uint64

// incube computes fishes population.
func incube(a []uint64) {
	i, n := len(a)-1, a[0]
	copy(a, a[1:])
	a[6], a[i] = a[6]+n, n
}

func sum(p popcnts) uint64 {
	var sum uint64
	for _, n := range p {
		sum += n
	}
	return sum
}

func main() {
	var fishes popcnts

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
		incube(fishes[:]) // pass slice
	}
	fmt.Println(sum(fishes)) // part2
}
