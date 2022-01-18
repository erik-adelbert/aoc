package main

import (
	"fmt"
	"math"
)

// Axis
const (
	X = iota
	Y
)

type (
	point [2]int
	speed [2]int
)

// target area: x=150..193, y=-136..-86
var (
	min = point{150, -136}
	max = point{193, -86}
)

func hit(v speed) bool {
	var p point
	for { // shoot
		p[X] += v[X]
		p[Y] += v[Y]

		if v[X] != 0 {
			v[X]--
		}
		v[Y]--

		if p[X] > max[X] || p[Y] < min[Y] { // over/under shoot
			return false
		}

		λ := func(a int) bool {
			return min[a] <= p[a] && p[a] <= max[a]
		}

		if λ(X) && λ(Y) {
			return true
		}
	}
}

func main() {
	vmin := speed{
		int(math.Sqrt(float64(2 * min[X]))), // FPU rules!
		min[Y],
	}

	vmax := speed{
		max[X],
		int(math.Abs(float64(min[Y] + 1))), // since we use math...
	}

	fmt.Println((vmax[Y] + 1) * vmax[Y] / 2) // part1

	n := 0
	for vx := vmin[X]; vx <= vmax[X]; vx++ {
		for vy := vmin[Y]; vy <= vmax[Y]; vy++ {
			if hit(speed{vx, vy}) {
				n++
			}
		}
	}
	fmt.Println(n) // part2
}
