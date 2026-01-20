package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/warriorguo/vocabulary/internal/handlers"
	"github.com/warriorguo/vocabulary/internal/repository"
	"github.com/warriorguo/vocabulary/internal/services"
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/vocabulary?sslmode=disable"
	}

	// Get server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Run migrations
	if err := runMigrations(ctx, pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize layers
	repo := repository.New(pool)
	dictSvc := services.NewDictionaryService(repo)
	handler := handlers.New(repo, dictSvc)

	// Setup Gin
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	handler.SetupRoutes(r)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS wordbook_entries (
			id BIGSERIAL PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL DEFAULT 'default',
			word VARCHAR(128) NOT NULL,
			short_definition TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, word)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_wordbook_user_created ON wordbook_entries(user_id, created_at DESC)`,
		`CREATE TABLE IF NOT EXISTS dictionary_cache (
			word VARCHAR(128) PRIMARY KEY,
			data JSONB NOT NULL,
			source VARCHAR(64) NOT NULL,
			fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_cache_expires ON dictionary_cache(expires_at)`,
	}

	for i, migration := range migrations {
		if _, err := pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}
