package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// sha256 hashes are frequently used to compute short identities for binary or text blobs, TLS/SSL certificates.
func GenClientId(peerName string) string {
	h := sha256.New()
	_, err := h.Write([]byte(peerName))
	CheckError(err)
	hash := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(hash))
}
