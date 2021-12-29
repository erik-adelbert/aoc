package main

import (
	"fmt"
	"strings"
)

type board [11]string // hhRhRhRhRhh h(allway), R(oom)

var costs = map[byte]int{
	'A': 1, 'B': 10, 'C': 100, 'D': 1000,
}

func goal(p byte) int {
	return 2 + 2*int(p-'A') // 'A': 2, 'B': 4, .. 'D': 8
}

func room(r int) bool {
	return 1 < r && r < 9 && r&1 == 0 // true for 2, 4, 6, 8
}

func (b board) free(s, t int) bool { // s(rc), t(arget)
	start, s, t := s, min(s, t), max(s, t)
	for i := s; i <= t; i++ { // check if hallway cells between s,t are free
		switch {
		case i == start || room(i): // skip starting cell and rooms
		case b[i] != ".": // hallway is not clear
			return false
		}
	}
	return true
}

func (b board) granted(r int, p byte) bool { // r(oom), p(awn)
	return len(b[r]) == strings.Count(b[r], ".")+strings.Count(b[r], string(p))
}

func (b board) pawn(r int) (byte, bool) { // r(oom)
	room := b[r]
	for _, c := range room {
		if c != '.' {
			return byte(c), true
		}
	}
	return 0, false
}

func (b board) pop(r int) (string, int) { // r(oom)
	var cell []rune
	dist, first := 0, true
	for _, c := range b[r] {
		dist++
		switch {
		case c == '.':
			cell = append(cell, c)
		case first:
			first = false
			cell = append(cell, '.')
		default:
			cell = append(cell, c)
			dist--
		}
	}
	return string(cell), dist
}

func (b board) put(r int, p byte) (string, int) { // r(oom), p(awn)
	bytes := []byte(b[r])
	if i := strings.Count(b[r], "."); i != 0 { // room has free slots
		bytes[i-1] = p // take the deeper one
		return string(bytes), i
	}
	return string(bytes), 0
}

func (b board) moves(r int) []int { // r(oom)
	p, _ := b.pawn(r) // pawn to move

	if r == goal(p) && b.granted(r, p) { // pawn already at destination
		return nil
	}

	if !room(r) { // pawn in the hallway, moving to goal is the only move
		if b.free(r, goal(p)) && b.granted(goal(p), p) {
			return []int{goal(p)} // move only if way is free and room is open
		}
		return nil
	}

	moves := make([]int, 0)
	for i := 0; i < len(b); i++ {
		switch {
		case i == r: // skip starting room
		case i != goal(p) && room(i): // enter no room except for...
		case i == goal(p) && !b.granted(i, p): // ...the goal one, if not closed
		case b.free(r, i):
			moves = append(moves, i)
		}
	}
	return moves
}

func (b board) move(s, t int) (board, int) { // s(ource), t(arget) -> board, cost
	nxt := b
	p, _ := b.pawn(s) // p(awn)
	n, dist := 0, 0
	nxt[s] = "."
	if room(s) {
		nxt[s], n = b.pop(s)
		dist += n
	}
	nxt[t] = string(p)
	if room(t) {
		nxt[t], n = b.put(t, p)
		dist += n
	}
	dist += abs(s - t)
	return nxt, dist * costs[p]
}

func (b board) solve() map[board]int {
	stack := []board{b}
	push := func(b board) {
		stack = append(stack, b)
	}
	pop := func() board {
		b := stack[len(stack)-1]
		stack, stack[len(stack)-1] = stack[:len(stack)-1], board{}
		return b
	}
	empty := func() bool {
		return len(stack) == 0
	}

	costs := map[board]int{b: 0}
	push(b)        // from start...
	for !empty() { // ...play all possible games
		b := pop()         // pop a (sub)game
		for i := range b { // for all cells
			if _, ok := b.pawn(i); !ok { // empty cell, nothing to do
				continue
			}
			for _, j := range b.moves(i) { // for all legal moves from cell i...
				sub, ncost := b.move(i, j) // ...play one
				ncost += costs[b]
				if costs[sub] == 0 || costs[sub] > ncost {
					costs[sub] = ncost // if it's a better move so far...
					push(sub)          // ...send subgame to resolution
				}
			}
		}
	}
	return costs
}

func main() {
	goal := board{
		".", ".", "AA", ".", "BB", ".", "CC", ".", "DD", ".", ".",
	}
	p1 := board{
		".", ".", "AB", ".", "DC", ".", "BA", ".", "DC", ".", ".",
	}
	fmt.Println(p1.solve()[goal])

	goal = board{
		".", ".", "AAAA", ".", "BBBB", ".", "CCCC", ".", "DDDD", ".", ".",
	}
	p2 := board{
		".", ".", "ADDB", ".", "DCBC", ".", "BBAA", ".", "DACC", ".", ".",
	}
	fmt.Println(p2.solve()[goal])
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
