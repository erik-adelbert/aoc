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
)

func main() {

	XMAS := [][]string{
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

	MAS := [][]string{
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

	var matrix RuneMat
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		matrix = append(matrix, []rune(input.Text()))
	}

	count1 := 0
	for _, subMatrix := range XMAS {
		found := matrix.findAllSubMatrices(toRuneMat(subMatrix))
		// fmt.Println(toRuneMat(subMatrix), found)
		count1 += len(found)
	}

	count2 := 0
	for _, subMatrix := range MAS {
		found := matrix.findAllSubMatrices(toRuneMat(subMatrix))
		// fmt.Println(toRuneMat(subMatrix), found)
		count2 += len(found)
	}

	fmt.Println(count1, count2)

}

type RuneMat [][]rune

func toRuneMat(s []string) RuneMat {
	m := make(RuneMat, len(s))
	for i, row := range s {
		m[i] = []rune(row)
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

// findAllSubMatrices searches for all occurrences of a sub-matrix in the larger matrix,
// allowing jokers as wildcards in the sub-matrix.
func (m RuneMat) findAllSubMatrices(sm RuneMat) [][2]int {
	H, W := len(m), len(m[0])
	h, w := len(sm), len(sm[0])

	matchWithJokers := func(matrixChar, subMatrixChar rune) bool {
		return subMatrixChar == '*' || matrixChar == subMatrixChar
	}

	matches := [][2]int{}

	// slide through the larger matrix
	for j := 0; j <= H-h; j++ {
	HSCAN:
		for i := 0; i <= W-w; i++ {
			// check if the sub-matrix matches at this position
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if !matchWithJokers(m[j+y][i+x], sm[y][x]) {
						continue HSCAN
					}
				}
			}
			// if we get here, the sub-matrix matches
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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
