# Timings

fastest end-to-end timing minus `cat` time of 100+ runs for part1&2 in ms - mbair M1/16GB - darwin 23.6.0 - go version go1.23.3 darwin/arm64 - hyperfine 1.19.0 - 2024-12

## Installation and benchmark

0. optionally install [gocyclo](https://github.com/fzipp/gocyclo)
1. install [hyperfine](https://github.com/sharkdp/hyperfine)
2. `git clone` this repository somewhere in your `$GOPATH`
3. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
4. `$ cd 2024`
5. `$ make`
6. `$ make runtime && cat runtime.md`
7. explore the other `Makefile` goals

## Day 1: [The Tyranny of the Rocket Equation](https://adventofcode.com/2019/day/1)

This marks my journey with implementing IntCode.

For the first challenge, I focused on making the code mathematically readable aloud. Because, why not?

## Day 2: [1202 Program Alarm](https://adventofcode.com/2019/day/2)

I have put the basics in place, now the code supports:

- IntCode programs
- IntCode cpu with μcode for ADD, MUL, HLT, HCF and tracing

I'm solving part 1&2 statically. Let `f(n, v)` be the result of running our add/mul program:

| n\v | 0  | 1  |
|-----|----|----|
|  0  | f0 | δv |
|  1  | δn |    |

`f0 = f(0, 0)`, `δn = f(1, 0)` and `δv = f(0, 1)`
`f(n, v) = f(0, 0) + n*f(1, 0) + v*f(0, 1) ∀ n, v ∈ [0, len(code)[`

From there it is easy to answer parts 1 & 2.

## Day 5: [Sunny with a Chance of Asteroids](https://adventofcode.com/2019/day/5)

The IntCode CPU now has `input`, `output`, `parameter modes` (part 1), `jit`, `jif`, `lt` and `eq` (part2).

```bash
❯ make run
go run ./aoc5.go < input.txt
9938601 4283952
```

## Day 7: [Amplification Circuit](https://adventofcode.com/2019/day/7)

Ah wow! The cpu went to an overhaul, it supports:

- `id` and `tagged trace`
- ~~multi-channel input~~ (useless)
- output channel

Parts 1 & 2 cannot be solved with the same code because the parameters and amplifier cabling differ. The solution for Part 2 maps one IntCode CPU per amplifier and connects them all in a loop.

It is possible to output the disassembled trace by setting the environment variable $TRACE to true.

For now, the computing throughput is: `29.2ms / 20401 instr  = 1,43μs/instr <=> ~700kops/s`

The trace is interesting: we can see the amplifiers initializing concurrently and then sync themselves in a serialized loop.

It is easy now to build a pool of concurrent workers if needed.

```bash
❯ TRACE=true make run
go run ./aoc7.go < input.txt
cpu4: INP $8 9   <- 9
cpu1: INP $8 6   <- 6
cpu0: INP $8 5   <- 5
cpu4: ADD $8 $8 10       <- 9, 10
cpu1: ADD $8 $8 10       <- 6, 10
cpu0: ADD $8 $8 10       <- 5, 10
cpu4: JIT $19 1  <- 1, 426
.
.
.
cpu2: HLT 3, pc: 344
cpu3: INP $9 11943924    <- 11943924
cpu3: MUL $9 $9 2        <- 11943924, 2
cpu3: OUT $9 23887848    <- 23887848
cpu3: HLT 3, pc: 263
cpu4: INP $9 23887848    <- 23887848
cpu4: ADD $9 $9 1        <- 23887848, 1
cpu4: OUT $9 23887849    <- 23887849
cpu4: HLT 3, pc: 182
89603079
```

## Day 9: [Sensor Boost](https://adventofcode.com/2019/day/9)

Now the cpu now supports `vmem` and `relative base addressing`

```bash
cpu0: MUL $63 34463338 34463338  <- 34463338, 34463338
cpu0: LT $63 $63 34463338        <- 1187721666102244, 34463338
cpu0: JIT $63 53         <- 0, 53
cpu0: MUL $1000 1 3      <- 1, 3
cpu0: RBO 988    <- 988 (988)
cpu0: RBO @12    <- 3 (991)
cpu0: RBO $1000  <- 3 (994)
cpu0: RBO @6     <- 3 (997)
cpu0: RBO @3     <- 3 (1000)
cpu0: INP $1000 2        <- 2
```

## Day 11: [Space Police](https://adventofcode.com/2019/day/11)

Just plain fun in launching intcode cpus concurrently and interacting with them.
Goroutines are tailored to the task.

## Day 13: [Care Package](https://adventofcode.com/2019/day/13)

Multiple inputs aren’t happening here. More fun with IntCode!

## Day 15: [Oxygen System](https://adventofcode.com/2019/day/15)

The solution performs a DFS with backtracking to find the shortest path to the goal, then uses BFS to calculate distances from the goal to all other cells.

## Day 17: [Set and Forget](https://adventofcode.com/2019/day/17)

This challenge is an usual AoC grid problem but is augmented by the asynchonicity of the intcode CPU. It is also a straightforward way to ensure we can have a decent terminal session with it.

```bash
❯ make run
go run ./aoc17.go < input.txt
Main: A,A,B,C,C,A,B,C,A,B
Function A: L,12,L,12,R,12
Function B: L,8,L,8,R,12,L,8,L,8
Function C: L,10,R,8,R,12
Continuous video feed? n
........................#########..............
........................#.......#..............
........................#.......#..............
........................#.......#..............
........................#.......#..............
........................#.......#..............
........................#.......#..............
........................#.......#..............
........................v.......#..............
................................#..............
................................#..............
................................#..............
................................#########......
........................................#......
............#########...................#......
............#.......#...................#......
............#.......#...................#......
............#.......#...................#......
............#.......#...................#......
............#.......#...................#......
........#########...#.......#############......
........#...#...#...#.......#..................
#############...#...#.......#..................
#.......#.......#...#.......#..................
#.......#.......#...#.......#..................
#.......#.......#...#.......#..................
#.......#.......#...#############..............
#.......#.......#...........#...#..............
#.......#.......#.........#############........
#.......#.......#.........#.#...#..............
#########.......#.........#.#...#..............
................#.........#.#...#..............
................#############...#..............
..........................#.....#..............
..........................#.....#..............
..........................#.....#..............
..........................#.....#..............
..........................#.....#..............
....................#############..............
....................#.....#....................
..............#############....................
..............#.....#..........................
..............#.....#..........................
..............#.....#..........................
..............#.....#..........................
..............#.....#..........................
..............#.....###########................
..............#...............#................
..............#...............#.......#########
..............#...............#.......#.......#
..............#...............#.......#.......#
..............#...............#.......#.......#
..............#############...#.......#.......#
..........................#...#.......#.......#
..........................#...#.......#.......#
..........................#...#.......#.......#
..........................#...#...#############
..........................#...#...#...#........
..........................#...#########........
..........................#.......#............
..........................#.......#............
..........................#.......#............
..........................#.......#............
..........................#.......#............
..........................#########............
9544 1499679
```

## Day 19: [Tractor Beam](https://adventofcode.com/2019/day/19)

The IntCode CPU sustains a consistent `2μs/iop` (`~500k iops`), with `50ns/iop` spent on addressing. This is also an example of an always valid `intcode` function: we don't need to clone it over and over again, resetting the cpu works perfectly fine.

## Day 21: [Springdroid Adventure](https://adventofcode.com/2019/day/21)

The interactive sessions are incredible!

## Day 23: [Category Six](https://adventofcode.com/2019/day/23)

The solution is concurrent, race free and fast. I have wrapped intcode cpus into concurrent network machines that maintain packet queues as described in the challenge.
The result feels strong. After careful review, the VM delivers `~12.86 Miops` <-> `77ns/op`. It has a new internal opcode `100` for an `EOT` while waiting for an input. The design favors an always correct auto-ordering of related ops.

## Day 25: [Cryostasis](https://adventofcode.com/2019/day/25)

The IntCode program powers a full-fledged old-school text adventure game. I have managed to write an almost bug-free CLI that should work on any entry and is ok for now.

The solution is a game CLI with some solving builtins like:

- `a`, `automap`
- `go`, `go in`, `go <room>`
- `b`, `breakin`

Drawing the starship is not an easy challenge and I'm still figuring it out.

## How was it?
