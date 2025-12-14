# Writing AoC solutions in Go

![We need guns. Lots of guns.](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne//lots-of-guns.jpg)

I chose [Go](https://go.dev) mainly because **it appeals to me visually**, and because I enjoy communicating with 21st-century computers using it. As I see it, Go sits right in the middle of the path between me and the machine: it gives me the freedom of expression I crave when coding, while at the same time [encouraging](https://go.dev/doc/effective_go) very focused, usually simple, machine-oriented code. Finally, having been born from [C](https://en.wikipedia.org/wiki/C_(programming_language)) and [Newsqueak](https://swtch.com/~rsc/thread/newsqueak.pdf), it refines and improves upon techniques that have already proven their worth.

As soon as I discovered the [Commander](https://www.youtube.com/watch?v=yE5Tpp2BSGw)’s [work](http://www.r-5.org/files/books/computers/internals/unix/Francisco_Ballesteros-Notes_on_the_Plan_9_Kernel_Source-EN.pdf), I knew I was interested in the [same side aspects](https://theswissbay.ch/pdf/Gentoomen%20Library/Software%20Engineering/B.W.Kernighan%2C%20R.Pike%20-%20The%20Practice%20of%20Programming.pdf) of coding that he was exploring. And he and his [fellow coders](https://developers.googleblog.com/en/go-a-new-programming-language/), being who they were, brought their entire world to us—again!

Today’s [entry](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go) is proof that Go is built around [CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes), so [concurrency](https://go.dev/blog/waza-talk) is part of the language. It basically means that our programs can seamlessly transition between linear and concurrent execution. Here, if we choose to treat [`reader()`](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go#L58-L67) and [`dialer()`](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go#L70-L82) as black boxes, the code [spawns and pipelines](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go#L41-L47) as many [goroutines](https://go.dev/tour/concurrency/1) (stages) as there are commands in the input—each one [receiving the current dialer state](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go#L71-L71), being wired to a single command, computing the next state, and [passing it along](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go#L102-L104) to the next stage.

The ultimate stage is reached in `main()`, when the program [receives the results back from the pipeline and outputs them](https://github.com/erik-adelbert/aoc/blob/main/2025/1/RedditOne/aoc1_reddit_one.go#L50-L54). It is a [canonical implementation](https://go.dev/test/chan/sieve1.go) of a (inefficient but illustrative) Go pipeline.

    ┌────────────┐
    │    main    │
    │            │
    │         ┌────────┐
    │         │ reader │──┐
    │         └────────┘  │
    │            │        │
    │         ┌─────────┐ │
    │  init → │ dialer1 │-|
    │         └─────────┘ │
    │            │  ↓     │
    │         ┌─────────┐ │
    │         │ dialer2 │-|
    │         └─────────┘ │
    │            │  ↓     │
    .            .  .     .
    .            .  .     .
    │            │  ↓     │
    │         ┌─────────┐ │
    │         │ dialerx │─┘
    │         └─────────┘
    │            │  │
    │  results ←────┘
    │            │
    └────────────┘

It does everything in 71 lines of standard Go and runs in under 2.3 ms.

Concurrency aside, Go is also on par with Rust (and similar languages) in terms of its purely numerical capabilities, flexibility, efficiency, and performance—and this is usually the kind of support I’m looking for when doing AoC.

This is my very first submission to the AoC Community Fun contest, and no matter what, I’m glad I can present it to you all. Many virtual hugs, and many thanks, to the AoC team and this community!

Happy coding!

**P.S.** As a matter of fact, this submission has an accompanying [soundtrack](https://www.youtube.com/watch?v=LxRgU1NoDCI).
