package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	blocks := make(blocks, 0, 1500)

	input := bufio.NewScanner(os.Stdin)
	for id := 0; input.Scan(); id++ {
		args := split(input.Text(), "~")

		α := split(args[0], ",")
		β := split(args[1], ",")

		b := newBlock(id)
		for i := range α {
			b.α[i] = atoi(α[i])
			b.β[i] = atoi(β[i])
		}
		blocks = append(blocks, b)
	}

	fmt.Println(blocks.graph())
}

type blocks []*block

func (bs blocks) graph() (int, int) {
	bs.sort(blkZ)

	type base struct {
		blk, dep int
	}

	unsafe := make([]bool, len(bs))
	bases := make([]base, 0, len(bs))

	hmap := make([]int, 100)
	idxs := make([]int, 100)
	for i := range idxs {
		idxs[i] = MaxInt
	}

	for i, b := range bs {
		α, β := b.α, b.β
		s, e := 10*α[Y]+α[X], 10*β[Y]+β[X]
		h := b.β[Z] - b.α[Z] + 1

		step := 1
		if β[Y] > α[Y] {
			step = 10
		}

		top, old := 0, MaxInt
		ndown, base := 0, base{}

		for j := s; j <= e; j += step {
			top = max(top, hmap[j])
		}

		for j := s; j <= e; j += step {
			if i := idxs[j]; hmap[j] == top && i != old {
				old = i

				if ndown++; ndown == 1 {
					base = bases[old]
				} else {
					a, b := base, bases[old]
					for a.dep > b.dep {
						a = bases[a.blk]
					}
					for b.dep > a.dep {
						b = bases[b.blk]
					}
					for a.blk != b.blk {
						a, b = bases[a.blk], bases[b.blk]
					}
					base = a
				}
			}
			hmap[j] = top + h
			idxs[j] = i
		}

		if ndown == 1 {
			unsafe[old] = true
			base.blk, base.dep = old, bases[old].dep+1
		}

		bases = append(bases, base)
	}

	nsafe := 0
	for i := range unsafe {
		if !unsafe[i] {
			nsafe++
		}
	}

	nfall := 0
	for i := range bases {
		nfall += bases[i].dep
	}

	return nsafe, nfall
}

func (bs blocks) sort(cmp func(a, b *block) int) {
	slices.SortFunc(bs, cmp)
}

const (
	X = iota
	Y
	Z
)

type vec3 [3]int

func (v *vec3) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%v", *v)
	return sb.String()
}

type AABB struct {
	α, β *vec3
}

func (b AABB) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%v~%v", b.α, b.β)
	return sb.String()
}

type block struct {
	AABB
	id int
}

func blkZ(a, b *block) int {
	return a.α[Z] - b.α[Z]
}

func blkId(a, b *block) int {
	return a.id - b.id
}

func newBlock(id int) *block {
	return &block{
		AABB{new(vec3), new(vec3)},
		id,
	}
}

var split = strings.Split

const MaxInt = int(^uint(0) >> 1)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
