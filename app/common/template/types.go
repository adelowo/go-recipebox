package template

import (
	"github.com/adelowo/RecipeBox/app/common/error"
	"html/template"
)

type Page struct {
	Title string
}

type PageWithForm struct {
	Page
	CsrfField template.HTML
}

//This struct is for pages with redirects
//This allows an easy way to get the errors unto the screen
type ErrorBagAwarePage struct {
	PageWithForm
	Errors *error.ValidatorErrorBag
}
