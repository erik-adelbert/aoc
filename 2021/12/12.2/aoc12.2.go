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

var paths map[string]bool

func (g graph) all(a, b string) {
	visits := make(map[*node]int)
	path := make([]*node, 0, len(g))

	var reall func(*node, *node, map[*node]int, []*node)
	reall = func(u, t *node, visits map[*node]int, path []*node) {
		seen := func(n *node) bool {
			return !n.big() && visits[n] >= n.limit
		}

		visits[u]++
		path = append(path, u)

		if u == t {
			var sb strings.Builder
			for _, n := range path {
				sb.WriteString(n.name)
			}
			if !paths[sb.String()] {
				paths[sb.String()] = true
			}
		} else {
			for _, v := range u.links {
				if !seen(v) {
					reall(v, t, visits, path)
				}
			}
		}

		visits[u]--
		i := len(path) - 1
		path, path[i] = path[:i], nil

		return
	}

	reall(g[a], g[b], visits, path)
}

func main() {
	paths = make(map[string]bool)
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
	fmt.Println(g)

	for _, n := range g {
		if n.name != "start" && n.name != "end" && !n.big() {
			n.limit = 2
			g.all("start", "end")
			n.limit = 1
		}
	}

	fmt.Println(len(paths))
}
