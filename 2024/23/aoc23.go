// aoc23.go --
// advent of code 2024 day 23
//
// https://adventofcode.com/2024/day/23
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-23: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Graph map[string]map[string]bool

func main() {
	nodes := make(Graph, 520)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		edge := strings.Split(input.Text(), "-")
		from, to := edge[0], edge[1]

		if nodes[from] == nil {
			nodes[from] = make(map[string]bool, 13)
		}
		nodes[from][to] = true

		if nodes[to] == nil {
			nodes[to] = make(map[string]bool, 13)
		}
		nodes[to][from] = true
	}

	count1 := nodes.tris()
	clique2 := nodes.largest()
	fmt.Println(count1, clique2) // part 1 & 2
}

// brute force the triangle count
func (g Graph) tris() int {
	sort3 := func(x ...string) [3]string {
		if x[0] > x[1] {
			x[0], x[1] = x[1], x[0]
		}
		if x[1] > x[2] {
			x[1], x[2] = x[2], x[1]
		}
		if x[0] > x[1] {
			x[0], x[1] = x[1], x[0]
		}
		return [3]string(x)
	}

	tris := make(map[[3]string]struct{})
	for a, edges := range g {
		if a[0] == 't' {
			for b := range edges {
				for c := range g[b] {
					if _, ok := edges[c]; ok {
						set := sort3(a, b, c)
						tris[set] = struct{}{}
					}
				}
			}
		}
	}
	return len(tris)
}

// get the the largest clique by expanding from each node
func (g Graph) largest() string {
	seen := make(map[string]bool, len(g))
	var clique, largest []string

	for n0, edges := range g {
		if !seen[n0] {
			clique = []string{n0} // start with a single node

			// expand the clique
			for n1 := range edges {
				// for each node n1 connected to n0
				neighs := g[n1]

				// check if it is connected to all nodes in the clique
				match := true
				for _, n2 := range clique {
					if _, ok := neighs[n2]; !ok {
						match = false
						break
					}
				}

				// if yes, add it to the clique
				if match {
					seen[n1] = true
					clique = append(clique, n1)
				}
			}

			if len(clique) > len(largest) {
				largest = append([]string(nil), clique...)
			}
		}
	}

	slices.Sort(largest)
	return strings.Join(largest, ",")
}
