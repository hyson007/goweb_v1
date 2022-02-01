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
	// by default schema will error when it sees unknown field in struct
	// 2022/01/31 11:03:39 schema: invalid path "gorilla.csrf.Token"
	// to disable this behavior
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
