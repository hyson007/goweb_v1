package middleware

import (
	"goweb_v1/context"
	"goweb_v1/models"
	"net/http"
	"strings"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// if the user is requesting static assert or image
		// we no need to lookup the current user
		if strings.HasPrefix(path, "/assets/") ||
			strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("remember_token")

		// if unable to find the user cookie, we assume user never login
		// no need to redirect, we just run the next
		if err != nil {
			next(w, r)
			return
		}

		user, err := mw.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		// give me the current context from request
		ctx := r.Context()
		// apply user to that context
		ctx = context.WithUser(ctx, user)
		// update request with new context
		r = r.WithContext(ctx)
		// fmt.Println("User Found: ", user)
		// the new context will be pass to next
		next(w, r)
	})
}

// RequireUser assume that user has been run otherwise it will not work
type RequireUser struct {
	// models.UserService
	User
}

// Apply assumes that the middle has already been run otherwise it will not work
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn assumes that user middleware has been already run, otherwise it will not work
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	ourHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})

	return mw.User.Apply(ourHandler)

	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	//this is an example of what we can do in middlware
	// 	//t := time.Now()
	// 	//fmt.Println("Request start at :", t)
	// 	//next(w, r)
	// 	//fmt.Println("Request ended at:", time.Since(t))

	// 	// if the user is logged in ...
	// 	cookie, err := r.Cookie("remember_token")

	// 	// if unable to find the user cookie, then we redirect them
	// 	if err != nil {
	// 		http.Redirect(w, r, "/login", http.StatusFound)
	// 		return
	// 	}

	// 	user, err := mw.ByRemember(cookie.Value)
	// 	if err != nil {
	// 		http.Redirect(w, r, "/login", http.StatusFound)
	// 		return
	// 	}
	// 	// give me the current context from request
	// 	ctx := r.Context()
	// 	// apply user to that context
	// 	ctx = context.WithUser(ctx, user)
	// 	// update request with new context
	// 	r = r.WithContext(ctx)
	// 	fmt.Println("User Found: ", user)
	// 	// the new context will be pass to next
	// 	next(w, r)
	// })
}
