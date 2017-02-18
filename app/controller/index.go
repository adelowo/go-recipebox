package controller

import (
	"fmt"
	//	"github.com/gorilla/mux"
	"github.com/adelowo/RecipeBox/app/common/session"
	"github.com/adelowo/RecipeBox/app/common/template"
	"net/http"
)

type indexPage struct {
	template.Page
	Username string
}

func Index(w http.ResponseWriter, r *http.Request) {

	sess, err := session.GetSession(r, "user")

	if err != nil {
		InternalError(w, err)
		return
	}

	flash := sess.Flashes("recipe.created")

	fmt.Println(flash)

	data := indexPage{template.NewPageStruct("View your recipeboard"), sess.Values["username"].(string)}

	template.HomeTemplate.Execute(w, data)
}
