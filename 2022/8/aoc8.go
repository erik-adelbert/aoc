package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	Part1 = iota
	Part2
)

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

func main() {
	counts := [2]int{0, 0}

	// store all axis
	M := make([][]byte, 0, 128)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		M = append(M, []byte(input.Text()))
	}

	MM := mirror(M)
	T := transpose(M)
	MT := mirror(T)

	views := func(x, y int) [4][]byte {
		U := MT[x][len(MT[0])-y:] // up
		L := MM[y][len(MM[0])-x:] // left
		R := M[y][x+1:]           // right
		D := T[x][y+1:]           // down

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
			count := 1    // part2
			seen := false // part1

			for _, v := range views(x, y) {

				d, h := dist(o, v)
				// fmt.Println(d, h, o, v)

				// part1
				if !seen && o > h {
					counts[Part1]++
					seen = true
				}

				// part2
				count *= d
			}

			// part2
			if counts[Part2] < count {
				counts[Part2] = count
			}
		}
	}

	fmt.Println(counts)
}
