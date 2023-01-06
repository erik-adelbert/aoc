## Timings

| day | time |
|-----|-----:|
| 15 | 1.1 |
| 2 | 1.1 |
| 3 | 1.1 |
| 1 | 1.2 |
| 10 | 1.2 |
| 5 | 1.2 |
| 6 | 1.2 |
| 7 | 1.2 |
| 4 | 1.3 |
| 12 | 1.4 |
| 22 | 1.4 |
| 8 | 1.6 |
| 13 | 1.7 |
| 21 | 1.7 |
| 14 | 2.0 |
| 25 | 2.1 |
| 17 | 2.3 |
| 9 | 3.2 |
| 18 | 5.0 |
| 19 | 5.6 |
| 11 | 6.0 |
| 23 | 13.5 |
| 16 | 16.3 |
| 24 | 58.0 |
| 20 | 118.8 |
| total | 251.2 |

end-to-end timing for part1&2 in ms - mbair M1/16GB - go1.19.4 darwin/arm64 - hyperfine 1.15.0

## Day 1
For this 2022 edition first day, I have written a simple and fast solution:
The logic is like a fragment of [insertion sort](https://en.wikipedia.org/wiki/Insertion_sort).
There's not much to say here, once again it's about pure composing speed. 
Nonetheless, it reminded me the 
[first challenge from last year](https://github.com/erik-adelbert/aoc/blob/main/2021/en_notes.md): 
I chose (again) to use 3 variables. I favored the `switch/case` form because the `if/else if/..`
lacked a final `else` clause (ie. balance, pedantic though).
Finally, I like the way max3, as a closure, captures global vars and gives the
resulting code a vintage-ish look and feel.

## Day 2
There is two efficient ways to solve today challenge: either the solution should precompute all
possible outcomes and match inputs against them or devise a scoring formula that is fast enough 
to be computed on the fly while parsing inputs.

But what would be the score anyway?

With R(ock), P(aper), S(cissors), let me consider the round where my O(pponent) plays R and I play P, 
as paper wins rock, my score is:

2 + 6 = 8 that is `m + r` with `m`, my move score and `r`, the outcome score

Move scores are:
| move | m |
|:----:|:-:|
|   R  | 1 |
|   P  | 2 |
|   S  | 3 |

With o, any O moves, i, any of my moves, D(raw), L(ost) and W(in), outcomes are [widely known](https://en.wikipedia.org/wiki/Rock_paper_scissors) to be:
|o\i| R | P | S |
|:-:|:-:|:-:|:-:|
| R | D | W | L |
| P | L | D | W |
| S | W | L | D |

When replacing `D=3`, `L=0`, `R=0`, `P=1`, `S=2` and `W=6`, I build the outcome scoring scale:
|o\i| 0 | 1 | 2 |
|:-:|:-:|:-:|:-:|
| 0 | 3 | 6 | 0 |
| 1 | 0 | 3 | 6 |
| 2 | 6 | 0 | 3 |

Now, when divided by `3`, it comes:
|o\i| 0 | 1 | 2 |
|:-:|:-:|:-:|:-:|
| 0 | 1 | 2 | 0 |
| 1 | 0 | 1 | 2 |
| 2 | 2 | 0 | 1 |

That is to say, I have:

`m = i + 1` with `0 <= i <= 2`

and with *an always positive* modulo:

`r = 3 * ((i-o+1) % 3)`, `0 <= o <= 2`

Finally, the *first part formula* is:

`(i + 1) + 3 * ((i-o+1) % 3)`

And it's easy enough to compute `i` and `o` values from input:

`o = ('A'|'B'|'C') - 'A'` and `i = ('X'|'Y'|'Z') - 'X'`

For part 2 and given the success of the previous approach, I'm
inclined to look for another matrix that summarizes the problem.
Here, given an opponent move, I have to follow a g(oal).
There's always a unique move to do that:

|o\g| L | D | W |
|:-:|:-:|:-:|:-:|
| R | P | R | S |
| P | S | P | R |
| S | R | S | P |

Applying the same substitutions as before, it comes:

|o\g| 0 | 1 | 2 |
|:-:|:-:|:-:|:-:|
| 0 | 1 | 0 | 2 |
| 1 | 2 | 1 | 0 |
| 2 | 0 | 2 | 1 |

It's not too hard to see that this matrix is symetric to the first one,
here's the *second part formula*:

`(((g+o-1) % 3) + 1) + 3*g` with `g = ('X'|'Y'|'Z') - 'X')`

Back to the practical design of my solution, I now have my answer:
I have found fast enough formulas to compute scores on the fly. But
why bother? I have also found a static scaling matrix for an even 
faster score computation!

## Day 3
Today my solution closely follows the challenge story and it's pretty boring 
but fast. First, the program splits any input line in two halves and maps 
each of them, in turn, with either `Head` or `Tail`, *two strictly positive
values*, in a *single map*. This amounts to build a set from the input line
while intersecting its head and tail.

ex: 

`ABCabc -> { A: Head, B: Head, C: Head, a: Tail, b: Tail, c: Tail }`

When updating the set with symbols from the tail, if there's already an 
item seen in the head, it is *counted for Part1*.
Once done, the map is memoized.
Moving forward, every 3 lines, the program scans the last 3 maps for a 
common item. Once found it is *counted for Part2*.

Using a sparse `[128]int` (no symbol is higher than `z = 122`), instead of, say,
a hashmap is trading space for speed. 
Only the 2 previous maps are buffered. 

If item `c` is common to all 3 input lines, then:

`mapbuf[0][c] * mapbuf[1][c] > 0` and `c` is in the current line

`count` closure is written as "prioritize no matter what for a lowercase symbol
and amend the priority if it is an uppercase one". This form was once known to 
ease branch prediction and here, I fancy the code layout.

As a final word and just for fun, I've tried my best to balance the naming.

## Day 4
My solution is exceptionnally boring but fast. This challenge plays against
Go fortes but the difficulty is really low.
It is about [1D geometry](https://en.wikipedia.org/wiki/One-dimensional_space).

There's a fact about 1D segments that translates elegantly in Go: *when one 
segment contains the other one, they also intersect*. It nicely becomes a 
`fallthrough` in the 2 branches `switch/case` of the main loop.

`<EDIT>` I've even removed the `fallthrough` because it simply means that part2 
score is part1 score + something!

## Day 5
Today I wasn't in the mood for parsing the challenge initial state: I hard-coded 
it into the program. From there, most of the pain was gone. Part1, is a classical
stack operations problem while Part2 means slice operations instead. Implementations 
of `atoi`, `last`, `push`, `pop`, `move` operations are unsafe and trust inputs.

Finally, to fast build the result string, I use a `bytes.Buffer` while scanning
crate stacks.

## Day 6
My solution grows a [sliding window](https://itnext.io/sliding-window-algorithm-technique-6001d5fbe8b3) over the input while scanning it. At any point, it ensures this window to be the largest possible with no repeating symbols. The underlying algorithm is:

    for each index, symbol in the input:
        if symbol is not in the current window:
            add it to the window
        else: 
            reset the window accordingly

        if window length is maximum:
            print index + 1
            stop

        update seen map with current symbol/index pair
        loop

This algorithm `runtime` `complexity` is `O(n)` with `n` the number of symbols. _It is the fastest possible_.

Now, I just have to define how a symbol is unique in a window, for this purpose I use a symbol indexed array (say a faster map) that records for every occurring symbol its last index in the input (ie. one place back scan). Let me dig in a little more:

_What is a locally unique symbol window-wised ?_

- if it has not been seen before or
- if it has been seen outside the current window

`<EDIT>` Fixed thanks to [@nicl271](https://www.reddit.com/r/adventofcode/comments/zdw0u6/comment/iz6sfv3/?utm_source=share&utm_medium=web2x&context=3)

_What if it is not unique?_

- to avoid repetition, the current window has to be shrunk to start just after the previous symbol index

In practice, only the current window length needs to be maintained.

_Did you notice (I didn't at first)?_

The proposed solution builds [_the longest substring without repeating characters_](https://leetcode.com/problems/longest-substring-without-repeating-characters/solutions/) that means it can *generally solve* the question.
Let me rephrase this idea: The _same code_ solves part 1&2. At runtime, if we _watch_ the window len and display the first index after a 4 non-repeating chars and later on the first one after 14 such chars, we're done! 

Last but not least the internal memory size is fixed, the solution also has `O(1)` `space` `complexity` `n-wise`:
*It is one of the best to solve the task at hand*.

`<EDIT>` I have an [ongoing discussion](https://www.reddit.com/r/adventofcode/comments/zdw0u6/comment/iz6e67e/?utm_source=share&utm_medium=web2x&context=3) about the space complexity that I may have not gotten right on this... ~more to come~!

## Day 7
Today, the challenge is an other kind of beast: the program has to 1) parse a shell dump, 2) rebuild a filesystem and
3) help decide what to do because of a storage shortage.

The shell session itself presents a [preorder traversal](https://en.wikipedia.org/wiki/Tree_traversal#Pre-order_implementation) of
the filesystem. Here is the today sample directory layout:

    ❯ grep cd sample.txt
    $ cd /
    $ cd a
    $ cd e
    $ cd ..
    $ cd ..
    $ cd d

As I said, preorder traversal of:
```
    /
    ├── a
    │   └── e
    └── d
```    
    

This part of the challenge is a classical question about the filesystem graph: I've looked for a recursive solution 
from the start because it's the easiest to develop and fix in this context.

There is not much room from improvement here and `part1` is pretty linear. But `part2` is a tad more interesting to design, the program has to memorize `subdir` sizes and fast scan them afterward. The trick here is to keep their array sorted as we build it as it will speed up search when scanning them later on. Fortunately, there's a single useful tool to help us for both building and scanning a sorted array: it's [binary search](https://en.wikipedia.org/wiki/Binary_search_algorithm). This algorithm is so handy that it lives in many standard libraries. Here, in `Go`, it is available as [`sort.SearchInts`](https://pkg.go.dev/sort#SearchInts).

Finally, some remarks about inputs: we don't care the `dir` or `ls` lines, they are noise here:

    ❯ grep -v dir sample.txt | grep -v ls
    $ cd /
    14848514 b.txt
    8504156 c.dat
    $ cd a
    29116 f
    2557 g
    $ cd e
    584 i
    $ cd ..
    $ cd ..
    $ cd d
    4060174 j
    8033020 d.log
    5626152 d.ext
    7214296 k

See? Without them it's even easier to answer today's questions!

~~Once again my solution runs bounded by `O(n)`: input lines are accessed only once.~~
The programs run in `O(n + log d)` with `n` the input lines count and `d` the subdirs
count. Thanks to the memoization for `part2`, the filesystem tree is never retraversed.

The programs runs in about `1ms` so I will leave it there for now but be warned:
there's a way to throw out everything except for the subdirs size calculation and
to get away with it. If I ever was to look for more total speed, I'll be coming back
for this.

`<EDIT>` I couldn't resist so I simplified the program. If you want to follow my notes, please pull the previous version.

## Day 8
My solution runs in ~3ms on my mbair M1/16GB but I'm not satisfied. The runtime complexity feels too high.

The program follows the naive approach. To speed things up, it precomputes the 4-axis field rotations matrices and scanning the 4-axis views becomes easy: it's a matter of slicing the 
right matrix at the right place and scanning this slice.

As trees are *counted* from *distances* and all distances are `chars` (offsetted by `'0'`), I really don't care bringing them back into integers: the `'0'` offset is auto-cancelled during computations. *The problem is a `pure` `byte` one*.

~~I will eventually rework this one to use a `monotonic` `stack` and I'm sure that will bring the complexity down to `n^2` instead of `n^3`. That is cutting the runtime by 1/3 in this case.~~

`<EDIT>` I've realised that `dist(o, v)` which counts the viewing distance from `o` needed to also output `h` the highest height. From there I was able to remove the call to `max(v)`. And the program runtime went down to ~1.6ms which I'm happy with.  

## Day 9
I find this one funny. The challenge teasing adds a lot of complications and by the end of reading it, all the necessary operations are split into too many small pieces. I have not find a better way than to follow the text: the program is very straightforward, it turns the challenge at hand into a vector problem.

`dir()` returns a translation vector that is used to move a tail knot closer to its head one. This simplifies the workflow. The other trick I've used is to compute `dist(A,B)^2` instead of `dist(A,B)` this also eases the flow of control.

## Day 10
Today's challenge reminds me of [last year day13](https://github.com/erik-adelbert/aoc/blob/2022/2021/13/aoc13.go), I have *carefully* implemented the description: offsets are everywhere and magic numbers such as `20`, `21`, `39`, `40` keep popping from seemingly nowhere.
In this case, I had to be torough writing every single case separately and later on I went on simplifying the code. As you can easily guess, there's a `+1` offset on `20` and a `-1` on `40`. `40` is obviously known from the challenge but `20` is interesting:

At first, I needed to relate `{20, 60, 100, 140, 180, 220, ...}` form `part1` with a closed formula, I went to [Sloane Sequence Encyclopedia](https://oeis.org) and discovered `f(x) = 40 x − 20`. This is where `20` was hidden. Reversing `f(x)` gave `(cycle+20)%40 == 0` as a starting point.

    10760
    ���� ���   ��  ���  �  � ����  ��  �  � 
    �    �  � �  � �  � �  � �    �  � �  � 
    ���  �  � �    �  � ���� ���  �    ���� 
    �    ���  � �� ���  �  � �    � �� �  � 
    �    �    �  � �    �  � �    �  � �  � 
    �    �     ��� �    �  � �     ��� �  �

I've updated my go toolchain from 1.17 to 1.19: it spared me 1ms. As of today my programs from day1-10 run all parts collectively under 15ms!

## Day 11
AoC 2022 is on! Today challenge describes an interesting dispatching or routing mechanism based on [modular arithmetic](https://en.wikipedia.org/wiki/Modular_arithmetic).

My solution is a straightforward implementation of the described design. 
*Monkeys* are modelised as parameterized states capable to self-update.
One of these parameters is an update function description: there's an embed minimal interpret that parse, tokenize and finally evaluates to an integer.

This works great for `part1` but coding it put me in a slow state of mind.
`part2` caught me off guard. First, I thought, well `math/big` could 
do the trick... well not really feasible given the `10_000` loop.

After a while, and a lot of circular thinking, I realised that it was about 
*modular arithmetic*: to keep the numbers checked while preserving speed, 
I needed a way to reduce the number in a way invisible to `updates` *and* 
`data routing`. Updates are not subject to number cutting side-effects: they 
are dumb single arithmetic operations. Routing is done modular-wise, the number 
that is invisible to all routing tests done in the network is the *least common 
multiple* of all modulos! 

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `cat input.txt` | 0.6 ± 0.2 | 0.2 | 2.1 | 1.00 |
| `cat input.txt \| ./aoc11` | 6.5 ± 0.3 | 6.0 | 13.0 | 10.52 ± 3.45 |

PS. ahah, program runtimes have been doubling for the last 3 days, `dijkstra` 
is probably coming on day 16! 

## Day 12
See? This is day 16 already! My last year [`day` `15`](https://github.com/erik-adelbert/aoc/blob/2022/2021/15/aoc15.go) comment still stands: you can't say must when using [`dijkstra`]((https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)).

~~The design encloses dijkstra in a loop to set the starting points. This allows the same code to solve `part1` given a singleton and `part2` a list of starting points.~~ 

There's much to say: we can solve this problem backward. By working out the solution from then end, there's a unique run that brings all the answers efficiently. 

<div style="text-align:center">
  <img src="https://upload.wikimedia.org/wikipedia/commons/2/23/Dijkstras_progress_animation.gif" />
</div>

## Day 13
Today's challenge describes some `packet` numbers. They are somewhat related to `snailfish` numbers from [last year day 18](https://github.com/erik-adelbert/aoc/blob/main/2021/18/aoc18.go). These numbers parsing is the crux of today's challenge, my parser is recursive and pretty straightforward. I had though times putting it together though and the tyniest mistake here is fatal. ~~The `cmp` function is nicely described: it was easy to have the standard Go library sorting the numbers.~~ No need to sort the packets!

I've chosen to represent a `packet` as `struct{ val int, list []packet }`, with the added convention that if `p.val < 0` then the number is a `list`. It's an `integer` otherwise. The following have been modified to display the internal representation of `[1, 1, 3, 1, 1]`:

    ❯ make sample
    cat sample.txt | go run ./aoc13.go
    {-1 [{-1 [{1 []} {1 []} {3 []} {1 []} {1 []}]}]}
    13
    140

## Day 14
There's so much to say about this challenge! I built my solution through many reworks and here it is running in less than 2ms!!
You'll see much more wizardry in this program than I intended at first. But the first iteration of this program was running in the 150ms realm, drowning in map accesses... Then I decided to translate and resize the world to fit it into a byte array: 50ms. 

Casually talking about this challenge with a friend, he was telling me about is plan to solve part2 by mapping hollows instead of walls.
Thinking of it, I realised that this world aisles would *always be the same* and *easy to compute*. I devised right away a slicing mechanism: an AABB defines the rectangle of interest where the action is. Having cut more than 50% of the world space, I was getting 
there: 15ms.

Finally, I added DFS backtracking to sand grains motion: I ended up assembling a nice dfs iterator.
All of it is about 200 LoC tailored to the problem but yet very generic.

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `cat input.txt` | 0.5 ± 0.2 | 0.2 | 1.4 | 1.00 |
| `cat input.txt \| ./aoc14` | 2.5 ± 0.2 | 2.0 | 4.1 | 4.61 ± 1.55 |

## Day 15
I have been working this challenge for almost 6 hours today. I've tried (number of) different ways and my solution
time was down to ~10ms and going. When things go sideways when coding, I always take a break and then work my way 
out using a paper and a pencil. It was when drawing corner cases for a friend that I realised it was possible to 
solve this problem from another prospective: it is possible (and [easy](https://www.reddit.com/r/adventofcode/comments/zmw9d8/comment/j0dhu1w/?utm_source=share&utm_medium=web2x&context=3)) to track the gap that could enclose the missing beacon.

And then I saw this [reddit post](https://www.reddit.com/r/adventofcode/comments/zmcn64/comment/j0cdi3j/?utm_source=share&utm_medium=web2x&context=3). The solution is so beautiful and balanced, that it would have been a waste of my time to finish mine (same idea anyway). Instead, I studied this one and adapted my work to become a port of it.

## Day 16
The dropout dilemma all over again, for now I'm not finished with this one! I did manage to get the stars but I'm not satisfied with the performance. ~~I'll figure it out later.~~ Finally! I've tried many techniques and settled for a floyd-warshall computing of all pairs shortest paths. Before a slightly modified `A*`. And then I saw the same ideas [here](https://github.com/orlp/aoc2022/blob/master/src/bin/day16.rs) published before mine. Consider my work to be a port of this `rust` solution.


## Day 17
This one is kind of fun! The program has to simulate a very bad [tetrish](https://en.wikipedia.org/wiki/Tetris) player that can't even rotate the pieces. `part1` is quiet easy to simulate and [`tetrominoes`](https://en.wikipedia.org/wiki/Tetromino) are [well known](https://gamedevelopment.tutsplus.com/tutorials/implementing-tetris-collision-detection--gamedev-852).

The main pitfall is in `part2`: we can't just simulate everything because the number of tetrominoes to drop `10^12` seems beyond comprehension. But the huge size of this number is also the key to this problem: we have `5` tetrominoes and `2k+` jets and they cycle.
so, _at least_, the all sequence seen up to `lcm(5, 2k+)` is repeating. We just have to simulate the play until we find a cycle. Then it's easy to statically fast forward all the cycles and to simulate the rest of the play until we have dropped the required number of pieces.

But how to detect a cycle? 

There are at least two good algorithms to solve the [general problem](https://en.wikipedia.org/wiki/Cycle_detection) but, here, they don't fit well... Let's try the naive approach for once: a cycle appears when we are about to drop the *same tetromino* with the *same jet* as before. Wait! Is that *all*? No it isn't, we also have to garantee that the new *tetromino* follows the same *path* as before. To this end we could *record* for each *tetromino*, the initial *jet* and the resulting *skyline*. And from there, we could naively (but efficiently) detect cycles!

The resulting computations are a little tricky but definitely manageable.

PS. It's funny to see `part2` computed faster than `part1` because it has fewer remaining moves.

## Day 18
Today's solution is pretty naive, `part1` scans every cube side `x` in the input: if any of the other sides is missing from input, then `x` is added to the area. `part2` is a classical [flood fill](https://en.wikipedia.org/wiki/Flood_fill#Stack-based_recursive_implementation_(four-way)) algorithm: It starts outside of all cubes and eventually moves toward an external side. From there, it surfs the surface while
keeping track of the outside area.

## Day 19
~~Just like Day 16: stars but no joy for now (TPSORT+).~~
Finally! I've managed to rework this challenge and the solution is surprisingly simple but hard to get right. The program runs a `DFS` search on possible moves with a cost heuristic to `cut` non-promising world states. Building a robot skips time forward sparing a lot of non interesting states in the process.

    Benchmark 1: cat input.txt
    Time (mean ± σ):       0.5 ms ±   0.2 ms    [User: 0.2 ms, System: 0.2 ms]
    Range (min … max):     0.2 ms …   2.0 ms    1160 runs
    
    Benchmark 2: cat input.txt | ./aoc19
    Time (mean ± σ):       6.1 ms ±   0.3 ms    [User: 4.8 ms, System: 1.9 ms]
    Range (min … max):     5.6 ms …   7.5 ms    1000 runs

## Day 20 
Today's about [cryptography](https://en.wikipedia.org/wiki/Key_(cryptography)) in a box!

The solution program runs an `array-based circular doubly linked list` with
two special ops: `shuffle` and `key`. Although it means writing some other basic ops from scratch, this choice prevents data to move around by updating their indices instead.

But today sample is not up to the task: an ill-designed code can pass the sample in various ways and obfuscate the crux of this challenge: offsets!! By the way, the one that can drive any programmer in endless circles comes from the imposed order of ops during `shuffle`:

1) If the program *removes* the current element
2) and then, finds it a new insertion point into the list
3) *therefore*, this point is lying between the *remaining elements*

By that time, the list contains only `n-1` items of the `n` from the sequence. Once found, there's another pitfall around this point: when going `forward`, `i` should be inserted `after` but it should go `before` when going `backward`. Fortunately, inserting *before* an item is the same as inserting *after* its predecessor.

`part1` & `part2` are solved the same way. Except for, `part2` is injected a fairly big `prime` [`salt`](https://en.wikipedia.org/wiki/Salt_(cryptography)) before being shuffled *ten* times. This is too scrambled for me. I can't see a faster way to solve today's problem other than handling the tedious `shuffling` task. It amounts for `90%` of the running time: `~159ms`.

Finally, the way I've coded this enabled me to use the central idea to [`Knuth's Algorithm X`](https://en.wikipedia.org/wiki/Dancing_Links#Main_ideas) (aka `DLX`):
the `cover/uncover` ops. On the funny side, my solution is akin to a *Step Dancing Subkeys*. As I've also recycled the `path halving` technique exposed by
[Sedgewick](https://en.wikipedia.org/wiki/Robert_Sedgewick_(computer_scientist)) during his [`Quick Union-Find`](https://sedgewick.io/wp-content/themes/sedgewick/slides/algorithms/Algs01-UnionFind.pdf) study, here comes the `Quick Step Dancing SubMonkeys` or more on point `Algorithm QSDSk`!

Yes, today I was inspired by [`Stanford`'s CS](https://web.stanford.edu/class/archive/cs/cs106b/cs106b.1126/lectures/17/Slides17.pdf)!

`<EDIT>` thanks to the reddit community and [u/azzal07](https://www.reddit.com/r/adventofcode/comments/1046aia/2022_all_daysgo_fast_solutions_291ms_total_runtime/) the runtime is down to `119ms` but there's an even better solution provided by u/CountableFiber, stay tuned!

## Day 21
Last year's day 23 was so painful to me: The challenge was about `compiler analysis`, I was unprepared. I tackled the challenge by writing it down and working out my solution with a pencil!
Then I saw [Russ Cox](https://www.youtube.com/watch?v=hmq6veCFo0Y) solve the challenge and learn
many things.

I was waiting for today's challenge to reclaim vengeance for last year! I have simplified the technique. First of all and without even giving it a thought I defined a `val` type that supports all supported instructions. Then I wrote an `eval()` for `val`. That was all for `part1`.

`part2` is about `computer algebra`: We have to solve `humn` value to make our program work. Although I'm aware that a [`bisection`](https://en.wikipedia.org/wiki/Binary_search_algorithm) would perfectly do the trick here, I decided to go for the algebra (vengeance!).

The idea, here, is truly basic, first all the symbolic instructions of the program are marked while descending its tree (ie. top to bottom). Then the solution solves for `humn` by *forcing* the values along this path to be equal to the *evaluated* values that don't depends on `humn` (the other side of the equality).

All in all, I could go for more speed by `bisecting` out the value. But for now, my solution runs in `1.7ms` (thanks to the small size of input) which I'm happy with!

## Day 22
~~First star but I'm unable to undertake `part2` because 1) I'm tired and 2) It will make my brain swells not in the good way... I'm taking a break and I will take care of it in due time.~~

Challenge is about [Cube Mapping](https://en.wikipedia.org/wiki/Cube_mapping), the main problem is to get the input cube right. For the rest, the current position is stored as a vector relative to `O` the origin of this worldmap and converted back and forth each time it enters/go out of a side.

Right now, I've got 3 days to finish reworking (`16`, `19` et `22.2`). I'm quite sure I won't beat last year `380ms` for all parts/all days: day20 is the culprit I guess. I'm still hoping to be under `500ms` but there's no real room to jiggle.

Let's see what's coming!

`<EDIT>` solution is ~~coming soon~~ here!

## Day 23
It's a multi-valued GoL, I've got the stars by writting a straight forward `python` script because my Go design for it will sureley takes a lot of time but runs really fast ~~(compared to naïve solutions, even the packed ones)~~. My solution is a packed simulation built upon a `u256` custom type derived from a `u128` custom type. It is akin to [u/SLiV9's](https://www.reddit.com/r/adventofcode/comments/zt6xz5/comment/j1f9cz2/?utm_source=share&utm_medium=web2x&context=3) but faster.

## Day 24
For this day challenge, the program `precomputes` all `wind conditions`: they `cycle` every `lcm(H, W)` with `H, W` the dimensions of this world. From there, it `floods` the resulting maze from `t0` and `start` to every `reachable cell` at a given time. The first time the `goal` cell is reached is garanteed to be the smallest possible (`dijkstra`-ish).

Whenever the goal is `not reachable` (there's no way to get through), the solution is to restart the flooding from `later` than `t0`.
## Day 25
For the last day of this AoC, the program defines a new number type `snafu` alongside the addition. The solution's core is a `digit adder` with `carry propagation` that can operate on `bytes`.

## What was it like?

This year, I wanted to compose fast and simple programs every day from the start. I failed for days `16`, `19` and `22.2` because, for me, they required thorough studies. I wasn't sure until the very last moment (day19 rework) I would be able to break my [last year record](https://www.reddit.com/r/adventofcode/comments/rzvsjq/2021_all_daysgo_fast_solutions_under_a_second/) (380ms) and I agree it makes little sense to try: challenges aren't even the same. Anyway, it felt like the right way of [`upping the ante`](https://www.reddit.com/r/adventofcode/comments/zaumkz/whats_up_with_upping_the_ante/) for me.

**Finally, here it is, this year collection runs all parts for all days in less than 291ms!!!**

I am so happy with this result! 
Feedback is welcome on reddit [u/erikade](https://www.reddit.com/user/erikade/).

Happy new year and Happy coding to you all!!


