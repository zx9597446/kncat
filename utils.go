package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
)

func randomBytes(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func hmacSha256(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func validateChallenge(challenge, reply, key []byte) bool {
	return hmac.Equal(hmacSha256(challenge, key), reply)
}
