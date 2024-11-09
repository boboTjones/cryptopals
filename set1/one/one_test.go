package main

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

// hexToBase64 converts a hexadecimal string to a base64 encoded string
func hexToBase64(hexStr string) (string, error) {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(decoded), nil
}

func TestHexToBase64(t *testing.T) {
	// Test case from the comments in one.go
	hexStr := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	expectedBase64 := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	result, err := hexToBase64(hexStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != expectedBase64 {
		t.Errorf("Expected %s but got %s", expectedBase64, result)
	}
}
