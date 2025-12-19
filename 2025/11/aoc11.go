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

	// mapping from 3-letter tags to unique integer IDs
	var IDs [26 * 26 * 26]int // map [a..z][a..z][a..z] to integer ID --- 17576 possible IDs
	var nextID int            // next available ID

	// id returns the unique integer ID for a 3-letter tags
	id := func(s string) int {
		k := idh(s) // compute hash key

		// lazily assign ID
		if IDs[k] == NullID { // check if not assigned yet
			nextID++        // advance next ID
			IDs[k] = nextID // assign new ID
		}

		return IDs[k]
	}

	// read graph edges as adjacency lists
	edges := make([][]int, MaxID)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// parse line: "src: dst1 dst2 dst3 ..."
		parts := strings.Split(input.Text(), ":")

		// src and [dst1 dst2 ...]
		src, dsts := parts[0], strings.Fields(parts[1]) // car and cons

		// add edges
		for _, d := range dsts {
			sk, dk := id(src), id(d)

			edges[sk] = append(edges[sk], dk)
		}
	}

	// part 1 use DFS to count all paths from 'you' to 'out'
	you, out := id("you"), id("out")

	stack := make([]int, 0, MaxID)

	stack = append(stack, you)
	for len(stack) > 0 {
		cur := stack[len(stack)-1] // pop
		stack = stack[:len(stack)-1]

		if cur == out {
			acc1++
			continue // don't expand further
		}

		// expand
		for _, nxt := range edges[cur] {
			stack = append(stack, nxt)
		}
	}

	// part 2 use DP to count all paths from 'svr' to 'out' that contain both 'dac' and 'fft'
	svr, dac, fft := id("svr"), id("dac"), id("fft")

	dp := make(map[uint32]int)

	// presence map to avoid cycles
	seen := make([]uint8, MaxID)

	var recount func(cur int, hasDac, hasFft bool) int

	recount = func(cur int, hasDac, hasFft bool) int {
		k := dpk(cur, hasDac, hasFft) // unique key

		// cycle detection
		if seen[cur] == Seen {
			return 0
		}

		// memoization check
		if count, ok := dp[k]; ok {
			return count
		}

		// base case
		if cur == out {
			if hasDac && hasFft {
				return 1 // valid path
			}
			return 0
		}

		seen[cur] = Seen                      // mark current node as seen
		defer func() { seen[cur] = Unseen }() // backtrack on return

		// explore neighbors
		count := 0
		for _, nxt := range edges[cur] {
			if seen[nxt] == Unseen { // recurse only if not already visited
				count += recount(
					nxt,
					hasDac || nxt == dac,
					hasFft || nxt == fft,
				)
			}
		}

		// memoize and return
		dp[k] = count
		return count
	}

	acc2 = recount(svr, svr == dac, svr == fft)

	fmt.Println(acc1, acc2, time.Since(t0))
}

const (
	// sugars for presence map
	Unseen = iota
	Seen
)

const (
	// MaxID is the maximum number of unique 3-letter IDs expected in input
	MaxID  = 616
	NullID = 0 // null ID value is invalid
)

// idh computes a hash for a 3-letter string
func idh(s string) int {
	return int(s[0]-'a')*26*26 + int(s[1]-'a')*26 + int(s[2]-'a')
}

// dpk computes a unique key for DP memoization
func dpk(cur int, hasDac, hasFft bool) uint32 {
	k := uint32(cur) << 2

	if hasDac {
		k |= 1 << 1
	}

	if hasFft {
		k |= 1
	}

	return k
}
