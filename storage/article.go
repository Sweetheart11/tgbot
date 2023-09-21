package storage

import (
	"context"
	"time"

	"github.com/Sweetheart11/tgbot/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type ArticlePostgresStore struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStore {
	return &ArticlePostgresStore{db: db}
}

func (s *ArticlePostgresStore) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at )
    VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStore) AllNotPosted(
	ctx context.Context,
	since time.Time,
	limit uint64,
) ([]model.Article, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var dbArticles []dbArticle
	rows, err := conn.QueryContext(
		ctx,
		`SELECT id, source_id, title, link, summary, published_at, created_at FROM articles
    WHERE posted_at IS NULL AND published_at >= $1::timestamp ORDER BY published_at DESC LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dbArticle dbArticle
		if err := rows.Scan(
			&dbArticle.ID,
			&dbArticle.SourceID,
			&dbArticle.Title,
			&dbArticle.Link,
			&dbArticle.Summary,
			&dbArticle.PublishedAt,
			&dbArticle.CreatedAt,
		); err != nil {
			return nil, err
		}
		dbArticles = append(dbArticles, dbArticle)
	}

	return lo.Map(dbArticles, func(dbArticle dbArticle, _ int) model.Article {
		return model.Article{
			ID:          dbArticle.ID,
			SourceID:    dbArticle.SourceID,
			Title:       dbArticle.Title,
			Link:        dbArticle.Link,
			Summary:     dbArticle.Summary,
			PublishedAt: dbArticle.PublishedAt,
			CreatedAt:   dbArticle.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStore) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		"UPDATE articles SET posted_at = $1::timestamp WHERE id = $2",
		time.Now().UTC().Format(time.RFC3339),
		id)

	return err
}

type dbArticle struct {
	ID          int64     `json:"id"`
	SourceID    int64     `json:"source_id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Summary     string    `json:"summary"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	PostedAt    time.Time `json:"posted_at"`
}
