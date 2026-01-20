import { useState, useRef, useEffect } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { api } from '../services/api';
import type { DictionaryEntry } from '../types';

export function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialWord = searchParams.get('word') || '';

  const [inputValue, setInputValue] = useState(initialWord);
  const [entry, setEntry] = useState<DictionaryEntry | null>(null);
  const [inWordbook, setInWordbook] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const audioRef = useRef<HTMLAudioElement>(null);

  const handleSearch = async (word?: string) => {
    const searchWord = word || inputValue.trim();
    if (!searchWord) return;

    setLoading(true);
    setError(null);
    setSearchParams({ word: searchWord });

    try {
      const response = await api.lookupWord(searchWord);
      setEntry(response.entry);
      setInWordbook(response.in_wordbook);
    } catch (err) {
      setEntry(null);
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to lookup word');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  const handleAddToWordbook = async () => {
    if (!entry) return;

    const shortDef = entry.meanings[0]?.definitions[0]?.definition || '';
    try {
      await api.addToWordbook({
        word: entry.word,
        short_definition: shortDef,
      });
      setInWordbook(true);
    } catch (err) {
      console.error('Failed to add to wordbook:', err);
    }
  };

  const handleRemoveFromWordbook = async () => {
    if (!entry) return;

    try {
      await api.removeFromWordbook(entry.word);
      setInWordbook(false);
    } catch (err) {
      console.error('Failed to remove from wordbook:', err);
    }
  };

  const playAudio = () => {
    if (audioRef.current) {
      audioRef.current.play();
    }
  };

  const audioUrl = entry?.phonetics.find(p => p.audio)?.audio;
  const phoneticText = entry?.phonetics.find(p => p.text)?.text;

  // Load initial word if present
  useEffect(() => {
    if (initialWord) {
      handleSearch(initialWord);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="search-page">
      <header className="header">
        <h1>English Wordbook</h1>
        <Link to="/wordbook" className="nav-link">My Wordbook</Link>
      </header>

      <div className="search-container">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Enter a word..."
          className="search-input"
        />
        <button onClick={() => handleSearch()} disabled={loading} className="search-button">
          {loading ? 'Searching...' : 'Search'}
        </button>
      </div>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {entry && (
        <div className="result-card">
          <div className="word-header">
            <h2 className="word-title">{entry.word}</h2>
            {phoneticText && <span className="phonetic">{phoneticText}</span>}
            {audioUrl && (
              <>
                <button onClick={playAudio} className="audio-button" aria-label="Play pronunciation">
                  🔊
                </button>
                <audio ref={audioRef} src={audioUrl} />
              </>
            )}
          </div>

          <button
            onClick={inWordbook ? handleRemoveFromWordbook : handleAddToWordbook}
            className={`wordbook-button ${inWordbook ? 'remove' : 'add'}`}
          >
            {inWordbook ? '✓ In Wordbook' : '+ Add to Wordbook'}
          </button>

          <div className="meanings">
            {entry.meanings.map((meaning, idx) => (
              <div key={idx} className="meaning">
                <h3 className="part-of-speech">{meaning.partOfSpeech}</h3>
                <ol className="definitions">
                  {meaning.definitions.map((def, defIdx) => (
                    <li key={defIdx} className="definition">
                      <p>{def.definition}</p>
                      {def.example && (
                        <p className="example">"{def.example}"</p>
                      )}
                    </li>
                  ))}
                </ol>
                {meaning.synonyms && meaning.synonyms.length > 0 && (
                  <p className="synonyms">
                    <strong>Synonyms:</strong> {meaning.synonyms.slice(0, 5).join(', ')}
                  </p>
                )}
              </div>
            ))}
          </div>

          {entry.sourceUrl && (
            <a href={entry.sourceUrl} target="_blank" rel="noopener noreferrer" className="source-link">
              Source
            </a>
          )}
        </div>
      )}
    </div>
  );
}
