# Timings

| day | time |
|-----|-----:|
| 8 | 0.6 |
| 13 | 0.7 |
| 2 | 0.8 |
| 5 | 0.8 |
| 1 | 0.9 |
| 7 | 1.0 |
| 10 | 1.0 |
| 15 | 1.3 |
| 3 | 1.5 |
| 17 | 1.5 |
| 6 | 1.6 |
| 4 | 1.7 |
| 14 | 2.5 |
| 12 | 5.5 |
| 9 | 9.1 |
| 11 | 10.1 |
| 16 | 17.4 |
| total | 59.0 |

fastest end-to-end timing minus `cat` time of 100+ runs for part1&2 in ms - mbair M1/16GB - darwin 23.6.0 - go version go1.23.3 darwin/arm64 - hyperfine 1.19.0 - 2024-12

## Installation and benchmark

0. optionnally install [gocyclo](https://github.com/fzipp/gocyclo)
1. install [hyperfine](https://github.com/sharkdp/hyperfine)
2. `git clone` this repository somewhere in your `$GOPATH`
3. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
4. `$ cd 2024`
5. `$ make`
6. `$ make runtime && cat runtime.md`
7. explore the other `Makefile` goals

## Day 1: Historian Hysteria

*This year, I’m freelancing and available to take on projects—preferably in Go or Python. Please help spread the word!*

On this first day of AoC 2024, the challenge seems reasonable. For today’s solution, I’m using `sort` instead of `slices` because the problem invites presorted integer data. This choice allows the code to perform [binary searches](https://en.wikipedia.org/wiki/Binary_search) on the right dataset using the left dataset. When factoring in the presorting, the overall [runtime complexity](https://en.wikipedia.org/wiki/Time_complexity) is [O(n log n)](https://go.dev/src/sort/sort.go).

## Day 2: Red-Nosed Reports

My solution to today’s problem is fairly straightforward. It involves a left-to-right scan to ensure the safety constraints, evolving into a [tail-recursive call](https://en.wikipedia.org/wiki/Tail_call) to tolerate exactly one misplaced element. The only tricky part is when the first item is the misplaced one. In that case, we can simply check whether the report starting from the second element is safe. This approach applies anytime there’s a misplaced element, significantly simplifying the flow of control.

It’s interesting to see that the vast majority of other coders opted to generate all possible reports and brute-force the solution. In contrast, my approach requires at most two generated reports to validate or invalidate any given report and only one of them needs extra memory allocation and data copy.

## Day 3: Mull It Over

Go should be renowned for its blazing-fast [regular expression](https://en.wikipedia.org/wiki/Regular_expression) matching [engine](https://swtch.com/~rsc/regexp/). I thoroughly enjoyed today’s problem as it reminded me of [Exercice 12](https://clc-wiki.net/wiki/K%26R2_solutions:Chapter_1:Exercise_12) from *[The C Programming Language](https://en.wikipedia.org/wiki/The_C_Programming_Language)*. Combining the two concepts was a delight—and a satisfying way to wrap up the day.

Here's a fun fact collected while analysing my input:

```sh
❯ grep -oE '\w+\([0-9]+,[0-9]+\)' input.txt | sed -E 's/\(([0-9]+,[0-9]+)\)//' | sort | uniq | tr '\n' ' '
354who 953why from how mul select what when where who why
❯ grep -oE '[a-z]+' input.txt | sed -E 's/\(([0-9]+,[0-9]+)\)//' | sort | uniq | tr '\n' ' '
bin do don from how mul mulselect mulwhen mulwho mulwhy perl select t usr what when where who why
❯ grep -oE '.{0,10}perl.{0,10}' input.txt
!/usr/bin/perl@~mul?what
```

Today's problem seems to be a tribute to Perl by [u/topaz](https://x.com/ericwastl/status/1465082878073753600?lang=en), the creator of AoC.

Note: I may revisit this solution at some point.

## Day 4: Ceres Search

For today's solution, I opted for a sub-matrix matcher that makes multiple passes over the original matrix. It identifies 8 sub-matrices for 'XMAS' and 4 for a valid X-MAS. On the plus side, there's no need for boundary checks, and the matching process is fast. However, the downside is that it requires 12 passes over the original data. That said, the performance is acceptable for now.

```C
var MAS = [][]string{
    {
        "M*M",
        "*A*",
        "S*S",
    },
    {
        "S*M",
        "*A*",
        "S*M",
    },
    ...
```

## Day 5: Print Queue

Today's problem is certainly a brain teaser, but a straightforward approach can still be surprisingly effective. The key insight is that to move from one page number, a, to another, b, it must hold that `a ∈ rules[b]` which is `b is greater than a`. This is the heart of the challenge:

- Part 1 requires verifying that every page number in a given set satisfies this rule.
- Part 2 involves sorting the page numbers so that this rule holds true throughout.

```bash
❯ cat rules.txt | cut -d "|" -f 1 | sort | uniq | wc -l
      49
❯ cat rules.txt | cut -d "|" -f 1 | sort | uniq -c
  24 12
  24 13
  24 14
  24 15
  24 17
  24 18
  24 19
```

The 49 numbers in the set `{ x ∈ I | 12 <= x <= 99 ∧ x%10 != 0 }` each appear exactly 24 times as the left and right sides of all the rules. What’s incredible is that, despite the existence of a global cycle, this structure is sliced in a way that induces [a total ordering](https://en.wikipedia.org/wiki/Total_order). I got lucky this time—I didn't knew about the cycle but didn’t feel like diving into [DAG](https://en.wikipedia.org/wiki/Directed_acyclic_graph), [topological sorting](https://en.wikipedia.org/wiki/Topological_sorting) and stuff so early in AoC (and at breakfast). Instead, I grabbed a good cup of coffee, took a chance, and just sorted it, trusting the relationship a < b ⇔ a ∈ rules[b]. It was so satisfying to see it worked!

## My take on Go 1.23.3

~~The introduction of the iterators and the perimeter of `slices` are somewhat unsatisfactory.
I don't believe this will evolve positively.~~ `<EDIT>` I'm RTFMing.

## Day 6: Guard Gallivant

This problem has been the most demanding challenge so far. I managed a `20ms+` runtime previously, but now I'm seeing results like:

```bash
counts: 41 8 16.625µs
```

```bash
counts: 4883 1390 1.442291ms
```

It feels like I'm close—but not quite there yet!

`<EDIT>` i'm commiting the version i'm currently working on. It is not done yet (but what is done is blazing fast) and I don't know *for sure* if it can be done this way.

## Day 7: Bridge Repair

Today's solution is an elegant recursive, multi-branched [DFS](https://en.wikipedia.org/wiki/Depth-first_search). The key insight is to start from the target value and work *backward*, deconstructing it step by step. This approach naturally prunes certain branches—like divisions or concatenations—when they become impossible.

[![DFS search animated](https://img.youtube.com/vi/NUgMa5coCoE/0.jpg)](https://www.youtube.com/watch?v=NUgMa5coCoE)

## Day 8: Resonant Collinearity

Given the size of today's input, [brute-forcing](https://en.wikipedia.org/wiki/Brute-force_search) the solution did the trick.

`<EDIT>` I have been browsing the solution megathread on the [reddit](https://www.reddit.com/r/adventofcode/) and I have a tip for all the LLM-prompters out there. A [cartesian product](https://en.wikipedia.org/wiki/Cartesian_product) to generate all *non repeating* pairs of a set can be expressed as:

```C
for i, a := range set {
    for _, b := range set[i+1:] {
        blah(a, b)
    }
}
```

## Day 9: Disk Fragmenter

Today, I implemented a compact [filesystem](https://en.wikipedia.org/wiki/File_system) with a [File Allocation Table](https://en.wikipedia.org/wiki/File_Allocation_Table) (FAT) to support file handling primitives, tailored to meet the problem's requirements. [Simulations](https://en.wikipedia.org/wiki/Simulation) are great but surely require a lot of editing to bring their basic concepts to life.

```bash
❯ make lines
      57 ./3/aoc3.go
      72 ./1/aoc1.go
      91 ./8/aoc8.go
      98 ./5/aoc5.go
     101 ./2/aoc2.go
     104 ./7/aoc7.go
     167 ./4/aoc4.go
     234 ./9/aoc9.go
     308 ./6/aoc6.go
    1232 total
```

## Day 10: Hoof It

As on day 7, DFS proves to be the tool of choice for tackling the problem. Starting from every `0` in the grid, each traversal efficiently identifies all possible paths and their respective goals (`9`). While optimization tricks exist, the small input size means these tweaks wouldn't significantly impact performance.

Got lucky today! Initially, I mistakenly counted the paths instead of the reached goals. I quickly fixed it to focus on goals, only to find that Part 2 required the path count! Both numbers are natural outputs of DFS, so the adjustment was easy.

```bash
      --------Part 1--------   --------Part 2--------
Day       Time   Rank  Score       Time   Rank  Score
 10   00:26:19   4169      0   00:27:48   3435      0
 ```

Today is the *10th* day of the *10th* anniversary of [Advent Of Code](https://en.wikipedia.org/wiki/Advent_of_Code) and my calendar surely looks like it.

## Day 11: Plutonian Pebbles

In today’s solution, the `Counter` works as a frequency map for integers, making it easy to add and retrieve counts. It processes stones in batches, avoiding unnecessary loops. Last but not least, it [minimizes redundant operations by using direct map access](https://en.wikipedia.org/wiki/Memoization). As a result and for my input, the final and longest line consists of `3799` unique numbers, each appearing one or more times for a total of `223M+` stones.

I really appreciate the minimalism and simplicity of this solution—it is about [design](https://en.wikipedia.org/wiki/Dieter_Rams)!

![Dieter Ram's chair](https://upload.wikimedia.org/wikipedia/commons/0/00/1952_Dieter_Rams_Chair_from_Studienarbeit_drawing_%28resized%29.jpg)

## Day 12: Garden Groups

I nearly fried my brain today! The first part of the problem was a breeze—it was similar to Day 10, and adapting my solution didn’t take any effort. But then came Part 2, and it completely blindsided me. My brain froze. Eventually, I decided to take a long break and enjoy a nice breakfast. Only after that was I able to focus, isolate the solution for Part 2, and come up with a simple routine to handle the computation. I’m not finished with this problem yet!

```bash
      --------Part 1--------   --------Part 2--------
Day       Time   Rank  Score       Time   Rank  Score
 12   00:07:21    436      0   00:50:48   1506      0
```

PS. #436 is my personal best. Having detailed day10 was a huge payoff!

## Day 13: Claw Contraption

Today's problem has a straightforward mathematical solution, specifically solving [systems of two linear equations](https://en.wikipedia.org/wiki/System_of_linear_equations). The key and easy challenge is to avoid [integer overflows](https://en.wikipedia.org/wiki/Integer_overflow). My solution will work out of the box on 64-bit machines, and if needed, a [simple adaptation](https://go.dev/ref/spec#Numeric_types) for 32-bit systems is left for you to handle.

Today's parsing might seem tedious, but it's not! Here's my advice for events like AoC: embrace laziness. The more regular the input format, the easier it is to extract data with quick and efficient string routines instead of, say, regular expressions. Especially during AoC, I always trust the input and [never waste time validating](https://www.youtube.com/watch?v=_oNOTknRTSU) it.

It's Friday, and we've just passed the middle of December. I think today's problem is more of a fun break before the weekend. Stay tuned—I'm betting on some special weekend topics!

PS. So far so good! 3.74 of [cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) on average is ok:

```bash
14 main (Maze).run 6/aoc6.go:183:1
14 main solve 10/aoc10.go:43:1
13 main safe 2/aoc2.go:49:1
12 main decompose 12/aoc12.go:49:1
10 main antinodes 8/aoc8.go:59:1
10 main check 7/aoc7.go:47:1
9 main shape 12/aoc12.go:107:1
8 main scan 6/aoc6.go:78:1
8 main main 5/aoc5.go:26:1
7 main (RuneMat).findAll 4/aoc4.go:143:1
Average: 3.74
```

## Day 14: Restroom Redoubt

What a day! Today's challenge is easy—Part 1 just checks if we’re moving the robots correctly. But Part 2? That’s a whole different game. I had to generate and manually scan 10K+ images with my own eyes to figure it out. The good news? I nailed it and learned something interesting about the easter egg along the way. My advice: don’t peek at the solution unless you’re okay with spoiling the fun!

I have included a visualization routine that outputs a png file in the working directory.

![Advent Of Code 2024 - 10 years edition - day14 easter egg - a framed xmas tree in the snow](./images/aoc14-6516.png)

`<EDIT>` I've discovered what seems to be a very fast way to automatically find the Easter egg: simply subsample half of the points and calculate the [standard deviation](https://en.wikipedia.org/wiki/Standard_deviation) along the X and Y axes. There are two keyframes between [0–103] that will have significantly lower deviation on either axis. Once the keyframe times are identified, we can determine the first cycle's coincidence using the [Chinese Remainder Theorem](https://en.wikipedia.org/wiki/Chinese_remainder_theorem). which is kind of a [recurring joke](https://www.reddit.com/r/adventofcode/comments/1hgegw6/comment/m2izyvv/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button) in AoC.

## Day 15: Warehouse Woes

I tackled today’s challenge in place, using the maze itself as the data storage! On the plus side, this approach comes with built-in visualization. However, the trade-off lies in how the maze needs to be updated. Thankfully, deferring all moves and then filtering out unnecessary clears turned out to be just fine. The code reminds me a weird old but precise clockwork.

```bash
❯ make sample
go run ./aoc15.go < sample.txt
10092 9021
####################
##[].......[].[][]##
##[]...........[].##
##[]........[][][]##
##[]......[]....[]##
##..##......[]....##
##..[]............##
##.........[].[][]##
##......[][]..[]..##
####################
```

## Day 16: Reindeer Maze

~~I'm not completely satisfied with today's solution.~~ I often say there's not much to discuss when it comes to Dijkstra—it’s like handing the problem over to a computational mycelium, after all: it will eventually be decomposed but it takes time. That said, I ended up using [Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm) twice (RIP problem!):

- first, from the start node to the end node,
- and then from the end node to all other cells.

 By combining the distance matrices, I was able to reconstruct the best path scores for all traversed cells, selecting only the lowest values. Why go through all that effort? After all, it's straightforward to memorize the paths during the first run, right? Well, while that’s easy to implement, the core loop would end up allocating over 57,000 small chunks (~300 words each) and stall there as a result. The proposed approach shifts that allocation overhead into computing time, making it not only faster but also inherently parallel!

[![Dijkstra's algorithm animated](https://img.youtube.com/vi/EFg3u_E6eHU/0.jpg)](https://www.youtube.com/watch?v=EFg3u_E6eHU)

## Day 17: Chronospatial Computer

Building a computer is a recurring theme in AoC, and once we’ve built a bytecode machine, the next challenge often involves disassembling and reversing a program. Today, the goal was to find a configuration that turns the code into a [quine](https://en.wikipedia.org/wiki/Quine_(computing)). *What an idea!*

Two years ago, I did it painfully with pen and paper. Last year, I went all in with symbolic resolution and variable binding. This year? I kept it simple—I’m trying to save my energy for the challenges ahead.
First, I built a reliable, decently fast emulator. Then, I reversed the program. Finally, I went with a trial-and-error approach: browsing Reddit, checking out different solutions, and eventually settling on [this one in the code](https://www.reddit.com/r/adventofcode/comments/1hg38ah/comment/m2go5os/). The method uses insights from reversing the code to force digits one at a time and then packs the forcing values into a single initial value.

Who knows—maybe I’ll try dynamically forcing the value later on.

```bash
go run ./aoc17.go < input.txt
00 BST 4 A 66245665 00000001 00000000   out:
02 BXL 7 ? 66245665 00000006 00000000   out:
04 CDV 5 B 66245665 00000006 01035088   out:
06 BXL 7 ? 66245665 00000001 01035088   out:
08 BXC 6 C 66245665 01035089 01035088   out:
10 ADV 3 3 08280708 01035089 01035088   out:
12 OUT 5 B 08280708 01035089 01035088   out: 1
14 JNZ 0 0 08280708 01035089 01035088   out: 1
00 BST 4 A 08280708 00000004 01035088   out: 1
02 BXL 7 ? 08280708 00000003 01035088   out: 1
04 CDV 5 B 08280708 00000003 01035088   out: 1
06 BXL 7 ? 08280708 00000004 01035088   out: 1
08 BXC 6 C 08280708 01035092 01035088   out: 1
10 ADV 3 3 01035088 01035092 01035088   out: 1
```
