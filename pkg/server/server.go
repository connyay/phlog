package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func Listen(addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Env": os.Environ(),
		}

		t.ExecuteTemplate(w, "index.html.tmpl", data)
	})

	log.Printf("listening on http://%s", addr)
	return http.ListenAndServe(addr, nil)
}
