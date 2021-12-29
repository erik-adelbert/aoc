package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type coor struct {
	x, y int
}

type card struct {
	row [5]int
	col [5]int
	val map[int]coor
}

func NewCard() *card {
	c := new(card)
	c.val = make(map[int]coor)
	return c
}

func (c *card) add(n, y, x int) {
	c.row[y] += n
	c.col[x] += n
	c.val[n] = coor{x, y}
}

func (c *card) biff(n int) bool {
	if rc, ok := c.val[n]; ok {
		c.row[rc.y] -= n
		c.col[rc.x] -= n

		if c.row[rc.y] == 0 || c.col[rc.x] == 0 { // win
			return true
		}
	}
	return false
}

func (c *card) sum() int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += c.row[i]
	}
	return sum
}

func main() {
	draw := make([]int, 0, 128)
	deck := make([]*card, 0, 128)
	cur, row := NewCard(), 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := input.Text()
		data := strings.Fields(args)
		switch len(data) {
		case 0: // sep
			if row > 0 {
				deck = append(deck, cur)
				cur, row = NewCard(), 0
			}
		case 5: // cardboard row
			for col, s := range data {
				n, _ := strconv.Atoi(s)
				cur.add(n, row, col)
			}
			row++
		default: // first line
			data = strings.Split(data[0], ",")
			for _, s := range data {
				n, _ := strconv.Atoi(s)
				draw = append(draw, n)
			}
		}
	}
	if row > 0 { // no newline at eof
		deck = append(deck, cur)
	}

	// stack!! heavy but easy & reliable
	type item struct {
		n int
		c *card
	}

	stack := make([]item, 0, 128)

	push := func(n int, c *card) {
		stack = append(stack, item{n, c})
	}

	pop := func() (int, *card) {
		if len(stack) == 0 {
			return 0, nil
		}
		top := stack[len(stack)-1]                                // pop
		stack, stack[len(stack)-1] = stack[:len(stack)-1], item{} // remove
		return top.n, top.c
	}

	for _, n := range draw {
		i := 0
		for _, card := range deck {
			if card.biff(n) { // win
				push(n, card) // save
				continue
			}
			deck[i] = card
			i++
		}
		deck = deck[:i]
	}
	n, c := pop() // last to win
	fmt.Println(n * c.sum())
}
