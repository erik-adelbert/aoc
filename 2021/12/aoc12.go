package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type id int

func (i id) String() string {
	return names[i]
}

type edges []int

func (e edges) empty() bool {
	return len(e) == 0
}

func (e edges) contains(n int) bool {
	for _, v := range e {
		if n == v {
			return true
		}
	}
	return false
}

var (
	graph [][]id
	names []string
	adjs  []edges
	bigs  []bool
	nids  map[string]id
)

func init() {
	graph = make([][]id, 0, 16)
	names = make([]string, 0, 16)
	bigs = make([]bool, 16)
	nids = make(map[string]id, 16)
	node("end")   // nid 0
	node("start") // nid 1
}

func node(s string) id {
	if nid, ok := nids[s]; ok {
		return nid
	}

	nid := id(len(nids))
	nids[s] = nid
	bigs[nid] = unicode.IsUpper(rune(s[0]))

	names = append(names, s)
	graph = append(graph, make([]id, 0, 8))

	return nid
}

func link(args []string) {
	i, j := node(args[0]), node(args[1])
	if j != 0 {
		graph[i] = append(graph[i], j)
	}
	if i != 0 {
		graph[j] = append(graph[j], i)
	}
}

func adjacency() {
	adjs = make([]edges, len(names))

	for _, nid := range nids {
		if bigs[nid] {
			continue
		}
		for _, i := range graph[nid] {
			if adjs[nid] == nil {
				adjs[nid] = make(edges, 12)
			}
			if bigs[i] { // bigs act like teleports
				for _, j := range graph[i] {
					adjs[nid][j]++
				}
			} else {
				adjs[nid][i]++
			}
		}
	}
}

func dfs(part bool) int {
	type item struct {
		s, t id  // source, target
		d    int // dist
		m    bool
	}
	stack := make([]item, 0)

	empty := func() bool {
		return len(stack) == 0
	}

	push := func(s, t id, d int, m bool) {
		stack = append(stack, item{s, t, d, m})
	}

	pop := func() (id, id, int, bool) {
		i := len(stack) - 1
		pop := stack[i]
		stack = stack[:i]
		return pop.s, pop.t, pop.d, pop.m
	}

	hops := make(edges, 12)
	for i := range hops {
		hops[i] = 1
	}
	push(node("end"), node("start"), 1, part)

	count := 0
	for !empty() {
		s, t, n, m := pop() // source, target, npath, multipaths (part1/part2)
		hops[int(t)] = int(s)
		adj := adjs[s]
		for i, npath := range adj {
			if npath == 0 {
				continue
			}
			switch i {
			case 1:
				count += n * npath
			default:
				if unseen := i == 0 || !hops[2:t+1].contains(i); unseen || m {
					push(id(i), t+1, n*npath, m && unseen)
				}
			}
		}
	}
	return count
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		link(strings.Split(input.Text(), "-"))
	}
	adjacency()

	const (
		part1 = false
		part2 = !part1
	)

	fmt.Println(dfs(part1))
	fmt.Println(dfs(part2))
}
