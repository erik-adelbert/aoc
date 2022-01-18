package main

import (
	"container/heap"
	_ "embed"
	"fmt"
	"strings"
	"time"
	"unicode"
)

//go:embed input.txt
var input_day string

const hallwayY int = 1

var hallwayPos = []Pos{{1, 1}, {2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1}, {11, 1}}

type Pos struct {
	x, y int
}

type MoveCost struct {
	src  Pos
	dest Pos
	cost int
}

const empty byte = '.'

type World struct {
	maxY int
	grid []byte
	// grid string
}

func createWorld(lines []string) World {
	// WARNING: we assume no occupant is at home
	world := World{
		// grid: strings.Join(lines, "\n"),
		grid: []byte(strings.Join(lines, "\n")),
		maxY: len(lines),
	}
	return world
}

func (w World) String() string {
	return string(w.grid)
}

func index(p Pos) int {
	return p.y*14 + p.x
}

func replace(s string, p Pos, c byte) string {
	i := index(p)
	return s[:i] + string(c) + s[i+1:]
}

// An occupant can be '.', 'A', 'B', 'C', 'D', or 'a', 'b', 'c', 'd'
func (w World) occupant(p Pos) byte {
	return w.grid[index(p)]
}

func (w World) occupied(p Pos) bool {
	return w.grid[index(p)] != empty
}

func (w World) move(src, dest Pos, home bool) World {
	if w.occupied(dest) {
		panic("dest occupied")
	}

	s := make([]byte, len(w.grid))
	copy(s, w.grid)
	if !home {
		s[index(dest)] = w.occupant(src)
	} else {
		s[index(dest)] = byte(unicode.ToLower(rune(w.occupant(src))))
	}
	s[index(src)] = empty

	// s := w.grid
	// s = replace(s, src, empty)
	// if !home {
	// 	s = replace(s, dest, w.occupant(src))
	// } else {
	// 	s = replace(s, dest, byte(unicode.ToLower(rune(w.occupant(src)))))
	// }
	return World{
		grid: s,
		maxY: w.maxY,
	}
}

func (w World) moveHome(src, dest Pos) World {
	return w.move(src, dest, true)
}

func (w World) atHome(p Pos) bool {
	return unicode.IsLower(rune(w.occupant(p)))
}

// Returns true when [src+1..dest-1] is no occupied
func (w World) accessibleHallway(srcX, destX int) bool {
	if srcX == destX {
		return true
	}

	if srcX < destX {
		for x := srcX + 1; x <= destX; x++ {
			if w.occupied(Pos{x, hallwayY}) {
				return false
			}
		}
	}
	for x := srcX - 1; x >= destX; x-- {
		if w.occupied(Pos{x, hallwayY}) {
			return false
		}
	}

	return true
}

// List of free hallway positions accessible from column roomX
func (w World) accessiblePos(roomX int) []Pos {
	var res []Pos
	for _, h := range hallwayPos {
		if h.x < roomX {
			if w.occupied(h) {
				res = nil
			} else {
				res = append(res, h)
			}
		}
		if h.x > roomX {
			if w.occupied(h) {
				return res
			}
			res = append(res, h)
		}
	}
	return res
}

func (w World) blockedHallway1() bool {
	// If we find two elements in the hallway that have to pass through
	for x1 := 4; x1 <= 8; x1 += 2 {
		occupant1 := w.occupant(Pos{x: x1, y: hallwayY})
		for x2 := x1 + 2; x2 <= 8; x2 += 2 {
			occupant2 := w.occupant(Pos{x: x2, y: hallwayY})
			if occupant1 != occupant2 && occupant1 != empty && occupant2 != empty {
				if roomX(occupant1) > x2 && roomX(occupant2) < x1 {
					// fmt.Println("blocked1")
					return true
				}
			}
		}
	}
	return false
}

func (w World) blockedHallway2() bool {
	if w.occupant(Pos{8, 1}) == 'D' {
		homeX := 9
		freeSpace := 0
		if !w.occupied(Pos{10, 1}) {
			freeSpace += 1
			if !w.occupied(Pos{11, 1}) {
				freeSpace += 1
			}
		}
		nbForeign := 0
		for homeY := w.maxY - 2; homeY >= 2 && w.occupied(Pos{homeX, homeY}); homeY-- {
			if !w.atHome(Pos{homeX, homeY}) {
				nbForeign += 1
			}
		}
		if nbForeign > freeSpace {
			return true
		}
	}
	if w.occupant(Pos{4, 1}) == 'A' {
		homeX := 3
		freeSpace := 0
		if !w.occupied(Pos{2, 1}) {
			freeSpace += 1
			if !w.occupied(Pos{1, 1}) {
				freeSpace += 1
			}
		}
		nbForeign := 0
		for homeY := w.maxY - 2; homeY >= 2 && w.occupied(Pos{homeX, homeY}); homeY-- {
			if !w.atHome(Pos{homeX, homeY}) {
				nbForeign += 1
			}
		}
		if nbForeign > freeSpace {
			return true
		}
	}
	return false
}

// Returns first available home position
func (w World) freeHomeY(roomX int) (int, bool) {
	for y := w.maxY - 2; y >= 2; y-- {
		if !w.occupied(Pos{roomX, y}) {
			return y, true
		}
		if !w.atHome(Pos{roomX, y}) {
			return y, false
		}
	}
	return 0, false
}

func roomX(b byte) int {
	switch b {
	case 'A', 'a':
		return 3
	case 'B', 'b':
		return 5
	case 'C', 'c':
		return 7
	case 'D', 'd':
		return 9
	default:
		return 0
	}
}

func costMove(b byte) int {
	switch b {
	case 'A', 'a':
		return 1
	case 'B', 'b':
		return 10
	case 'C', 'c':
		return 100
	case 'D', 'd':
		return 1000
	default:
		return 0
	}
}

func (w World) moveHallwayToHome() (World, int) {
	cost := 0
	stop := false
	for !stop {
		stop = true
		for _, p := range hallwayPos {
			if w.occupied(p) {
				occupant := w.occupant(p)
				homeX := roomX(occupant)
				if w.accessibleHallway(p.x, homeX) {
					if homeY, ok := w.freeHomeY(homeX); ok {
						home := Pos{homeX, homeY}
						cost += manhattanDistance(p, home) * costMove(occupant)
						w = w.moveHome(p, home)
						stop = false
					}
				}
			}
		}
	}
	return w, cost
}

func (w World) moveRoomToHome(x int) (World, int) {
	cost := 0
	for roomY := 2; roomY <= w.maxY-2; roomY++ {
		p := Pos{x, roomY}
		if w.occupied(p) {
			occupant := w.occupant(p)
			homeX := roomX(occupant)

			if w.atHome(p) || p.x == homeX {
				return w, cost
			}

			homeY, ok := w.freeHomeY(homeX)
			if ok && w.accessibleHallway(x, homeX) {
				distance := manhattanDistance(p, Pos{homeX, hallwayY}) + manhattanDistance(Pos{homeX, hallwayY}, Pos{homeX, homeY})
				cost += distance * costMove(occupant)
				w = w.moveHome(p, Pos{homeX, homeY})
			} else {
				return w, cost
			}
		}
	}
	return w, cost
}

func (w World) moveRoomToHallway(roomX int) []MoveCost {
	var res []MoveCost

	for roomY := 2; roomY <= w.maxY-2; roomY++ {
		p := Pos{roomX, roomY}
		if w.occupied(p) {
			if w.atHome(p) {
				return res
			}
			occupant := w.occupant(p)
			for _, h := range w.accessiblePos(roomX) {
				cost := manhattanDistance(p, h) * costMove(occupant)
				res = append(res, MoveCost{src: p, dest: h, cost: cost})
			}
			// for _, h := range hallwayPos {
			// 	if w.accessibleHallway(roomX, h.x) {
			// 		cost := manhattanDistance(p, h) * costMove(occupant)
			// 		res = append(res, MoveCost{src: p, dest: h, cost: cost})
			// 	}
			// }
			return res
		}
	}
	return res
}

type State struct {
	world World
	cost  int
}

func (w World) step() []State {
	var res []State
	var cost int
	var c int

	// This is an optimization, not necessary
	if w.blockedHallway1() || w.blockedHallway2() {
		return res
	}

	w, c = w.moveHallwayToHome()
	cost += c

	// This is an optimization, not necessary
	for roomX := 3; roomX <= 9; roomX += 2 {
		w, c = w.moveRoomToHome(roomX)
		cost += c
	}

	for roomX := 3; roomX <= 9; roomX += 2 {
		for _, m := range w.moveRoomToHallway(roomX) {
			res = append(res, State{world: w.move(m.src, m.dest, false), cost: cost + m.cost})
		}
	}

	if len(res) == 0 && cost > 0 {
		res = append(res, State{w, cost})
	}

	return res
}

type node struct {
	World
	priority int
	index    int
}

func manhattanDistance(from, to Pos) int {
	absX := from.x - to.x
	if absX < 0 {
		absX = -absX
	}
	absY := from.y - to.y
	if absY < 0 {
		absY = -absY
	}
	return absX + absY
}

func heuristicCost(w World) int {
	var res int
	cpt := [10]int{}

	for x := 3; x <= 9; x += 2 {
		for y := w.maxY - 2; y >= 2; y-- {
			p := Pos{x, y}
			if w.occupied(p) && !w.atHome(p) {
				occupant := w.occupant(p)
				cpt[roomX(occupant)] += 1
				distance := cpt[roomX(occupant)] + manhattanDistance(p, Pos{roomX(occupant), hallwayY})
				res += distance * costMove(occupant)
			}
		}
	}
	for _, p := range hallwayPos {
		if w.occupied(p) {
			occupant := w.occupant(p)
			cpt[roomX(occupant)] += 1
			distance := cpt[roomX(occupant)] + manhattanDistance(p, Pos{roomX(occupant), hallwayY})
			res += distance * costMove(occupant)
		}
	}
	// fmt.Println(w, res)
	return res
}

func signature(w World) uint64 {
	var res uint64
	for _, p := range hallwayPos {
		o := w.occupant(p)
		res = res * 5
		if o == empty {
			res += 0
		} else if o >= 'a' {
			res += uint64(1 + o - 'a')
		} else if o >= 'A' {
			res += uint64(1 + o - 'A')
		}
	}
	for y := 2; y <= w.maxY-2; y++ {
		for x := 3; x <= 9; x += 2 {
			o := w.occupant(Pos{x, y})
			res = res * 5
			if o == empty {
				res += 0
			} else if o >= 'a' {
				res += uint64(1 + o - 'a')
			} else if o >= 'A' {
				res += uint64(1 + o - 'A')
			}
		}
	}
	return res
}

func path(start, to World) (path []World, distance int) {
	toSignature := signature(to)
	startSignature := signature(start)

	frontier := &PriorityQueue{}
	heap.Init(frontier)
	heap.Push(frontier, &node{World: start, priority: 0})

	cameFrom := map[uint64]World{startSignature: start}
	costSoFar := map[uint64]int{startSignature: 0}

	for {
		if frontier.Len() == 0 {
			// There's no path, return found false.
			return
		}
		var current World = heap.Pop(frontier).(*node).World
		currentSignature := signature(current)

		if currentSignature == toSignature {
			// Found a path to the goal.
			var path []World
			currentSignature = signature(current)
			for currentSignature != startSignature {
				path = append(path, current)
				current = cameFrom[currentSignature]
				currentSignature = signature(current)
			}
			return path, costSoFar[toSignature]
		}

		var next []State = current.step()
		for _, neighbor := range next {
			newCost := costSoFar[currentSignature] + neighbor.cost
			neighborSignature := signature(neighbor.world)
			if _, ok := costSoFar[neighborSignature]; !ok || newCost < costSoFar[neighborSignature] {
				costSoFar[neighborSignature] = newCost
				priority := newCost + heuristicCost(neighbor.world)
				heap.Push(frontier, &node{World: neighbor.world, priority: priority})
				cameFrom[neighborSignature] = current
			}
		}
	}
}

// A PriorityQueue implements heap.Interface and holds Items.
// Code copied from https://pkg.go.dev/container/heap@go1.17.5
type PriorityQueue []*node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// #############
// #...........#
// ###D#A#B#C###
//   #B#A#D#C#
//   #########
func Part1(input string) int {
	var d int
	input = strings.TrimSuffix(input, "\n")
	l := strings.Split(input, "\n")
	w := createWorld(l)
	t := "  #a#b#c#d#  "
	g := createWorld([]string{l[0], l[1], t, t, l[4]})
	_, d = path(w, g)
	return d
}

func Part2(input string) int {
	input = strings.TrimSuffix(input, "\n")
	l := strings.Split(input, "\n")
	l1 := "  #D#C#B#A#  "
	l2 := "  #D#B#A#C#  "
	lines := []string{l[0], l[1], l[2], l1, l2, l[3], l[4]}
	w := createWorld(lines)
	t := "  #a#b#c#d#  "
	g := createWorld([]string{l[0], l[1], t, t, t, t, l[4]})
	_, d := path(w, g)
	return d
	// return 0
}

func main() {
	fmt.Println("--2021 day 23 solution--")
	start := time.Now()
	fmt.Println("part1: ", Part1(string(input_day)))
	fmt.Println(time.Since(start))

	start = time.Now()
	fmt.Println("part2: ", Part2(string(input_day)))
	fmt.Println(time.Since(start))
}
