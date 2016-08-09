package ggdb

import (
	"database/sql"
	"log"
	"os"
	"os/user"

	"github.com/leerumler/gengar/ggconf"

	// Import the SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

// GGDB holds a connection to gengar's database.
var GGDB *sql.DB

// findGGDB locates the gengar database file.
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
	dbfile := confdir + "/gg.db"
	return &dbfile
}

// ConnectGGDB establishes a connection to Gengar's database.
func connectGGDB() *sql.DB {

	// Find the database and open it.
	dbfile := findGGDB()
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}

	// Return a pointer to that database connection.
	return db
}

// CleanSlate creates an empty SQLite database to act as gengar's central data store.
func CleanSlate() {

	// Get a pointer to the database connection.
	db := connectGGDB()
	defer db.Close()

	// Drop existing tables and (re-)create them.
	createTables := `
	DROP TABLE IF EXISTS expansions;
	DROP TABLE IF EXISTS phrases;
	CREATE TABLE expansions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expansion TEXT
		);
	CREATE TABLE phrases (
		phrase TEXT PRIMARY KEY,
		exp_id INTEGER,
		FOREIGN KEY (exp_id) REFERENCES expansions(id)
	);
	`
	_, err := db.Exec(createTables)

	// Die on error.
	if err != nil {
		log.Fatal("Couldn't create CleanSlate:", err)
	}
}

// findExpansionID finds and returns the expansion ID from gengar's database.
func findExpansionID(exp *ggconf.Expander) int {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Check the expansion ID in the database.
	var expID int
	err := db.QueryRow("SELECT id FROM expansions WHERE expansion=$1", exp.Expansion).Scan(&expID)
	switch {
	case err == sql.ErrNoRows:
		log.Fatal("No matching expansions found.")
	case err != nil:
		log.Fatal(err)
	}

	// Return the ID.
	return expID
}

// InsertExpansion inserts an Expansion into gengar's database.
func InsertExpansion(exp *ggconf.Expander) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the expansion.
	_, err := db.Exec("INSERT INTO expansions (expansion) VALUES ($1);", exp.Expansion)
	if err != nil {
		log.Fatal(err)

	}
}

// MapPhrase maps a phrase to an expansion in gengar's database.
func MapPhrase(exp *ggconf.Expander) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the expansion.
	_, err := db.Exec("INSERT INTO phrases (phrase, exp_id) VALUES ($1, $2)", exp.Phrase, exp.ID)
	if err != nil {
		log.Fatal(err)

	}
}

// ReadExpanders reads expansions from the database and returns a slice of ggconf.Expanders.
func ReadExpanders() *[]ggconf.Expander {
	var exps []ggconf.Expander

	// Get pointer to database connection.
	db := connectGGDB()
	defer db.Close()

	// Query the database for the expansions.
	rows, err := db.Query("SELECT exp_id, phrase, expansion FROM phrases JOIN expansions ON phrases.exp_id = expansions.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Read through the results and populate exps with returned info.
	for rows.Next() {
		var exp ggconf.Expander
		err = rows.Scan(&exp.ID, &exp.Phrase, &exp.Expansion)
		exps = append(exps, exp)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Return pointer to exps.
	return &exps
}

// CreateTestDB creates a test database.
func CreateTestDB() {

	// Create some testing expansions.
	var exps []ggconf.Expander
	exps = append(exps, ggconf.Expander{Phrase: "test1", Expansion: "This is test 1!", ID: 1})
	exps = append(exps, ggconf.Expander{Phrase: "test2", Expansion: "this is test 2", ID: 2})
	exps = append(exps, ggconf.Expander{Phrase: "test3", Expansion: "this is test 3", ID: 3})

	// Wipe/create a blank gengar database.
	CleanSlate()

	// Insert each of our testing expansions.
	for _, exp := range exps {
		InsertExpansion(&exp)
		MapPhrase(&exp)
	}
}
