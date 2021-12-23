package main

import (
	"encoding/base64"
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
		postStore store.Store
		err       error
	)
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		log.Println("Using PG store")
		postStore, err = store.NewPG(dsn)
		if err != nil {
			log.Fatalf("Initializing store %v", err)
		}
	} else {
		log.Println("Using mem store")
		postStore = &store.Mem{}
	}
	// err = postStore.AddPost(store.Post{
	// 	Title:       "Hello World!",
	// 	Attachment:  `iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII`,
	// 	ContentType: "image/png",
	// })
	// if err != nil {
	// 	log.Fatalf("Writing to store store %v", err)
	// }
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
					post.Attachment = base64.RawStdEncoding.EncodeToString(part.Body)
					post.ContentType = part.Type
				}
			}
			log.Printf("Storing post title=%q type=%q len=%d", post.Title, post.ContentType, len(post.Attachment))
			err := postStore.AddPost(post)
			if err != nil {
				log.Printf("Failed storing post %v %v", post, err)
			}
		}
	}()
	log.Fatal(server.ListenHTTP(httpAddr, postStore))
}
