package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/warriorguo/vocabulary/internal/models"
)

// Mock repository for testing
type mockRepo struct {
	entries      []models.WordbookEntry
	wordExists   bool
	returnError  error
}

func (m *mockRepo) GetWordbookEntries(ctx context.Context, userID string) ([]models.WordbookEntry, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.entries, nil
}

func (m *mockRepo) AddWordbookEntry(ctx context.Context, userID, word, shortDef string) (*models.WordbookEntry, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	entry := &models.WordbookEntry{
		ID:              1,
		UserID:          userID,
		Word:            word,
		ShortDefinition: shortDef,
		CreatedAt:       time.Now(),
	}
	m.entries = append(m.entries, *entry)
	return entry, nil
}

func (m *mockRepo) DeleteWordbookEntry(ctx context.Context, userID, word string) error {
	return m.returnError
}

func (m *mockRepo) WordExistsInWordbook(ctx context.Context, userID, word string) (bool, error) {
	if m.returnError != nil {
		return false, m.returnError
	}
	return m.wordExists, nil
}

// Mock dictionary service
type mockDictSvc struct {
	entry       *models.DictionaryEntry
	returnError error
}

func (m *mockDictSvc) LookupWord(ctx context.Context, word string) (*models.DictionaryEntry, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.entry, nil
}

// Test handler with mocks
type testHandler struct {
	repo    *mockRepo
	dictSvc *mockDictSvc
}

func newTestHandler() *testHandler {
	return &testHandler{
		repo:    &mockRepo{entries: []models.WordbookEntry{}},
		dictSvc: &mockDictSvc{},
	}
}

func setupTestRouter(th *testHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api")
	{
		api.GET("/wordbook", func(c *gin.Context) {
			entries, err := th.repo.GetWordbookEntries(c.Request.Context(), defaultUserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if entries == nil {
				entries = []models.WordbookEntry{}
			}
			c.JSON(http.StatusOK, gin.H{"entries": entries})
		})

		api.POST("/wordbook", func(c *gin.Context) {
			var req models.AddWordRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			entry, err := th.repo.AddWordbookEntry(c.Request.Context(), defaultUserID, req.Word, req.ShortDefinition)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"entry": entry})
		})

		api.DELETE("/wordbook/:word", func(c *gin.Context) {
			word := c.Param("word")
			if word == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "word parameter is required"})
				return
			}
			err := th.repo.DeleteWordbookEntry(c.Request.Context(), defaultUserID, word)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "word removed from wordbook"})
		})
	}

	return r
}

func TestGetWordbook(t *testing.T) {
	th := newTestHandler()
	th.repo.entries = []models.WordbookEntry{
		{ID: 1, UserID: "default", Word: "hello", ShortDefinition: "a greeting", CreatedAt: time.Now()},
		{ID: 2, UserID: "default", Word: "world", ShortDefinition: "the earth", CreatedAt: time.Now()},
	}

	router := setupTestRouter(th)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/wordbook", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Entries []models.WordbookEntry `json:"entries"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(response.Entries))
	}
}

func TestGetWordbookEmpty(t *testing.T) {
	th := newTestHandler()
	router := setupTestRouter(th)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/wordbook", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Entries []models.WordbookEntry `json:"entries"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(response.Entries))
	}
}

func TestAddToWordbook(t *testing.T) {
	th := newTestHandler()
	router := setupTestRouter(th)

	body := `{"word": "hello", "short_definition": "a greeting"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/wordbook", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response struct {
		Entry models.WordbookEntry `json:"entry"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Entry.Word != "hello" {
		t.Errorf("expected word 'hello', got '%s'", response.Entry.Word)
	}
}

func TestAddToWordbookBadRequest(t *testing.T) {
	th := newTestHandler()
	router := setupTestRouter(th)

	tests := []struct {
		name string
		body string
	}{
		{"missing word", `{"short_definition": "a greeting"}`},
		{"missing definition", `{"word": "hello"}`},
		{"invalid json", `{invalid`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/wordbook", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestDeleteFromWordbook(t *testing.T) {
	th := newTestHandler()
	router := setupTestRouter(th)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/wordbook/hello", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
