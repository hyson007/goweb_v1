package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir string = "views/layout/"
	LayoutExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	// using a function to get all gohtml under views/layout folder
	filesList, err := filepath.Glob(LayoutDir + "*" + LayoutExt)
	if err != nil {
		panic(err)
	}

	files = append(files, filesList...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{Template: t,
		Layout: layout}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}
