package controllers

import (
	"fmt"
	"goweb_v1/context"
	"goweb_v1/models"
	"goweb_v1/views"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte
)

func NewGallery(gs models.GalleryService, is models.ImageService, r *mux.Router) *Gallery {
	return &Gallery{
		NewView:   views.NewView("bootstrap", "views/galleries/new.gohtml"),
		ShowView:  views.NewView("bootstrap", "views/galleries/show.gohtml"),
		EditView:  views.NewView("bootstrap", "views/galleries/edit.gohtml"),
		IndexView: views.NewView("bootstrap", "views/galleries/index.gohtml"),
		gs:        gs,
		r:         r,
		is:        is,
	}
}

type Gallery struct {
	NewView   *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	r         *mux.Router
	is        models.ImageService
}

func (g *Gallery) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, r, nil)

}

type GalleryForm struct {
	Title string `schema:"title"`
}

// this will show every single gallery use has access to
// GET /galleries/
func (g *Gallery) Index(w http.ResponseWriter, r *http.Request) {
	// first get current user, then get galleries by userID
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	var vd views.Data
	vd.Yield = galleries
	// fmt.Println("from edit", vd)
	// fmt.Fprintln(w, galleries)
	g.IndexView.Render(w, r, vd)
	// fmt.Fprintln(w, gallery)
}

// GET /galleries/:id
// majority of this show code is same for the edit code
// hence we write a new func galleryByID
func (g *Gallery) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	// fmt.Println("gallery Show", gallery, err)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
	// fmt.Fprintln(w, gallery)
}

// GET /galleries/:id/edit
func (g *Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)

	if err != nil {
		return
	}

	// verify if the user can access the gallery or not
	user := context.User(r.Context())
	// fmt.Println("gallery Edit user", user)
	// fmt.Println("gallery Edit", gallery.ID, err)
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery Not found", http.StatusNotFound)
		return
	}

	// if found the gallery then we render a form for users to edit
	var vd views.Data
	vd.Yield = gallery
	// make vd users aware
	// vd.User = user
	// fmt.Println("from edit", vd)
	g.EditView.Render(w, r, vd)
	// fmt.Fprintln(w, gallery)
}

// POST /galleries/:id/update
func (g *Gallery) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// verify if the user can access the gallery or not
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery Not found", http.StatusNotFound)
		return
	}

	// if found the gallery then we render a form for users to edit
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseFormHelper(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvSuccess,
		Message: "Gallery successfully updated",
	}
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/delete
func (g *Gallery) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// verify if the user can access the gallery or not
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery Not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}
	// TODO: redirect to index page
	fmt.Fprintln(w, "successfully delete! ")
}

// POST /galleries/:id/images
func (g *Gallery) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// verify if the user can access the gallery or not
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery Not found", http.StatusNotFound)
		return
	}

	// we need to parse the multi-part form
	err = r.ParseMultipartForm(maxMultipartMem)
	var vd views.Data
	vd.Yield = gallery
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	//create folder, only need to run one time
	// galleryPath := fmt.Sprintf("images/galleries/%v/", gallery.ID)
	// err = os.MkdirAll(galleryPath, 0755)
	// if err != nil {
	// 	vd.SetAlert(err)
	// 	g.EditView.Render(w, r, vd)
	// 	return
	// }

	files := r.MultipartForm.File["images"]
	for _, f := range files {
		// open file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()
		err = g.is.CreateImage(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}

		http.Redirect(w, r, "/galleries", http.StatusFound)

	}
}

// POST /galleries/:id/images/:filename/delete
// this way no need to post data in payload
func (g *Gallery) ImageDelete(w http.ResponseWriter, r *http.Request) {
	//Get gallery
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// verify if the user can access the gallery or not
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery Not found", http.StatusNotFound)
		return
	}

	FileName := mux.Vars(r)["filename"]
	// fmt.Println(FileName)

	i := models.Image{
		Filename:  FileName,
		GalleryID: gallery.ID,
	}

	_ = i
	var vd views.Data
	// delete from os
	err = g.is.Delete(&i)
	if err != nil {
		fmt.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	path := fmt.Sprintf("%v", gallery.ID)
	http.Redirect(w, r, "/galleries/"+path, http.StatusFound)
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
		g.NewView.Render(w, r, vd)
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
		g.NewView.Render(w, r, vd)
		return
	}

	// redirect to gallery after user post
	// url, err := g.r.Get("show_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	// this err should not occur
	// if err != nil {
	// 	http.Redirect(w, r, "/", http.StatusFound)
	// 	return
	// }

	// finally redirect user
	http.Redirect(w, r, "/galleries", http.StatusFound)
	// fmt.Fprintln(w, gallery)
}

func (g *Gallery) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	// this id is string now, we need to convert it to int
	id, err := strconv.Atoi(idStr)
	// fmt.Println("before lookup, galleryByID, id is:", id, "error is: ", err)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	// otherwise we look up this gallery id
	gallery, err := g.gs.ByID(uint(id))
	// fmt.Println("after lookup, galleryByID gallery is:", gallery, "error is: ", err)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "whoops! something went wrong", http.StatusInternalServerError)

		}
		return nil, err
	}
	images, err := g.is.ByGalleryID(gallery.ID)
	if err != nil {
		gallery.Images = []models.Image{}
	} else {
		gallery.Images = images
	}

	return gallery, nil
}
