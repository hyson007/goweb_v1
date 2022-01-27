package middleware

import (
	"fmt"
	"goweb_v1/context"
	"goweb_v1/models"
	"net/http"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//this is an example of what we can do in middlware
		//t := time.Now()
		//fmt.Println("Request start at :", t)
		//next(w, r)
		//fmt.Println("Request ended at:", time.Since(t))

		// if the user is logged in ...
		cookie, err := r.Cookie("remember_token")

		// if unable to find the user cookie, then we redirect them
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// give me the current context from request
		ctx := r.Context()
		// apply user to that context
		ctx = context.WithUser(ctx, user)
		// update request with new context
		r = r.WithContext(ctx)
		fmt.Println("User Found: ", user)
		// the new context will be pass to next
		next(w, r)
	})
}
