package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

//small letter
func parseFormHelper(r *http.Request, dst interface{}) error {
	// return errors.New("balh")
	if err := r.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	if err := decoder.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
