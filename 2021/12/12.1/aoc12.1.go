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
	links []*node
}

func newNode(s string) *node {
	links := make([]*node, 0, 16)
	return &node{s, links}
}

func (a *node) link(b *node) {
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
	var sb strings.Builder
	sb.WriteString(n.name)
	return sb.String()
}

type graph map[string]*node

var npath int

func (g graph) paths(a, b string) {
	visits := make(map[*node]int)
	path := make([]*node, 0, len(g))

	var repaths func(*node, *node, map[*node]int, []*node)
	repaths = func(u, t *node, visits map[*node]int, path []*node) {
		seen := func(n *node) bool {
			return !n.big() && visits[n] > 0
		}

		visits[u]++
		path = append(path, u)

		if u == t {
			npath++
		} else {
			for _, v := range u.links {
				if !seen(v) {
					repaths(v, t, visits, path)
				}
			}
		}

		visits[u]--
		i := len(path) - 1
		path, path[i] = path[:i], nil

		return
	}

	repaths(g[a], g[b], visits, path)
}

func main() {
	g := make(graph)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), "-")
		if _, ok := g[args[0]]; !ok {
			g[args[0]] = newNode(args[0])
		}
		if _, ok := g[args[1]]; !ok {
			g[args[1]] = newNode(args[1])
		}
		g[args[0]].link(g[args[1]])
	}
	g.paths("start", "end")

	fmt.Println(npath)
}
