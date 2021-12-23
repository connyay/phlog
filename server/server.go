package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/connyay/phlog/store"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func ListenHTTP(addr string, postStore store.Store) error {
	postStore.AddPost(store.Post{
		Title:       "Hello World!",
		Attachment:  `iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII`,
		ContentType: "image/png",
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Env": os.Environ(),
		}

		t.ExecuteTemplate(w, "index.html.tmpl", data)
	})
	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		posts, err := postStore.GetPosts("")
		if err != nil {
			http.Error(w, "failed getting posts", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Posts": posts,
		}
		t.ExecuteTemplate(w, "posts.html.tmpl", data)
	})

	log.Printf("listening on http://%s", addr)
	return http.ListenAndServe(addr, nil)
}
