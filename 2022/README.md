| day | time (ms) |
|-----|-----:|
| 6 | 0.6 |
| 3 | 1.2 |
| 1 | 1.3 |
| 2 | 1.3 |
| 4 | 1.3 |
| 5 | 1.3 |
| 7 | 1.3 |
| total | 8.3 |

mbair M1/16GB - go1.17.5 darwin/arm64

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
values*, in a *single map*. This amounts to build a set from the input while
intersecting its head and tail.

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

_What is a unique symbol window-wised ?_

- if it has not been seen before, it is unique or
- if it has been seen outside the current window, it is locally unique

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

There is not much room from improvement here and `part1` is pretty linear. But `part2` is a tad more interesting to design, the program has to memorize `subdir` sizes and fast scan them afterward. The trick here is to keep their array sorted as we build it as it will speed up search when scanning them later on. Fortunately, there's a single usefull tool to help us for both building and scanning a sorted array: it's [binary search](https://en.wikipedia.org/wiki/Binary_search_algorithm). This algorithm is so handy that it lives in many standard libraries. Here, in `Go`, it is available as [`sort.SearchInts`](https://pkg.go.dev/sort#SearchInts).

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

Once again my solution runs bounded by `O(n)`: input lines are accessed only once.
Thanks to the memoization for `part2`, the filesystem tree is never retraversed.

The programs runs in about `1ms` so I will leave it there for now but be warned:
there's a way to throw out everything except for the subdirs size calculation and
to get away with it. If I ever was to look for more total speed, I'll be coming back
for this.

`<EDIT>` I couldn't resist so I simplified the program. If you want to follow my notes, please pull the previous version.