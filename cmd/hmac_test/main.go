package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"goweb_v1/hash"
	"goweb_v1/models"
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

	// test the new user service
	fmt.Println("test the new user service...")
	const (
		host     = "localhost"
		port     = "5432"
		user     = "baloo"
		password = "junglebook"
		dbname   = "lenslocked"
	)

	psqlinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	us, err := models.NewUserService(psqlinfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()

	userobj := models.User{
		Name:     "Jon",
		Email:    "abc@abc.com",
		Password: "jon",
		Remember: "abc123",
	}

	if err := us.Create(&userobj); err != nil {
		panic(err)
	}
	fmt.Println("testing original users...")
	fmt.Printf("%+v\n", userobj)

	user2, err := us.ByRemember("abc123")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", user2)

	// testing original users...
	// the reason why first userobj password is set to empty is due to create function has override that!!!!

	// {Model:{ID:1 CreatedAt:2022-01-23 12:42:30.843174 +0800 +08 m=+0.090328443 UpdatedAt:2022-01-23 12:42:30.843174 +0800 +08 m=+0.090328443 DeletedAt:<nil>} Name:Jon Email:abc@abc.com Password: PasswordHash:$2a$10$sQV75kXTJyeE2L5zobEKL.uXsVVeriDPYhrL2Hud4FqKVqwOZyucO Remember:abc123 RememberHash:as2txooZwlx7QUKz3vmKEpEcrRDZfR_U9Ikev3WRhos=}

	// &{Model:{ID:1 CreatedAt:2022-01-23 04:42:30.843174 +0000 UTC UpdatedAt:2022-01-23 04:42:30.843174 +0000 UTC DeletedAt:<nil>} Name:Jon Email:abc@abc.com Password: PasswordHash:$2a$10$sQV75kXTJyeE2L5zobEKL.uXsVVeriDPYhrL2Hud4FqKVqwOZyucO Remember: RememberHash:as2txooZwlx7QUKz3vmKEpEcrRDZfR_U9Ikev3WRhos=}

}
