package main

import (
	"fmt"
	"goweb_v1/rand"
)

func main() {
	token, _ := rand.RememberToken()
	fmt.Println(token)
}
