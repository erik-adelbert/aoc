package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type node struct {
	name  string
	links []*node // could be []string
	limit int
}

func newNode(s string) *node {
	links := make([]*node, 0, 16)
	return &node{s, links, 1}
}

func link(a, b *node) {
	a.links = append(a.links, b)
	b.links = append(b.links, a)
}

func (n *node) big() bool {
	for _, r := range n.name {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func (n *node) String() string {
	return n.name
}

type nodes []*node

func (n nodes) push(x *node) nodes {
	return append(n, x)
}

func (n nodes) pop() (nodes, *node) {
	i := len(n) - 1

	pop := n[i]
	n, n[i] = n[:i], nil
	return n, pop
}

type graph map[string]*node

func (g graph) add(nodes []string) {
	for _, n := range nodes { // len(nodes) == 2
		if _, ok := g[n]; !ok {
			g[n] = newNode(n)
		}
	}
	link(g[nodes[0]], g[nodes[1]])
}

func (g graph) paths(a, b string) {
	visits := make(map[*node]int, 31)
	path := make(nodes, 0, len(g)) // stack as path!

	var repaths func(*node, *node)
	repaths = func(s, t *node) {
		seen := func(n *node) bool {
			return !n.big() && visits[n] >= n.limit
		}

		visits[s]++
		path = path.push(s)

		if s == t {
			var sb strings.Builder
			for _, n := range path {
				sb.WriteString(n.name[:2]) // cheat mode on
			}
			paths[sb.String()] = true
		} else {
			for _, v := range s.links {
				if !seen(v) {
					repaths(v, t)
				}
			}
		}

		path, s = path.pop()
		visits[s]--

		return
	}

	repaths(g[a], g[b])
}

var paths map[string]bool

func init() {
	paths = make(map[string]bool, 130513)
}

func main() {
	g := make(graph, 31)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		g.add(strings.Split(input.Text(), "-"))
	}
	g.paths("start", "end")
	fmt.Println(len(paths)) // part1

	for _, n := range g {
		if n.name != "start" && n.name != "end" && !n.big() {
			n.limit = 2
			g.paths("start", "end")
			n.limit = 1
		}
	}
	fmt.Println(len(paths)) // part2
}
