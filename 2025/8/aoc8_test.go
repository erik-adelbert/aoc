package main

import (
	"bytes"
	"testing"
)

var sink []byte

func BenchmarkCut(b *testing.B) {
	s := []byte("1234567890-1234567890")
	var x, y []byte

	for b.Loop() {
		x, y, _ = bytes.Cut(s, []byte("-"))
	}

	sink = x // prevent optimization
	sink = y
}

func BenchmarkSplitN(b *testing.B) {
	s := []byte("1234567890-1234567890")

	var x [][]byte

	for b.Loop() {
		x = bytes.SplitN(s, []byte("-"), 2)
	}

	sink = x[0] // prevent optimization
}
