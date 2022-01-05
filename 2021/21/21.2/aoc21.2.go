package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type game struct { // unique game state
	c1, s1, c2, s2 uint64 // cell1, score1, cell2, score2
}

type wins [2]uint64

var cache map[game]wins

func init() {
	cache = make(map[game]wins, 17317)
}

var rolls = [...]uint64{ // all dice rolls
	3, 4, 5, 4, 5, 6, 5, 6, 7,
	4, 5, 6, 5, 6, 7, 6, 7, 8,
	5, 6, 7, 6, 7, 8, 7, 8, 9,
}

// solve is a minimax-like game space resolver
func solve(g game) wins {
	switch {
	case g.s1 >= 21:
		return wins{1, 0}
	case g.s2 >= 21:
		return wins{0, 1}
	}

	if _, ok := cache[g]; !ok { // new game!
		var count wins
		for _, r := range rolls { // play all
			c1 := (g.c1+r-1)%10 + 1 // one move at a time
			s1 := g.s1 + c1
			sub := solve(game{g.c2, g.s2, c1, s1}) // swap players
			count = wins{
				count[0] + sub[1], // update with swapped back players
				count[1] + sub[0],
			}
		}
		cache[g] = count
	}
	return cache[g]
}

const (
	p1 = iota // player1
	p2        // player2
)

func main() {
	var c [2]uint64 // player start cells

	i, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		args := strings.Split(line, ": ")
		n, _ := strconv.ParseUint(args[1], 10, 64)
		c[i] = n
		i++
	}

	stats := solve(game{c[p1], 0, c[p2], 0}) // solve all games
	fmt.Println(max(stats[p1], stats[p2]))
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
