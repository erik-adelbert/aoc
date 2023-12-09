package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

var (
	links    []node
	AAA, ZZZ int
)

func init() {
	links = make([]node, ZZZ+1)
	AAA, ZZZ = hash("AAA"), hash("ZZZ")
}

func main() {
	var input *bufio.Scanner

	var cmds string

	getcmds := func() {
		input.Scan()        // advance scanner
		cmds = input.Text() // first line
		input.Scan()        // consume empty line
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

	roots := make([]int, 0, 8)

	addnode := func() {
		input := input.Text()

		h, left, right := hash(input[0:3]), hash(input[7:10]), hash(input[12:15])
		if isroot(h) {
			roots = append(roots, h)
		}
		links[h] = mknode(left, right)
	}

	browse := func(start int, isgoal func(int) bool) int {
		cmd, step := cmdstepper()

		for node := start; !isgoal(node); step() {
			switch cmd() {
			case 'R':
				node = right(node)
			default:
				node = left(node)
			}
		}
		return step() - 1
	}

	// read input
	input = bufio.NewScanner(os.Stdin)
	getcmds()
	for input.Scan() {
		addnode()
	}

	// part1
	browseAAA := func() int {
		return browse(AAA, func(node int) bool {
			return node == ZZZ
		})
	}

	// part2
	browseAll := func() int {
		cycles := make([]int, len(roots))

		for i := range roots {
			cycles[i] = browse(roots[i], func(node int) bool {
				return isgoal(node)
			})
		}
		return lcm(cycles)
	}

	fmt.Println(browseAAA(), browseAll())
}

type node int

func mknode(left, right int) node {
	return node(right<<15 | left)
}

func left(h int) int {
	return int(links[h] & 0x7fff)
}

func right(h int) int {
	return int((links[h] >> 15) & 0x7fff)
}

func isroot(h int) bool {
	return (h>>10)&0x1f == 0 // last car is 'A'
}

func isgoal(h int) bool {
	return (h>>10)&0x1f == 25 // last car is 'Z'
}

func hash(s string) int {
	ctoi := func(b byte) int {
		return int(b - 'A')
	}

	return ctoi(s[2])<<10 + ctoi(s[1])<<5 + ctoi(s[0])
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

	if x == 0 {
		return "AAA"
	}

	var sb strings.Builder
	for n := int(x); n > 0; n >>= 5 {
		sb.WriteByte('A' + byte(n&0x1f))
	}

	return sb.String()
}
