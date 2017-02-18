package controller

import (
	"github.com/adelowo/RecipeBox/app/common/template"
	"net/http"
)

type indexPage struct {
	template.Page
	Username string
}

func Index(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/recipes", http.StatusFound)
}
