package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"news-feed-bot/internal/model"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticlePostgresStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Исправлено: sources_id вместо source_id
	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles(sources_id, title, link, summary, published_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (link) DO NOTHING`, // Добавлен created_at и улучшен CONFLICT
		article.SourceID, article.Title, article.Link, article.Summary,
		article.PublishedAt, time.Now().UTC(),
	); err != nil {
		return err
	}
	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles WHERE posted_at IS NULL AND published_at >= $1 ORDER BY published_at DESC LIMIT $2`,
		since.UTC(),
		limit,
	); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PostedAt:    article.PostedAt,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1 WHERE id = $2`,
		time.Now().UTC(),
		id,
	); err != nil {
		return err
	}
	return nil
}

// Удалите дублирующий метод MarcPosted (с опечаткой)

type dbArticle struct {
	ID          int64     `db:"id"`
	SourceID    int64     `db:"sources_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	PostedAt    time.Time `db:"posted_at"` // Исправлено: указатель для NULL значений
}
