package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/warriorguo/vocabulary/internal/models"
	"github.com/warriorguo/vocabulary/internal/repository"
)

const (
	freeDictAPIURL = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	cacheTTL       = 7 * 24 * time.Hour // 7 days
	sourceFreeDic  = "freedictionaryapi"
)

type DictionaryService struct {
	repo   *repository.Repository
	client *http.Client
}

func NewDictionaryService(repo *repository.Repository) *DictionaryService {
	return &DictionaryService{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FreeDictAPIResponse represents the raw API response
type FreeDictAPIResponse []struct {
	Word      string `json:"word"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string   `json:"definition"`
			Example    string   `json:"example"`
			Synonyms   []string `json:"synonyms"`
			Antonyms   []string `json:"antonyms"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []string `json:"antonyms"`
	} `json:"meanings"`
	SourceUrls []string `json:"sourceUrls"`
}

func (s *DictionaryService) LookupWord(ctx context.Context, word string) (*models.DictionaryEntry, error) {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	// Check cache first
	cached, err := s.repo.GetCachedDictionary(ctx, word)
	if err != nil {
		return nil, fmt.Errorf("cache lookup failed: %w", err)
	}
	if cached != nil {
		var entry models.DictionaryEntry
		if err := json.Unmarshal(cached.Data, &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return &entry, nil
	}

	// Fetch from API
	entry, err := s.fetchFromAPI(ctx, word)
	if err != nil {
		return nil, err
	}

	// Cache the result
	data, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %w", err)
	}

	if cacheErr := s.repo.SetCachedDictionary(ctx, word, data, sourceFreeDic, cacheTTL); cacheErr != nil {
		// Log but don't fail - caching is optional
		fmt.Printf("Warning: failed to cache dictionary entry: %v\n", cacheErr)
	}

	return entry, nil
}

func (s *DictionaryService) fetchFromAPI(ctx context.Context, word string) (*models.DictionaryEntry, error) {
	url := freeDictAPIURL + word
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("word not found: %s", word)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp FreeDictAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(apiResp) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	// Normalize to our format
	return s.normalizeResponse(apiResp), nil
}

func (s *DictionaryService) normalizeResponse(apiResp FreeDictAPIResponse) *models.DictionaryEntry {
	first := apiResp[0]

	entry := &models.DictionaryEntry{
		Word:      first.Word,
		Phonetics: make([]models.Phonetic, 0),
		Meanings:  make([]models.Meaning, 0),
	}

	if len(first.SourceUrls) > 0 {
		entry.SourceURL = first.SourceUrls[0]
	}

	// Process phonetics - prefer ones with audio
	for _, p := range first.Phonetics {
		phonetic := models.Phonetic{
			Text:      p.Text,
			Audio:     p.Audio,
			SourceURL: p.SourceURL,
		}
		entry.Phonetics = append(entry.Phonetics, phonetic)
	}

	// Process meanings
	for _, m := range first.Meanings {
		meaning := models.Meaning{
			PartOfSpeech: m.PartOfSpeech,
			Definitions:  make([]models.Definition, 0),
			Synonyms:     m.Synonyms,
			Antonyms:     m.Antonyms,
		}

		for _, d := range m.Definitions {
			def := models.Definition{
				Definition: d.Definition,
				Example:    d.Example,
				Synonyms:   d.Synonyms,
				Antonyms:   d.Antonyms,
			}
			meaning.Definitions = append(meaning.Definitions, def)
		}

		entry.Meanings = append(entry.Meanings, meaning)
	}

	return entry
}
