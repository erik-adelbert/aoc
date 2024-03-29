// aoc14.go --
// advent of code 2021 day 14
//
// https://adventofcode.com/2021/day/14
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-14: initial commit

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

func popcnt(rule string, depth int) histo {
	if len(cache[rule][depth]) > 0 {
		return cache[rule][depth]
	}

	cache[rule][depth] = histo{rules[rule]: 1} // cache current rule byte product

	if depth > 1 { // subsequent rules
		l := string([]byte{rule[0], rules[rule]}) // left
		r := string([]byte{rules[rule], rule[1]}) // right
		cache[rule][depth] = merge(cache[rule][depth], popcnt(l, depth-1))
		cache[rule][depth] = merge(cache[rule][depth], popcnt(r, depth-1))
	}

	return cache[rule][depth]
}

func main() {
	var seed []byte

	const (
		depth1 = 10
		depth2 = 40
	)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if args := strings.Split(line, " -> "); len(args) == 2 {
			rules[args[0]] = args[1][0]
			cache[args[0]] = make([]histo, depth2+1) // allocate cache space to accommodate for new rule
		} else if line != "" {
			seed = []byte(line)
		}
	}

	extent := func(depth int) int64 {
		popcnts := make(histo)
		for _, b := range seed {
			popcnts[b]++
		}
		for i := range seed[:len(seed)-1] { // extract and count pairs from seed
			popcnts = merge(popcnts, popcnt(string(seed[i:i+2]), depth))
		}

		min, max := extrema(popcnts)
		return max - min
	}

	fmt.Println(extent(depth1)) // part1
	fmt.Println(extent(depth2)) // part2
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
