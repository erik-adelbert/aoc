package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

const (
	LO = iota
	HI
)

type sigcnt [2]int

type module struct {
	id   string
	kind byte

	acc int

	srcs map[*module]int
	dsts []*module

	cnt *sigcnt
}

type pulse struct {
	dst, src *module
	val      int
}

func (p pulse) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "(%v %d -> %v)", p.src, p.val, p.dst)
	return sb.String()
}

func newModule(s string, cnt *sigcnt) *module {
	m := new(module)

	if m.kind == 0 {
		switch {
		case s == "broadcaster":
			m.kind = 'b'
		case s[0] == '%', s[0] == '&':
			m.kind, s = s[0], s[1:]
		default:
			m.kind = 's'
		}
	}

	m.id = s
	m.srcs = make(map[*module]int)
	m.dsts = make([]*module, 0, 8)
	m.cnt = cnt

	return m
}

func (m *module) link(u *module) {
	u.srcs[m] = 0
	m.dsts = append(m.dsts, u)
}

func (m *module) emit(v int) (pulses []pulse) {
	pulses = make([]pulse, len(m.dsts))

	i := 0
	for _, d := range m.dsts {
		pulses[i] = pulse{dst: d, src: m, val: v}
		i++
	}

	return
}

func (m *module) broadcast(p pulse) []pulse {
	val := p.val
	m.cnt[val]++

	return m.emit(val)
}

func (m *module) and(p pulse) []pulse {
	src, val := p.src, p.val
	m.cnt[val]++

	m.srcs[src] = val

	m.acc = HI
	for _, v := range m.srcs {
		m.acc &= v
	}

	return m.emit(HI - m.acc)
}

func (m *module) flip(p pulse) []pulse {
	val := p.val
	m.cnt[val]++

	if val == LO {
		m.acc = 1 - m.acc
		return m.emit(m.acc)
	}

	return []pulse{}
}

func (m *module) sink(p pulse) []pulse {
	m.cnt[p.val]++
	return []pulse{}
}

func (m *module) run(p pulse) []pulse {
	module := []func(pulse) []pulse{
		'b': m.broadcast,
		'%': m.flip,
		'&': m.and,
		's': m.sink,
	}[m.kind]

	return module(p)
}

func (m *module) String() string {
	return string(m.kind) + m.id
}

type modules map[string]*module

func (ms modules) button(rxinv *module) {
	cur := make([]pulse, 0, 64)
	nxt := make([]pulse, 0, 64)

	probe := make(map[*module]int, 4)

	N := 1000
	if rxinv != nil {
		N = MaxInt
	}

	var npresses int
	for !(rxinv != nil && len(probe) == len(rxinv.srcs)) && npresses < N { // for sample and input alike
		cur = append(cur, pulse{
			src: nil,
			dst: ms["broadcaster"],
			val: LO,
		})

		for len(cur) > 0 {
			nxt = nxt[:0]

			for _, p := range cur {
				// part2
				if rxinv != nil && p.val == HI {
					if _, ok := rxinv.srcs[p.src]; ok && probe[p.src] == 0 {
						probe[p.src] = npresses + 1
					}
				}
				nxt = append(nxt, p.dst.run(p)...)
			}
			cur, nxt = nxt, cur
		}

		if npresses == 999 {
			// part1
			cnt := ms["broadcaster"].cnt
			part1 = cnt[LO] * cnt[HI]
		}
		npresses++
	}
	part2 = lcm(values(probe))
}

var part1, part2 int

func main() {

	var cnt sigcnt
	mods := make(modules, 64)

	links := make([][]string, 0, 64)

	// parse network
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := split(input.Text(), " -> ")

		m := newModule(args[0], &cnt)
		mods[m.id] = m

		links = append(links, append(split(args[1], ", "), m.id))
	}

	// relink ands find probe
	rxinv := ""
	for i := range links {
		ω := len(links[i]) - 1
		m := links[i][ω]
		for _, u := range links[i][:ω] {
			if u == "rx" {
				rxinv = m
			}
			if _, ok := mods[u]; !ok {
				x := newModule(u, &cnt)
				mods[x.id] = x
			}
			mods[m].link(mods[u])
		}
	}

	// press button and probe network
	mods.button(mods[rxinv])

	fmt.Println(part1, part2)
}

var split = strings.Split

func keys[T comparable, V any](m map[T]V) []T {
	list := make([]T, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	return list
}

func values[T comparable, V any](m map[T]V) []V {
	list := make([]V, 0, len(m))
	for _, v := range m {
		list = append(list, v)
	}
	return list
}

func lcm(nums []int) (Π int) {
	Π = 1
	for i := range nums {
		Π *= nums[i] / gcd(Π, nums[i])
	}
	return Π
}

// https://en.wikipedia.org/wiki/Binary_GCD_algorithm
func gcd(a, b int) int {
	u, v := uint(a), uint(b)

	if u == 0 {
		return b
	}

	if v == 0 {
		return a
	}

	// `|` is bitwise OR. `TrailingZeros` quickly counts a binary number's
	// trailing zeros, giving its prime factorization's exponent on two.
	log2 := bits.TrailingZeros(u | v)

	// `>>=` divides the left by two to the power of the right, storing that in
	// the left variable. `u` divided by its prime factorization's power of two
	// turns it odd.
	u >>= bits.TrailingZeros(u)
	v >>= bits.TrailingZeros(v)

	for u != v {
		if u < v {
			u, v = v, u
		}
		u -= v
		u >>= bits.TrailingZeros(u)
	}

	// `<<` multiplies the left by two to the power of the right.
	return int(u << log2)
}

const MaxInt = int(^uint(0) >> 1)

const DEBUG = false

func debug(format string, a ...any) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}
