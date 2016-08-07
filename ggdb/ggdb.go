package ggdb

import (
	"database/sql"
	"log"
	"os"
	"os/user"

	// Import the SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

// expander defines a struct that holds text epansion information.
type expander struct {
	phrase, expansion string
}

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
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	_, err = db.Exec(createTables)
	if err != nil {
		log.Printf("%q: %s\n", err, createTables)
		return
	}
}

// InsertExpansion inserts an Expansion into gengar's database.
func InsertExpansion(exp expander) {

	// Find the database and open it.
	dbfile := findGGDB()
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert the expansion.
	_, err = db.Exec("insert into expansions (expansion) values ($1);", exp.expansion)
	if err != nil {
		log.Fatal(err)

	}
}

func findExpansionID(exp expander) int {
	// Find the database and open it.
	dbfile := findGGDB()
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var expID int
	err = db.QueryRow("select id from expansions where expansion=$1", exp.expansion).Scan(&expID)
	switch {
	case err == sql.ErrNoRows:
		log.Fatal("No matching expansions found.")
	case err != nil:
		log.Fatal(err)
	default:
		log.Println("Expansion ID is", expID)
	}

	return expID

}

// MapPhrase maps a phrase to an expansion in gengar's database.
func MapPhrase(exp expander) {

}
