package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	sum1, sum2 := 0, 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		input := Fields(input.Text())

		blocks := Split(input[1], ",")

		springs1 := input[0]
		blocks1 := make([]int, len(blocks))
		for i := range blocks {
			blocks1[i] = atoi(blocks[i])
		}

		springs2 := Join([]string{
			springs1, springs1, springs1, springs1, springs1,
		}, "?")
		blocks2 := make([]int, 5*len(blocks1))
		for i := 0; i < len(blocks2); i += len(blocks1) {
			copy(blocks2[i:], blocks1)
		}

		sum1 += solve(springs1, blocks1)
		sum2 += solve(springs2, blocks2)

	}
	fmt.Println(sum1, sum2)
}

func solve(springs string, blocks []int) int {
	springs = "." + strings.TrimRight(springs, ".")

	// rolling block arrangement counts
	// ie. dp table last 2 rows
	cur := make([]int, len(springs)+1)
	nxt := make([]int, len(springs)+1)

	cur[0] = 1 // match at start
	for i, c := range springs {
		if c == '#' {
			break
		}
		cur[i+1] = 1
	}

	for _, block := range blocks {
		build := 0
		// build current block in all possible locations
		for i, c := range springs {
			build++
			if c == '.' {
				build = 0 // restart
			}

			// propagate current arrangement count
			if c != '#' {
				nxt[i+1] = nxt[i]
			}

			if build >= block && springs[i-block] != '#' {
				// current block is buildable in nxt[i+1] ways
				// accumulate arrangements of i springs into all blocks up to current
				nxt[i+1] += cur[i-block]
			}
		}

		// current is base counts for next block
		cur, nxt = nxt, cur
		clear(nxt)
	}

	// arrangements of all springs into all blocks
	return cur[len(cur)-1]
}

var (
	Fields    = strings.Fields
	Join      = strings.Join
	Split     = strings.Split
	TrimRight = strings.TrimRight
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
