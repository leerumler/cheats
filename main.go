package main

import (
	"github.com/leerumler/gengar/gengar"
	"github.com/leerumler/gengar/ggdb"
)

func main() {

	// Populate a test database.
	ggdb.CreateTestDB()

	// Start gengar.
	gengar.ListenClosely()

}
