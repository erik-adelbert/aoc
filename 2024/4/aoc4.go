// aoc4.go --
// advent of code 2024 day 4
//
// https://adventofcode.com/2024/day/4
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-4: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// "*" (star) is a wildcard character that can match any letter.
var XMAS = [][]string{
	{
		"XMAS",
	},
	{
		"X***",
		"*M**",
		"**A*",
		"***S",
	},
	{
		"X",
		"M",
		"A",
		"S",
	},
	{
		"***X",
		"**M*",
		"*A**",
		"S***",
	},
	{
		"SAMX",
	},
	{
		"S***",
		"*A**",
		"**M*",
		"***X",
	},
	{
		"S",
		"A",
		"M",
		"X",
	},
	{
		"***S",
		"**A*",
		"*M**",
		"X***",
	},
}

var MAS = [][]string{
	{
		"M*M",
		"*A*",
		"S*S",
	},
	{
		"S*M",
		"*A*",
		"S*M",
	},
	{
		"S*S",
		"*A*",
		"M*M",
	},
	{
		"M*S",
		"*A*",
		"M*S",
	},
}

func main() {
	var matrix RuneMat

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		matrix = append(matrix, []rune(input.Text()))
	}

	var wg sync.WaitGroup

	count1 := 0
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, sub := range XMAS {
			matches := matrix.findAll(toRuneMat(sub))
			count1 += len(matches)
		}
	}()

	count2 := 0
	for _, sub := range MAS {
		matches := matrix.findAll(toRuneMat(sub))
		count2 += len(matches)
	}

	wg.Wait()

	fmt.Println(count1, count2)
}

type RuneMat [][]rune

func toRuneMat(s []string) RuneMat {
	m := make(RuneMat, len(s))
	for j, row := range s {
		m[j] = []rune(row)
	}
	return m
}

func (m RuneMat) String() string {
	var sb strings.Builder

	for _, row := range m {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

// findAll searches for all occurrences of a sub-matrix in the larger matrix,
// allowing wildcards in the sub-matrix.
func (m RuneMat) findAll(sm RuneMat) [][2]int {
	H, W := len(m), len(m[0])
	h, w := len(sm), len(sm[0])

	matches := make([][2]int, 0, 600) // pre-allocate for 600 matches (arbitrary)

	// slide through the larger matrix
	for j := 0; j <= H-h; j++ {
	HSCAN:
		for i := 0; i <= W-w; i++ {
			// check if the sub-matrix matches at this position
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if sm[y][x] != '*' && m[j+y][i+x] != sm[y][x] {
						continue HSCAN // mismatch!
					}
				}
			}
			// the sub-matrix matches
			matches = append(matches, [2]int{j, i})
		}
	}
	return matches
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
