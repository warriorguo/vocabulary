package models

import (
	"time"
)

// WordbookEntry represents a word saved in the user's wordbook
type WordbookEntry struct {
	ID              int64     `json:"id"`
	UserID          string    `json:"user_id"`
	Word            string    `json:"word"`
	ShortDefinition string    `json:"short_definition"`
	CreatedAt       time.Time `json:"created_at"`
}

// DictionaryCache represents cached dictionary data
type DictionaryCache struct {
	Word      string    `json:"word"`
	Data      []byte    `json:"data"`
	Source    string    `json:"source"`
	FetchedAt time.Time `json:"fetched_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Phonetic represents pronunciation information
type Phonetic struct {
	Text      string `json:"text,omitempty"`
	Audio     string `json:"audio,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
}

// Definition represents a single definition
type Definition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example,omitempty"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
}

// Meaning represents a meaning with part of speech
type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
	Synonyms     []string     `json:"synonyms,omitempty"`
	Antonyms     []string     `json:"antonyms,omitempty"`
}

// DictionaryEntry represents the normalized dictionary response
type DictionaryEntry struct {
	Word      string     `json:"word"`
	Phonetics []Phonetic `json:"phonetics"`
	Meanings  []Meaning  `json:"meanings"`
	SourceURL string     `json:"sourceUrl,omitempty"`
}

// AddWordRequest represents the request body for adding a word
type AddWordRequest struct {
	Word            string `json:"word" binding:"required"`
	ShortDefinition string `json:"short_definition" binding:"required"`
}
