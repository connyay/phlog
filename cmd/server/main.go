package main

import (
	"log"
	"os"

	"github.com/connyay/phlog/server"
	"github.com/connyay/phlog/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	store := &store.Mem{}
	log.Fatal(server.ListenHTTP("0.0.0.0:"+port, store))
}
