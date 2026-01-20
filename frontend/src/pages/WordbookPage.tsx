import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import type { WordbookEntry } from '../types';

export function WordbookPage() {
  const [entries, setEntries] = useState<WordbookEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    loadWordbook();
  }, []);

  const loadWordbook = async () => {
    try {
      setLoading(true);
      const data = await api.getWordbook();
      setEntries(data);
    } catch (err) {
      setError('Failed to load wordbook');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (word: string) => {
    try {
      await api.removeFromWordbook(word);
      setEntries(entries.filter(e => e.word !== word));
    } catch (err) {
      console.error('Failed to delete word:', err);
    }
  };

  const handleWordClick = (word: string) => {
    navigate(`/?word=${encodeURIComponent(word)}`);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <div className="wordbook-page">
      <header className="header">
        <h1>My Wordbook</h1>
        <Link to="/" className="nav-link">Search</Link>
      </header>

      {loading && <div className="loading">Loading...</div>}

      {error && <div className="error-message">{error}</div>}

      {!loading && !error && entries.length === 0 && (
        <div className="empty-state">
          <p>Your wordbook is empty.</p>
          <p>Search for words and add them to build your vocabulary!</p>
          <Link to="/" className="start-button">Start Searching</Link>
        </div>
      )}

      {!loading && entries.length > 0 && (
        <div className="wordbook-list">
          {entries.map((entry) => (
            <div key={entry.id} className="wordbook-item">
              <div
                className="word-info"
                onClick={() => handleWordClick(entry.word)}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => e.key === 'Enter' && handleWordClick(entry.word)}
              >
                <span className="word">{entry.word}</span>
                <span className="definition">{entry.short_definition}</span>
                <span className="date">{formatDate(entry.created_at)}</span>
              </div>
              <button
                className="delete-button"
                onClick={() => handleDelete(entry.word)}
                aria-label={`Delete ${entry.word}`}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
