package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	part1 int   // small subdirs total size
	part2 []int // subdirs sizes for part2
)

// sort insert
func record(s int) {
	i := sort.SearchInts(part2, s)
	part2 = append(part2, 0)
	copy(part2[i+1:], part2[i:])
	part2[i] = s
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

					// part1 counting
					if subdir <= 100000 {
						part1 += subdir
					}

					// part2 memoization
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

	// part2 binsearch
	i := sort.SearchInts(part2, root-40000000)

	fmt.Println(part1, part2[i])
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
