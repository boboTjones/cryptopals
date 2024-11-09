package main

import (
	"encoding/hex"
	"testing"
)

// Controlled ECB and non-ECB test data
var controlledECBHex = "abababababababababababababababababababababababababababababababababab"    // Repeated "ab" pattern simulates ECB
var controlledNonECBHex = "abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefab" // Non-repeating pattern

// Helper function to count repeated 16-byte blocks in a byte slice
func countRepeats(data []byte, blockSize int) int {
	blockCount := make(map[string]int)
	for i := 0; i+blockSize <= len(data); i += blockSize {
		block := string(data[i : i+blockSize])
		blockCount[block]++
	}

	// Count blocks that are repeated
	repeatCount := 0
	for _, count := range blockCount {
		if count > 1 {
			repeatCount++
		}
	}
	return repeatCount
}

func TestDetectECB(t *testing.T) {
	// Decode the controlled ECB and non-ECB hex strings
	ecbData, _ := hex.DecodeString(controlledECBHex)
	nonEcbData, _ := hex.DecodeString(controlledNonECBHex)

	ecbRepeats := countRepeats(ecbData, 16)
	nonEcbRepeats := countRepeats(nonEcbData, 16)

	if ecbRepeats <= nonEcbRepeats {
		t.Errorf("Expected ECB data to have more repeated 16-byte blocks than non-ECB data")
	}
}
