package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>welcome to my awesome site </h1>")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	// router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	// http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":3000", router)

}
