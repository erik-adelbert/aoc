// aoc6.go --
// advent of code 2021 day 6
//
// https://adventofcode.com/2021/day/6
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-6: initial commit

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

func popcnt(p popcnts) uint64 {
	var popcnt uint64
	for _, n := range p {
		popcnt += n
	}
	return popcnt
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
			fmt.Println(popcnt(fishes)) // part1
		}
		incube(fishes[:]) // pass slice
	}
	fmt.Println(popcnt(fishes)) // part2
}
