package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// setup database
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// setup server
	portNumber := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server := NewApiServer(portNumber, store)
	server.Run()
}
