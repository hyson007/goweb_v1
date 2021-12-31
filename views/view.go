package views

import (
	"html/template"
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
