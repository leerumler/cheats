package ggdb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// CreateGengarDB creates an SQLite database to act as gengar's central data store.
func CreateGengarDB() {

	ggdb, err := sql.Open("sqlite3", "~/.config/gengar/ggdb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer ggdb.Close()

	sqlStmt := "create table expansions (id integer primary key autoincrement, expansion text);"
	_, err = ggdb.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
