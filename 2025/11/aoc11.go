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

	var IDs [26 * 26 * 26]int // map string hash to integer ID
	var nextID = 1

	id := func(s string) int {
		if IDs[h(s)] == 0 {
			IDs[h(s)] = nextID // assign new ID
			nextID++
		}

		return IDs[h(s)]
	}

	edges := make([][]int, MaxIDHint)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		parts := strings.Split(input.Text(), ":")

		src, dst := parts[0], strings.Fields(parts[1])

		for _, d := range dst {
			edges[id(src)] = append(edges[id(src)], id(d))
		}
	}

	// part 1 use DFS to count all paths
	you, out := id("you"), id("out")

	q := make([]int, 0, MaxIDHint)

	q = append(q, you)
	for len(q) > 0 {
		cur := q[len(q)-1]
		q = q[:len(q)-1]

		if cur == out {
			acc1++
			continue
		}

		for _, nxt := range edges[cur] {
			q = append(q, nxt)
		}
	}

	// part 2 use DP to count all paths from 'svr' to 'out' that contain both 'dac' and 'fft'
	svr, dac, fft := id("svr"), id("dac"), id("fft")

	dp := make(map[uint32]int)

	// presence map to avoid cycles
	seen := make([]int, MaxIDHint)

	var recount func(cur int, hasDac, hasFft bool) int

	recount = func(cur int, hasDac, hasFft bool) int {
		k := hk(cur, hasDac, hasFft)

		if seen[cur] > 0 {
			return 0
		}

		if v, ok := dp[k]; ok {
			return v
		}

		if cur == out {
			if hasDac && hasFft {
				return 1
			}
			return 0
		}

		seen[cur]++
		defer func() { seen[cur]-- }() // backtrack on return

		count := 0
		for _, nxt := range edges[cur] {
			if seen[nxt] == 0 {
				count += recount(
					nxt,
					hasDac || nxt == dac,
					hasFft || nxt == fft,
				)
			}
		}

		dp[k] = count
		return count
	}

	acc2 = recount(svr, svr == dac, svr == fft)

	fmt.Println(acc1, acc2, time.Since(t0))
}

const MaxIDHint = 616

func h(s string) int {
	return int(s[0]-'a')*26*26 + int(s[1]-'a')*26 + int(s[2]-'a')
}

func hk(cur int, hasDac, hasFft bool) uint32 {
	k := uint32(cur) << 2
	if hasDac {
		k |= 1 << 1
	}
	if hasFft {
		k |= 1
	}
	return k
}
