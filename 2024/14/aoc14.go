// aoc14.go --
// advent of code 2024 day 14
//
// https://adventofcode.com/2024/day/14
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-14: initial commit

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
)

const (
	H = 103
	W = 101

	// for sample input
	// H = 7
	// W = 11
	T0     = 100
	MAXDIM = 500
)

type Vec struct {
	x, y int
}

type Robots struct {
	pos Vec
	mov Vec
}

func main() {
	robots := make([]Robots, 0, MAXDIM) // arbitrary but educated guess

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fields := strings.Fields(input.Text())
		argp := strings.Split(fields[0], ",")
		argv := strings.Split(fields[1], ",")

		robots = append(robots, Robots{
			pos: Vec{atoi(argp[0][2:]), atoi(argp[1])},
			mov: Vec{atoi(argv[0][2:]), atoi(argv[1])},
		})
	}

	// part 1
	robots = move(robots, T0)
	prod1 := check(robots)

	// part 2 search
	time2 := T0
	for !easter(robots) {
		robots = move(robots, 1)
		time2++
	}
	fmt.Println(prod1, time2) // part 1 & 2
	// output(time2, H, W, robots)
}

func move(robots []Robots, t int) []Robots {
	for i := range robots {
		robots[i].pos.x = mod(robots[i].pos.x+t*robots[i].mov.x, W)
		robots[i].pos.y = mod(robots[i].pos.y+t*robots[i].mov.y, H)
	}
	return robots
}

func check(robots []Robots) int {
	ox, oy := W/2, H/2

	quads := [4]int{}
	for _, r := range robots {
		p := r.pos
		switch {
		case p.x < ox && p.y > oy:
			quads[0]++
		case p.x < ox && p.y < oy:
			quads[1]++
		case p.x > ox && p.y < oy:
			quads[2]++
		case p.x > ox && p.y > oy:
			quads[3]++
		}
	}

	prod1 := 1
	for _, n := range quads {
		prod1 *= n
	}

	return prod1
}

func output(timestamp int, H, W int, robots []Robots) {
	fname := fmt.Sprintf("aoc14-%d.png", timestamp)

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	black := color.RGBA{0, 0, 0, 255}
	for _, r := range robots {
		img.Set(r.pos.x, r.pos.y, black)
	}

	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	fmt.Printf("PNG file written to %s\n", fname)
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

// strconv.Atoi modified loop
// s is ^-?\d+.*$
// the suffix part .* is ditched
func atoi(s string) (n int) {
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}

	for i := 0; i < len(s) && s[i]-'0' < 10; i++ {
		n = 10*n + int(s[i]-'0')
	}
	return n * neg
}

func easter(robots []Robots) bool {
	n := len(robots)

	pos := make(map[Vec]struct{}, n)
	for _, r := range robots {
		pos[r.pos] = struct{}{}
	}

	return len(pos) == n
}
