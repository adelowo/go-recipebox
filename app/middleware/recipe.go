package middleware

import (
	"github.com/adelowo/RecipeBox/app/common/session"
	"github.com/adelowo/RecipeBox/app/controller"
	"github.com/adelowo/RecipeBox/app/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//Middleware to ensure that only the original owner of a recipe can delete it
func RecipeOwner(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := mux.Vars(r)["id"]

		//Force convert the id to an integer.
		//If it is not an integer, send a 404
		//If it is an integer,
		//make sure the recipe is owned by the user before any action is carried out
		//If it is not owned by the user,
		//Send a 404

		if recipeId, err := strconv.Atoi(id); err == nil {

			sess, err := session.GetSession(r, "user")

			if err != nil {
				controller.InternalError(w, err)
				return
			}

			currentUser, err := model.FindUserByEmail(sess.Values["email"].(string))

			if err != nil {
				controller.InternalError(w, err)
				return
			}

			if model.IsRecipeOwnedBy(currentUser, recipeId) {
				r.Form.Add("id", id)
				h.ServeHTTP(w, r)
				return
			}

			http.NotFound(w, r)
			return
		}

		http.NotFound(w, r)
		return
	})
}
