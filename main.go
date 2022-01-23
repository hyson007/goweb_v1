package main

import (
	"fmt"
	"goweb_v1/controllers"
	"goweb_v1/models"
	"net/http"

	"github.com/gorilla/mux"
)

// var homeTemplate *template.Template
// var contactTemplate *template.Template

// var homeView *views.View
// var contactView *views.View

// var signupView *views.View

// func home(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	// we need to change from Execute to ExecuteTemplate to indicate the name of template.
// 	// homeView.Layout in this case is set to "bootstrap", which matches the boostrap template name.
// 	must(homeView.Render(w, nil))
// }

// func contact(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(contactView.Render(w, nil))
// }

// func signup(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	must(signupView.Render(w, nil))
// }

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Frequently asked questions</h1><p>Here is a list of commmonly asked questions.</p>")
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "baloo"
	password = "junglebook"
	dbname   = "lenslocked"
)

func main() {
	// reason to initate err here is to ensure homeTemplate is using the global variable

	// homeView = views.NewView("bootstrap", "views/home.gohtml")
	// contactView = views.NewView("bootstrap", "views/contact.gohtml")

	// signupView = controllers.NewUsers().NewView
	psqlinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	us, err := models.NewUserService(psqlinfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()
	// us.ResetDB()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	//this one can be replaced with
	//r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	//note the two different ways to do this
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	http.ListenAndServe(":3000", r)

}
