package model

import (
	"errors"
	"github.com/adelowo/RecipeBox/app/common/database"
	"github.com/adelowo/RecipeBox/app/common/hasher"
)

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func FindUserByEmail(email string) (User, error) {

	var user User

	row := database.Db.QueryRowx("SELECT * FROM users WHERE email=?", email)

	err := row.StructScan(&user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

//This wraps/delegates to the FindUserByEmail.
//It justs allows for expressive usage
//Do not know if this is idiomatic go or PHP land stuff
func DoesUserExist(email string) bool {

	_, err := FindUserByEmail(email)

	return err == nil
}

func CreateUser(username, email, password string) error {

	hashedPassword, err := hasher.NewBcryptHasher().Hash(password)

	if err != nil {
		return err
	}

	stmt, err := database.Db.Preparex("INSERT INTO users(username, email, password) VALUES(?,?,?)")

	if err != nil {
		return err
	}

	result := stmt.MustExec(username, email, hashedPassword)

	if x, _ := result.RowsAffected(); x == 1 {
		return nil
	}

	return errors.New("An error occured while trying to save your details")
}
