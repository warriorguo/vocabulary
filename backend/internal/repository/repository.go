package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/warriorguo/vocabulary/internal/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Wordbook operations

func (r *Repository) GetWordbookEntries(ctx context.Context, userID string) ([]models.WordbookEntry, error) {
	query := `
		SELECT id, user_id, word, short_definition, created_at
		FROM wordbook_entries
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.WordbookEntry
	for rows.Next() {
		var entry models.WordbookEntry
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Word, &entry.ShortDefinition, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (r *Repository) AddWordbookEntry(ctx context.Context, userID, word, shortDef string) (*models.WordbookEntry, error) {
	query := `
		INSERT INTO wordbook_entries (user_id, word, short_definition)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, word) DO UPDATE SET short_definition = EXCLUDED.short_definition
		RETURNING id, user_id, word, short_definition, created_at`

	var entry models.WordbookEntry
	err := r.db.QueryRow(ctx, query, userID, word, shortDef).Scan(
		&entry.ID, &entry.UserID, &entry.Word, &entry.ShortDefinition, &entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *Repository) DeleteWordbookEntry(ctx context.Context, userID, word string) error {
	query := `DELETE FROM wordbook_entries WHERE user_id = $1 AND word = $2`
	_, err := r.db.Exec(ctx, query, userID, word)
	return err
}

func (r *Repository) WordExistsInWordbook(ctx context.Context, userID, word string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM wordbook_entries WHERE user_id = $1 AND word = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, word).Scan(&exists)
	return exists, err
}

// Cache operations

func (r *Repository) GetCachedDictionary(ctx context.Context, word string) (*models.DictionaryCache, error) {
	query := `
		SELECT word, data, source, fetched_at, expires_at
		FROM dictionary_cache
		WHERE word = $1 AND expires_at > NOW()`

	var cache models.DictionaryCache
	err := r.db.QueryRow(ctx, query, word).Scan(
		&cache.Word, &cache.Data, &cache.Source, &cache.FetchedAt, &cache.ExpiresAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

func (r *Repository) SetCachedDictionary(ctx context.Context, word string, data []byte, source string, ttl time.Duration) error {
	query := `
		INSERT INTO dictionary_cache (word, data, source, fetched_at, expires_at)
		VALUES ($1, $2, $3, NOW(), NOW() + $4::interval)
		ON CONFLICT (word) DO UPDATE SET
			data = EXCLUDED.data,
			source = EXCLUDED.source,
			fetched_at = EXCLUDED.fetched_at,
			expires_at = EXCLUDED.expires_at`

	_, err := r.db.Exec(ctx, query, word, data, source, ttl.String())
	return err
}

func (r *Repository) CleanExpiredCache(ctx context.Context) error {
	query := `DELETE FROM dictionary_cache WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, query)
	return err
}
