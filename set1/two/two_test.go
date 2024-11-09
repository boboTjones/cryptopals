package main

import (
	"encoding/hex"
	"testing"
)

func TestFixedXOR(t *testing.T) {
	// Test case: example from comments
	hexStr1 := "1c0111001f010100061a024b53535009181c"
	hexStr2 := "686974207468652062756c6c277320657965"
	expectedResult := "746865206b696420646f6e277420706c6179"

	// Decode hex strings
	strOne := deFunk(hexStr1)
	strTwo := deFunk(hexStr2)

	// XOR the two decoded byte slices
	result := make([]byte, len(strOne))
	for i := range strOne {
		result[i] = strOne[i] ^ strTwo[i]
	}

	// Encode the result to hex for comparison
	resultHex := hex.EncodeToString(result)

	// Verify the result matches the expected output
	if resultHex != expectedResult {
		t.Errorf("Expected %s but got %s", expectedResult, resultHex)
	}
}
