package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numSamples = 10000000

// Hash function provided
func hash(a int) int {
	a ^= (a << 6) & 0xFFFFFF
	a ^= (a >> 5) & 0xFFFFFF
	return a ^ (a<<11)&0xFFFFFF
}

func measureSpread() {
	hashes := make(map[int]int)
	n := rand.Intn(1000000)
	for i := 0; i < 2000; i++ {
		n = hash(n)
		hashes[n]++
	}

	n = rand.Intn(1000000)
	for i := 0; i < 2000; i++ {
		n = hash(n)
		hashes[n]++
	}

	n = rand.Intn(1000000)
	for i := 0; i < 2000; i++ {
		n = hash(n)
		hashes[n]++
	}

	// Print distribution statistics
	totalHashes := 0
	for _, count := range hashes {
		totalHashes += count
	}

	fmt.Println("Total loops:", 2000)
	fmt.Println("Distinct hash values:", len(hashes))
}

// Measure the uniformity and collision resistance of the hash function
func measureUniformity() {
	hashes := make(map[int]int)
	for i := 0; i < numSamples; i++ {
		// Generate random input between 0 and 1,000,000
		input := rand.Intn(1000000)
		h := hash(input)

		// Count the occurrence of each hash value
		hashes[h]++
	}

	// Print distribution statistics
	totalHashes := 0
	for _, count := range hashes {
		totalHashes += count
	}

	fmt.Println("Total samples:", numSamples)
	fmt.Println("Distinct hash values:", len(hashes))

	// Calculate mean and variance
	var sum, sumOfSquares int
	for _, count := range hashes {
		sum += count
		sumOfSquares += count * count
	}
	mean := float64(sum) / float64(len(hashes))
	variance := float64(sumOfSquares)/float64(len(hashes)) - mean*mean
	fmt.Printf("Mean: %f, Variance: %f\n", mean, variance)
}

// Measure the avalanche effect of the hash function (how much it changes when a bit is flipped)
func measureAvalanche(input int) int {
	originalHash := hash(input)
	changedInput := input ^ (1 << rand.Intn(32)) // Flip a random bit
	changedHash := hash(changedInput)

	// XOR the original and changed hashes, then count the number of differing bits
	difference := originalHash ^ changedHash
	bitCount := 0
	for difference != 0 {
		bitCount++
		difference &= difference - 1 // Count set bits
	}
	return bitCount
}

// Measure how well the range of hash outputs is covered
func measureRangeCoverage() {
	hashes := make(map[int]struct{})
	for i := 0; i < numSamples; i++ {
		input := rand.Intn(1000000)
		hashVal := hash(input)
		hashes[hashVal] = struct{}{}
	}

	coverage := float64(len(hashes)) / float64(0xFFFFFF)
	fmt.Printf("Coverage: %.6f%%\n", coverage*100)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Measure uniformity and collision resistance
	fmt.Println("Uniformity & Collision Resistance:")
	measureUniformity()
	measureSpread()

	// Measure avalanche effect
	fmt.Println("\nAvalanche Effect:")
	input := rand.Intn(1000000) // Random input
	avalancheEffect := measureAvalanche(input)
	fmt.Printf("Avalanche effect (bit difference): %d\n", avalancheEffect)

	// Measure range coverage
	fmt.Println("\nRange Coverage:")
	measureRangeCoverage()
}
