import { describe, it, expect } from 'vitest';
import type {
  Phonetic,
  Definition,
  Meaning,
  DictionaryEntry,
  WordbookEntry,
  LookupResponse,
  WordbookResponse,
  AddWordRequest,
} from './index';

describe('Type definitions', () => {
  it('should allow valid Phonetic objects', () => {
    const phonetic: Phonetic = {
      text: '/həˈloʊ/',
      audio: 'https://example.com/audio.mp3',
      sourceUrl: 'https://example.com',
    };
    expect(phonetic.text).toBe('/həˈloʊ/');
  });

  it('should allow Phonetic with optional fields', () => {
    const phonetic: Phonetic = {};
    expect(phonetic.text).toBeUndefined();
  });

  it('should allow valid Definition objects', () => {
    const definition: Definition = {
      definition: 'used as a greeting',
      example: 'hello there!',
      synonyms: ['hi', 'hey'],
      antonyms: ['goodbye'],
    };
    expect(definition.definition).toBe('used as a greeting');
  });

  it('should allow valid Meaning objects', () => {
    const meaning: Meaning = {
      partOfSpeech: 'noun',
      definitions: [{ definition: 'a test' }],
      synonyms: ['trial'],
      antonyms: [],
    };
    expect(meaning.partOfSpeech).toBe('noun');
    expect(meaning.definitions).toHaveLength(1);
  });

  it('should allow valid DictionaryEntry objects', () => {
    const entry: DictionaryEntry = {
      word: 'hello',
      phonetics: [{ text: '/həˈloʊ/' }],
      meanings: [
        {
          partOfSpeech: 'exclamation',
          definitions: [{ definition: 'a greeting' }],
        },
      ],
      sourceUrl: 'https://example.com',
    };
    expect(entry.word).toBe('hello');
    expect(entry.meanings).toHaveLength(1);
  });

  it('should allow valid WordbookEntry objects', () => {
    const entry: WordbookEntry = {
      id: 1,
      user_id: 'default',
      word: 'hello',
      short_definition: 'a greeting',
      created_at: '2024-01-15T10:00:00Z',
    };
    expect(entry.id).toBe(1);
    expect(entry.word).toBe('hello');
  });

  it('should allow valid LookupResponse objects', () => {
    const response: LookupResponse = {
      entry: {
        word: 'test',
        phonetics: [],
        meanings: [],
      },
      in_wordbook: true,
    };
    expect(response.in_wordbook).toBe(true);
  });

  it('should allow valid WordbookResponse objects', () => {
    const response: WordbookResponse = {
      entries: [
        {
          id: 1,
          user_id: 'default',
          word: 'test',
          short_definition: 'a test',
          created_at: '2024-01-15T10:00:00Z',
        },
      ],
    };
    expect(response.entries).toHaveLength(1);
  });

  it('should allow valid AddWordRequest objects', () => {
    const request: AddWordRequest = {
      word: 'hello',
      short_definition: 'a greeting',
    };
    expect(request.word).toBe('hello');
    expect(request.short_definition).toBe('a greeting');
  });
});
