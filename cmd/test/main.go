package main

import (
	"fmt"
	"io"
	"os"
)

func Contents(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close() // f.Close will run when we're finished.

	var result []byte
	buf := make([]byte, 10)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...) // append is discussed later.
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err // f will be closed if we return here.
		}
		fmt.Println("have read so far: ", n)
	}
	return string(result), nil // f will be closed if we return here.
}

func test() {
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}
}

func untrace(s string) { fmt.Println("leaving:", s) }

func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}

func un(s string) {
	fmt.Println("leaving:", s)
}

func a() {
	defer un(trace("a"))
	fmt.Println("in a")
}

func b() {
	// defer un(trace("b"))
	fmt.Println("in b")
	a()
}

func main() {
	type MyString string

	func (m MyString) String() string {
		return fmt.Sprintf("MyString=%s", m) // Error: will recur forever.
	}
}
