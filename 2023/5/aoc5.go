// aoc5.go --
// advent of code 2023 day 5
//
// https://adventofcode.com/2023/day/5
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2023-12-5: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	seeds1, seeds2 spans // part 1 & 2
)

var world [7]spans

func init() {
	seeds1 = mkSpans() // part1
	seeds2 = mkSpans() // part2

	for i := range world {
		world[i] = mkSpans()
	}
}

func main() {
	state := SEED
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		input := input.Text()
		switch {
		case len(input) == 0:
			state++
		case state == SEED:
			fields := Fields(input)[1:]
			for i := range fields {
				// tricks locate() by making a small interval for part1
				// no need to deal with 0 length intervals this way
				n := atoi(fields[i])
				seeds1 = append(seeds1,
					span{n, n + 1, 0}, // len == 1, slightly incorrect but ok
				)

				// build part2 every 2nd line
				if isodd(i) {
					seeds2 = append(seeds2,
						span{
							seeds1[i-1].src,
							seeds1[i-1].src + seeds1[i].src,
							0,
						},
					)
				}
			}
		case Contains(input, ":"): // discard header
		default:
			fields := Fields(input)
			world[state] = append(world[state],
				span{
					atoi(fields[1]),
					atoi(fields[1]) + atoi(fields[2]),
					atoi(fields[0]),
				},
			)
		}
	}

	fmt.Println(locate(seeds1), locate(seeds2)) // parts 1&2
}

func locate(seeds spans) (minloc int) {
	cur := mkSpans()

	// cur has a stack interface
	// it is convenient to add arbitray split intervals for the ones
	// that could match many ranges in a single mapping
	push := func(s span) {
		cur = append(cur, s)
	}

	pop := func() span {
		i := len(cur) - 1
		pop := cur[i]
		cur, cur[i] = cur[:i], span{}
		return pop
	}

	minloc = MaxInt
	for _, seed := range seeds {
		cur = cur[:0]
		push(seed)

		// remap seed ranges iteratively
		for _, maps := range world {
			nxt := mkSpans() // mapped span stack

		SPLITMAP:
			for len(cur) > 0 {
				br := pop()               // base range
				for _, cm := range maps { // current map
					// match by intersecting
					x := (span{max(cm.src, br.src), min(cm.end, br.end), 0})
					if x.src < x.end { // valid intersection

						// remap range for next step
						off := cm.dst - cm.src
						mapped := span{x.src + off, x.end + off, 0}
						nxt = append(nxt, mapped)

						// conditional LR split
						if br.src < x.src {
							left := span{br.src, x.src, 0}
							push(left)
						}

						if x.end < br.end {
							right := span{x.end, br.end, 0}
							push(right)
						}

						continue SPLITMAP // deal with new split ranges
					}
				}
				nxt = append(nxt, br) // no match yet, keep unchanged for next step
			}
			cur = nxt // prepare next mapping step
		}

		// get the lowest location from the last mappimg step
		for i := range cur {
			minloc = min(minloc, cur[i].src)
		}
	}

	return
}

const (
	// parsing DFA states
	SOL  = iota // soil
	FRT         // fertilizer
	WTR         // water
	LIG         // light
	TMP         // temperature
	HUM         // humidity
	LOC         // location
	SEED = -1   // seeds

	// default initial size
	SIZE = 16
)

type span struct { // range is a reserved word
	src, end, dst int
}

type spans []span

// make wrapper
func mkSpans() spans {
	return make(spans, 0, SIZE)
}

const MaxInt = int(^uint(0) >> 1)

// Go strings package wrappers/sugar
var Contains, Fields = strings.Contains, strings.Fields

func isodd(n int) bool {
	return n*(n&1) > 0
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
