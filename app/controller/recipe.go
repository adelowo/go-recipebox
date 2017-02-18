package controller

import (
	"github.com/adelowo/RecipeBox/app/common/error"
	"github.com/adelowo/RecipeBox/app/common/session"
	"github.com/adelowo/RecipeBox/app/common/template"
	"github.com/adelowo/RecipeBox/app/model"
	"github.com/gorilla/mux"
	h "html/template"
	"net/http"
	"strconv"
	"strings"
)

//Bind form values to a struct
//This helps to prevent the user from having to retype/remember the previous entry if there was a failed validation
type formData struct {
	FormTitle, FormDescription, FormIngredients string //Prefixed with Form to model an HTML Form
}

type createRecipe struct {
	template.PageWithForm
	formData
	TitleErr, DescriptionErr, IngredientsErr string //suffixed Titl since <PageWithForm> already has a Title field
}

type recipePage struct {
	template.PageWithForm
	Recipes model.RecipeBox
}

func ShowAllRecipes(w http.ResponseWriter, r *http.Request) {

	sess, err := session.GetSession(r, "user")

	if err != nil {
		InternalError(w, err)
		return
	}

	currentUser, err := model.FindUserByEmail(sess.Values["email"].(string))

	if err != nil {
		InternalError(w, err)
		return
	}

	allRecipes, _ := model.FetchAllRecipesFor(currentUser)

	p := template.NewPageStruct("View all of your recipes")

	template.AllRecipeTemplate.Execute(w, recipePage{template.NewPageWithFormStruct(p, getCsrfTemplate(r)), allRecipes})
}

func AddRecipe(w http.ResponseWriter, r *http.Request) {
	p := template.NewPageStruct("Add a new recipe to your board")

	r.ParseForm()

	formData := newFormDataStruct(r.Form.Get("title"), r.Form.Get("description"), r.Form.Get("ingredients"))

	data := withCreateRecipeStruct(getCsrfTemplate(r), p, error.NewValidatorErrorBag(), formData)

	template.CreateRecipeTemplate.Execute(w, data)
}

func withCreateRecipeStruct(c h.HTML, p template.Page, eb *error.ValidatorErrorBag, formData formData) createRecipe {
	title, _ := eb.Get("title")
	description, _ := eb.Get("description")
	ingredients, _ := eb.Get("ingredients")

	return createRecipe{template.NewPageWithFormStruct(p, c), formData, title, description, ingredients}
}

func newFormDataStruct(title, description, ingredients string) formData {
	return formData{title, description, ingredients}
}

func SaveRecipe(w http.ResponseWriter, r *http.Request) {

	validationErrorBag := error.NewValidatorErrorBag()

	r.ParseForm()

	title := r.Form.Get("title")

	//Ingredients are seperated by "," in the GUI. Let's make sure we don't have any at index 0 and/or n
	ingredients := strings.Trim(r.Form.Get("ingredients"), ",")

	description := r.Form.Get("description")

	//Validate length ? Some sick recipes might be 2,3 letters long though
	//Because of the above, we only make sure we don't get a non empty title

	if title == "" {
		validationErrorBag.Add("title", "Please provide a title for your recipe")
	}

	if ingredients == "" {
		validationErrorBag.Add("ingredients", "Please provide the ingredients needed to make this recipe")
	}

	if description == "" {
		validationErrorBag.Add("description", "Please provide a description for this recipe")
	}

	if validationErrorBag.Count() != 0 {
		sendCreateRecipeFailureResponse(w, r, validationErrorBag)
		return
	}

	sess, err := session.GetSession(r, "user")

	if err != nil {
		InternalError(w, err)
		return
	}

	currentUser, err := model.FindUserByEmail(sess.Values["email"].(string))

	if err != nil {
		//This shouldn't happen though. Let's just double check
		InternalError(w, err)
		return
	}

	if model.DoesRecipeExist(title) {
		validationErrorBag.Add("title", "You already have a recipe with the title, "+title)
		sendCreateRecipeFailureResponse(w, r, validationErrorBag)
		return
	}

	err = model.CreateRecipe(currentUser, title, ingredients, description)

	if err != nil {
		validationErrorBag.Add("title", "An error occured while we tried saving your recipe")
		sendCreateRecipeFailureResponse(w, r, validationErrorBag)
		return
	}

	sess.AddFlash("Your recipe, "+title+" was created successfully", "recipe.created")

	sess.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func sendCreateRecipeFailureResponse(w http.ResponseWriter, r *http.Request, eb *error.ValidatorErrorBag) {
	w.WriteHeader(http.StatusFound)
	csrf := getCsrfTemplate(r)
	p := template.NewPageStruct("Add a new recipe to your board")

	formData := newFormDataStruct(r.Form.Get("title"), r.Form.Get("description"), r.Form.Get("ingredients"))

	template.CreateRecipeTemplate.Execute(w, withCreateRecipeStruct(csrf, p, eb, formData))
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id, _ := strconv.Atoi((r.Form.Get("id")))

	err := model.DeleteRecipe(id)

	if err != nil {
		InternalError(w, err)
	}

	http.Redirect(w, r, "/recipes", http.StatusFound)

}

type singleRecipePage struct {
	template.Page
	Recipe model.Recipe
}

func ViewRecipe(w http.ResponseWriter, r *http.Request) {

	//We don't care about the error value here since we defined a regex to capture only the numerical part of the route
	//We can thus be certain that there wouldn't be any error.
	//Am i missing something here ?
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	recipe, err := model.FindRecipeById(id)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := singleRecipePage{template.NewPageStruct("View this recipe"), recipe}

	template.SingleRecipePage.Execute(w, data)
}

func EditRecipe(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		postEditRecipe(w, r)
		return
	}

	getEditRecipe(w, r)
}

func postEditRecipe(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	recipe, err := model.FindRecipeById(id)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	validationErrorBag := error.NewValidatorErrorBag()

	//Duplication and in some other places though

	r.ParseForm()

	title := r.Form.Get("title")

	ingredients := strings.Trim(r.Form.Get("ingredients"), ",")

	description := r.Form.Get("description")

	if title == "" {
		validationErrorBag.Add("title", "Please provide a title for your recipe")
	}

	if ingredients == "" {
		validationErrorBag.Add("ingredients", "Please provide the ingredients needed to make this recipe")
	}

	if description == "" {
		validationErrorBag.Add("description", "Please provide a description for this recipe")
	}

	if validationErrorBag.Count() != 0 {
		sendEditRecipeFailure(w, r, validationErrorBag, recipe)
		return
	}

	//End duplication. 8)

	err = model.Update(recipe, title, description, ingredients)

	if err != nil {
		InternalError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)

}

func sendEditRecipeFailure(w http.ResponseWriter, r *http.Request, eb *error.ValidatorErrorBag, recipe model.Recipe) {
	w.WriteHeader(http.StatusFound)

	formData := newFormDataStruct(r.Form.Get("title"), r.Form.Get("description"), r.Form.Get("ingredients"))

	d := withCreateRecipeStruct(getCsrfTemplate(r), template.NewPageStruct("Edit this recipe"), eb, formData)

	data := editPage{recipe, d}

	template.EditRecipeTemplate.Execute(w, data)
}

func getEditRecipe(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	recipe, err := model.FindRecipeById(id)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	formData := newFormDataStruct(r.Form.Get("title"), r.Form.Get("description"), r.Form.Get("ingredients"))

	d := withCreateRecipeStruct(getCsrfTemplate(r), template.NewPageStruct("Edit this recipe"), error.NewValidatorErrorBag(), formData)

	data := editPage{recipe, d}

	template.EditRecipeTemplate.Execute(w, data)
}

type editPage struct {
	Recipe model.Recipe
	createRecipe
}
