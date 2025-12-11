// aoc11.go --
// advent of code 2025 day 11
//
// https://adventofcode.com/2025/day/11
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-11: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	t0 := time.Now()   // start timer
	var acc1, acc2 int // parts 1 and 2 accumulators

	var nextID int
	IDs := make(map[string]int, MaxIDHint)

	edges := make([][]int, MaxIDHint)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		parts := strings.Split(line, ":")
		src := parts[0]
		dst := strings.Fields(parts[1])

		if _, ok := IDs[src]; !ok {
			IDs[src] = nextID
			nextID++
		}

		for _, d := range dst {
			if _, ok := IDs[d]; !ok {
				IDs[d] = nextID
				nextID++
			}

			edges[IDs[src]] = append(edges[IDs[src]], IDs[d])
		}
	}

	// part 1 use DFS to count all paths
	you, out := IDs["you"], IDs["out"]

	var recount func(curr int, seen []int) int
	recount = func(curr int, seen []int) int {
		if curr == out {
			return 1
		}

		seen[curr]++
		defer func() { seen[curr]-- }()

		count := 0
		for _, next := range edges[curr] {
			if seen[next] == 0 {
				count += recount(next, seen)
			}
		}

		return count
	}

	// part 2 use DP to count all paths from 'svr' to 'out' that contain both 'dac' and 'fft'
	svr, dac, fft := IDs["svr"], IDs["dac"], IDs["fft"]

	type key struct {
		cur    int
		hasDac bool
		hasFft bool
	}
	dp := make(map[key]int)

	var recountWithNodes func(cur int, seen []int, hasDac, hasFft bool) int

	recountWithNodes = func(cur int, seen []int, hasDac, hasFft bool) int {
		key := key{cur, hasDac, hasFft}

		if seen[cur] > 0 {
			return 0
		}

		if v, ok := dp[key]; ok {
			return v
		}

		if cur == out {
			if hasDac && hasFft {
				return 1
			}
			return 0
		}

		seen[cur]++
		defer func() { seen[cur]-- }()

		count := 0
		for _, next := range edges[cur] {
			if seen[next] == 0 {
				count += recountWithNodes(
					next,
					seen,
					hasDac || next == dac,
					hasFft || next == fft,
				)
			}
		}

		dp[key] = count
		return count
	}

	acc1 = recount(you, make([]int, MaxIDHint))
	acc2 = recountWithNodes(svr, make([]int, MaxIDHint), svr == dac, svr == fft)

	fmt.Println(acc1, acc2, time.Since(t0))
}

const MaxIDHint = 615
