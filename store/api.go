package store

import (
	"bytes"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Store interface {
	AddBlob(data []byte, ext string) (string, error)
	GetBlobByRef(ref string) (blob io.Reader, ext string, err error)
	AddPost(Post) error
	GetPosts(category string) ([]Post, error)
	GetPostByID(id string) (Post, error)
}

type Post struct {
	ID       string
	Title    string
	Category string
	Blobs    []string
}

type blob struct {
	data []byte
	ext  string
}

type Mem struct {
	posts map[string]Post
	blobs map[string]blob
}

func (m *Mem) AddPost(p Post) error {
	p.ID = uuid.NewString()
	if m.posts == nil {
		m.posts = map[string]Post{}
	}
	m.posts[p.ID] = p
	return nil
}
func (m *Mem) GetPosts(category string) ([]Post, error) {
	posts := make([]Post, 0, len(m.posts))
	// This range order is not stable.
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	return posts, nil
}
func (m *Mem) GetPostByID(id string) (Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return Post{}, fmt.Errorf("post %q not found", id)
	}
	return p, nil
}
func (m *Mem) AddBlob(data []byte, ext string) (string, error) {
	if m.blobs == nil {
		m.blobs = map[string]blob{}
	}
	ref := uuid.NewString()
	m.blobs[ref] = blob{data, ext}
	return ref, nil
}
func (m *Mem) GetBlobByRef(ref string) (io.Reader, string, error) {
	blob, ok := m.blobs[ref]
	if !ok {
		return nil, "", fmt.Errorf("blob %q not found", ref)
	}
	return bytes.NewReader(blob.data), blob.ext, nil
}
