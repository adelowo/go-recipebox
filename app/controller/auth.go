package controller

import (
	"github.com/adelowo/RecipeBox/app/common/error"
	"github.com/adelowo/RecipeBox/app/common/hasher"
	"github.com/adelowo/RecipeBox/app/common/session"
	"github.com/adelowo/RecipeBox/app/common/template"
	"github.com/adelowo/RecipeBox/app/model"
	v "github.com/asaskevich/govalidator"
	"github.com/gorilla/csrf"
	h "html/template"
	"log"
	"net/http"
)

type signUpError struct {
	template.PageWithForm
	Username, Password, Email string
}

type loginError struct {
	template.PageWithForm
	Email, Password string
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := session.GetSession(r, "user")

	if err == nil {
		delete(session.Values, "username")
		delete(session.Values, "active")
		session.Save(r, w)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
	return
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		postLogin(w, r)

		return
	}

	getLogin(w, r)

	return
}

func postLogin(w http.ResponseWriter, r *http.Request) {

	validatorErrorBag := error.NewValidatorErrorBag()

	r.ParseForm()

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	//Validate

	if email != "" {
		if !v.IsEmail(email) {
			validatorErrorBag.Add("email", "Please provide a valid email address")
		}
	} else {
		validatorErrorBag.Add("email", "Please provide an email address")
	}

	if password == "" {
		validatorErrorBag.Add("password", "Please provide your password")
	}

	if validatorErrorBag.Count() != 0 {
		sendLoginFailureResponse(w, r, validatorErrorBag)
		return
	}

	currentUser, err := model.FindUserByEmail(email)

	if err != nil {
		validatorErrorBag.Add("password", "Invalid password/email combination")
		sendLoginFailureResponse(w, r, validatorErrorBag)
	}

	bcryptHasher := hasher.NewBcryptHasher()

	if bcryptHasher.Verify(currentUser.Password, password) {
		sess, _ := session.GetSession(r, "user")

		sess.Values["username"] = currentUser.Username
		sess.Values["active"] = true

		err := sess.Save(r, w)

		if err != nil {
			validatorErrorBag.Add("email", "An error occured while we try to log you in. Please try again")
			sendLoginFailureResponse(w, r, validatorErrorBag)
			return
		}

		log.Printf("%s just logged in \n", currentUser.Username)

		http.Redirect(w, r, "/", http.StatusFound)

		return

	}

	//Password verification failed.

	validatorErrorBag.Add("password", "Invalid email/password. Try again")

	sendLoginFailureResponse(w, r, validatorErrorBag)
	return
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	csrfField := getCsrfTemplate(r)

	p := template.NewPageStruct("Login to edit your recipeboard")

	template.LoginTemplate.Execute(w, withLoginErrorStruct(csrfField, p, error.NewValidatorErrorBag()))
}

func Signup(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		postSignUp(w, r)
		return
	}

	getSignUp(w, r)
	return
}

func getSignUp(w http.ResponseWriter, r *http.Request) {

	csrfField := getCsrfTemplate(r)

	p := template.NewPageStruct("Signup to own a RecipeBoard")

	template.SignUpTemplate.Execute(w, withSignUpErrorStruct(csrfField, p, error.NewValidatorErrorBag()))
}

func postSignUp(w http.ResponseWriter, r *http.Request) {

	errorBag := error.NewValidatorErrorBag()

	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")

	if username != "" {
		if len(username) < 6 {

			errorBag.Add("username", "Your username should at least 6 characters")
		}

		if len(username) > 30 {
			errorBag.Add("username", "Your username should not be greater than 30 characters")
		}
	} else {
		errorBag.Add("username", "Please provide your username")
	}

	if email != "" {
		if !v.IsEmail(email) {
			errorBag.Add("email", "Please provide a valid email address")
		}
	} else {
		errorBag.Add("email", "Please provide your email address")
	}

	if password == "" {
		errorBag.Add("password", "Please provide a password")
	} else {
		if len(password) < 10 {
			errorBag.Add("password", "Your password should not be lesser than 10 characters")
		}
	}

	if errorBag.Count() != 0 {
		sendSignUpFailureResponse(w, r, errorBag)
		return
	}

	if model.DoesUserExist(email) {
		errorBag.Add("email", "This email address have been registered by someone else. Early bird wins")
	}

	if errorBag.Count() != 0 {
		sendSignUpFailureResponse(w, r, errorBag)
		return
	}

	//Save the user to the database

	if err := model.CreateUser(username, email, password); err != nil {
		errorBag.Add("email", "An error occured while trying to save your details")
		sendSignUpFailureResponse(w, r, errorBag)
		return
	}

	//Save to session, then redirect
	//Don't care about the err since it would be a new session if it does not exist anyways
	session, _ := session.GetSession(r, "user")

	// If we get here, we can successfully persist to session and keep the user logged in

	session.Values["username"] = username
	session.Values["active"] = true

	//Redirect to the index page

	err := session.Save(r, w)

	if err != nil {
		errorBag.Add("email", "Your account have been registered but an error occured while trying to log you in. Please visit the login page to continue")
		sendSignUpFailureResponse(w, r, errorBag)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func getCsrfTemplate(r *http.Request) h.HTML {
	return csrf.TemplateField(r)
}

func sendSignUpFailureResponse(w http.ResponseWriter, r *http.Request, e *error.ValidatorErrorBag) {
	w.WriteHeader(http.StatusFound)
	csrfField := getCsrfTemplate(r)

	p := template.NewPageStruct("Signup to own a RecipeBoard")

	template.SignUpTemplate.Execute(w, withSignUpErrorStruct(csrfField, p, e))

}

//Where did Generics go ?

func sendLoginFailureResponse(w http.ResponseWriter, r *http.Request, eb *error.ValidatorErrorBag) {
	w.WriteHeader(http.StatusFound)
	csrfField := getCsrfTemplate(r)
	p := template.NewPageStruct("Login to edit your recipeboard")

	template.LoginTemplate.Execute(w, withLoginErrorStruct(csrfField, p, eb))

}

func withSignUpErrorStruct(c h.HTML, p template.Page, eb *error.ValidatorErrorBag) signUpError {
	username, _ := eb.Get("username")
	password, _ := eb.Get("password")
	email, _ := eb.Get("email")

	return signUpError{template.NewPageWithFormStruct(p, c), username, password, email}
}

//Again
func withLoginErrorStruct(c h.HTML, p template.Page, eb *error.ValidatorErrorBag) loginError {
	password, _ := eb.Get("password")
	email, _ := eb.Get("email")

	return loginError{template.NewPageWithFormStruct(p, c), email, password}

}
