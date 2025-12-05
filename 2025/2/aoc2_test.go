package main

import (
	"testing"
)

func TestAllSplitSpans(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected [][2]int
	}{
		{
			name:     "single digit range",
			a:        1,
			b:        5,
			expected: [][2]int{{1, 5}},
		},
		{
			name:     "spans across 10 boundary",
			a:        5,
			b:        15,
			expected: [][2]int{{5, 9}, {10, 15}},
		},
		{
			name:     "spans across multiple boundaries",
			a:        50,
			b:        1500,
			expected: [][2]int{{50, 99}, {100, 999}, {1000, 1500}},
		},
		{
			name:     "exact boundary start",
			a:        100,
			b:        200,
			expected: [][2]int{{100, 200}},
		},
		{
			name:     "exact boundary end",
			a:        50,
			b:        100,
			expected: [][2]int{{50, 99}, {100, 100}},
		},
		{
			name:     "single value",
			a:        42,
			b:        42,
			expected: [][2]int{{42, 42}},
		},
		{
			name: "large range spanning many boundaries",
			a:    1,
			b:    1000000000,
			expected: [][2]int{
				{1, 9}, {10, 99}, {100, 999}, {1000, 9999}, {10000, 99999},
				{100000, 999999}, {1000000, 9999999}, {10000000, 99999999},
				{100000000, 999999999}, {1000000000, 1000000000},
			},
		},
		{
			name:     "no split needed - high values",
			a:        2000000000,
			b:        3000000000,
			expected: [][2]int{{2000000000, 3000000000}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result [][2]int
			for span := range allSplitSpans(tt.a, tt.b) {
				result = append(result, span)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d spans, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("span %d: expected %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

func TestSumMultiples(t *testing.T) {
	tests := []struct {
		name     string
		a, b, x  int
		expected int
	}{
		{
			name:     "no multiples in range",
			a:        1,
			b:        5,
			x:        10,
			expected: 0,
		},
		{
			name:     "single multiple",
			a:        10,
			b:        10,
			x:        10,
			expected: 10,
		},
		{
			name:     "multiple multiples of 3 in range 1-10",
			a:        1,
			b:        10,
			x:        3,
			expected: 18, // 3 + 6 + 9 = 18
		},
		{
			name:     "multiples of 5 in range 5-20",
			a:        5,
			b:        20,
			x:        5,
			expected: 50, // 5 + 10 + 15 + 20 = 50
		},
		{
			name:     "multiples of 2 in range 3-9",
			a:        3,
			b:        9,
			x:        2,
			expected: 18, // 4 + 6 + 8 = 18
		},
		{
			name:     "range starts after first multiple",
			a:        12,
			b:        25,
			x:        7,
			expected: 35, // 14 + 21 = 35
		},
		{
			name:     "range ends before next multiple",
			a:        8,
			b:        13,
			x:        5,
			expected: 10, // only 10
		},
		{
			name:     "large multiples",
			a:        100,
			b:        300,
			x:        101,
			expected: 303, // 101 + 202 = 303
		},
		{
			name:     "x equals 1",
			a:        5,
			b:        8,
			x:        1,
			expected: 26, // 5 + 6 + 7 + 8 = 26
		},
		{
			name:     "exact range boundaries are multiples",
			a:        15,
			b:        45,
			x:        15,
			expected: 90, // 15 + 30 + 45 = 90
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sumMultiples(tt.a, tt.b, tt.x)
			if result != tt.expected {
				t.Errorf("sumMultiples(%d, %d, %d) = %d, expected %d",
					tt.a, tt.b, tt.x, result, tt.expected)
			}
		})
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "lcm of 1 and any number",
			a:        1,
			b:        5,
			expected: 5,
		},
		{
			name:     "lcm of same numbers",
			a:        7,
			b:        7,
			expected: 7,
		},
		{
			name:     "lcm of coprime numbers",
			a:        3,
			b:        5,
			expected: 15,
		},
		{
			name:     "lcm of numbers with common factor",
			a:        6,
			b:        9,
			expected: 18,
		},
		{
			name:     "lcm where one divides the other",
			a:        4,
			b:        12,
			expected: 12,
		},
		{
			name:     "lcm of larger numbers",
			a:        15,
			b:        20,
			expected: 60,
		},
		{
			name:     "lcm with zero (edge case)",
			a:        0,
			b:        5,
			expected: 0,
		},
		{
			name:     "lcm of powers of 2",
			a:        8,
			b:        16,
			expected: 16,
		},
		{
			name:     "lcm of large primes",
			a:        97,
			b:        101,
			expected: 9797,
		},
		{
			name:     "lcm order independence",
			a:        12,
			b:        8,
			expected: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lcm(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("lcm(%d, %d) = %d, expected %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
