package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const MaxInt = int(^uint(0) >> 1)

func main() {
	olds := []int{MaxInt, MaxInt, MaxInt} // 3 last depths window

	n, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		depth, _ := strconv.Atoi(input.Text())
		if olds[0] < depth { // increase detected!
			n++
		}
		last := copy(olds, olds[1:]) // shift left
		olds[last] = depth           // update window
	}
	fmt.Println(n)
}
