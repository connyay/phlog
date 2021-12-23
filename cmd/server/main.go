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
	postStore := &store.Mem{}
	// postStore.AddPost(store.Post{
	// 	Title:       "Hello World!",
	// 	Attachment:  `iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII`,
	// 	ContentType: "image/png",
	// })
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
