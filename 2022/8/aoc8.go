// aoc8.go --
// advent of code 2022 day 8
//
// https://adventofcode.com/2022/day/8
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-8: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

// part indices
const (
	Part1 = iota
	Part2
)

func main() {
	counts := [2]int{0, 0}

	// store all axis
	M := make([][]byte, 0, 128)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		M = append(M, []byte(input.Text()))
	}

	W, Σ := mirror(M), transpose(M)
	Ͻ := mirror(Σ) // Ͻ reads reverse (lunate) Σ

	views := func(x, y int) [4][]byte {
		U := Ͻ[x][len(Ͻ[0])-y:] // up
		L := W[y][len(W[0])-x:] // left
		R := M[y][x+1:]         // right
		D := Σ[x][y+1:]         // down

		return [4][]byte{U, L, R, D}
	}

	// part2
	dist := func(o byte, axis []byte) (int, byte) {
		var (
			acc int
			x   byte
		)
		for _, x = range axis {
			acc++
			if x >= o {
				break
			}
		}
		return acc, x
	}

	for y, r := range M {
		for x, o := range r {
			seen := false // part1
			count := 1    // part2

			for _, v := range views(x, y) {
				d, h := dist(o, v)

				// part1
				if !seen && o > h {
					counts[Part1]++
					seen = true
				}

				// part2
				count *= d
			}

			// part2
			counts[Part2] = max(counts[Part2], count)
		}
	}

	fmt.Println(counts[Part1], counts[Part2])
}

func transpose(m [][]byte) [][]byte {
	t := mkmat(len(m[0]), len(m))
	for i := 0; i < len(t); i++ {
		r := t[i]
		for j := 0; j < len(r); j++ {
			r[j] = m[j][i]
		}
	}
	return t
}

func mirror(m [][]byte) [][]byte {
	t := mkmat(len(m[0]), len(m))
	for i := 0; i < len(t); i++ {
		r := t[i]
		for j := 0; j < len(r); j++ {
			r[j] = m[i][len(r)-(j+1)]
		}
	}
	return t
}

func mkmat(h, w int) [][]byte {
	r := make([]byte, h*w)
	m := make([][]byte, h)
	lo, hi := 0, w
	for i := range m {
		m[i] = r[lo:hi:hi]
		lo, hi = hi, hi+w
	}
	return m
}
