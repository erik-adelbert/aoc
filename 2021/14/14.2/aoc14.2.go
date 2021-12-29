package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type histo map[byte]int64

var (
	rules map[string]byte
	cache map[string][]histo
)

func init() {
	rules = make(map[string]byte)
	cache = make(map[string][]histo)
}

func merge(a, b histo) histo {
	for k, v := range b {
		a[k] += v
	}
	return a
}

func count(rule string, depth int) histo {
	if len(cache[rule][depth]) > 0 {
		return cache[rule][depth]
	}

	cache[rule][depth] = histo{rules[rule]: 1} // cache current rule byte product

	if depth > 1 { // subsequent rules
		l := string([]byte{rule[0], rules[rule]}) // left
		r := string([]byte{rules[rule], rule[1]}) // right
		cache[rule][depth] = merge(cache[rule][depth], count(l, depth-1))
		cache[rule][depth] = merge(cache[rule][depth], count(r, depth-1))
	}

	return cache[rule][depth]
}

func main() {
	const depth = 40

	var seed []byte

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if args := strings.Split(line, " -> "); len(args) == 2 {
			rules[args[0]] = args[1][0]
			cache[args[0]] = make([]histo, depth+1) // allocate cache space to accomodate for new rule
		} else if line != "" {
			seed = []byte(line)
		}
	}

	counts := make(histo)
	for _, b := range seed {
		counts[b]++
	}
	for i := range seed[:len(seed)-1] { // extract and count pairs from seed
		counts = merge(counts, count(string(seed[i:i+2]), depth))
	}

	fmt.Println(extrema(counts))
}

func extrema(m histo) (int64, int64) {
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
