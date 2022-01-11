package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	histo map[byte]int64
	rules map[string]byte
)

func main() { // suboptimal but easy
	var seed []byte
	rules := make(rules)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if args := strings.Split(line, " -> "); len(args) == 2 {
			rules[args[0]] = args[1][0]
		} else if line != "" {
			seed = []byte(line)
		}
	}

	const depth = 10

	count := make(histo)
	for _, b := range seed {
		count[b]++
	}

	bufs := [][]byte{ // double buffering
		make([]byte, (1<<depth)*(len(seed)-1)+1),
		make([]byte, (1<<depth)*(len(seed)-1)+1),
	}

	// init slices
	cur, nxt := bufs[0], bufs[1]
	cur = cur[:len(seed)]
	copy(cur, seed)

	for j := 0; j < depth; j++ {
		n := len(cur) - 1 // last car
		nxt = nxt[:0]
		for i := range cur[:n] {
			cur := cur[i : i+2] // current pair from current seed
			new := rules[string(cur)]
			nxt = append(nxt, []byte{cur[0], new}...)
			count[new]++
		}
		cur, nxt = append(nxt, cur[n]), cur
	}

	min, max := extrema(count)
	fmt.Println(max - min)
}

func extrema(m map[byte]int64) (int64, int64) {
	const (
		MaxInt64 = int64(^uint64(0) >> 1)
		MinInt64 = -MaxInt64 - 1
	)

	min, max := MaxInt64, MinInt64
	for _, v := range m {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
