package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	idsum, pwsum := 0, 0

	input := bufio.NewScanner(os.Stdin)
	for i := 1; input.Scan(); i++ {
		valid, power := true, [3]int{}

		// input is ^Game \d+:\s(.*;\s)+(.*)$
		input := input.Text()
		draws := Split(input[Index(input, ": ")+2:], "; ") // ditch "^Game X: " prefix, split tail
		for j := range draws {
			draws := Split(draws[j], ", ") // split game draws
			for jj := range draws {
				datas := Fields(draws[jj])                   // split RGB component and count
				count, color := atoi(datas[0]), datas[1][:1] // get values, single letter color r, g, b

				colid := Index("bgr", color) // R, G, B

				// check for validity, part1
				valid = valid && count <= 14-colid

				// record power RGB component, part2
				power[colid] = max(power[colid], count)
			}
		}
		if valid {
			idsum += i // part1
		}
		pwsum += power[B] * power[G] * power[R] // part2
	}
	fmt.Println(idsum, pwsum)
}

// sugars
var Index, Fields, Split = strings.Index, strings.Fields, strings.Split

const (
	B = iota
	G
	R
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

var DEBUG = true

func debug(a ...any) (int, error) {
	if DEBUG {
		return fmt.Println(a...)
	}
	return 0, nil
}
