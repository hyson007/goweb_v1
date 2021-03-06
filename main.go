package main

import (
	"fmt"
	"goweb_v1/controllers"
	"goweb_v1/middleware"
	"goweb_v1/models"
	"goweb_v1/rand"
	"net/http"

	"github.com/gorilla/csrf"
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

	// us, err := models.NewUserService(psqlinfo)
	// if err != nil {
	// 	panic(err)
	// }

	svc, err := models.NewService(psqlinfo)
	if err != nil {
		panic(err)
	}

	defer svc.Close()
	svc.AutoMigrate()

	//csrf token, which works like middleware but it's a function
	//csrfMW := csrf.Protect([]byte("32-byte-long-auth-key"))
	isProd := false
	randByte, _ := rand.Bytes(32)
	csrfMW := csrf.Protect(randByte, csrf.Secure(isProd))

	// for the name route to work, we have to declear the gallery controller
	// first before then we define mux new router
	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(svc.User)

	userMw := middleware.User{
		UserService: svc.User,
	}
	requireUseMw := middleware.RequireUser{
		User: userMw,
	}

	// pass in the r into gallery
	galleryC := controllers.NewGallery(svc.Gallery, svc.Image, r)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	//this one can be replaced with
	//r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	//note the two different ways to do this
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/logout",
		requireUseMw.ApplyFn(usersC.Logout)).Methods("POST")

	r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	//set up routes for css
	cssHandler := http.FileServer(http.Dir("./assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", cssHandler))

	//images routes
	// we can get by by just using http.Dir("./") as the image path matches
	// with our FS
	// but we do this just to show strip prefix which works like a middleware
	imageHandler := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	//Gallery routes
	// galleryNew := requireUseMw.Apply(galleryC.NewView)
	r.Handle("/galleries", requireUseMw.ApplyFn(galleryC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUseMw.Apply(galleryC.NewView)).Methods("GET")
	r.HandleFunc("/galleries", requireUseMw.ApplyFn(galleryC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUseMw.ApplyFn(galleryC.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUseMw.ApplyFn(galleryC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUseMw.ApplyFn(galleryC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleryC.Show).Methods("GET").Name("show_gallery")

	//images POST
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUseMw.ApplyFn(galleryC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUseMw.ApplyFn(galleryC.ImageDelete)).Methods("POST")
	http.ListenAndServe(":3000", csrfMW(userMw.Apply(r)))

}
