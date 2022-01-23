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
	jiao
}

type jiao interface {
	speak()
}

type Cat struct{}
type Dog struct{}

func (c Cat) speak() {
	fmt.Println("cat")
}

func (d Dog) speak() {
	fmt.Println("dog")
}

func main() {
	d := Animal{Cat{}}
	d.speak()

	d = Animal{Dog{}}
	d.speak()

}
