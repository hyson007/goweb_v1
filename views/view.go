package views

import (
	"bytes"
	"errors"
	"goweb_v1/context"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
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

	// t, err := template.ParseFiles(files...)
	// template.New("") is same as without new, except can put in a name
	// Funcs takes in a func map, which key is a field gets pass in template
	// example csrfField can be used in template {{csrfField}}
	// t, err := template.New("").Funcs(template.FuncMap{
	// 	"csrfField": func() template.HTML {
	// 		return "<h1> test csrf </h1>"
	// 	},
	// }).ParseFiles(files...)

	// we can also let the func returns html template and error
	// we shall see this kind of error
	// 2022/01/31 10:53:28 template: edit.gohtml:98:2: executing "deleteImagesForm" at <csrfField>: error calling csrfField: csrfField is not implemented
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
		},
	}).ParseFiles(files...)

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
	// implement alert awareness for render
	if alert := getAlert(r); alert != nil {
		vd.Alert = alert
		clearAlert(w)
	}

	vd.User = context.User(r.Context())

	// csrf template part
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	// write to buffer first instead of directly send to w
	var buf bytes.Buffer
	// fmt.Printf("%+v\n", vd.Yield)
	// if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "something went wrong, pls contact support",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}
