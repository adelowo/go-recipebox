package middleware

import (
	"github.com/adelowo/RecipeBox/app/common/session"
	"net/http"
)

//Middleware that makes sure authenticated users cannot view pages meant for guests
//Example : Once logged in, you shouldn't be able to use the login or sign up again till you are logout'd
func Guest(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if session.IsUserLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)

	})
}
