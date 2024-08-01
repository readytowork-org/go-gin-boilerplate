package utils

import (
	"crypto/rand"
	"io"
	"strconv"
	"time"
)

var digitsAndNumbers = [...]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '|', '}', '{', '[', ']'}

// GenerateRandomCode generate random string
func GenerateRandomCode(codeLength int) string {
	b := make([]byte, codeLength)
	n, err := io.ReadAtLeast(rand.Reader, b, codeLength)
	if n != codeLength {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = digitsAndNumbers[int(b[i])%len(digitsAndNumbers)]
	}
	return string(b)
}

// GenerateRandomFileName genrates the fileName with unique time
func GenerateRandomFileName() string {
	_time := time.Now().UnixNano()
	return strconv.FormatInt(_time, 10)
}

// GenerateRandomDigitSequence generates random digit sequence
func GenerateRandomDigitSequence(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
