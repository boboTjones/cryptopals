package main

import (
	"bytes"
	"reflect"
	"sort"
	"testing"
)

func TestMapMaker(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantFreqs map[byte]float64
		wantTotal int
	}{
		{
			name:      "empty input",
			input:     []byte{},
			wantFreqs: map[byte]float64{},
			wantTotal: 0,
		},
		{
			name:      "single character",
			input:     []byte("A"),
			wantFreqs: map[byte]float64{'A': 1.0},
			wantTotal: 1,
		},
		{
			name:      "repeated character",
			input:     []byte("AAA"),
			wantFreqs: map[byte]float64{'A': 1.0},
			wantTotal: 3,
		},
		{
			name:  "multiple characters",
			input: []byte("HELLO"),
			wantFreqs: map[byte]float64{
				'H': 0.2,
				'E': 0.2,
				'L': 0.4,
				'O': 0.2,
			},
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFreqs, gotTotal := mapMaker(tt.input)

			if gotTotal != tt.wantTotal {
				t.Errorf("mapMaker() total = %v, want %v", gotTotal, tt.wantTotal)
			}

			if !reflect.DeepEqual(gotFreqs, tt.wantFreqs) {
				t.Errorf("mapMaker() freqs = %v, want %v", gotFreqs, tt.wantFreqs)
			}
		})
	}
}

func TestDeFunk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []byte
	}{
		{
			name:  "empty string",
			input: "",
			want:  []byte{},
		},
		{
			name:  "valid hex string",
			input: "4142",
			want:  []byte("AB"),
		},
		{
			name:  "known cipher text",
			input: "1b37373331363f78",
			want:  []byte{0x1b, 0x37, 0x37, 0x33, 0x31, 0x36, 0x3f, 0x78},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deFunk(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("deFunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeXor(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		key   byte
		want  []byte
	}{
		{
			name:  "empty input",
			input: []byte{},
			key:   'X',
			want:  []byte{},
		},
		{
			name:  "single byte",
			input: []byte{0x1b},
			key:   'X',
			want:  []byte{0x1b ^ 'X'},
		},
		{
			name:  "multiple bytes",
			input: []byte{0x1b, 0x37},
			key:   'X',
			want:  []byte{0x1b ^ 'X', 0x37 ^ 'X'},
		},
		{
			name:  "zero key",
			input: []byte("Hello"),
			key:   0,
			want:  []byte("Hello"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deXor(tt.input, tt.key)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("deXor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultsSorting(t *testing.T) {
	results := Results{
		&Result{key: 'A', score: 0.5, total: 34},
		&Result{key: 'B', score: 0.1, total: 34},
		&Result{key: 'C', score: 0.3, total: 34},
	}

	expected := Results{
		&Result{key: 'B', score: 0.1, total: 34},
		&Result{key: 'C', score: 0.3, total: 34},
		&Result{key: 'A', score: 0.5, total: 34},
	}

	sort.Sort(sortByScore{results})

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Sorting failed\ngot: %+v\nwant: %+v", results, expected)
	}
}

// This test verifies the end-to-end decryption process
func TestKnownDecryption(t *testing.T) {
	// Known test case where 'X' is the key
	cipherHex := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	key := byte('X')

	// Decode the hex
	cipherBytes := deFunk(cipherHex)

	// Decrypt with known key
	decrypted := deXor(cipherBytes, key)

	// Verify properties of the decrypted text
	if len(decrypted) != len(cipherBytes) {
		t.Errorf("Decrypted length = %d, want %d", len(decrypted), len(cipherBytes))
	}

	// Check if decrypted text contains only printable ASCII
	for _, b := range decrypted {
		if b < 32 || b > 126 {
			t.Errorf("Decrypted text contains non-printable character: %d", b)
		}
	}

	// Calculate frequency distribution
	freqs, total := mapMaker(decrypted)
	if total != len(cipherBytes) {
		t.Errorf("Frequency analysis total = %d, want %d", total, len(cipherBytes))
	}

	// Verify some expected frequency properties
	for char, freq := range freqs {
		if freq < 0 || freq > 1.0 {
			t.Errorf("Invalid frequency for char %c: %f", char, freq)
		}
	}
}
