package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type positions []int

func (p positions) mean() int {
	sum := 0
	for _, x := range p {
		sum += x
	}
	return int(math.Round(float64(sum) / float64(len(p))))
}

func (p positions) sumdist() int {
	g := func(x int) int { // gauss sum
		return (x * (x + 1)) / 2
	}

	m, sum := p.mean(), 0
	for _, x := range p {
		sum += g(abs(x - m))
	}
	return sum
}

func main() {
	var crabs positions

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for _, arg := range strings.Split(input.Text(), ",") {
			x, _ := strconv.Atoi(arg)
			crabs = append(crabs, x)
		}
	}
	fmt.Println(crabs.sumdist())
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
