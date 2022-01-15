// https://www.youtube.com/watch?v=JyrNC74r2SI&list=PLrwpzH1_9ufMLOB6BAdzO08Qx-9jHGfGg&index=25
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, _ := os.ReadFile("input.txt")
	text := strings.Fields(r.Replace(string(data)))

	type dim struct {
		min, max int
	}

	type flip struct {
		on int
		p  [3]dim
	}

	var prog []*flip

	const Off = 100_000
	for ; len(text) > 0; text = text[7:] {
		var f flip
		f.on,
			f.p[0].min, f.p[0].max,
			f.p[1].min, f.p[1].max,
			f.p[2].min, f.p[2].max =
			atoi(text[0]),
			atoi(text[1])+Off,
			atoi(text[2])+Off+1,
			atoi(text[3])+Off,
			atoi(text[4])+Off+1,
			atoi(text[5])+Off,
			atoi(text[6])+Off+1
		prog = append(prog, &f)
	}

	const N = 850

	remap := new([3][2 * Off]int)
	width := new([3][N]int)
	for _, f := range prog {
		for i, d := range f.p {
			remap[i][d.min] = 1
			remap[i][d.max] = 1
		}
	}

	for i := range remap {
		t := 0
		for j, v := range &remap[i] {
			t += v
			remap[i][j] = t
			width[i][t]++
		}
	}
	for _, f := range prog {
		for i, d := range f.p {
			f.p[i] = dim{remap[i][d.min], remap[i][d.max]}
		}
	}

	sw := new([N][N][N]byte)
	for _, f := range prog {
		for x := f.p[0].min; x < f.p[0].max; x++ {
			for y := f.p[1].min; y < f.p[1].max; y++ {
				for z := f.p[2].min; z < f.p[2].max; z++ {
					sw[x][y][z] = byte(f.on)
				}
			}
		}
	}
	total := 0
	for x := range sw {
		for y := range sw {
			for z := range sw {
				total += width[0][x] * width[1][y] * width[2][z] * int(sw[x][y][z])
			}
		}
	}
	fmt.Println(total)
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

var r = strings.NewReplacer(
	"x=", "",
	"y=", "",
	"z=", "",
	"on", "1",
	"off", "0",
	",", " ",
	"..", " ",
)
