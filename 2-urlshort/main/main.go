package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercises/urlshort"
	_ "github.com/mattn/go-sqlite3"
)

func createTable(db *sql.DB) {
	createSQL := `
		BEGIN;
		CREATE TABLE mapping (
			"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"path" TEXT,
			"url" TEXT
		);
		CREATE INDEX path_idx ON mapping (path);
		COMMIT;
	`

	_, err := db.Exec(createSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func insertMapping(db *sql.DB, path string, url string) {
	insertSQL := `INSERT INTO mapping(path, url) VALUES (?, ?)`
	statement, err := db.Prepare(insertSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(path, url)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func setupDatabase() *sql.DB {
	databaseName := "urlshort.db"

	os.Remove(databaseName)
	file, err := os.Create(databaseName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	file.Close()

	db, _ := sql.Open("sqlite3", databaseName)

	createTable(db)
	insertMapping(db, "/urlshort-godoc", "https://godoc.org/github.com/gophercises/urlshort")
	insertMapping(db, "/yaml-godoc", "https://godoc.org/gopkg.in/yaml.v3")

	return db
}

func main() {
	db := setupDatabase()
	defer db.Close()

	mux := defaultMux()

	dbHandler := urlshort.DBHandler(db, mux)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
