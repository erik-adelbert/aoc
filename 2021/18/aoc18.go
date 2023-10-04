// aoc18.go --
// advent of code 2021 day 18
//
// https://adventofcode.com/2021/day/18
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-18: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type tokens []string

func (t *tokens) shift() string {
	if len(*t) == 0 {
		return ""
	}

	head := (*t)[0]
	(*t), (*t)[0] = (*t)[1:], ""
	return head
}

func tokenize(s string) *tokens { // old school wow inside!
	r := strings.NewReplacer(
		",", " ",
		"[", " [ ",
		"]", " ] ",
	)
	fields := tokens(strings.Fields(r.Replace(s)))
	return &fields
}

type snum struct {
	vals []int
	deps []int
}

// SNum constructs a snailfish number from variadic args
func SNum(args ...interface{}) snum {
	var sn snum
	switch len(args) {
	case 0:
		return snum{}
	case 1:
		toks := args[0].(*tokens)
		// fmt.Println(*toks)
		sn = snum{
			make([]int, 0, len(*toks)),
			make([]int, 0, len(*toks)),
		}

		var depth int
		for _, t := range *toks {
			t = strings.Trim(t, " ")
			switch {
			case t[0] == '[':
				depth++
			case t[0] == ']':
				depth--
			default:
				n, _ := strconv.Atoi(t)
				sn.vals = append(sn.vals, n)
				sn.deps = append(sn.deps, depth-1)
			}
		}
	case 2:
		a, b := args[0].(snum), args[1].(snum)
		if reflect.DeepEqual(a, snum{}) {
			sn.vals = b.vals
			sn.deps = b.deps
		} else {
			sn.vals = append(a.vals, b.vals...)
			sn.deps = append(a.deps, b.deps...)
			for i := range sn.deps {
				sn.deps[i]++
			}
		}
	default:
		panic("illegal SNum() call")
	}
	return sn
}

func clone(sn snum) snum {
	// tedious but faster than append(sn.x[:0:0], sn.x...)
	cn := snum{
		make([]int, len(sn.vals)),
		make([]int, len(sn.deps)),
	}
	copy(cn.vals, sn.vals)
	copy(cn.deps, sn.deps)
	return cn
}

func mag(sn snum) int {
	vals, deps := sn.vals, sn.deps
	for len(vals) > 1 {
		for i := range deps {
			if deps[i] == deps[i+1] {
				vals[i] = 3*vals[i] + 2*vals[i+1]
				vals, _ = remove(vals, i+1)
				deps, _ = remove(deps, i+1)

				if deps[i] > 0 {
					deps[i]--
				}
				break
			}
		}
	}
	return vals[0]
}

func explode(sn snum) (snum, bool) {
	for i := range sn.deps {
		if sn.deps[i] != 4 {
			continue
		}

		if i > 0 {
			sn.vals[i-1] += sn.vals[i]
		}

		if i+2 < len(sn.vals) {
			sn.vals[i+2] += sn.vals[i+1]
		}

		sn.vals[i], sn.deps[i] = 0, 3
		if i+1 < len(sn.deps) {
			sn.vals, _ = remove(sn.vals, i+1)
			sn.deps, _ = remove(sn.deps, i+1)
		}
		return sn, true
	}
	return sn, false
}

func reduce(sn snum) snum {
	const (
		xflag = iota + 1 // 1
		sflag            // 2
		both             // xflag | sflag == 3
	)

	done := 0
	for done != both {
		var more bool
		done = both // reset
		if sn, more = explode(sn); more {
			done &^= xflag // unflag
			continue
		}
		if sn, more = split(sn); more {
			done &^= sflag // unflag
			continue
		}
	}
	return sn
}

func split(sn snum) (snum, bool) {
	ok := false

	for i, v := range sn.vals {
		if v < 10 {
			continue
		}

		l, r := v/2, v-v/2
		sn.vals[i] = l
		sn.deps[i]++
		sn.vals = insert(sn.vals, i+1, r)
		sn.deps = insert(sn.deps, i+1, sn.deps[i])
		ok = true
		break
	}
	return sn, ok
}

// String outputs a representation of a snum.
// With n the count of numbers in the snum, it regrows the snum into a b-tree
// in O(n^2) and recursivly prints the tree in O(n). It's slow and heavy but
// reliable and somehow easy to understand.
// The overall O(n^2) is acceptable: the sole purpose here is to help us
// compose, debug and tune the code. This is also why we have regrow and
// reprint instead of a single but intricate function.
func (sn snum) String() string {
	type bnode struct {
		v           int
		left, right *bnode
	}

	new := func() *bnode {
		return &bnode{v: -1}
	}

	var regrow func(r *bnode, v, d int) bool
	regrow = func(r *bnode, v, d int) bool {
		if d == 0 {
			if *r == *new() { // *new() is null
				r.v = v
				return true
			}
			return false
		}
		if r.v != -1 { // not a leaf
			return false
		}
		if r.left == nil { // grow one level
			r.left = new()
		}
		if regrow(r.left, v, d-1) {
			return true
		}
		if r.right == nil { // grow one level
			r.right = new()
		}
		return regrow(r.right, v, d-1)
	}

	var reprint func(*bnode) string
	reprint = func(r *bnode) string {
		switch {
		case r == nil:
			return ""
		case -1 < r.v:
			return fmt.Sprint(r.v)
		default:
			return fmt.Sprintf("[%s,%s]", reprint(r.left), reprint(r.right))
		}
	}

	root := &bnode{v: -1}
	for i, v := range sn.vals {
		regrow(root, v, sn.deps[i]+1)
	}

	return reprint(root)
}

func main() {
	args := make([]snum, 0, 128)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		args = append(args, SNum(tokenize(line)))
	}

	// fmt.Println(args[0])

	jobs := make(chan snum)
	mags := make(chan int)

	go func() { // producer
		defer close(jobs)

		for i, a := range args {
			for j, b := range args {
				if i != j {
					jobs <- SNum(clone(a), clone(b))
				}
			}
		}
	}()

	for i := 0; i < 4; i++ { // consumers
		go func() {
			for sn := range jobs {
				mags <- mag(reduce(sn))
			}
		}()
	}

	num := SNum()
	for _, sn := range args {
		num = reduce(SNum(num, clone(sn)))
	}
	fmt.Println(mag(num)) // part1

	i, max := 0, 0
	for n := range mags {
		if n > max {
			max = n
		}
		if i++; i >= len(args)*(len(args)-1) {
			break
		}
	}
	close(mags)
	fmt.Println(max) // part2
}

func remove(a []int, i int) ([]int, bool) {
	if i >= len(a) {
		return a, false
	}
	return append(a[:i], a[i+1:]...), true
}

func insert(a []int, i int, v int) []int {
	if len(a) == i { // nil or empty slice or after last element
		return append(a, v)
	}
	a = append(a[:i+1], a[i:]...) // i < len(a)
	a[i] = v
	return a
}
