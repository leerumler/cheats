package main

import (
	"flag"

	"github.com/leerumler/gengar/gengar"
	"github.com/leerumler/gengar/ggdb"
	"github.com/leerumler/gengar/ggui"
)

func main() {

	//
	listen := flag.Bool("run", false, "Tell gengar to start listening.")
	reset := flag.Bool("reset", false, "Reset Gengar's database. This will wipe out everything.")
	flag.Parse()

	//
	switch {
	case *listen:

		// Start gengar.
		gengar.ListenClosely()

	case *reset:

		// Reset the test database.
		ggdb.CreateTestDB()

	default:

		// Open gengar's configuration menu.
		ggui.GengarMenu()

	}
}

// Oh no.  She is crazy and needs to go down.
// 								~ Uncle Iroh.
