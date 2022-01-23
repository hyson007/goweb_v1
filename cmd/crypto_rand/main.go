package main

import (
	"fmt"
	"goweb_v1/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
