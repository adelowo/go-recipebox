package template

import (
	"github.com/adelowo/RecipeBox/app/common/error"
	"html/template"
	"io/ioutil"
	"log"
)

var templatesParsed bool

var templateHandler *template.Template

var (
	SignUpTemplate,
	LoginTemplate,
	HomeTemplate,
	CreateRecipeTemplate,
	AllRecipeTemplate *template.Template
)

func ParseTemplates(templates ...string) {

	if templatesParsed {
		panic("Templates have been parsed already.. Program exiting....")
	}

	templatesDir := "./templates/"

	var allFiles []string

	for _, file := range templates {
		_, err := ioutil.ReadFile(templatesDir + file)

		if err != nil {
			log.Fatal("Template not found or cannot be read")
		}

		allFiles = append(allFiles, templatesDir+file)
	}

	allFiles = append(allFiles, getPartials()...)

	log.Println(allFiles)

	templateHandler = template.Must(template.ParseFiles(allFiles...))

	SignUpTemplate = lookTemplateUp("signup.html")
	LoginTemplate = lookTemplateUp("login.html")
	HomeTemplate = lookTemplateUp("index.html")
	CreateRecipeTemplate = lookTemplateUp("add_recipe.html")
	AllRecipeTemplate = lookTemplateUp("recipes.html")

	templatesParsed = true

}

func lookTemplateUp(t string) *template.Template {
	template := templateHandler.Lookup(t)

	if template == nil {
		log.Fatalf("Couldn't find the specified template, %s", t)
	}

	return template
}

func getPartials() []string {
	return []string{"./templates/_partial/_head.html", "./templates/_partial/_footer.html"}
}

func NewPageStruct(t string) Page {

	return Page{t}
}

func NewErrorPageAwareStruct(p Page, c template.HTML, v *error.ValidatorErrorBag) ErrorBagAwarePage {
	return ErrorBagAwarePage{NewPageWithFormStruct(p, c), v}
}

func NewPageWithFormStruct(p Page, c template.HTML) PageWithForm {
	return PageWithForm{p, c}
}
