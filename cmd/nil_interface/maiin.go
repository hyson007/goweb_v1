package main

import (
	"fmt"
)

// type Animal interface {
// 	speak()
// }

// func jiao(a Animal) {
// 	a.speak()
// }

// type Cat struct {
// }

// func (c Cat) speak() {
// 	fmt.Println("cat")
// }

// type Dog struct {
// }

// func (d Dog) speak() {
// 	fmt.Println("Dog")
// }

// func main() {
// 	c := Cat{}
// 	d := Dog{}
// 	jiao(c)
// 	jiao(d)
// }

type Animal struct {
	speaker
}

type speaker interface {
	speak()
}

type Cat struct {
}

type Fenail struct {
	speaker
}

func (c Cat) speak() {
	fmt.Println("cat")
}

func (f Fenail) speak() {
	fmt.Print("Prefix: ....")
	f.speaker.speak()
}

func main() {
	d := Animal{&Fenail{Cat{}}}
	d.speak()

}
