// aoc23.go --
// advent of code 2019 day 25
//
// https://adventofcode.com/2019/day/25
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <intcode_file>")
		return
	}

	raw, err := loadIntCode(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading IntCode program: %v\n", err)
		return
	}
	code := newIC(raw)

	in := make(chan int, 1)
	cpu := newIntCodeCPU(0, in)

	go cpu.run(code)

	// Interactive shell
	input := bufio.NewScanner(os.Stdin)

	var room *room
	var starship dungeon

	var skiproom bool
	for {
		if !skiproom {
			r, err := getroom(cpu)
			switch {
			case err != nil:
				fmt.Printf("%s\n\n", err)

				if strings.Contains(err.Error(), "hello") {
					fmt.Println("Exiting shell...")
					return
				}
			case r.name != "":
				room = r
			}
		}
		fmt.Printf("\n%v\n> ", room)

		if !input.Scan() {
			fmt.Println("\nError reading input. Exiting.")
			break
		}

		switch input.Text() {
		case "exit", "q":
			fmt.Println("Exiting shell...")
			return

		case "n", "north", "s", "south", "e", "east", "w", "west":
			d := input.Text()[0:1]
			if strings.Contains(room.doors, d) {
				writeln(in, DIRS[d[0]], false)
			} else {
				fmt.Println("You can't go that way.")
				skiproom = true
			}

		case "a", "automap":
			starship = automap(room, cpu, in)
			writeln(in, "inv", false)

		case "b", "breakin":
			skiproom = true
			if room.name != "Security Checkpoint" {
				fmt.Println("Not in Security Checkpoint")
				break
			}
			// find the cockpit door and try to enter
			for _, x := range room.doors {
				if starship[room.name][x] != "" {
					continue
				}
				if breakin(cpu, in, x) {
					fmt.Println("Exiting shell...")
					return
				}
			}

		case "i", "inv":
			writeln(in, "inv", false)

		case "r":
			skiproom = true
			fmt.Println("read:", cpu.readln())

		default:
			line := input.Text()
			switch {
			case strings.HasPrefix(line, "take"), line == "take", line == "t":
				skiproom = true
				if len(room.items) > 0 {
					available := strings.Join(room.items, " ")
					blacklisted := strings.Join(blacklist, " ")

					item := room.items[0]
					if len(line) > 5 {
						item = line[5:]
					}

					if strings.Contains(available, item) && !strings.Contains(blacklisted, item) {
						writeln(in, "take "+item, false)
						skiproom = false
					}
				}
			case strings.HasPrefix(input.Text(), "go"):
				skiproom = true
				if len(starship) == 0 {
					fmt.Println("Automap first")
					break
				}

				dst := strings.Join(strings.Fields(input.Text())[1:], " ")
				if dst == "" || dst == "in" {
					dst = "Security Checkpoint"
				}
				if _, ok := starship[dst]; !ok {
					fmt.Println("Invalid destination")
				} else if rr := autogo(cpu, in, starship, room.name, dst); rr != nil {
					room = rr
				}
			default:
				writeln(in, input.Text(), false)
			}
		}
	}
}

func loadIntCode(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	code, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(code), nil
}

type dungeon map[string]map[rune]string

func (d dungeon) add(name string) {
	if d[name] == nil {
		d[name] = make(map[rune]string, 4)
	}
}

type room struct {
	name  string
	desc  string
	items []string
	doors string
}

func getroom(cpu *IntCodeCPU) (*room, error) {
	room := &room{}

	const (
		PROLOG = iota
		ROOM
		TEAZ
		DOORS
		EPILOG
		ITEMS
	)

	var err error
	empty := 0
	state := PROLOG
	for {
		line := cpu.readln()
		// fmt.Println(line)
		switch state {
		case PROLOG:
			switch line {
			case "":
				if empty += 1; empty == 3 {
					state = ROOM
					empty = 0
				}
			case "Command?":
				return room, err
			default:
				// print error and ignore
				fmt.Printf("%s\n", line)
				empty = 0
			}
		case ROOM:
			room.name = strings.Trim(strings.Split(line, "==")[1], " ")
			state = TEAZ
		case TEAZ:
			room.desc = line
			state = DOORS
		case DOORS:
			switch line {
			case "":
				empty += 1
				switch empty {
				case 1:
					continue
				case 2:
					state = EPILOG
					empty = 0
				}
			case "Doors here lead:":
				continue
			default:
				room.doors += line[2:3]
			}
		case EPILOG:
			switch line {
			case "Items here:":
				state = ITEMS
			case "Command?":
				return room, err
			default:
				err = fmt.Errorf("%s", line)

				for _, pat := range []string{"hello", "checkpoint"} {
					if strings.Contains(line, pat) {
						return room, err
					}
				}

			}
		case ITEMS:
			switch line {
			case "":
				state = EPILOG
			default:
				room.items = append(room.items, line[2:])
			}
		}
	}
}

var blacklist = []string{"escape po\x17", "escape pod", "giant electromagnet", "photons", "molten lava", "mutex", "infinite loop"}

func writeln(in chan int, s string, trace bool) {
	if trace {
		fmt.Println(">", s)
	}
	for _, c := range s {
		in <- int(c)
	}
	in <- '\n'
}

func automap(root *room, cpu *IntCodeCPU, in chan int) dungeon {
	fmt.Println("Automapping...")
	starship := make(dungeon, 19)

	var reexplore func(*room)
	reexplore = func(r *room) {
		starship.add(r.name)

		// writeln := func(s string) {
		// 	// fmt.Println(">", s)
		// 	for _, c := range s {
		// 		in <- int(c)
		// 	}
		// 	in <- '\n'
		// }

		take := func(i string) {
			if slices.Index(blacklist, i) != -1 {
				return
			}

			// fmt.Println("> take", i)
			writeln(in, "take "+i, false)
			for i := 0; i < 4; i++ {
				if line := cpu.readln(); i == 1 {
					fmt.Println(line)
				}
			}
		}

		fmt.Println("-", r.name)
		for _, i := range r.items {
			fmt.Println("  -", i)
			take(i)
		}

		move := func(d rune) *room {
			writeln(in, DIRS[d], false)
			rr, _ := getroom(cpu)
			return rr
		}

		backward := []rune{'n': 's', 's': 'n', 'e': 'w', 'w': 'e'}

		unmove := func(d rune) *room {
			return move(backward[d])
		}

		for _, d := range r.doors {
			if starship[r.name][d] == "" {
				rr := move(d)
				starship[r.name][d] = rr.name

				starship.add(rr.name)
				starship[rr.name][backward[d]] = r.name

				if rr.name == "Security Checkpoint" {
					fmt.Printf("- %s (ignoring)\n", rr.name)
				} else {
					reexplore(rr)
				}
				unmove(d)
			}
		}
	}

	now := time.Now()
	reexplore(root)
	fmt.Println("Time:", time.Since(now))

	return starship
}

func autogo(cpu *IntCodeCPU, in chan int, starship dungeon, src, dst string) *room {
	fmt.Printf("Autogo: %s -> %s\n", src, dst)

	move := func(d rune) *room {
		writeln(in, DIRS[d], true)
		r, _ := getroom(cpu)
		return r
	}

	var path []rune
	seen := make(map[string]bool)

	var dfs func(string, string) *room
	dfs = func(src, dst string) *room {
		var rr *room

		if src == dst {
			return rr
		}

		seen[src] = true
		for d, r := range starship[src] {
			if !seen[r] {
				seen[r] = true
				path = append(path, d)

				if rr := dfs(r, dst); rr != nil {
					return rr
				}

				if r == dst {
					for _, d := range path {
						rr = move(d)
					}
					return rr
				}
				path = path[:len(path)-1]
			}
		}
		return nil
	}
	return dfs(src, dst)
}

func breakin(cpu *IntCodeCPU, in chan int, out rune) bool {
	inventory := make([]string, 0, 8)
	writeln(in, "inv", false)
	cpu.readln()
	cpu.readln()
INVENTORY:
	for {
		line := cpu.readln()
		switch line {
		case "Command?":
			break INVENTORY
		case "":
			continue
		}
		inventory = append(inventory, line[2:])
	}
	slices.Sort(inventory)
	// fmt.Println("inventory:", inventory, len(inventory))

	take := func(i string) {
		writeln(in, "take "+i, false)
		for i := 0; i < 4; i++ {
			if line := cpu.readln(); i == 1 {
				fmt.Println(line)
			}
		}
	}

	drop := func(i string) {
		writeln(in, "drop "+i, false)
		for i := 0; i < 4; i++ {
			if line := cpu.readln(); i == 1 {
				fmt.Println(line)
			}
		}
	}

	now := time.Now()
	tryexit := func() bool {
		writeln(in, DIRS[out], false)

		line := cpu.readln()
		for i := 0; i < 9; i++ {
			line = cpu.readln()
		}
		if strings.Contains(line, "proceed") {
			cpu.readln()
			fmt.Println(cpu.readln())
			fmt.Println("Time:", time.Since(now))
			return true
		}
		_, _ = getroom(cpu)

		return false
	}

	for _, x := range inventory[4:] {
		drop(x)
	}

	combs := mkcombs(len(inventory), 4)
	cur, combs := combs[0], combs[1:]

	for _, new := range combs[1:] {
		if tryexit() {
			return true
		}

		for _, x := range cur {
			if slices.Index(new, x) == -1 {
				drop(inventory[x])
				break
			}
		}

		for _, x := range new {
			if slices.Index(cur, x) == -1 {
				take(inventory[x])
				break
			}
		}

		cur = new
	}

	return false
}

func mkcombs(n, k int) [][]int {
	if k == 0 {
		// Base case: yield the empty combination
		return [][]int{{}}
	} else if k == n {
		// Base case: yield the combination (0, 1, 2, ..., n-1)
		res := make([][]int, 1)
		res[0] = make([]int, n)
		for i := 0; i < n; i++ {
			res[0][i] = i
		}
		return res
	} else {
		// Recursive case: yield combinations from combinationsGraycode(n-1, k)
		var res [][]int
		res = append(res, mkcombs(n-1, k)...)

		// For each combination from combinationsGraycode(n-1, k-1), append n-1 to it
		subs := mkcombs(n-1, k-1)
		for i := len(subs) - 1; i >= 0; i-- {
			comb := subs[i]
			// Create a new combination by adding n-1
			new := append(comb, n-1)
			res = append(res, new)
		}

		return res
	}
}

func (r room) String() string {
	var sb strings.Builder
	sb.WriteString("" +
		"== " + r.name + " ==\n" +
		r.desc + "\n" +
		"\n" +
		"Doors here lead:\n",
	)
	for _, d := range r.doors {
		sb.WriteString("- " + DIRS[d] + "\n")
	}
	if len(r.items) > 0 {
		sb.WriteString("\nItems here:\n")
		for _, i := range r.items {
			sb.WriteString("  - " + i + "\n")
		}
	}
	return sb.String()
}

var DIRS = []string{'n': "north", 's': "south", 'e': "east", 'w': "west"}
