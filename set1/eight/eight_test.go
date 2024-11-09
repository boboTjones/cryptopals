package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestTwoStep(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: []string{},
		},
		{
			name:     "single byte",
			input:    []byte{0x41},
			expected: []string{"A"},
		},
		{
			name:     "two bytes",
			input:    []byte{0x41, 0x42},
			expected: []string{"AB"},
		},
		{
			name:     "three bytes",
			input:    []byte{0x41, 0x42, 0x43},
			expected: []string{"AB", "C"},
		},
		{
			name:     "four bytes",
			input:    []byte{0x41, 0x42, 0x43, 0x44},
			expected: []string{"AB", "CD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := twoStep(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("twoStep(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		size     int
		expected [][]string
	}{
		{
			name:     "empty input",
			input:    []string{},
			size:     2,
			expected: [][]string{},
		},
		{
			name:     "nil input",
			input:    nil,
			size:     2,
			expected: [][]string{},
		},
		{
			name:     "zero size",
			input:    []string{"A", "B"},
			size:     0,
			expected: [][]string{},
		},
		{
			name:     "negative size",
			input:    []string{"A", "B"},
			size:     -1,
			expected: [][]string{},
		},
		{
			name:     "single item",
			input:    []string{"A"},
			size:     2,
			expected: [][]string{{"A"}},
		},
		{
			name:     "exact chunks",
			input:    []string{"A", "B", "C", "D"},
			size:     2,
			expected: [][]string{{"A", "B"}, {"C", "D"}},
		},
		{
			name:     "partial last chunk",
			input:    []string{"A", "B", "C"},
			size:     2,
			expected: [][]string{{"A", "B"}, {"C"}},
		},
		{
			name:     "size larger than input",
			input:    []string{"A", "B"},
			size:     3,
			expected: [][]string{{"A", "B"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chunk(tt.input, tt.size)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("chunk(%v, %d) = %v, want %v",
					tt.input, tt.size, result, tt.expected)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		size     int
		expected int
	}{
		{
			name:     "empty input",
			input:    []string{},
			size:     2,
			expected: 0,
		},
		{
			name:     "no matches",
			input:    []string{"A", "B", "C", "D"},
			size:     2,
			expected: 0,
		},
		{
			name:     "all matches",
			input:    []string{"A", "B", "A", "B"},
			size:     2,
			expected: 2, // Both elements in the chunks match
		},
		{
			name:     "partial matches",
			input:    []string{"A", "B", "A", "C"},
			size:     2,
			expected: 1, // Only the first element matches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare(tt.input, tt.size)
			if result != tt.expected {
				t.Errorf("compare(%v, %d) = %d, want %d",
					tt.input, tt.size, result, tt.expected)
			}
		})
	}
}

func TestLinesSorting(t *testing.T) {
	lines := Lines{
		{Number: 1, Score: 10},
		{Number: 2, Score: 30},
		{Number: 3, Score: 20},
	}

	expected := Lines{
		{Number: 2, Score: 30},
		{Number: 3, Score: 20},
		{Number: 1, Score: 10},
	}

	sort.Sort(sortLinesByScore{lines})

	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("Sorting failed. Got %v, want %v", lines, expected)
	}
}
