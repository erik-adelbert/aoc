package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type pos []int

func (p pos) mean() int {
	sum := 0
	for _, x := range p {
		sum += x
	}
	return int(math.Floor(float64(sum) / float64(len(p)))) // Round doesn't work on my input
}

func (p pos) median() int {
	sort.Ints(p)
	return p[len(p)/2]
}

// sumdist does parallel summation of distances for part1 & part2.
func (p pos) sumdist() (int, int) {
	m1, m2 := p.median(), p.mean()
	g := func(x int) int { return (x * (x + 1)) / 2 } // gauss sum

	sum1, sum2 := 0, 0
	for _, x := range p {
		sum1 += abs(x - m1)    // dist1
		sum2 += g(abs(x - m2)) // dist2
	}
	return sum1, sum2
}

func main() {
	crabs := make(pos, 0, 1000)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for _, arg := range strings.Split(input.Text(), ",") {
			x, _ := strconv.Atoi(arg)
			crabs = append(crabs, x)
		}
	}
	fmt.Println(crabs.sumdist()) // part1 & part2
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
