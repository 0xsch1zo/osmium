package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/handlers"
)

func main() {
	port := flag.Int("port", 8080, "Specify the port that the server will run on.")

	http.HandleFunc("POST /register", handlers.Register)

	log.Print("Starting listening on port: " + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
