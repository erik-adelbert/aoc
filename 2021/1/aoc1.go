package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// MaxInt is defined in the idiomatic way
const MaxInt = int(^uint(0) >> 1)

func main() {
	old1, old2, old3 := MaxInt, MaxInt, MaxInt // 3 last depths window

	n1, n2 := 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cur, _ := strconv.Atoi(input.Text())
		if old1 < cur { // increase!
			n1++
		}
		if old3 < cur { // increase!
			n2++
		}
		old1, old2, old3 = cur, old1, old2 // shift/update window
	}
	fmt.Println(n1) // part1
	fmt.Println(n2) // part2
}
