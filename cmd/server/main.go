package main

import (
	"log"
	"os"

	"github.com/connyay/phlog/pkg/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(server.Listen("localhost:" + port))
}
