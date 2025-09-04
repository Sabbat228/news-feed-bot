package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"news-feed-bot/internal/model"
	"time"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticlePostgresStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {}
func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
}
func (s *ArticlePostgresStorage) MarcPosted(ctx context.Context, id int64) error {}

type dbArticle struct {
	ID          int64     `db:"id"`
	SourceId    int64     `db:"source_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
}
