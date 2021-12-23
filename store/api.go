package store

import "github.com/google/uuid"

type Store interface {
	AddPost(Post) error
	GetPosts(category string) ([]Post, error)
}

type Post struct {
	ID          string
	Title       string
	Category    string
	Attachment  string
	ContentType string
}

type Mem struct {
	posts map[string]Post
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
