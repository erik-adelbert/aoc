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

type nodes []*node

func (n *nodes) push(x *node) {
	*n = append(*n, x)
}

func (n *nodes) pop() *node {
	i := len(*n) - 1

	pop := (*n)[i]
	*n, (*n)[i] = (*n)[:i], nil
	return pop
}

type graph map[string]*node

var npath int

func (g graph) paths(a, b string) {
	visits := make(map[*node]int)
	path := make(nodes, 0, len(g))

	var repaths func(*node, *node)
	repaths = func(u, t *node) {
		seen := func(n *node) bool {
			return !n.big() && visits[n] > 0
		}

		visits[u]++
		path.push(u)

		if u == t {
			npath++
		} else {
			for _, v := range u.links {
				if !seen(v) {
					repaths(v, t)
				}
			}
		}

		path.pop()
		visits[u]--

		return
	}

	repaths(g[a], g[b])
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
