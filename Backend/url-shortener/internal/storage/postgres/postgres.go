package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.Ping()
	fmt.Println(err)
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) (int64, error) {
	const op = "storage.postgres.urlToSave"

	stmt, err := s.db.Prepare(`insert into url(url, alias) values ($1, $2) returning id`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	var id int64

	err = stmt.QueryRowContext(ctx, url, alias).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) URL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.URL"

	stmt, err := s.db.Prepare(`select url from url where alias=$1`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	defer stmt.Close()

	var url string
	row := stmt.QueryRowContext(ctx, alias)
	err = row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.postgres.DeleteURL"

	stmt, err := s.db.Prepare(`delete from url where alias=$1`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_ = stmt.QueryRowContext(ctx, alias)
	return nil
}
