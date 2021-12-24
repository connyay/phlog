package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/connyay/phlog/mail"
	"github.com/connyay/phlog/server"
	"github.com/connyay/phlog/store"
)

func main() {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "0.0.0.0:8080"
	}
	mailAddr := os.Getenv("MAIL_ADDR")
	if mailAddr == "" {
		mailAddr = ":8081"
	}
	var (
		storage store.Store
		err     error
	)
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		log.Println("Using PG store")
		storage, err = store.NewPG(dsn)
		if err != nil {
			log.Fatalf("Initializing store %v", err)
		}
	} else {
		log.Println("Using mem store")
		storage = &store.Mem{}
	}
	if _, seed := os.LookupEnv("SEED_DB"); seed {
		log.Println("Seeding data")
		if err := seedData(storage); err != nil {
			log.Fatalf("seeding data %v", err)
		}
	}

	mailReader := mail.RawDogReader{
		ListenAddr: mailAddr,
	}

	go func() {
		for msg := range mailReader.Messages() {
			var post store.Post
			for _, part := range msg.Parts {
				if strings.HasPrefix(part.Type, "text/plain") {
					post.Title = string(part.Body)
				}
				if strings.HasPrefix(part.Type, "image/") {
					ref, err := storage.AddBlob(part.Body, part.Type)
					if err != nil {
						// FIXME(cjh): Handle this somehow. Retry?
						panic(fmt.Errorf("failed storing %q", err))
					}
					post.Blobs = append(post.Blobs, ref)
				}
			}
			log.Printf("Storing post title=%q blobs=%q", post.Title, post.Blobs)
			err := storage.AddPost(post)
			if err != nil {
				// FIXME(cjh): Handle this somehow. Retry?
				log.Printf("Failed storing post %v %v", post, err)
			}
		}
	}()
	log.Fatal(server.ListenHTTP(httpAddr, storage))
}

func seedData(storage store.Store) error {
	blob, err := base64.RawStdEncoding.DecodeString(`iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII`)
	if err != nil {
		return err
	}
	ref, err := storage.AddBlob(blob, ".png")
	if err != nil {
		return err
	}
	err = storage.AddPost(store.Post{
		Title: "Hello World!",
		Blobs: []string{ref},
	})
	if err != nil {
		return err
	}
	return nil
}
