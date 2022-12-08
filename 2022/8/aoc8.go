package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

const (
	Part1 = iota
	Part2
)

func mat(h, w int) [][]byte {
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
	t := mat(len(m[0]), len(m))
	for i := 0; i < len(t); i++ {
		r := t[i]
		for j := 0; j < len(r); j++ {
			r[j] = m[j][i]
		}
	}
	return t
}

func mirror(m [][]byte) [][]byte {
	t := mat(len(m[0]), len(m))
	for i := 0; i < len(t); i++ {
		r := t[i]
		for j := 0; j < len(r); j++ {
			r[j] = m[i][len(r)-(j+1)]
		}
	}
	return t
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	counts := [2]int{0, 0}

	// store all axis
	M := make([][]byte, 0, 128)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		M = append(M, []byte(input.Text()))
	}

	T := transpose(M)
	MM := mirror(M)
	MT := mirror(T)

	_, _, _, _ = M, T, MM, MT

	views := func(x, y int) [][]byte {
		U := MT[x][len(MT[0])-y:] // up
		L := MM[y][len(MM[0])-x:] // left
		R := M[y][x+1:]           // right
		D := T[x][y+1:]           // down

		return [][]byte{U, L, R, D}
		// return [][]byte{}
	}

	// part2
	dist := func(o byte, axis []byte) int {
		acc := 0
		for _, x := range axis {
			acc++
			if x >= o {
				break
			}
		}
		return acc
	}

	for y, r := range M {
		for x, o := range r {
			count := 1    // part2
			seen := false // part1
			for _, v := range views(x, y) {
				// part1
				if !seen && o > max(v) {
					counts[Part1]++
					seen = true
				}

				// part2
				count *= dist(o, v)
			}

			// part2
			if counts[Part2] < count {
				counts[Part2] = count
			}
		}
	}

	fmt.Println(counts)
}

func max(b []byte) byte {
	var m byte = 0
	for _, v := range b {
		if v > m {
			m = v
		}
	}
	return m
}
