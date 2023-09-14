package main

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestParseCmd(t *testing.T) {
	tests := []struct {
		name         string
		strCmd       string
		expectedCmd  string
		expectedArgs []string
	}{
		{
			name:         "No arguments",
			strCmd:       "command",
			expectedCmd:  "command",
			expectedArgs: []string{},
		},
		{
			name:         "Single argument",
			strCmd:       "command -- argument",
			expectedCmd:  "command",
			expectedArgs: []string{"argument"},
		},
		{
			name:         "Multiple arguments",
			strCmd:       "command -- argument1 argument2 argument3",
			expectedCmd:  "command",
			expectedArgs: []string{"argument1", "argument2", "argument3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := parseCmd(tt.strCmd)
			if cmd != tt.expectedCmd {
				t.Errorf("got %s, want %s", cmd, tt.expectedCmd)
			}
			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("got %v, want %v", args, tt.expectedArgs)
			}
		})
	}
}
func TestValidateChallenge(t *testing.T) {
	challenge := []byte("challenge")
	reply := []byte("reply")
	key := []byte("key")

	// Test case: valid reply
	validReply := hmacSha256(challenge, key)
	if !validateChallenge(challenge, validReply, key) {
		t.Errorf("validateChallenge() returned false for valid reply")
	}

	// Test case: invalid reply
	invalidReply := []byte("invalid")
	if validateChallenge(challenge, invalidReply, key) {
		t.Errorf("validateChallenge() returned true for invalid reply")
	}

	// Test case: empty challenge
	emptyChallenge := []byte("")
	if validateChallenge(emptyChallenge, reply, key) {
		t.Errorf("validateChallenge() returned true for empty challenge")
	}

	// Test case: empty reply
	emptyReply := []byte("")
	if validateChallenge(challenge, emptyReply, key) {
		t.Errorf("validateChallenge() returned true for empty reply")
	}

	// Test case: empty key
	emptyKey := []byte("")
	if validateChallenge(challenge, reply, emptyKey) {
		t.Errorf("validateChallenge() returned true for empty key")
	}
}

func TestRandomBytes(t *testing.T) {
	// Test case: size is 0
	t.Run("SizeZero", func(t *testing.T) {
		result := randomBytes(0)
		if len(result) != 0 {
			t.Errorf("Expected length of result to be 0, but got %d", len(result))
		}
	})

	// Test case: size is positive
	t.Run("SizePositive", func(t *testing.T) {
		size := 10
		result := randomBytes(size)
		if len(result) != size {
			t.Errorf("Expected length of result to be %d, but got %d", size, len(result))
		}
	})
}

func TestHmacSha256(t *testing.T) {
	data := []byte("hello")
	key := []byte("secret")

	s := "88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b"
	expected, _ := hex.DecodeString(s)

	// Positive test case
	result := hmacSha256(data, key)
	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %x, but got %x", expected, result)
	}

	// Negative test case
	wrongKey := []byte("wrong")
	result = hmacSha256(data, wrongKey)
	if bytes.Equal(result, expected) {
		t.Errorf("Expected different result, but got %x", result)
	}

	// Additional test cases can be added to cover different scenarios
}
