package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var c [2]int // player cells

	i, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), ": ")
		c[i], _ = strconv.Atoi(args[1])
		i++
	}

	const (
		p1 = iota // p(layer)1
		p2
	)

	other := func(p int) int {
		return (p + 1) & 1
	}

	var s [2]int        // scores
	p, n, d := p1, 0, 0 // player, roll count, dice value
	for {
		for i := 0; i < 3; i++ { // 3 dice rolls
			n++
			d = d%100 + 1
			c[p] = (c[p]+d-1)%10 + 1 // update player location
		}
		if s[p] += c[p]; s[p] >= 1000 { // update and check score
			fmt.Println(n * s[other(p)])
			break
		}
		p = other(p)
	}
}
