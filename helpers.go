package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"time"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// sha256 hashes are frequently used to compute short identities for binary or text blobs, TLS/SSL certificates.
func GenSHA256(peerName string) string {
	h := sha256.New()
	_, err := h.Write([]byte(peerName))
	CheckError(err)
	hash := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(hash))
}

func ValidatePort(port int) {
	if port < 1024 || port > math.MaxUint16 {
		panic(errors.New("Invalid port range."))
	}
}

func TrimWhitespaces(str string) string {
	const cutset = " \n\t\r\f\a\b\v"
	return strings.Trim(str, cutset)
}

// Server specific.
func Echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func OneOfMany(str string, many ...string) bool {
	for _, m := range many {
		if str == m {
			return true
		}
	}
	return false
}
