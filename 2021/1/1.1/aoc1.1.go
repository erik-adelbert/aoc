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
	old := MaxInt
	n, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cur, _ := strconv.Atoi(input.Text())
		if cur > old { // increase!
			n++
		}
		old = cur
	}
	fmt.Println(n)
}
