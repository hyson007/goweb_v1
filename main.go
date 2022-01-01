package main

import (
	"fmt"
	"goweb_v1/views"
	"net/http"

	"github.com/gorilla/mux"
)

// var homeTemplate *template.Template
// var contactTemplate *template.Template

var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// we need to change from Execute to ExecuteTemplate to indicate the name of template.
	// homeView.Layout in this case is set to "bootstrap", which matches the boostrap template name.
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Frequently asked questions</h1><p>Here is a list of commmonly asked questions.</p>")
}

func main() {
	// reason to initate err here is to ensure homeTemplate is using the global variable

	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	http.ListenAndServe(":3000", r)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
