package main

import (
	"bufio"
	"fmt"
	"os"
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
					record(subdir)
				}
			}
			// discard ls
		default:
			root += file(line)
		}
	}
	return root
}

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

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int
	for _, c := range []byte(s) {
		n = 10*n + int(c-'0')
	}
	return n
}
