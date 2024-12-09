// aoc8.go --
// advent of code 2024 day 8
//
// https://adventofcode.com/2024/day/8
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-8: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	ASCIIMAX = 123 // [0-9A-Za-z] -> 48-122
	MAXDIM   = 50
)

type Point struct {
	y, x int
}

type City struct {
	H, W     int
	antennas [][]Point // antennas map
}

func (c City) inbounds(p Point) bool {
	return p.y >= 0 && p.x >= 0 && p.y < c.H && p.x < c.W
}

func main() {
	antennas := make([][]Point, ASCIIMAX)

	input := bufio.NewScanner(os.Stdin)
	h, w := 0, 0
	for input.Scan() {
		line := input.Text()
		for i, c := range line {
			if c != '.' {
				antennas[c] = append(antennas[c], Point{h, i})
			}
		}
		w = len(line)
		h++
	}

	city := City{h, w, antennas}
	count1 := antinodes(city, 1, 2)
	count2 := antinodes(city, 0, max(city.H, city.W))
	fmt.Println(count1, count2) // part 1 & 2
}

func antinodes(city City, dmin, dmax int) int { // min distance, max distance
	H, W := city.H, city.W

	antinodes := make([]int, H*W)
	for _, set := range city.antennas['0':ASCIIMAX] {
		for i, a := range set {
			for _, b := range set[i+1:] {
				for d := dmin; d < dmax; d++ { // distance factor
					δy, δx := b.y-a.y, b.x-a.x
					p1 := Point{a.y - δy*d, a.x - δx*d}
					p2 := Point{b.y + δy*d, b.x + δx*d}

					if city.inbounds(p1) {
						antinodes[p1.y*W+p1.x] = 1
					}

					if city.inbounds(p2) {
						antinodes[p2.y*W+p2.x] = 1
					}

					if !city.inbounds(p1) && !city.inbounds(p2) {
						break
					}
				}
			}
		}
	}
	count := 0
	for _, n := range antinodes {
		count += n
	}
	return count
}
