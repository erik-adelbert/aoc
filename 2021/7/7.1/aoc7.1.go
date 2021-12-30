package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type pos []int

func (p pos) median() int {
	sort.Ints(p)
	return p[len(p)/2]
}

func (p pos) sumdist() int {
	m, sum := p.median(), 0
	for _, x := range p {
		sum += abs(x - m)
	}
	return sum
}

func main() {
	var crabs pos

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
