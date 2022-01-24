package rand

import (
	"encoding/base64"
	"math/rand"
)

const RememberTokenByte = 32

// there is another math/rand

// Bytes will help us generate n random bytes or will return error
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, err
}

// String will generate a byte slice of size nByte and return
// a string thats the base64 URL encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {

		// note we can't return nil here, as the func expect to return a string
		// if the func return []byte, like the previous one, we can return nil
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RemeberToken is a helper func that design to generate remember token at
// a predefined byte size
func RememberToken() (string, error) {
	return String(RememberTokenByte)
}

//NBytes return the number of bytes used in the base64 URL encoded string.
func NBytes(base64string string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64string)
	if err != nil {
		return -1, err
	}
	return len(b), nil

}
