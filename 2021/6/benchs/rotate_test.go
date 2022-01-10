package rotate

import (
	"math/rand"
	"testing"
)

var sample []int

func init() {
	sample = rand.Perm(1000)
}

func BenchmarkCopyRotate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CopyRotateLeft(sample)
	}
}

func BenchmarkDKRotate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DKRotateLeft(sample)
	}
}
