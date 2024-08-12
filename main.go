package main

import (
	"log"

	"github.com/AdamDomagalsky/gon/api"
)

func main() {
	config := &api.Config{
		Port: "8080",
	}

	server := api.NewServer(config)
	if err := server.Start(); err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
