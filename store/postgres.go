package store

import (
	"context"
	"embed"
	"io"
	"strconv"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Create migration source
//go:embed postgres/migrations
var embedFS embed.FS

var embedSource = &migration.EmbedMigrationSource{
	EmbedFS: embedFS,
	Dir:     "postgres/migrations",
}

// NewPG returns an initialized Store backed by postgres based on the given
// connection string.
func NewPG(dsn string) (Store, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	s := &pgstore{dsn, pool}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

type pgstore struct {
	dsn  string
	pool *pgxpool.Pool
}

func (pg *pgstore) AddBlob(data []byte, ext string) (string, error) {
	panic("not impl")
}
func (pg *pgstore) GetBlobByRef(ref string) (blob io.Reader, ext string, err error) {
	panic("not impl")
}

func (pg *pgstore) AddPost(post Post) error {
	pg.pool.BeginTx(context.Background(), pgx.TxOptions{})
	_, err := pg.pool.Exec(context.Background(), `
	INSERT INTO posts (
		title
	)
	VALUES ($1)
	`,
		post.Title,
	)
	if err != nil {
		return err
	}
	return nil
}

func (pg *pgstore) GetPosts(category string) (posts []Post, err error) {
	sql := `
	SELECT
		id, title
	FROM posts`
	rows, err := pg.pool.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id   int
			post Post
		)
		if err := rows.Scan(&id, &post.Title); err != nil {
			return nil, err
		}
		post.ID = strconv.Itoa(id)
		posts = append(posts, post)
	}
	return posts, nil
}

func (pg *pgstore) GetPostByID(id string) (post Post, err error) {
	sql := `
	SELECT
		id, title
	FROM posts where id = $1`
	row := pg.pool.QueryRow(context.Background(), sql, id)
	var intID int
	if err := row.Scan(&intID, &post.Title); err != nil {
		return post, err
	}
	post.ID = strconv.Itoa(intID)
	return post, nil
}

func (pg *pgstore) migrate() error {
	driver, err := postgres.New(pg.dsn)
	if err != nil {
		return err
	}
	defer driver.Close()
	_, err = migration.Migrate(driver, embedSource, migration.Up, 0)
	if err != nil {
		return err
	}
	return nil
}
