package main

import (
	"flag"
	"log"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/config"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/http"
)

func main() {
	configPath := flag.String("config", "", "Explicitly specify conifg location")
	flag.Parse()

	var serverConfig *config.Config
	var err error
	if len(*configPath) != 0 {
		serverConfig, err = config.ParseConfig(*configPath)
	} else {
		serverConfig, err = config.ParseDefaultConfig()
	}
	if err != nil {
		log.Fatal(err)
	}

	err = config.ValidateConfig(serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase("teamserver.db")
	if err != nil {
		log.Fatal(err)
	}

	server, err := http.NewServer(serverConfig, db)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.ListenAndServe())
}
