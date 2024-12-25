// aoc25.go --
// advent of code 2024 day 25
//
// https://adventofcode.com/2024/day/25
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-25: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	schems := make([]int, 0, 500)

	input := bufio.NewScanner(os.Stdin)
	var hash int
	for input.Scan() {
		row := input.Bytes()
		switch {
		case len(row) == 0:
			schems, hash = append(schems, hash), 0
		default:
			for _, c := range row {
				hash <<= 1
				if c == '#' {
					hash |= 1
				}
			}
		}
	}
	schems = append(schems, hash)

	count1 := 0
	for i, a := range schems[:len(schems)-1] {
		for _, b := range schems[i+1:] {
			if a^b == a+b {
				count1++
			}
		}
	}

	fmt.Println(count1) // part 1
}

//         ____  ___            _____      _____    _________ _______________   ________    _____
//         \   \/  /           /     \    /  _  \  /   _____/ \_____  \   _  \  \_____  \  /  |  |
//          \     /   ______  /  \ /  \  /  /_\  \ \_____  \   /  ____/  /_\  \  /  ____/ /   |  |_
//          /     \  /_____/ /    Y    \/    |    \/        \ /       \  \_/   \/       \/    ^   /
//         /___/\  \         \____|__  /\____|__  /_______  / \_______ \_____  /\_______ \____   |
//               \_/                 \/         \/        \/          \/     \/         \/    |__|
//          ___________   __  .__         _____         _________
//         /_   \   _  \_/  |_|  |__     /  _  \   ____ \_   ___ \
//          |   /  /_\  \   __\  |  \   /  /_\  \ /  _ \/    \  \/
//          |   \  \_/   \  | |   Y  \ /    |    (  <_> )     \____
//          |___|\_____  /__| |___|  / \____|__  /\____/ \______  /
//                     \/          \/          \/               \/
