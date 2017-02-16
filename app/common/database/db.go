package database

import (
	c "github.com/adelowo/RecipeBox/app/common/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Db *sqlx.DB

func init() {

	var config c.Configuration

	config, err := c.ReadConfig("config/config.json")

	if err != nil {
		log.Fatal("An error occured while trying to read the configuration value")
	}

	Db = sqlx.MustConnect("sqlite3", config.Database)

}
