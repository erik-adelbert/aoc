package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Axis
const (
	X = iota
	Y
	Z
)

type vec [3]int

func (v vec) sign(s vec) vec {
	return vec{
		s[X] * v[X], s[Y] * v[Y], s[Z] * v[Z],
	}
}

func (v vec) sub(u vec) vec {
	return vec{
		v[X] - u[X], v[Y] - u[Y], v[Z] - u[Z],
	}
}

func (v vec) manh() int {
	return abs(v[X]) + abs(v[Y]) + abs(v[Z])
}

type reading []vec

// Reading is a vector field abstraction
func Reading(points []vec) reading {
	return append(points[:0:0], points...)
}

type rotator func() reading // rotations iterator

func (r reading) π2rots() rotator {
	rots := []struct {
		s, a vec // sign, axis
	}{
		{vec{-1, -1, +1}, vec{X, Y, Z}},
		{vec{-1, +1, -1}, vec{X, Y, Z}},
		{vec{+1, -1, -1}, vec{X, Y, Z}},
		{vec{+1, +1, +1}, vec{X, Y, Z}},
		{vec{-1, -1, -1}, vec{X, Z, Y}},
		{vec{-1, +1, +1}, vec{X, Z, Y}},
		{vec{+1, -1, +1}, vec{X, Z, Y}},
		{vec{+1, +1, -1}, vec{X, Z, Y}},
		{vec{-1, -1, -1}, vec{Y, X, Z}},
		{vec{-1, +1, +1}, vec{Y, X, Z}},
		{vec{+1, -1, +1}, vec{Y, X, Z}},
		{vec{+1, +1, -1}, vec{Y, X, Z}},
		{vec{-1, -1, +1}, vec{Y, Z, X}},
		{vec{-1, +1, -1}, vec{Y, Z, X}},
		{vec{+1, -1, -1}, vec{Y, Z, X}},
		{vec{+1, +1, +1}, vec{Y, Z, X}},
		{vec{-1, -1, +1}, vec{Z, X, Y}},
		{vec{-1, +1, -1}, vec{Z, X, Y}},
		{vec{+1, -1, -1}, vec{Z, X, Y}},
		{vec{+1, +1, +1}, vec{Z, X, Y}},
		{vec{-1, -1, -1}, vec{Z, Y, X}},
		{vec{-1, +1, +1}, vec{Z, Y, X}},
		{vec{+1, -1, +1}, vec{Z, Y, X}},
		{vec{+1, +1, -1}, vec{Z, Y, X}},
	}

	i, turned := 0, make(reading, len(r))

	return func() reading { // rotator
		if i >= len(rots) {
			return reading(nil)
		}

		rot := rots[i]
		for j, v := range r {
			turned[j] = vec{
				v[rot.a[X]], v[rot.a[Y]], v[rot.a[Z]],
			}.sign(rot.s)
		}
		i++
		return turned
	}
}

var (
	reads []reading
	scans reading
	fixed map[vec]bool
)

func init() {
	reads = make([]reading, 0, 32)
	scans = append(make([]vec, 0, 28), vec{0, 0, 0})
	fixed = make(map[vec]bool, 337)
}

func list(m map[vec]bool) reading {
	list := make(reading, 0, len(m))
	for v := range m {
		list = append(list, v)
	}
	return list
}

func difs(r reading) reading {
	difs := make([]vec, len(r)-1)
	for i, v1 := range r[1:] { // v0, v1 = r[0], r[1] ...
		difs[i] = v1.sub(r[i]) // difs[0] == v1-v0, ...
	}

	return difs
}

func inter(a, b reading, first bool) reading {
	if len(a) > len(b) {
		a, b = b, a
	}

	m := make(map[vec]bool, len(a))
	for _, v := range a {
		m[v] = true
	}
	inter := make(reading, 0, len(m))
	for _, v := range b {
		if m[v] {
			if inter = append(inter, v); first {
				break
			}
		}
	}
	return inter
}

func rebase(r reading, o vec) reading {
	based := make(reading, len(r))
	for i, v := range r {
		based[i] = v.sub(o)
	}
	return based
}

func main() {
	points := make(reading, 0, 32)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch line := input.Text(); line != "" {
		case strings.Contains(line, "---"):
			if len(points) > 0 {
				reads = append(reads, Reading(points)) // clone points
				points = points[:0]                    // reset
			}
		default:
			args := strings.Split(line, ",")
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			z, _ := strconv.Atoi(args[2])
			points = append(points, vec{x, y, z})
		}
	}
	if points != nil {
		reads = append(reads, points) // last reading
	}

	for _, p := range reads[0] {
		fixed[p] = true
	}

	cur := reads[1:]
	for len(cur) > 0 {
		i := 0
		for _, r := range cur {
			if !rotal(r) {
				cur[i] = r
				i++
			}
		}
		cur = cur[:i]
	}

	fmt.Println(len(fixed)) // part1

	diam := 0
	for i, v0 := range scans {
		for _, v1 := range scans[i+1:] {
			if dist := v1.sub(v0).manh(); dist > diam {
				diam = dist
			}
		}
	}
	fmt.Println(diam) // part2
}

func rotal(r reading) bool {
	next := r.π2rots()

	rot := next()
	for rot != nil {
		if align(rot) {
			return true
		}
		rot = next()
	}
	return false
}

func align(r reading) bool {
	index := func(r reading, x vec) int {
		for i, v := range r {
			if x == v {
				return i
			}
		}
		return len(r)
	}

	known := list(fixed)
	for a := X; a <= Z; a++ { // X, Y, Z
		sort.Slice(r, func(i, j int) bool {
			return r[i][a] <= r[j][a]
		})

		sort.Slice(known, func(i, j int) bool {
			return known[i][a] <= known[j][a]
		})

		rdifs, kdifs := difs(r), difs(known)

		const (
			All   = false
			First = !All
		)

		if matches := inter(rdifs, kdifs, First); len(matches) > 0 {
			pivot := matches[0]
			i := index(rdifs, pivot)
			j := index(kdifs, pivot)
			o := r[i].sub(known[j])

			rebased := rebase(r, o)
			if len(inter(known, rebased, All)) >= 12 {
				for _, v := range rebased {
					fixed[v] = true
				}
				scans = append(scans, o)
				return true
			}
		}
	}
	return false
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
