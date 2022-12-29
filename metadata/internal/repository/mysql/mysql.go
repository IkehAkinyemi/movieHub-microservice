package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"moviehub.com/metadata/internal/repository"
	"moviehub.com/metadata/pkg/model"
)

// Repository defines a MySQL-based movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Get retrieves movie metadata by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	query := "SELECT title, description, director FROM movies WHERE id = ?"

	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}

		return nil, err
	}

	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	stmt := "INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, stmt, id, metadata.Title, metadata.Description, metadata.Director)
	return err
}
