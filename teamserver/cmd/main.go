package main

import (
	"errors"
	"flag"
	"log"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/http"
)

func checkIfArgDefined(argName string) (bool, error) {
	if !flag.Parsed() {
		return false, errors.New("flag.Parse() has not been yet called")
	}

	var found bool
	flag.Visit(func(flag *flag.Flag) {
		if flag.Name == argName {
			found = true
		}
	})

	return found, nil
}

func main() {
	port := flag.Int("port", 8080, "Specify the port that the server will run on")
	// Small hack to make the default look sane
	portFlag := flag.Lookup("port")
	portFlag.DefValue = "8080/8443"

	https := flag.Bool("https", false, "If https is enabled then both cert and key are required")
	certificate := flag.String("cert", "", "Specify the https certificate")
	key := flag.String("key", "", "Specify encryption key")

	flag.Parse()

	if *https {
		if len(*certificate) == 0 || len(*key) == 0 {
			log.Fatal("Key or certificate was not supplied")
		}

		portExplicitlyDefined, err := checkIfArgDefined("port")
		if err != nil {
			log.Fatal(err)
		}

		if !portExplicitlyDefined {
			*port = 8443
		}
	}

	db, err := database.NewDatabase("teamserver.db")
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServer(*port, db)
	if *https {
		server.ListenAndServeTLS(*certificate, *key)
	} else {
		server.ListenAndServe()
	}
}
