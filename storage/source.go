package storage

import (
	"context"
	"time"

	"github.com/Sweetheart11/tgbot/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func NewSourceStorage(db *sqlx.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{db: db}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource

	rows, err := conn.QueryContext(ctx, "SELECT * FROM sources")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var source dbSource
		if err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.FeedURL,
			&source.CreatedAt,
		); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source(source)
	}), nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	row := conn.QueryRowContext(ctx, "SELECT * FROM sources WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	if err := row.Scan(
		&source.ID,
		&source.Name,
		&source.FeedURL,
		&source.CreatedAt,
	); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int64
	row := conn.QueryRowContext(
		ctx,
		"INSERT INTO sources (name, feed_url, created_at) VALUES ($1, $2, $3) returning id",
		source.Name,
		source.FeedURL,
		source.CreatedAt,
	)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "DELETE FROM sources WHERE id = $1", id)
	return err
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}
