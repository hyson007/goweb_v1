package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
}

func main() {

	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, User{"jack"})
	if err != nil {
		panic(err)
	}
}
