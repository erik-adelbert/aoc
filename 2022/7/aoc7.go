package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var subdirs []int // subdir sizes for part2

// sort insert
func record(s int) {
	i := sort.SearchInts(subdirs, s)
	subdirs = append(subdirs, 0)
	copy(subdirs[i+1:], subdirs[i:])
	subdirs[i] = s
}

func file(line string) int {
	fields := strings.Fields(line)
	return atoi(fields[0])
}

func tree(input *bufio.Scanner) int {
	root := 0
	for input.Scan() {
		line := input.Text()

		switch line[0] {
		case 'd':
			// dir (discard)
		case '$':
			// ls is discarded
			fields := strings.Fields(line[2:])
			if fields[0] == "cd" {
				switch fields[1] {
				case "..":
					return root
				default:
					subdir := tree(input)
					root += subdir
					record(subdir)
				}
			}
		default:
			root += file(line)
		}
	}
	return root
}

func main() {
	var root int
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		_ = input.Text() // discard initial cd /
		root = tree(input)
	}

	// part1 sum
	smalls := 0
	for i := 0; subdirs[i] <= 100_000; i++ {
		smalls += subdirs[i]
	}

	// part2 binsearch
	i := sort.SearchInts(subdirs, root-40_000_000)

	fmt.Println(smalls, subdirs[i])
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int
	for _, c := range []byte(s) {
		n = 10*n + int(c-'0')
	}
	return n
}
