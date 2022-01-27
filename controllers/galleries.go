package controllers

import (
	"fmt"
	"goweb_v1/context"
	"goweb_v1/models"
	"goweb_v1/views"
	"log"
	"net/http"
)

func NewGallery(gs models.GalleryService) *Gallery {
	return &Gallery{
		NewView: views.NewView("bootstrap", "views/galleries/new.gohtml"),
		gs:      gs,
	}
}

type Gallery struct {
	NewView *views.View
	gs      models.GalleryService
}

func (g *Gallery) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, nil)

}

type GalleryForm struct {
	Title string `schema:"title"`
}

// POST /galleries
func (g *Gallery) Create(w http.ResponseWriter, r *http.Request) {
	var form GalleryForm
	var vd views.Data

	// here we start to render a better response page to use if certain func
	// fails
	if err := parseFormHelper(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	// this user should give us the user obj
	user := context.User(r.Context())
	fmt.Println("got user: ", user)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
