package main

import (
	"bytes"
	"reflect"
	"sort"
	"testing"
)

func TestDeFunk(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "simple hex string",
			input:    "4142",
			expected: []byte("AB"),
		},
		{
			name:     "complex hex string",
			input:    "48656c6c6f", // "Hello" in hex
			expected: []byte("Hello"),
		},
		{
			name:     "all zeros",
			input:    "0000",
			expected: []byte{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deFunk(tt.input)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("deFunk(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDeXor(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		key      byte
		expected []byte
	}{
		{
			name:     "empty input",
			input:    []byte{},
			key:      'A',
			expected: []byte{},
		},
		{
			name:     "single byte",
			input:    []byte{0x41}, // 'A'
			key:      0x01,
			expected: []byte{0x40},
		},
		{
			name:  "multiple bytes",
			input: []byte("ABC"),
			key:   'X',
			expected: []byte{
				'A' ^ 'X',
				'B' ^ 'X',
				'C' ^ 'X',
			},
		},
		{
			name:     "zero key",
			input:    []byte("Hello"),
			key:      0,
			expected: []byte("Hello"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deXor(tt.input, tt.key)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("deXor(%v, %v) = %v, want %v",
					tt.input, tt.key, result, tt.expected)
			}
		})
	}
}

func TestResultsSorting(t *testing.T) {
	results := Results{
		&Result{[]byte("test1"), 1, 'A', []byte("dec1"), 50},
		&Result{[]byte("test2"), 2, 'B', []byte("dec2"), 10},
		&Result{[]byte("test3"), 3, 'C', []byte("dec3"), 30},
	}

	expected := Results{
		&Result{[]byte("test2"), 2, 'B', []byte("dec2"), 10},
		&Result{[]byte("test3"), 3, 'C', []byte("dec3"), 30},
		&Result{[]byte("test1"), 1, 'A', []byte("dec1"), 50},
	}

	sort.Sort(sortByTotal{results})

	for i := range results {
		if !reflect.DeepEqual(results[i], expected[i]) {
			t.Errorf("Sorting failed at position %d\ngot: %+v\nwant: %+v",
				i, results[i], expected[i])
		}
	}
}

// Helper function to test a full decryption scenario
func TestDecryptionScenario(t *testing.T) {
	// Known plaintext encrypted with key 'X'
	input := "48656c6c6f" // "Hello" in hex
	key := byte('X')
	expectedPlain := []byte("Hello")

	// Test the full decryption process
	decoded := deFunk(input)
	if !bytes.Equal(decoded, expectedPlain) {
		t.Errorf("deFunk failed: got %v, want %v", decoded, expectedPlain)
	}

	// XOR with key
	encrypted := deXor(decoded, key)

	// XOR again to decrypt
	decrypted := deXor(encrypted, key)

	// Should get original plaintext back
	if !bytes.Equal(decrypted, expectedPlain) {
		t.Errorf("Full decryption failed: got %v, want %v",
			decrypted, expectedPlain)
	}
}
