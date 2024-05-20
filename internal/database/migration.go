package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func Migrate(db *sql.DB) {

	err := goose.SetDialect("postgres")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	path := "../../db/migrations" // path if started locally
	if _, err := os.Stat(path); os.IsNotExist(err) {
		old_path := path
		path = "migrations" // path if started in docker
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatalf("Migrations directory not found: %s OR %s", old_path, path)
			panic(err)
		}
	}

	err = goose.Up(db, path)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
