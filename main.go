package main

import (
	"flag"
	"os"

	"github.com/leerumler/gengar/gengar"
	"github.com/leerumler/gengar/ggdb"
	"github.com/leerumler/gengar/ggui"
)

func main() {

	//
	listen := flag.Bool("run", false, "Tell gengar to start listening.")
	scary := flag.Bool("scary", false, "Report all tracked events.  Kinda scary.")
	reset := flag.Bool("reset", false, "Reset Gengar's database. This will wipe out everything, so be careful.")
	flag.Parse()

	// Check if gengar's DB exists, and create it if it doesn't.
	dbfile := ggdb.FindGGDB()
	if _, err := os.Stat(*dbfile); os.IsNotExist(err) {
		ggdb.CleanSlate()
	}

	//
	if *scary {
		gengar.Scary = true
		*listen = true
	}

	//
	switch {
	case *listen:

		// Start gengar.
		gengar.ListenClosely()

	case *reset:

		// Create a blank database.
		ggdb.CleanSlate()

	default:

		// Open gengar's configuration menu.
		ggui.GengarMenu()

	}
}

// Oh no.  She is crazy and needs to go down.
// 								~ Uncle Iroh.
