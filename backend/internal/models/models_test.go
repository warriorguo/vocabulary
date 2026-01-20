package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWordbookEntryJSON(t *testing.T) {
	entry := WordbookEntry{
		ID:              1,
		UserID:          "test-user",
		Word:            "hello",
		ShortDefinition: "a greeting",
		CreatedAt:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded WordbookEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ID != entry.ID {
		t.Errorf("ID mismatch: got %d, want %d", decoded.ID, entry.ID)
	}
	if decoded.Word != entry.Word {
		t.Errorf("Word mismatch: got %s, want %s", decoded.Word, entry.Word)
	}
}

func TestDictionaryEntryJSON(t *testing.T) {
	entry := DictionaryEntry{
		Word: "test",
		Phonetics: []Phonetic{
			{Text: "/test/", Audio: "https://example.com/test.mp3"},
		},
		Meanings: []Meaning{
			{
				PartOfSpeech: "noun",
				Definitions: []Definition{
					{Definition: "a procedure", Example: "run a test"},
				},
			},
		},
		SourceURL: "https://example.com",
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded DictionaryEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Word != entry.Word {
		t.Errorf("Word mismatch: got %s, want %s", decoded.Word, entry.Word)
	}
	if len(decoded.Phonetics) != 1 {
		t.Errorf("Phonetics length mismatch: got %d, want 1", len(decoded.Phonetics))
	}
	if len(decoded.Meanings) != 1 {
		t.Errorf("Meanings length mismatch: got %d, want 1", len(decoded.Meanings))
	}
}

func TestAddWordRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     AddWordRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     AddWordRequest{Word: "hello", ShortDefinition: "a greeting"},
			wantErr: false,
		},
		{
			name:    "empty word",
			req:     AddWordRequest{Word: "", ShortDefinition: "a greeting"},
			wantErr: true,
		},
		{
			name:    "empty definition",
			req:     AddWordRequest{Word: "hello", ShortDefinition: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.req.Word == "" || tt.req.ShortDefinition == ""
			if hasErr != tt.wantErr {
				t.Errorf("validation mismatch: got error=%v, want error=%v", hasErr, tt.wantErr)
			}
		})
	}
}
