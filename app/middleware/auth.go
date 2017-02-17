package middleware

import (
	"github.com/adelowo/RecipeBox/app/common/session"
	"net/http"
)

//Protect routes from unauthenticated users
func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !session.IsUserLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		h.ServeHTTP(w, r)
	})
}
