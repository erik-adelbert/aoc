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
	DIM      = 50
)

type Point struct {
	y, x int
}

type City struct {
	H, W     int
	antennas [][]Point // antennas map
	terrain  [][]rune  // terrain map
}

func (c City) is_free(p Point) bool {
	return c.terrain[p.y][p.x] == '.'
}

func (c City) inbounds(p Point) bool {
	return p.y >= 0 && p.x >= 0 && p.y < c.H && p.x < c.W
}

func main() {
	terrain := make([][]rune, 0, DIM)
	antennas := make([][]Point, ASCIIMAX)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		j, line := len(terrain), input.Text()
		for i, c := range line {
			if c != '.' {
				antennas[c] = append(antennas[c], Point{j, i})
			}
		}
		terrain = append(terrain, []rune(line))
	}

	city := City{len(terrain), len(terrain[0]), antennas, terrain}
	count1 := antinodes(city, 1, 2)
	count2 := antinodes(city, 0, DIM)
	fmt.Println(count1, count2) // part 1 & 2
}

func antinodes(city City, m, M int) int {
	H, W := city.H, city.W

	antinodes := make([]int, H*W)
	for _, set := range city.antennas['0':ASCIIMAX] {
		for i, a := range set {
			for _, b := range set[i+1:] {
				for s := m; s < M; s++ {
					δy, δx := b.y-a.y, b.x-a.x
					p1 := Point{a.y - δy*s, a.x - δx*s}
					p2 := Point{b.y + δy*s, b.x + δx*s}

					if !city.inbounds(p1) && !city.inbounds(p2) {
						break
					}
					if city.inbounds(p1) {
						antinodes[p1.y*W+p1.x] = 1
					}
					if city.inbounds(p2) {
						antinodes[p2.y*W+p2.x] = 1
					}
				}
			}
		}
	}
	count := 0
	for _, v := range antinodes {
		count += v
	}
	return count
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
