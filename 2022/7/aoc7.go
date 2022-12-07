package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	Part1 = iota
	Part2
)

var answers [2]int

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int
	for _, c := range []byte(s) {
		n = 10*n + int(c-'0')
	}
	return n
}

type node struct {
	name string
	size int
	link []*node
}

func file(line string) *node {
	fields := strings.Fields(line)
	return &node{
		name: fields[1],
		size: atoi(fields[0]),
	}
}

// subdirs sizes for part2
var subdirs []int

// sort insert
func record(s int) {
	i := sort.SearchInts(subdirs, s)
	subdirs = append(subdirs, 0)
	copy(subdirs[i+1:], subdirs[i:])
	subdirs[i] = s
}

func tree(name string, input *bufio.Scanner) *node {
	root := new(node)
	root.name = name

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
					// subdir
					subdir := tree(fields[1], input)
					root.size += subdir.size
					root.link = append(root.link, subdir)

					// part1 counting
					if subdir.size <= 100000 {
						answers[Part1] += subdir.size
					}

					// part2 memoization
					record(subdir.size)
				}
			}
		default:
			// file
			leaf := file(line)
			root.size += leaf.size
			root.link = append(root.link, leaf)
		}
	}

	return root
}

func main() {
	var root *node
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		_ = input.Text() // discard initial cd /
		root = tree("/", input)
	}

	// part2 binsearch
	i := sort.SearchInts(subdirs, root.size-40000000)
	answers[Part2] = subdirs[i]

	fmt.Println(answers[Part1], answers[Part2])
}
