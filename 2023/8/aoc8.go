package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

var (
	graph    []node
	AAA, ZZZ int
)

func init() {
	AAA, ZZZ = hash("AAA"), hash("ZZZ")
	graph = make([]node, ZZZ+1)
}

func main() {
	var input *bufio.Scanner

	var cmds string

	roots := make([]int, 0, 8)

	getcmds := func() {
		input.Scan()
		cmds = input.Text()
		input.Scan()
	}

	addnode := func() {
		input := input.Text()

		h, left, right := hash(input[0:3]), hash(input[7:10]), hash(input[12:15])
		if isroot(h) {
			roots = append(roots, h)
		}
		graph[h] = node{left, right}
	}

	cmdstepper := func() (func() byte, func() int) {
		i := 0

		cmd := func() byte {
			// read from ring buffer
			return cmds[i%len(cmds)]
		}
		step := func() int {
			i++
			return i
		}

		return cmd, step
	}

	// part1
	browseAAA := func() int {
		cmd, step := cmdstepper()

		for node := AAA; node != ZZZ; step() {
			switch cmd() {
			case 'L':
				node = left(node)
			default:
				node = right(node)
			}
		}
		return step() - 1
	}

	// part2
	browseAll := func() int {
		cmd, step := cmdstepper()
		cycles := make([]int, len(roots))

		nstep, nroot := 0, 0
		for nroot < len(roots) {
			for i := range roots {
				if cycles[i] != 0 {
					continue
				}

				switch cmd() {
				case 'L':
					roots[i] = left(roots[i])
				default:
					roots[i] = right(roots[i])
				}
				if isgoal(roots[i]) { // first only
					cycles[i] = nstep + 1
					nroot++
				}

			}
			nstep = step()
		}
		return lcm(cycles)
	}

	// read input
	input = bufio.NewScanner(os.Stdin)
	getcmds()
	for input.Scan() {
		addnode()
	}

	fmt.Println(browseAAA(), browseAll())
}

type node struct {
	left, right int
}

func left(h int) int {
	return graph[h].left
}

func right(h int) int {
	return graph[h].right
}

func isroot(h int) bool {
	return h%26 == 0 // last car is 'A'
}

func isgoal(h int) bool {
	return h%26 == 25 // last car is 'Z'
}

func hash(s string) (h int) {
	ctoi := func(b byte) int {
		return int(b - 'A')
	}

	// encode s base 26
	h = (ctoi(s[0])*26+ctoi(s[1]))*26 + ctoi(s[2])
	return
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

	// `|` is bitwise OR. `trailing_zeros` quickly counts a binary number's
	// trailing zeros, giving its prime factorization's exponent on two.
	exp2 := bits.TrailingZeros(u | v)

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
	return int(u << exp2)
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

// debug print formatter
// ex:
//
//	fmt.Println(asString(hash("ZZZ")))
//	fmt.Println(asString(node.left))
type asString int

func (x asString) String() string {
	var sb strings.Builder

	buf := make([]byte, 0, 3)

	n := int(x)
	for n > 0 {
		buf = append(buf, byte(n%26)+'A')
		n /= 26
	}
	sb.WriteByte(buf[2])
	sb.WriteByte(buf[1])
	sb.WriteByte(buf[0])
	return sb.String()
}
