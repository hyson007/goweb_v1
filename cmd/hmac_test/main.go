package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"goweb_v1/hash"
)

func main() {
	// p := make([]byte, 100)
	toHash := []byte(" this is my strings to hash")
	h := hmac.New(sha256.New, []byte("my-secret-key"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(b)
	fmt.Println(string(b))

	//reset the data
	// each time we can print the same output
	h.Reset()
	h.Write(toHash)
	b = h.Sum(nil)
	fmt.Println(b)
	fmt.Println(string(b))

	//if we change the secret key, then the output is different
	h = hmac.New(sha256.New, []byte("my-secret-key111"))
	h.Reset()
	h.Write(toHash)
	b = h.Sum(nil)
	fmt.Println(b)
	fmt.Println(string(b))

	//calling the customized hmac wrapper
	toHash = []byte(" this is my strings to hash")
	h = hmac.New(sha256.New, []byte("my-secret-key"))
	h.Write(toHash)
	b = h.Sum(nil)
	fmt.Println("comparing the customized func")
	fmt.Println("original output")
	fmt.Println(base64.URLEncoding.EncodeToString(b))

	fmt.Println("customized func output")
	hmac := hash.NewHMAC("my-secret-key")
	fmt.Println(hmac.Hash(" this is my strings to hash"))

}
