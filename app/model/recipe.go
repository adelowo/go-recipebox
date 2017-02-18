package model

import (
	"errors"
	"github.com/adelowo/RecipeBox/app/common/database"
)

type Recipe struct {
	Id          int    `db:"id"`
	UserID      int    `db:"user_id"`
	Title       string `db:"title"`
	Ingredients string `db:"ingredients"`
	Description string `db:"description"`
}

type RecipeBox struct {
	Recipes []Recipe
}

func FindRecipeById(id int) (Recipe, error) {
	var recipe Recipe

	row := database.Db.QueryRowx("SELECT * FROM recipes WHERE id=?", id)

	err := row.StructScan(&recipe)

	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil

}

func IsRecipeOwnedBy(u User, id int) bool {

	recipe, err := FindRecipeById(id)

	if err != nil {
		return false
	}

	return recipe.UserID == u.Id
}

func FindRecipeByTitle(title string) (Recipe, error) {

	var recipe Recipe

	row := database.Db.QueryRowx("SELECT * FROM recipes WHERE title=?", title)

	err := row.StructScan(&recipe)

	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func CreateRecipe(u User, title, ingredients, description string) error {

	stmt, err := database.Db.Preparex("INSERT INTO recipes(user_id,title, ingredients, description) VALUES(?,?,?,?)")

	if err != nil {
		return err
	}

	result := stmt.MustExec(u.Id, title, ingredients, description)

	if x, _ := result.RowsAffected(); x == 1 {
		return nil
	}

	return errors.New("An error occured while we tried saving your recipe. Please try again")
}

func DoesRecipeExist(title string) bool {
	_, err := FindRecipeByTitle(title)

	return err == nil
}

func FetchAllRecipesFor(u User) (RecipeBox, error) {

	var recipeBox RecipeBox
	var recipe Recipe

	rows, err := database.Db.Queryx("SELECT * FROM recipes WHERE user_id=?", u.Id)

	defer rows.Close()

	if err != nil {
		return RecipeBox{}, err
	}

	for rows.Next() {
		rows.StructScan(&recipe)

		recipeBox.Recipes = append(recipeBox.Recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return RecipeBox{}, err
	}

	return recipeBox, nil
}

func DeleteRecipe(id int) error {

	stmt, err := database.Db.Preparex("DELETE FROM recipes WHERE id=?")

	if err != nil {
		return err
	}

	result := stmt.MustExec(id)

	if x, _ := result.RowsAffected(); x == 1 {
		return nil
	}

	return errors.New("An error occured while we tried deleting your recipe")
}

func Update(r Recipe, title, description, ingredients string) error {

	stmt, err := database.Db.Preparex("UPDATE recipes SET title=?,description=?, ingredients=? WHERE id=?")

	if err != nil {
		return err
	}

	result := stmt.MustExec(title, description, ingredients, r.Id)

	if x, _ := result.RowsAffected(); x == 1 {
		return nil
	}

	return errors.New("An error occured while your recipe was being updated")
}
