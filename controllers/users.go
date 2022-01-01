package controllers

import (
	"fmt"
	"goweb_v1/views"
	"net/http"
)

// NewUsers is used to create a new user controller
// this function will panic if template is unable to parse
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView    *views.View
	CreateView *views.View
}

// this is used to rend the form where a user can create a new account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// this is used to process the signup form where a user is submit in order to create a new account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseFormHelper(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
