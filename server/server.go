package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/connyay/phlog/store"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func ListenHTTP(addr string, postStore store.Store) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
