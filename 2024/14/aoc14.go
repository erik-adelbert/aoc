// aoc14.go --
// advent of code 2024 day 14
//
// https://adventofcode.com/2024/day/14
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-14: initial commit
// 2024-12-15: general automated solution

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
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
	THRESH = 25
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

	var i, tx0, ty0 int
	for i < T0 && (tx0 == 0 || ty0 == 0) {
		devx, devy := stddev2D(robots)
		switch {
		case devy < THRESH:
			tx0 = i
		case devx < THRESH:
			ty0 = i
		}
		robots = move(robots, 1)
		i++
	}

	robots = move(robots, T0-i)
	prod1 := check(robots)

	time2 := easter(tx0, ty0)
	robots = move(robots, time2-T0)
	output(time2, H, W, robots)

	fmt.Println(prod1, time2)
}

func move(robots []Robots, t int) []Robots {
	for i := range robots {
		robots[i].pos.x = mod(robots[i].pos.x+t*robots[i].mov.x, W)
		robots[i].pos.y = mod(robots[i].pos.y+t*robots[i].mov.y, H)
	}
	return robots
}

func check(robots []Robots) int {
	acc := 1
	for _, n := range quads(robots) {
		acc *= n
	}
	return acc
}

func quads(robots []Robots) [4]int {
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
	return quads
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
func atoi(s string) int {
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}

	n := 0
	for i := 0; i < len(s); i++ {
		n = 10*n + int(s[i]-'0')
	}
	return n * neg
}

func sample(robots []Robots) [][]int {
	binmat := make([][]int, H)
	for j := range binmat {
		binmat[j] = make([]int, W)
	}

	// sub-sample the first half of the robots
	for i := range robots[:len(robots)/2] {
		r := robots[i]
		binmat[r.pos.y][r.pos.x] = 1
	}

	return binmat
}

func stddev2D(robots []Robots) (float64, float64) {
	var xs, ys []int

	sample := sample(robots)

	for i := 0; i < len(sample); i++ {
		for j := 0; j < len(sample[0]); j++ {
			if sample[i][j] == 1 {
				ys = append(ys, i) // Y-coordinates
				xs = append(xs, j) // X-coordinates
			}
		}
	}

	μX, μY := mean(xs), mean(ys)
	return stddev(xs, μX), stddev(ys, μY)
}

// compute the mean of a never empty slice of ints
func mean(coords []int) float64 {
	var sum int
	for _, c := range coords {
		sum += c
	}
	return float64(sum) / float64(len(coords))
}

// compute the standard deviation of a never empty slice of ints
func stddev(coords []int, μ float64) float64 {
	var sum float64
	for _, c := range coords {
		sum += (float64(c) - μ) * (float64(c) - μ)
	}
	return math.Sqrt(sum / float64(len(coords)))
}

// function to find the first coincidence of two cycles
// tx0, ty0 are the starting time of the cycles with periods H, W
func easter(tx0, ty0 int) int {
	px, py := H, W // periods of the cycles

	// solve t = px * k + kx, substitute into second equation
	// px * k + kx ≡ ky (mod py)
	// simplify: (px % py) * k ≡ (ky - kx) % py
	m1 := px % py
	diff := (ty0 - tx0 + py) % py

	// compute k using modular inverse
	inv := modinv(m1, py)
	k := (diff * inv) % py

	// compute the first coincidence time
	t := k*px + tx0

	return t
}

// modular inverse by extended euclidean algorithm
func modinv(a, m int) int {
	m0, x0, x1 := m, 0, 1

	for a > 1 {
		q := a / m
		t := m
		m = a % m
		a = t
		t = x0
		x0 = x1 - q*x0
		x1 = t
	}

	if x1 < 0 {
		x1 += m0
	}

	return x1
}
