package ggdb

import (
	"database/sql"
	"log"
	"os"
	"os/user"

	// Import the SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

func findGGDB() *string {

	// Check the current user.
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Find the homedir and add the rest of the path.
	homedir := usr.HomeDir

	// Make the gengar config directory if it doesn't exist.
	confdir := homedir + "/.config/gengar"
	os.MkdirAll(confdir, 0700)

	// Append database name and return a pointer.
	dbfile := confdir + "/ggdb.db"
	return &dbfile
}

// CreateGengarDB creates an SQLite database to act as gengar's central data store.
func CreateGengarDB() {

	// Find the SQLite DB file.
	dbfile := findGGDB()

	// Open the database file.  This should create it if it doesn't exist.
	ggdb, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer ggdb.Close()

	createTables := `
	drop table if exists expansions;
	drop table if exists phrases;
	create table expansions (
		id integer primary key autoincrement,
		expansion text
		);
	create table phrases (
		phrase text primary key,
		exp_id integer,
		foreign key (exp_id) references expansions(id)
	);
	`

	_, err = ggdb.Exec(createTables)
	if err != nil {
		log.Printf("%q: %s\n", err, createTables)
		return
	}
}
