package server

import (
	"embed"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/connyay/phlog/store"
)

//go:embed templates/*
var resources embed.FS

func ListenHTTP(addr string, storage store.Store) error {
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseFS(resources, "templates/*")))
	r.GET("/p/:id", func(c *gin.Context) {
		post, err := storage.GetPostByID(c.Param("id"))
		if err != nil {
			c.Error(err)
			return
		}
		data := map[string]interface{}{
			"Post": post,
		}
		c.HTML(http.StatusOK, "post.html.tmpl", data)
	})
	r.GET("/b/:id", func(c *gin.Context) {
		blob, ext, err := storage.GetBlobByRef(c.Param("id"))
		if err != nil {
			c.Error(err)
			return
		}
		c.Header("Content-Type", mime.TypeByExtension(ext))
		_, err = io.Copy(c.Writer, blob)
		if err != nil {
			c.Error(err)
			return
		}
	})
	r.GET("/", func(c *gin.Context) {
		posts, err := storage.GetPosts("")
		if err != nil {
			c.Error(err)
			return
		}
		data := map[string]interface{}{
			"Posts": posts,
		}
		c.HTML(http.StatusOK, "posts.html.tmpl", data)
	})
	log.Printf("listening on http://%s", addr)
	return r.Run(addr)
}
