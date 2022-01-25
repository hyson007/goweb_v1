package views

import (
	"bytes"
	"html/template"
	"io"
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

func (v View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	// data is an empty interface, we want to ensure a check type in place

	switch data.(type) {
	case Data:
		// do nothing
	default:
		// otherwise we assume whatever data is under Yield type
		data = Data{
			Yield: data,
		}
	}
	// write to buffer first instead of directly send to w
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		http.Error(w, "something went wrong, pls contact support",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}
