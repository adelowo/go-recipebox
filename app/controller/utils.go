package controller

import (
	"github.com/gorilla/csrf"
	h "html/template"
	"net/http"
)

func getCsrfTemplate(r *http.Request) h.HTML {
	return csrf.TemplateField(r)
}
