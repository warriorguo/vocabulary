//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Run migrations
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS wordbook_entries (
			id BIGSERIAL PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL DEFAULT 'default',
			word VARCHAR(128) NOT NULL,
			short_definition TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, word)
		)`,
		`CREATE TABLE IF NOT EXISTS dictionary_cache (
			word VARCHAR(128) PRIMARY KEY,
			data JSONB NOT NULL,
			source VARCHAR(64) NOT NULL,
			fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMPTZ NOT NULL
		)`,
	}

	for _, m := range migrations {
		if _, err := pool.Exec(ctx, m); err != nil {
			t.Fatalf("failed to run migration: %v", err)
		}
	}

	cleanup := func() {
		pool.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return pool, cleanup
}

func TestRepositoryIntegration_WordbookCRUD(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()
	userID := "test-user"

	// Test empty wordbook
	entries, err := repo.GetWordbookEntries(ctx, userID)
	if err != nil {
		t.Fatalf("GetWordbookEntries failed: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}

	// Test add entry
	entry, err := repo.AddWordbookEntry(ctx, userID, "hello", "a greeting")
	if err != nil {
		t.Fatalf("AddWordbookEntry failed: %v", err)
	}
	if entry.Word != "hello" {
		t.Errorf("expected word 'hello', got '%s'", entry.Word)
	}
	if entry.ShortDefinition != "a greeting" {
		t.Errorf("expected definition 'a greeting', got '%s'", entry.ShortDefinition)
	}

	// Test get entries
	entries, err = repo.GetWordbookEntries(ctx, userID)
	if err != nil {
		t.Fatalf("GetWordbookEntries failed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}

	// Test word exists
	exists, err := repo.WordExistsInWordbook(ctx, userID, "hello")
	if err != nil {
		t.Fatalf("WordExistsInWordbook failed: %v", err)
	}
	if !exists {
		t.Error("expected word to exist")
	}

	// Test word doesn't exist
	exists, err = repo.WordExistsInWordbook(ctx, userID, "nonexistent")
	if err != nil {
		t.Fatalf("WordExistsInWordbook failed: %v", err)
	}
	if exists {
		t.Error("expected word to not exist")
	}

	// Test delete entry
	err = repo.DeleteWordbookEntry(ctx, userID, "hello")
	if err != nil {
		t.Fatalf("DeleteWordbookEntry failed: %v", err)
	}

	// Verify deletion
	entries, err = repo.GetWordbookEntries(ctx, userID)
	if err != nil {
		t.Fatalf("GetWordbookEntries failed: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after deletion, got %d", len(entries))
	}
}

func TestRepositoryIntegration_DictionaryCache(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()

	// Test cache miss
	cache, err := repo.GetCachedDictionary(ctx, "hello")
	if err != nil {
		t.Fatalf("GetCachedDictionary failed: %v", err)
	}
	if cache != nil {
		t.Error("expected nil cache for missing word")
	}

	// Test set cache
	testData := []byte(`{"word":"hello","meanings":[]}`)
	err = repo.SetCachedDictionary(ctx, "hello", testData, "test", 1*time.Hour)
	if err != nil {
		t.Fatalf("SetCachedDictionary failed: %v", err)
	}

	// Test cache hit
	cache, err = repo.GetCachedDictionary(ctx, "hello")
	if err != nil {
		t.Fatalf("GetCachedDictionary failed: %v", err)
	}
	if cache == nil {
		t.Fatal("expected cache to exist")
	}
	if string(cache.Data) != string(testData) {
		t.Errorf("cache data mismatch: got %s, want %s", cache.Data, testData)
	}

	// Test cache update (upsert)
	newData := []byte(`{"word":"hello","meanings":[{"partOfSpeech":"noun"}]}`)
	err = repo.SetCachedDictionary(ctx, "hello", newData, "test", 1*time.Hour)
	if err != nil {
		t.Fatalf("SetCachedDictionary update failed: %v", err)
	}

	cache, err = repo.GetCachedDictionary(ctx, "hello")
	if err != nil {
		t.Fatalf("GetCachedDictionary failed: %v", err)
	}
	if string(cache.Data) != string(newData) {
		t.Errorf("updated cache data mismatch: got %s, want %s", cache.Data, newData)
	}
}

func TestRepositoryIntegration_UpsertWordbook(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := New(pool)
	ctx := context.Background()
	userID := "test-user"

	// Add entry
	_, err := repo.AddWordbookEntry(ctx, userID, "test", "original definition")
	if err != nil {
		t.Fatalf("AddWordbookEntry failed: %v", err)
	}

	// Update entry (same word)
	entry, err := repo.AddWordbookEntry(ctx, userID, "test", "updated definition")
	if err != nil {
		t.Fatalf("AddWordbookEntry update failed: %v", err)
	}
	if entry.ShortDefinition != "updated definition" {
		t.Errorf("expected updated definition, got '%s'", entry.ShortDefinition)
	}

	// Verify only one entry exists
	entries, err := repo.GetWordbookEntries(ctx, userID)
	if err != nil {
		t.Fatalf("GetWordbookEntries failed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}
