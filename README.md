
# aoc

My advent of code programs

All programs have been carefully composed and run collectively in around 95ms(2022) and 254ms(2021).
Coding notes have proven useful for others on reddit.

Happy coding!

## Installation and benchmark

0. optionnally install [gocyclo](https://github.com/fzipp/gocyclo)
1. install [hyperfine](https://github.com/sharkdp/hyperfine)
2. `git clone` this repository somewhere in your `$GOPATH`
3. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
4. `$ cd 2022`
5. `$ make`
6. `$ make runtime && cat runtime.md`
7. explore the other `Makefile` goals
