package store

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type Store interface {
	AddBlob(data []byte, ext string) (string, error)
	GetBlobByRef(ref string) (blob io.ReadCloser, ext string, err error)
	AddPost(Post) error
	GetPosts(category string) ([]Post, error)
	GetPostByID(id string) (Post, error)
}

func New() (storage Store, err error) {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		log.Println("Using PG store")
		storage, err = NewPG(dsn)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Using mem store")
		storage = &Mem{}
	}
	if _, s3Creds := os.LookupEnv("AWS_S3_SECRET"); s3Creds {
		log.Println("Using s3 blob store")
		storage = S3BlobStore{storage}
	}
	return storage, nil
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

func (m *Mem) GetBlobByRef(ref string) (io.ReadCloser, string, error) {
	blob, ok := m.blobs[ref]
	if !ok {
		return nil, "", fmt.Errorf("blob %q not found", ref)
	}
	return io.NopCloser(bytes.NewReader(blob.data)), blob.ext, nil
}
