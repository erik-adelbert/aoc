package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() { // suboptimal but easy
	var seed []byte
	count := make(map[byte]int64)
	rules := make(map[string]byte)
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

	for _, b := range seed {
		count[b]++
	}
	nxt := make([]byte, 0, len(seed)*(1<<depth)+1)
	for j, cur := 0, seed; j < depth; j++ {
		nxt = nxt[:0] // reset
		for i := range cur[:len(cur)-1] {
			pair := cur[i : i+2]
			new := rules[string(pair)]
			nxt = append(nxt, []byte{pair[0], new}...)
			count[new]++
		}
		cur, nxt = append(nxt, cur[len(cur)-1]), []byte(nil)
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
