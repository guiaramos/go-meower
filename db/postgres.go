package db

import (
	"context"
	"database/sql"
	"github.com/guiaramos/go-meower/schema"
)

var (
	queryInsertMeow = "INSERT INTO meows(id, body, created_at) VALUES($1, $2, $3);"
	queryListMeows  = "SELECT * FROM meows ORDER BY id DESC OFFSET $1 LIMIT $2;"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r PostgresRepository) Close() {
	r.db.Close()
}

func (r PostgresRepository) InsertMeow(ctx context.Context, meow schema.Meow) error {
	_, err := r.db.Query(queryInsertMeow, meow.ID, meow.Body, meow.CreatedAt)
	return err
}

func (r PostgresRepository) ListMeows(ctx context.Context, skip uint64, take uint64) ([]schema.Meow, error) {
	rows, err := r.db.Query(queryListMeows, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meows := []schema.Meow{}
	for rows.Next() {
		meow := schema.Meow{}
		if err = rows.Scan(&meow.ID, &meow.Body, &meow.CreatedAt); err == nil {
			meows = append(meows, meow)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return meows, nil
}
