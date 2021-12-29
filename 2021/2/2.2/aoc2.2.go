package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	x, y, aim := 0, 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		arg, _ := strconv.Atoi(strings.Fields(line)[1])
		switch line[0] {
		case 'f': // forward
			x += arg
			y += aim * arg
		case 'u': // up
			aim -= arg
		case 'd': // down
			aim += arg
		}
	}
	fmt.Println(x * y)
}
