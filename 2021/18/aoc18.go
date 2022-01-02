package main

import (
	"bufio"
	"fmt"
	"os"
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
	s = strings.Replace(s, ",", " ", -1)
	s = strings.Replace(s, "[", " [ ", -1)
	s = strings.Replace(s, "]", " ] ", -1)
	fields := tokens(strings.Fields(s))
	return &fields
}

type pair struct {
	v           int
	left, right *pair
}

func (p *pair) String() string {
	var reprint func(p *pair) string

	reprint = func(p *pair) string {
		switch {
		case p == nil:
			return ""
		case p.leaf():
			return fmt.Sprint(p.v)
		default:
			return fmt.Sprintf("[ %s, %s ]", reprint(p.left), reprint(p.right))
		}
	}

	return reprint(p)
}

func (a *pair) eq(b *pair) bool {
	switch {
	case a == nil:
		return b == nil
	case b == nil:
		return a == nil
	case a == b:
		return true
	case a.leaf() != b.leaf():
		return false
	default:
		return a.left.eq(b.left) && a.right.eq(b.right)
	}
}

func (p *pair) null() bool {
	return p.v < -1 && p.left == nil && p.right == nil
}

func newPair(args ...interface{}) *pair {
	switch len(args) {
	case 0:
		return &pair{-1, nil, nil}
	case 1:
		if n, ok := args[0].(int); ok {
			return &pair{n, nil, nil}
		}
		in := args[0].(*tokens)
		for tok := in.shift(); tok != ""; tok = in.shift() {
			switch tok {
			case "[":
				return newPair(newPair(in), newPair(in))
			case "]":
				continue
			default:
				n, _ := strconv.Atoi(tok)
				return newPair(n)
			}
		}
		return &pair{-1, nil, nil}
	case 2:
		l, r := args[0].(*pair), args[1].(*pair)
		if l == nil || l.eq(&pair{-1, nil, nil}) {
			return r
		}
		return &pair{-1, l, r}
	}
	panic("newPair: unreachable")
}

func (p *pair) leaf() bool {
	return p.v > -1
}

func flatten(p *pair) []*pair {
	var flat []*pair
	switch {
	case p.null():
		break
	case p.leaf():
		flat = append(flat, p)
	default:
		flat = append(flat, append(flatten(p.left), flatten(p.right)...)...)
	}
	return flat
}

func explode(p *pair) *pair {
	// fmt.Println("explode")
	i, flat, done := 0, flatten(p), false

	var rexplode func(*pair, int) *pair
	rexplode = func(p *pair, depth int) *pair {
		if p.leaf() {
			i++
			return p
		}
		if p.left.leaf() && p.right.leaf() {
			if depth > 3 && !done {
				if i > 0 {
					flat[i-1].v += flat[i].v
				}
				if i+1 < len(flat)-1 {
					flat[i+2].v += flat[i+1].v
				}
				done = true
				return newPair(0)
			}
		}
		return newPair(
			rexplode(p.left, depth+1), rexplode(p.right, depth+1),
		)
	}
	return rexplode(p, 0)
}

func split(p *pair) *pair {
	done := false

	var resplit func(*pair) *pair
	resplit = func(p *pair) *pair {
		if p.leaf() {
			if p.v >= 10 && !done {
				done = true
				return newPair(newPair(p.v/2), newPair(p.v-p.v/2))
			}
			return p
		}
		return newPair(resplit(p.left), resplit(p.right))
	}
	return resplit(p)
}

func reduce(p *pair) *pair {
	const (
		xflag = iota + 1 // 1
		sflag            // 2
		both             // xflag | sflag == 3
	)

	cur, nxt, done := newPair(), p, 0
	for done != both {
		cur, done = nxt, both // reset, flag
		if nxt = explode(cur); !nxt.eq(cur) {
			done &= ^xflag // unflag
			continue
		}
		if nxt = split(cur); !nxt.eq(cur) {
			done &= ^sflag // unflag
			continue
		}
	}
	return cur
}

func mag(p *pair) int {
	if p.leaf() {
		return p.v
	}
	return 3*mag(p.left) + 2*mag(p.right)
}

func main() {
	lines := make([]string, 0, 128)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		lines = append(lines, line)
	}

	num := newPair()
	for _, line := range lines {
		num = reduce(newPair(num, newPair(tokenize(line))))
	}
	fmt.Println(mag(num)) // part1

	max := 0
	for i, a := range lines {
		for j, b := range lines {
			if i != j {
				pa, pb := newPair(tokenize(a)), newPair(tokenize(b))
				if n := mag(reduce(newPair(pa, pb))); n > max {
					max = n
				}
			}
		}
	}
	fmt.Println(max) // part2
}
