package views

import (
	"bytes"
	"goweb_v1/context"
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

// make the render become user status aware by passing the context
func (v View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data

	// data is an empty interface, we want to ensure a check type in place
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		// otherwise we assume whatever data is under Yield type
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	// write to buffer first instead of directly send to w
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		http.Error(w, "something went wrong, pls contact support",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}
