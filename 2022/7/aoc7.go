// aoc7.go --
// advent of code 2022 day 7
//
// https://adventofcode.com/2022/day/7
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-7: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

var subdirs []int // subdir sizes

func main() {
	var root int // root size
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	_ = input.Text() // discard initial cd /
	root = tree(input)

	// part1 sum
	smalls := 0
	for i := 0; subdirs[i] <= 100_000; i++ {
		smalls += subdirs[i]
	}

	// part2 binsearch
	i := sort.SearchInts(subdirs, root-40_000_000)

	fmt.Println(smalls, subdirs[i])
}

func tree(input *bufio.Scanner) int {
	root := 0
	for input.Scan() {
		line := input.Text()

		switch line[0] {
		case 'd':
			// discard dir
		case '$':
			fields := strings.Fields(line[2:])
			if fields[0] == "cd" {
				switch fields[1] {
				case "..":
					return root
				default:
					subdir := tree(input)
					root += subdir
					subdirs = append(subdirs, subdir)
				}
			}
			// discard ls
		default:
			root += file(line)
		}
	}
	slices.Sort(subdirs)
	return root
}

func file(line string) int {
	fields := strings.Fields(line)
	return atoi(fields[0])
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
