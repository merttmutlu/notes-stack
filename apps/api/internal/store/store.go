package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/merttmutlu/notes-stack/apps/api/internal/note"
)

type Dependencies struct {
	DatabaseURL string
}

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, deps Dependencies) (*Store, error) {
	if deps.DatabaseURL == "" {
		return nil, errors.New("database url is not configured")
	}

	pool, err := pgxpool.New(ctx, deps.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return &Store{
		pool: pool,
	}, nil
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.pool == nil {
		return errors.New("database pool is not initialized")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return s.pool.Ping(pingCtx)
}

func (s *Store) ListNotes(ctx context.Context) ([]note.Note, error) {
	if s == nil || s.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(queryCtx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]note.Note, 0)
	for rows.Next() {
		var item note.Note

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		notes = append(notes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Store) GetNote(ctx context.Context, id string) (note.Note, error) {
	if s == nil || s.pool == nil {
		return note.Note{}, errors.New("database pool is not initialized")
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var item note.Note
	err := s.pool.QueryRow(queryCtx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE id = $1
	`, strings.TrimSpace(id)).Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return note.Note{}, err
	}

	return item, nil
}

func (s *Store) CreateNote(ctx context.Context, req note.CreateRequest) (note.Note, error) {
	if s == nil || s.pool == nil {
		return note.Note{}, errors.New("database pool is not initialized")
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var item note.Note
	err := s.pool.QueryRow(queryCtx, `
		INSERT INTO notes (id, title, content, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
		RETURNING id, title, content, created_at, updated_at
	`,
		strings.TrimSpace(req.Title),
		strings.TrimSpace(req.Content),
	).Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return note.Note{}, err
	}

	return item, nil
}

func (s *Store) UpdateNote(ctx context.Context, id string, req note.UpdateRequest) (note.Note, error) {
	if s == nil || s.pool == nil {
		return note.Note{}, errors.New("database pool is not initialized")
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var item note.Note
	err := s.pool.QueryRow(queryCtx, `
		UPDATE notes
		SET title = $2, content = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, content, created_at, updated_at
	`,
		strings.TrimSpace(id),
		strings.TrimSpace(req.Title),
		strings.TrimSpace(req.Content),
	).Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return note.Note{}, err
	}

	return item, nil
}

func (s *Store) DeleteNote(ctx context.Context, id string) (bool, error) {
	if s == nil || s.pool == nil {
		return false, errors.New("database pool is not initialized")
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tag, err := s.pool.Exec(queryCtx, `
		DELETE FROM notes
		WHERE id = $1
	`, strings.TrimSpace(id))
	if err != nil {
		return false, fmt.Errorf("delete note: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (s *Store) Close() {
	if s == nil || s.pool == nil {
		return
	}

	s.pool.Close()
}
