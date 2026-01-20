package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warriorguo/vocabulary/internal/models"
	"github.com/warriorguo/vocabulary/internal/repository"
	"github.com/warriorguo/vocabulary/internal/services"
)

const defaultUserID = "default"

type Handler struct {
	repo    *repository.Repository
	dictSvc *services.DictionaryService
}

func New(repo *repository.Repository, dictSvc *services.DictionaryService) *Handler {
	return &Handler{
		repo:    repo,
		dictSvc: dictSvc,
	}
}

// LookupWord handles GET /api/dict?word={word}
func (h *Handler) LookupWord(c *gin.Context) {
	word := c.Query("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word parameter is required"})
		return
	}

	entry, err := h.dictSvc.LookupWord(c.Request.Context(), word)
	if err != nil {
		if err.Error() == "word not found: "+word {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if word is in wordbook
	inWordbook, err := h.repo.WordExistsInWordbook(c.Request.Context(), defaultUserID, word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entry":       entry,
		"in_wordbook": inWordbook,
	})
}

// GetWordbook handles GET /api/wordbook
func (h *Handler) GetWordbook(c *gin.Context) {
	entries, err := h.repo.GetWordbookEntries(c.Request.Context(), defaultUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if entries == nil {
		entries = []models.WordbookEntry{}
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

// AddToWordbook handles POST /api/wordbook
func (h *Handler) AddToWordbook(c *gin.Context) {
	var req models.AddWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.repo.AddWordbookEntry(c.Request.Context(), defaultUserID, req.Word, req.ShortDefinition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"entry": entry})
}

// RemoveFromWordbook handles DELETE /api/wordbook/:word
func (h *Handler) RemoveFromWordbook(c *gin.Context) {
	word := c.Param("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word parameter is required"})
		return
	}

	err := h.repo.DeleteWordbookEntry(c.Request.Context(), defaultUserID, word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "word removed from wordbook"})
}

// SetupRoutes configures all API routes
func (h *Handler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/dict", h.LookupWord)
		api.GET("/wordbook", h.GetWordbook)
		api.POST("/wordbook", h.AddToWordbook)
		api.DELETE("/wordbook/:word", h.RemoveFromWordbook)
	}
}
