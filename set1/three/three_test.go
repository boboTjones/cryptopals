package main

import (
	"encoding/hex"
	"testing"
)

func TestSingleCharacterXORDecryption(t *testing.T) {
	hexStr := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	expectedKey := byte('X')                                         // Assume 'X' is the correct key based on analysis.
	expectedDecryptedMessage := "Cooking MC's like a pound of bacon" // Replace with expected plaintext if known.

	// Decode the hex string
	cipherText, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}

	// XOR with the expected key to see if it produces the expected message
	decryptedText := deXor(cipherText, expectedKey)

	// Convert decrypted text to string for comparison
	decryptedMessage := string(decryptedText)

	// Check if the result matches the expected output
	if decryptedMessage != expectedDecryptedMessage {
		t.Errorf("Expected decrypted message %q but got %q", expectedDecryptedMessage, decryptedMessage)
	}
}
