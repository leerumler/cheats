package ggdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"

	// Import the SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

// Phrase holds information about phrases.
type Phrase struct {
	Name      string
	ID, ExpID int
}

// Expansion holds information about expansions
type Expansion struct {
	Name, Text string
	ID, CatID  int
}

// Category holds information about categories.
type Category struct {
	ID   int
	Name string
}

// Expander defines a struct that holds text epansion information.
type Expander struct {
	Phrase, Expansion string
	ID                int
}

// GGDB holds a connection to gengar's database.
var GGDB *sql.DB

// FindGGDB locates the gengar database file.
func FindGGDB() *string {

	// Check the current user.
	usr, err := user.Current()
	if err != nil {
		log.Panicln(err)
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
	dbfile := FindGGDB()
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Panicln(err)
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
	DROP TABLE IF EXISTS categories;
	DROP TABLE IF EXISTS expansions;
	DROP TABLE IF EXISTS phrases;
	CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE expansions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		text TEXT UNIQUE,
		cat_id INTEGER DEFAULT 1,
		FOREIGN KEY (cat_id) REFERENCES categories(id)
	);
	CREATE TABLE phrases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		exp_id INTEGER,
		FOREIGN KEY (exp_id) REFERENCES expansions(id)
	);

	`
	_, err := db.Exec(createTables)

	// Die on error.
	if err != nil {
		log.Panicln("Couldn't create CleanSlate:", err)
	}
}

// AddCategory inserts a new Category into gengar's database.
func AddCategory(cat *Category) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the category.
	_, err := db.Exec("INSERT INTO categories (name) VALUES ($1);", cat.Name)
	if err != nil {
		log.Panicln("Couldn't insert category: ", err)
	}
}

// AddExpansion inserts a new Expansion into gengar's database.
func AddExpansion(exp *Expansion) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the expansion.
	if _, err := db.Exec("INSERT INTO expansions (name, cat_id) VALUES ($1, $2);", exp.Name, exp.CatID); err != nil {
		log.Panicln("Couldn't insert expansion: ", err)
	}
}

// AddPhrase maps a phrase to an expansion in gengar's database.
func AddPhrase(phrase *Phrase) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the phrase.
	_, err := db.Exec("INSERT INTO phrases (name, exp_id) VALUES ($1, $2);", phrase.Name, phrase.ExpID)
	if err != nil {
		log.Panicln("Couldn't insert phrase: ", err)
	}
}

// UpdateCategory updates a category's name.
func UpdateCategory(cat *Category) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Update the category's name
	if _, err := db.Exec("UPDATE categories SET name = $1 WHERE id = $2;", cat.Name, cat.ID); err != nil {
		log.Panicln("Couldn't update expansion: ", err)
	}
}

// UpdateExpansionName updates an expansion's name.
func UpdateExpansionName(exp *Expansion) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Update the expansion.
	if _, err := db.Exec("UPDATE expansions SET name = $1 WHERE id = $2;", exp.Name, exp.ID); err != nil {
		log.Panicln("Couldn't update expansion: ", err)
	}
}

// UpdateExpansionText updates an expansion's text.
func UpdateExpansionText(exp *Expansion) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Update the expansion.
	if _, err := db.Exec("UPDATE expansions SET text = $1 WHERE id = $2;", exp.Text, exp.ID); err != nil {
		log.Panicln("Couldn't update expansion: ", err)
	}
}

// UpdatePhrase updates a phrase's name.
func UpdatePhrase(phrase *Phrase) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Update the phrase.
	if _, err := db.Exec("UPDATE phrases SET name = $1 WHERE id = $2;", phrase.Name, phrase.ID); err != nil {
		log.Panicln("Couldn't update phrase: ", err)
	}
}

// ReadCategories reads the available categories from the database.
func ReadCategories() []Category {

	// Create an empty slice of Categories to fill.
	var cats []Category

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for all available categories.
	rows, err := db.Query("SELECT id, name FROM categories ORDER BY name;")
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	// Load the query's results in to new Categories and append them to the slice.
	for rows.Next() {
		var cat Category
		err = rows.Scan(&cat.ID, &cat.Name)
		cats = append(cats, cat)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Panicln(err)
	}

	// Return the filled slice of Categories
	return cats
}

// ReadExpansions finds all of the Expansions within a given category.
func ReadExpansions(cat *Category) []Expansion {

	// Create an empty slice of Expansions to fill.
	var exps []Expansion

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for Expansions matching the category's ID.
	rows, err := db.Query("SELECT id, name, text, cat_id FROM expansions WHERE cat_id=$1 ORDER BY name;", cat.ID)
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	// Load the query's results in to new Expansions and append those to the slice.
	for rows.Next() {
		var exp Expansion
		err = rows.Scan(&exp.ID, &exp.Name, &exp.Text, &exp.CatID)
		exps = append(exps, exp)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Panicln(err)
	}

	// Return the populated slice of Expansions.
	return exps
}

// ReadPhrases finds all of the Phrases mapped to a given Expansion.
func ReadPhrases(exp *Expansion) []Phrase {

	// Create an empty slice of Phrases.
	var phrases []Phrase

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for Phrases matching the expansion's ID.
	rows, err := db.Query("SELECT id, name, exp_id FROM phrases WHERE exp_ID=$1 ORDER BY name;", exp.ID)
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	// Load the query's results in to the Phrase slice.
	for rows.Next() {
		var phrase Phrase
		err = rows.Scan(&phrase.ID, &phrase.Name, &phrase.ExpID)
		phrases = append(phrases, phrase)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Panicln(err)
	}

	// Return the populated slice of phrases.
	return phrases
}

// DeleteCategory deletes an entire category and all of its expansions and their mapped phrases.
func DeleteCategory(cat *Category) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Get a list of all expansions in that category.
	exps := ReadExpansions(cat)

	// Get a list of all phrases in each of those expansions.
	var phrases []Phrase
	for _, exp := range exps {
		phrases = append(phrases, ReadPhrases(&exp)...)
	}

	// Create a query that will delete all of the phrases.
	var delPhrase string
	for _, phrase := range phrases {
		delPhrase += "DELETE FROM phrases WHERE id = " + fmt.Sprint(phrase.ID) + "; "
	}

	// Delete the phrases.
	if _, err := db.Exec(delPhrase); err != nil {
		log.Panicln("Couldn't Delete Phrases: ", err)
	}

	// Create a query that will delete all of the expansions.
	var delExp string
	for _, exp := range exps {
		delExp += "DELETE FROM expansions WHERE id = " + fmt.Sprint(exp.ID) + "; "
	}

	// Delete the expansions.
	if _, err := db.Exec(delExp); err != nil {
		log.Panicln("Couldn't Delete Expansions: ", err)
	}

	// Delete the category.
	if _, err := db.Exec("DELETE FROM categories WHERE id = $1", cat.ID); err != nil {
		log.Panicln("Couldn't Delete Category: ", err)
	}
}

// DeleteExpansion deletes an expansion and all of its associated phrases.
func DeleteExpansion(exp *Expansion) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	//
	phrases := ReadPhrases(exp)

	// Create a query that will delete all of the phrases.
	var delPhrase string
	for _, phrase := range phrases {
		delPhrase += "DELETE FROM phrases WHERE id = " + fmt.Sprint(phrase.ID) + "; "
	}

	// Delete the phrases.
	if _, err := db.Exec(delPhrase); err != nil {
		log.Panicln("Couldn't Delete Phrases: ", err)
	}

	// Delete the expansion.
	if _, err := db.Exec("DELETE FROM expansions WHERE id = $1", exp.ID); err != nil {
		log.Panicln("Couldn't Delete Expansion: ", err)
	}
}

// DeletePhrase deletes a phrase from Gengar's database.
func DeletePhrase(phrase *Phrase) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Delete the phrase.
	if _, err := db.Exec("DELETE FROM phrases WHERE id = $1", phrase.ID); err != nil {
		log.Panicln("Couldn't Delete Phrase: ", err)
	}
}

// ReadExpanders reads expansions from the database and returns a slice of Expanders.
func ReadExpanders() []Expander {

	var exps []Expander

	// Get pointer to database connection.
	db := connectGGDB()
	defer db.Close()

	// Query the database for the expansions.
	rows, err := db.Query("SELECT phrases.name, text FROM phrases JOIN expansions ON phrases.exp_id = expansions.id;")
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	// Read through the results and populate exps with returned info.
	for rows.Next() {
		var exp Expander
		err = rows.Scan(&exp.Phrase, &exp.Expansion)
		exps = append(exps, exp)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Panicln(err)
	}

	// Return pointer to exps.
	return exps
}

// CreateTestDB creates a test database.
func CreateTestDB() {

	// Create a test category.
	cat := Category{Name: "category1", ID: 1}
	AddCategory(&cat)

	// Create some testing expansions.
	var exps []Expansion
	exps = append(exps, Expansion{Name: "Test 1", Text: "This is test 1!", ID: 1, CatID: 1})
	exps = append(exps, Expansion{Name: "Test 2", Text: "this is test 2?", ID: 2, CatID: 1})
	exps = append(exps, Expansion{Name: "Test 3", Text: "this is test 3!?@$", ID: 3, CatID: 1})

	// Create some test phrases.
	var phrases []Phrase
	phrases = append(phrases, Phrase{Name: "test1", ExpID: 1})
	phrases = append(phrases, Phrase{Name: "test2", ExpID: 2})
	phrases = append(phrases, Phrase{Name: "test3", ExpID: 3})

	// Wipe/create a blank gengar database.
	CleanSlate()

	// Insert each of our testing expansions and phrases.
	for _, exp := range exps {
		AddExpansion(&exp)
	}
	for _, exp := range exps {
		UpdateExpansionText(&exp)
	}
	for _, phrase := range phrases {
		AddPhrase(&phrase)
	}
}

// There shouldn't be much need to read all expansions or phrases, but...
// they're already written, so I have no reason to delete them yet.

// // ReadAllExpansions reads all of the available expansions from the database.
// func ReadAllExpansions() []Expansion {
//
// 	// Create an empty slice of Expansions to fill.
// 	var exps []Expansion
//
// 	// Connect to the database.
// 	db := connectGGDB()
// 	defer db.Close()
//
// 	// Query the database for all Expansions.
// 	rows, err := db.Query("SELECT id, name, text, cat_id FROM expansions;")
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer rows.Close()
//
// 	// Load the query's results in to new Expansions and append those to the slice.
// 	for rows.Next() {
// 		var exp Expansion
// 		err = rows.Scan(&exp.ID, &exp.Name, &exp.Text, &exp.CatID)
// 		exps = append(exps, exp)
// 	}
//
// 	// Die on error.
// 	if err = rows.Err(); err != nil {
// 		log.Panicln(err)
// 	}
//
// 	// Return the populated slice of expansions.
// 	return exps
// }
//
// // ReadAllPhrases reads all of the available phrases from the database.
// func ReadAllPhrases() []Phrase {
//
// 	// Create an empty slice of Phrases to fill.
// 	var phrases []Phrase
//
// 	// Connect to the database.
// 	db := connectGGDB()
// 	defer db.Close()
//
// 	// Get all of the available expansions.
// 	exps := ReadAllExpansions()
//
// 	// For each of the expansions, read the phrases and append them to the slice.
// 	for _, exp := range exps {
// 		newPhrases := ReadPhrases(&exp)
// 		phrases = append(phrases, newPhrases...)
// 	}
//
// 	// Return the populated slice of phrases.
// 	return phrases
// }
