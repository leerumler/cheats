package main

import (
	"flag"

	"github.com/leerumler/gengar/gengar"
	"github.com/leerumler/gengar/ggdb"
	"github.com/leerumler/gengar/ggui"
)

func main() {

	// Populate a test database.
	ggdb.CreateTestDB()

	//
	menu := flag.Bool("config", false, "Open gengar's configuration menu.")
	flag.Parse()

	//
	if !*menu {

		// Start gengar.
		gengar.ListenClosely()

	} else {

		// Open gengar's configuration menu.
		ggui.GengarMenu()
	}
}

// Oh no.  She is crazy and needs to go down.
// 								~ Uncle Iroh.
