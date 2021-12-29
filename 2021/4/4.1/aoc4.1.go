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

func newCard() *card {
	c := new(card)
	c.val = make(map[int]coor)
	return c
}

func (c *card) add(n, y, x int) {
	c.row[y] += n
	c.col[x] += n
	c.val[n] = coor{x, y}
}

func (c *card) biff(n int) bool { // biff a number, return true if it's a win
	if rc, ok := c.val[n]; ok {
		c.row[rc.y] -= n
		c.col[rc.x] -= n
		// delete(c.val, n)

		if c.row[rc.y] == 0 || c.col[rc.x] == 0 { // win
			return true
		}
	}
	return false
}

func (c *card) sum() int { // sum remaining numbers (kinda)
	sum := 0
	for i := 0; i < 5; i++ {
		sum += c.row[i]
	}
	return sum
}

func main() {
	draw := make([]int, 0, 128)
	deck := make([]*card, 0, 128)
	cur, row := newCard(), 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := input.Text()
		data := strings.Fields(args)
		switch len(data) {
		case 0: // sep
			if row > 0 {
				deck = append(deck, cur)
				cur, row = newCard(), 0
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

	for _, n := range draw {
		for _, card := range deck {
			if card.biff(n) { // first win
				fmt.Println(n * card.sum())
				return
			}
		}
	}
}
