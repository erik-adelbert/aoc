package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
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
	cn := snum{
		make([]int, len(sn.vals)),
		make([]int, len(sn.deps)),
	}
	copy(cn.vals, sn.vals)
	copy(cn.deps, sn.deps)
	return cn
}

func explode(sn snum) (snum, bool) {
	for i := 0; i < len(sn.deps); i++ {
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

func mag(sn snum) int {
	vals, deps := sn.vals, sn.deps
	for len(vals) > 1 {
		for i := 0; i < len(deps); i++ {
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

func (sn snum) String() string {
	var sb strings.Builder

	depth := 0
	sb.WriteRune('[')
	for i, v := range sn.vals {
		switch {
		case depth < sn.deps[i]:
			for depth < sn.deps[i] {
				sb.WriteRune('[')
				depth++
			}
		case depth > sn.deps[i]:
			for depth > sn.deps[i] {
				sb.WriteRune(']')
				depth--
			}
		case depth == sn.deps[i]:
			sb.WriteRune(',')
		}
		sb.WriteString(fmt.Sprintf(" %d ", v))
	}
	for depth > 0 {
		sb.WriteRune(']')
		depth--
	}
	sb.WriteRune(']')

	return sb.String()
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	args := make([]snum, 0, 128)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		args = append(args, SNum(tokenize(line)))
	}

	num := SNum()
	for _, sn := range args {
		num = reduce(SNum(num, clone(sn)))
	}
	fmt.Println(mag(num)) // part1

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
