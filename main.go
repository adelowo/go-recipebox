package main

import (
	"flag"
	"fmt"
	"github.com/adelowo/RecipeBox/app/common/config"
	_ "github.com/adelowo/RecipeBox/app/common/database"
	_ "github.com/adelowo/RecipeBox/app/common/session"
	"github.com/adelowo/RecipeBox/app/common/template"
	"github.com/adelowo/RecipeBox/app/controller"
	"github.com/adelowo/RecipeBox/app/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
)

const (
	CSRF_KEY = "RECIPEBOX_SECRET" //This value resides in the env
)

var isProduction bool

func main() {

	flag.BoolVar(&isProduction, "prod", true, "Allow only secure Csrf Tokens ? Set to false for dev env. True for production.")
	flag.Parse()

	config, err := config.ReadConfig("config/config.json")

	if err != nil {
		log.Fatal("Cannot read the config file")
	}

	csrfKey := os.Getenv(CSRF_KEY)

	if csrfKey == "" {
		message := fmt.Sprintf("The world is going to drop if you don't have a csrf key. \n You should set an env variable called %s", CSRF_KEY)
		panic(message)
	}

	template.ParseTemplates("auth/signup.html", "auth/login.html")

	r := mux.NewRouter()

	index := http.HandlerFunc(controller.Index)

	r.Handle("/", alice.New(middleware.Auth).ThenFunc(index)).Methods("GET")

	r.HandleFunc("/logout", controller.Logout).Methods("GET")

	r.HandleFunc("/login", controller.Login).Methods("GET", "POST")

	r.HandleFunc("/signup", controller.Signup).Methods("GET", "POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Printf("Server started on %s", config.Host+config.Port)

	opts := csrf.Secure(isProduction)

	http.ListenAndServe(config.Host+config.Port, csrf.Protect([]byte(csrfKey), opts)(r))
}
