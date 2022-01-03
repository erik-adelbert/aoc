package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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

func (p *pair) clone() *pair {
	t := *p
	if !p.null() && !p.leaf() {
		left, right := p.left.clone(), p.right.clone()
		t.left, t.right = left, right
	}
	return &t
}

func flatten(p *pair) []*pair {
	var flat []*pair
	switch {
	case p.null():
	case p.leaf():
		flat = append(flat, p)
	default:
		left, right := flatten(p.left), flatten(p.right)
		flat = append(flat, append(left, right...)...)
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

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// flag.Parse()
	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	pprof.StartCPUProfile(f)
	// 	defer pprof.StopCPUProfile()
	// }

	args := make([]*pair, 0, 128)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		args = append(args, newPair(tokenize(line)))
	}

	num := newPair()
	for _, x := range args {
		num = reduce(newPair(num, x.clone()))
	}
	fmt.Println(mag(num)) // part1

	jobs := make(chan [2]*pair)
	mags := make(chan int)

	go func() { // producer
		defer close(jobs)

		wp := &sync.WaitGroup{}
		wp.Add(len(args) / 5)
		for k := 0; k < len(args); k += len(args) / 5 {
			sub := args[k : k+len(args)/5]
			go func() { // sliced sub producer
				for i, a := range sub {
					for j, b := range args {
						if i != j {
							jobs <- [...]*pair{a.clone(), b.clone()}
						}
					}
				}
				wp.Done()
			}()
		}
		wp.Wait() // production done
	}()

	for i := 0; i < 16; i++ { // consumers
		go func() {
			for args := range jobs {
				a, b := args[0], args[1]
				mags <- mag(reduce(newPair(a, b)))
			}
		}()
	}

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
