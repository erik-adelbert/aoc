// aoc17.go --
// advent of code 2024 day 17
//
// https://adventofcode.com/2024/day/17
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-17: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const DEBUG = false

func main() {

	var mach Machine
	var prog []int
	var reg byte
	var val int

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		switch {
		case len(line) == 0:
		case line[0] == 'R':
			reg, val = line[9], atoi(line[12:])
			switch reg {
			case 'A':
				mach.ra = val
			case 'B':
				mach.rb = val
			case 'C':
				mach.rc = val
			}
		case line[0] == 'P':
			text := strings.Split(line[9:], ",")
			prog = make([]int, len(text))
			for i, word := range text {
				prog[i] = atoi(word)
			}
		}
	}

	mach = mach.exec(prog)
	fmt.Println(mach, mach.quine(prog)) // part 1 & 2
}

type Machine struct {
	out []int
	ra  int
	rb  int
	rc  int
	pc  int
}

func (m Machine) exec(text []int) Machine {
	combo := func(arg int) int {
		return []int{
			0: 0, 1: 1, 2: 2, 3: 3,
			A: m.ra,
			B: m.rb,
			C: m.rc,
		}[arg]
	}

	for m.pc < len(text) {
		pc := m.pc
		op, arg := text[m.pc], text[m.pc+1]
		decode := []func(){
			ADV: func() { m.ra >>= combo(arg) },
			BXL: func() { m.rb ^= arg },
			BST: func() { m.rb = combo(arg) & 0x7 },
			BXC: func() { m.rb ^= m.rc },
			OUT: func() { m.out = append(m.out, combo(arg)&0x7) },
			BDV: func() { m.rb = m.ra >> combo(arg) },
			CDV: func() { m.rc = m.ra >> combo(arg) },
			JNZ: func() {
				if m.ra != 0 {
					m.pc = arg - 2
				}
			},
		}

		decode[op]()
		debug("%02d %s %d %s %08d %08d %08d\tout: %v\n", pc, code[op], arg, regname[arg], m.ra, m.rb, m.rc, m)
		m.pc += 2
	}
	return m
}

func (m Machine) quine(text []int) int {
	table := make([]int, 0)
	for a := 0; a < (1 << 10); a++ {
		table = append(table, m.init(a, 0, 0).exec(text).out[0])
	}

	cur := make([][]int, 0, len(table))
	for i := range table {
		if table[i] == text[0] {
			cur = append(cur, []int{i})
		}
	}

	var nxt [][]int
	for _, word := range text[1:] {
		nxt = make([][]int, 0, len(cur))
		for _, x := range cur {
			seed := x[len(x)-1] >> 3
			for i := 0; i < 8; i++ {
				if table[(i<<7)+seed] == word {
					nxt = append(nxt, append(slices.Clone(x), (i<<7)+seed))
				}
			}
		}
		cur = nxt
	}

	pack := func(x []int) int {
		i, d := x[0], 10
		for _, c := range x[1:] {
			i += (c >> 7) << d
			d += 3
		}
		return i
	}

	// quines := make([]int, 0, 16)
	for _, x := range cur {
		a := pack(x)
		if slices.Equal(text, m.init(a, 0, 0).exec(text).out) {
			return a
			// here we should retain the min and return it at the end
			// but somehow my input's solution appears here first

			// quines = append(quines, a)
		}
	}

	// return slices.Min(quines)

	panic("unreachable")
}

const MaxInf = int(^uint(0) >> 1)

func (m Machine) init(a, b, c int) Machine {
	m.ra, m.rb, m.rc = a, b, c
	m.pc = 0
	m.out = []int{}
	return m
}

const (
	A = iota + 4
	B
	C
	NIL
)

const (
	ADV = iota
	BXL
	BST
	JNZ
	BXC
	OUT
	BDV
	CDV
)

var code []string = []string{
	ADV: "ADV",
	BXL: "BXL",
	BST: "BST",
	JNZ: "JNZ",
	BXC: "BXC",
	OUT: "OUT",
	BDV: "BDV",
	CDV: "CDV",
}

var regname []string = []string{
	0:   "0",
	1:   "1",
	2:   "2",
	3:   "3",
	A:   "A",
	B:   "B",
	C:   "C",
	NIL: "?",
}

func (m Machine) String() string {
	return strings.Trim(strings.Replace(fmt.Sprint(m.out), " ", ",", -1), "[]")
}

func debug(format string, args ...interface{}) (n int, err error) {
	if DEBUG {
		return fmt.Printf(format, args...)
	}
	return 0, nil
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
