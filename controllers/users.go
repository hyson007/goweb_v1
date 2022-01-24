package controllers

import (
	"fmt"
	"goweb_v1/models"
	"goweb_v1/rand"
	"goweb_v1/views"
	"net/http"
)

// NewUsers is used to create a new user controller
// this function will panic if template is unable to parse
// Update, change userservice to interface
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "views/users/new.gohtml"),
		LoginView: views.NewView("bootstrap", "views/users/login.gohtml"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	// CreateView *views.View
	us models.UserService
}

// this is used to rend the form where a user can create a new account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// this is used to process the signup form where a user is submit in order to create a new account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseFormHelper(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// here this user is not a pointer hence need &
	// once user created, we can just assign a cookie back
	err := u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintln(w, user)

	//redirect user to cookietest page
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify provided email and password
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseFormHelper(r, &form); err != nil {
		panic(err)
	}
	// Do something with login Form

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid user email address")
		case models.ErrInvalidPwd:
			fmt.Fprintln(w, "Invalid user password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// once login successfully, we can assign a cookie to user
	err = u.signIn(w, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintln(w, user)

	//redirect user to cookietest page
	http.Redirect(w, r, "/cookietest", http.StatusFound)

	// if no err and we are going to return a email cookie first
	// this setcookie must happen before Fprint
	// cookie := http.Cookie{
	// 	Name:  "email",
	// 	Value: user.Email,
	// }
	// http.SetCookie(w, &cookie)

	// optionally, call cookieTest to verify or we can create a new route in main
	//u.CookieTest(w, r)
}

//Create a new func just for setting cookie part
// func signIn(w http.ResponseWriter, user *models.User) {
// 	cookie := http.Cookie{
// 		Name:  "email",
// 		Value: user.Email,
// 	}
// 	http.SetCookie(w, &cookie)
// 	// fmt.Fprintln(w, user)
// }

//convert this func to receiver on user
//sign in is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {

	//this step ensures whatever user in DB has a tokenhash
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// fmt.Fprintln(w, user)
	return nil
}

//CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "remember_token cookie is", cookie.Value)
	fmt.Fprintln(w, user)
}
