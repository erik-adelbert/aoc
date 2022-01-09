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

func (v vec) sign(b vec) vec {
	return vec{
		v[X] * b[X], v[Y] * b[Y], v[Z] * b[Z],
	}
}

func (v vec) sub(a vec) vec {
	return vec{
		v[X] - a[X], v[Y] - a[Y], v[Z] - a[Z],
	}
}

func (v vec) manh() int {
	return abs(v[X]) + abs(v[Y]) + abs(v[Z])
}

type reading []vec

func Reading(points []vec) reading {
	return append(points[:0:0], points...)
}

func (r reading) π2rots() <-chan reading { // rotations iterator
	c := make(chan reading)

	rots := []struct {
		s, a vec
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

	go func() {
		defer close(c)
		for _, rot := range rots {
			turned := make([]vec, len(r))
			for i, v := range r {
				turned[i] = vec{
					v[rot.a[0]], v[rot.a[1]], v[rot.a[2]],
				}.sign(rot.s)
			}
			c <- turned
		}
	}()

	return c
}

var (
	reads []reading
	scans []vec
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
	for i, v1 := range r[1:] {
		difs[i] = v1.sub(r[i]) // v0 = r[i], difs[i] == v1-v0
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
		switch line := input.Text(); {
		case line == "":
		case strings.Contains(line, "---"):
			if len(points) > 0 {
				reads = append(reads, Reading(points))
				points = points[:0]
			}
		default:
			args := strings.Split(line, ",")
			x, _ := strconv.Atoi(args[0])
			y, _ := strconv.Atoi(args[1])
			z, _ := strconv.Atoi(args[2])
			points = append(points, vec{x, y, z})
		}
	}
	reads = append(reads, Reading(points)) // last reading

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
	for rot := range r.π2rots() {
		if align(rot) {
			return true
		}
	}
	return false
}

func align(r reading) bool {
	index := func(r reading, v vec) int {
		for i, x := range r {
			if x == v {
				return i
			}
		}
		return len(r)
	}

	sort := func(r *reading, a int) {
		sort.Slice(*r, func(i, j int) bool {
			return (*r)[i][a] <= (*r)[j][a]
		})
	}

	read, known := Reading(r), list(fixed)
	for a := X; a <= Z; a++ {
		sort(&read, a)
		sort(&known, a)

		rdifs, kdifs := difs(read), difs(known)

		if matches := inter(rdifs, kdifs, true); len(matches) > 0 {
			pivot := matches[0]
			i := index(rdifs, pivot)
			j := index(kdifs, pivot)
			o := read[i].sub(known[j])

			rebased := rebase(read, o)
			if len(inter(known, rebased, false)) >= 12 {
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
