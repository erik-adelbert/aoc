package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type packet struct {
	val  int
	list []packet
}

func (p packet) isint() bool {
	return p.val != -1
}

// If both values are integers, the lower integer should come first.
// If the left integer is lower, the inputs are right.
// If the left integer is higher,the inputs are not right.
// Otherwise, the inputs are the same integer; continue
//
// If both values are lists, compare the first value of each list, then the second value, and so on.
// If the left list runs out of items first, the inputs are right.
// If the right list runs out of items first, the inputs are not right.
// If the lists are the same length and no comparison makes a decision about the order, continue.
//
// If exactly one value is an integer, convert the integer to a list, then retry.
func cmp(a, b packet) int {
	switch {
	case a.isint() && b.isint():
		switch {
		case a.val < b.val:
			return -1
		case a.val > b.val:
			return 1
		}
	case !(a.isint() || b.isint()):
		for i := range a.list {
			if i >= len(b.list) {
				return 1
			}
			if r := cmp(a.list[i], b.list[i]); r != 0 {
				return r
			}
		}
		if len(b.list) > len(a.list) {
			return -1
		}
	case a.isint():
		return cmp(packet{-1, []packet{a}}, b)
	case b.isint():
		return cmp(a, packet{-1, []packet{b}})
	}
	return 0
}

func load(s []byte) packet {
	var rec func(int) (packet, int)

	rec = func(i int) (packet, int) {
		a := packet{val: -1}

		for ; i < len(s); i++ {
			switch s[i] {
			case '[':
				var b packet
				b, i = rec(i + 1)
				a.list = append(a.list, b)
				fallthrough
			case ',':
				continue
			case ']':
				return a, i
			}

			a.list = append(
				a.list, packet{val: atoi(s[i:])})
		}
		return a, i
	}

	a, _ := rec(0)
	return a
}

func main() {
	popcnt := 0
	packets := []packet{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bytes := input.Bytes()

		// part1
		if len(bytes) == 0 {
			a := packets[len(packets)-2]
			b := packets[len(packets)-1]

			if cmp(a, b) < 1 {
				popcnt += len(packets) / 2
			}
			continue
		}

		// part2
		packets = append(packets, load(bytes))
	}

	// part1
	fmt.Println(popcnt)

	// part2
	markers := []packet{
		load([]byte("[[2]]")),
		load([]byte("[[6]]")),
	}
	packets = append(packets, markers...)
	sort.Sort(byPacket(packets))

	key := 1
	for i, p := range packets {
		if cmp(p, markers[0])*cmp(p, markers[1]) == 0 {
			key *= i + 1
		}
	}
	fmt.Println(key)
}

// packet sorting interface
type byPacket []packet

func (a byPacket) Len() int           { return len(a) }
func (a byPacket) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPacket) Less(i, j int) bool { return cmp(a[i], a[j]) < 0 }

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) int {
	var n int
	for _, c := range s {
		if c < 48 || c > 57 {
			break
		}
		n = 10*n + int(c-'0')
	}
	return n
}
