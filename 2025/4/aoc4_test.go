package main

import (
	"slices"
	"testing"
)

// global variable to prevent compiler optimizations
var sink []byte

// BenchmarkCarelessAllocations demonstrates the performance impact of allocating
// in tight loops - this is what NOT to do in performance-critical code
func BenchmarkCarelessAllocations(b *testing.B) {
	for b.Loop() {
		// BAD: allocating inside the loop
		buf := make([]byte, 1024)

		// Fill buffer with 1s
		for i := range 1024 {
			buf[i] = 1
		}

		// Prevent optimization by assigning to global
		sink = buf
	}
}

// BenchmarkPreallocated shows the correct approach with pre-allocation
func BenchmarkPreallocated(b *testing.B) {
	// GOOD: allocate once outside the loop
	buf := make([]byte, 1024)

	for b.Loop() {
		// Reset and fill buffer with 1s
		for i := range 1024 {
			buf[i] = 1
		}

		// Prevent optimization by assigning to global
		sink = buf
	}
}

// BenchmarkCarelessAllocationsWithCopy demonstrates another bad pattern:
// allocating and copying in tight loops
func BenchmarkCarelessAllocationsWithCopy(b *testing.B) {
	// GOOD: allocate and fill once outside the loop
	src := slices.Repeat([]byte{1}, 1024)

	for i := range 1024 {
		src[i] = 1
	}

	for b.Loop() {
		// BAD: allocating and copying inside the loop
		buf := make([]byte, 1024)
		copy(buf, src)

		sink = buf // prevent optimization
	}
}

// BenchmarkPreallocatedWithCopy shows the efficient approach for copying
func BenchmarkPreallocatedWithCopy(b *testing.B) {
	src := slices.Repeat([]byte{1}, 1024)
	buf := make([]byte, 1024) // GOOD: allocate once

	for b.Loop() {
		// GOOD: just copy to pre-allocated buffer
		copy(buf, src)

		// Prevent optimization
		sink = buf
	}
}

// BenchmarkWorstCase demonstrates multiple bad practices combined
func BenchmarkWorstCase(b *testing.B) {
	for b.Loop() {
		// TERRIBLE: multiple allocations and unnecessary operations
		tmp := make([]byte, 512)
		buf := make([]byte, 1024)

		// Fill temp with 1s
		for i := range tmp {
			tmp[i] = 1
		}

		// Copy temp twice to fill buffer
		copy(buf[:512], tmp)
		copy(buf[512:], tmp)

		// More unnecessary allocation
		res := make([]byte, len(buf))
		copy(res, buf)

		sink = res // prevent optimization
	}
}

// BenchmarkOptimized shows how the same work can be done efficiently
func BenchmarkOptimized(b *testing.B) {
	// EXCELLENT: single allocation, reused across all iterations
	buf := make([]byte, 1024)

	for b.Loop() {
		// GOOD: direct filling without intermediate allocations
		for i := range 1024 {
			buf[i] = 1
		}

		sink = buf
	}
}

// BenchmarkRealWorldBad simulates a realistic bad pattern from parsing/processing
func BenchmarkRealWorldBad(b *testing.B) {
	for b.Loop() {
		// BAD: creating multiple temporary slices
		lines := make([][]byte, 0, 10)

		for range 10 {
			// BAD: allocating for each "line"
			line := make([]byte, 100)
			for i := range 100 {
				line[i] = byte(i % 256)
			}
			lines = append(lines, line)
		}

		// BAD: another allocation to flatten
		res := make([]byte, 0, 1000)
		for _, line := range lines {
			res = append(res, line...)
		}

		sink = res
	}
}

// BenchmarkRealWorldGood shows the efficient approach for the same work
func BenchmarkRealWorldGood(b *testing.B) {
	// GOOD: pre-allocate everything
	buf := make([]byte, 1000)

	for b.Loop() {
		// GOOD: direct writing to final buffer
		i := 0
		for range 10 {
			for j := range 100 {
				buf[i] = byte(j % 256)
				i++
			}
		}

		sink = buf // prevent optimization
	}
}

/*
BENCHMARK ANALYSIS:

The results clearly demonstrate the performance cost of careless allocations:

1. Basic Allocation Impact:
   - BenchmarkCarelessAllocations:  ~419 ns/op, 1024 B/op, 1 allocs/op
   - BenchmarkPreallocated:         ~333 ns/op, 0 B/op, 0 allocs/op
   → 25% performance improvement with pre-allocation

2. Copy Operations:
   - BenchmarkCarelessAllocationsWithCopy: ~101 ns/op, 1024 B/op, 1 allocs/op
   - BenchmarkPreallocatedWithCopy:        ~14 ns/op, 0 B/op, 0 allocs/op
   → 7x performance improvement with pre-allocation!

3. Real-World Scenario:
   - BenchmarkRealWorldBad:  ~712 ns/op, 2144 B/op, 11 allocs/op
   - BenchmarkRealWorldGood: ~454 ns/op, 0 B/op, 0 allocs/op
   → 57% performance improvement, eliminates all allocations

KEY TAKEAWAYS:
- Memory allocations have significant overhead beyond just the memory cost
- The garbage collector impact compounds with more allocations
- Pre-allocation can provide 25-700% performance improvements
- In tight loops, allocation overhead can dominate execution time
- This is why aoc4.go's pre-allocation strategy is so effective

This demonstrates why the aoc4.go solution achieves 1.6ms runtime -
zero allocations in the hot path is crucial for competitive programming performance.
*/
