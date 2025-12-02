package main

import (
	"bytes"
	"testing"
)

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{10, "10"},
		{99, "99"},
		{100, "100"},
		{123, "123"},
		{1000, "1000"},
		{12345, "12345"},
		{9876543210, "9876543210"},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.input)), func(t *testing.T) {
			result := itoa(tt.input)
			if !bytes.Equal(result, []byte(tt.expected)) {
				t.Errorf("itoa(%d) = %s, want %s", tt.input, string(result), tt.expected)
			}
		})
	}
}

func BenchmarkItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		itoa(9876543210)
	}
}
