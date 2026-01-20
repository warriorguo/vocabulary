export interface Phonetic {
  text?: string;
  audio?: string;
  sourceUrl?: string;
}

export interface Definition {
  definition: string;
  example?: string;
  synonyms?: string[];
  antonyms?: string[];
}

export interface Meaning {
  partOfSpeech: string;
  definitions: Definition[];
  synonyms?: string[];
  antonyms?: string[];
}

export interface DictionaryEntry {
  word: string;
  phonetics: Phonetic[];
  meanings: Meaning[];
  sourceUrl?: string;
}

export interface WordbookEntry {
  id: number;
  user_id: string;
  word: string;
  short_definition: string;
  created_at: string;
}

export interface LookupResponse {
  entry: DictionaryEntry;
  in_wordbook: boolean;
}

export interface WordbookResponse {
  entries: WordbookEntry[];
}

export interface AddWordRequest {
  word: string;
  short_definition: string;
}
