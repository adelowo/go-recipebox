package session

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

var store *sessions.CookieStore

func init() {

	secret := os.Getenv("RECIPEBOX_COOKIE_SECRET")

	if secret == "" {
		log.Fatal("You should set up the secret key for the session in an environment variable")
	}

	store = sessions.NewCookieStore([]byte(secret))
}

//Return a session instance from the session's Store Registry
func GetSession(r *http.Request, sessionName string) (*sessions.Session, error) {

	session, err := store.Get(r, sessionName)

	if err != nil {
		log.Println("Could not get store for the session")
		log.Println("New session used instead")
	}

	return session, err
}
