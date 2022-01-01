package controllers

import "goweb_v1/views"

type Static struct {
	Home    *views.View
	Contact *views.View
}

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
	}
}

// this is jack initial method, not as good
// func StaticHome() *views.View {
// 	return views.NewView("bootstrap", "views/static/home.gohtml")
// }

// func StaticContact() *views.View {
// 	return views.NewView("bootstrap", "views/static/contact.gohtml")
// }
