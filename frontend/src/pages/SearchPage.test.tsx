import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { SearchPage } from './SearchPage';
import { api } from '../services/api';

vi.mock('../services/api', () => ({
  api: {
    lookupWord: vi.fn(),
    addToWordbook: vi.fn(),
    removeFromWordbook: vi.fn(),
  },
}));

const mockedApi = vi.mocked(api);

const renderSearchPage = () => {
  return render(
    <BrowserRouter>
      <SearchPage />
    </BrowserRouter>
  );
};

describe('SearchPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render search input and button', () => {
    renderSearchPage();

    expect(screen.getByPlaceholderText('Enter a word...')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /search/i })).toBeInTheDocument();
  });

  it('should render navigation link to wordbook', () => {
    renderSearchPage();

    expect(screen.getByRole('link', { name: /my wordbook/i })).toBeInTheDocument();
  });

  it('should call lookupWord when search button is clicked', async () => {
    const mockEntry = {
      entry: {
        word: 'hello',
        phonetics: [{ text: '/həˈloʊ/', audio: 'https://example.com/hello.mp3' }],
        meanings: [
          {
            partOfSpeech: 'exclamation',
            definitions: [{ definition: 'used as a greeting', example: 'hello there!' }],
          },
        ],
      },
      in_wordbook: false,
    };
    mockedApi.lookupWord.mockResolvedValue(mockEntry);

    renderSearchPage();

    const input = screen.getByPlaceholderText('Enter a word...');
    const button = screen.getByRole('button', { name: /search/i });

    fireEvent.change(input, { target: { value: 'hello' } });
    fireEvent.click(button);

    await waitFor(() => {
      expect(mockedApi.lookupWord).toHaveBeenCalledWith('hello');
    });
  });

  it('should display word definition after successful search', async () => {
    const mockEntry = {
      entry: {
        word: 'hello',
        phonetics: [{ text: '/həˈloʊ/' }],
        meanings: [
          {
            partOfSpeech: 'exclamation',
            definitions: [{ definition: 'used as a greeting' }],
          },
        ],
      },
      in_wordbook: false,
    };
    mockedApi.lookupWord.mockResolvedValue(mockEntry);

    renderSearchPage();

    const input = screen.getByPlaceholderText('Enter a word...');
    const button = screen.getByRole('button', { name: /search/i });

    fireEvent.change(input, { target: { value: 'hello' } });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText('hello')).toBeInTheDocument();
      expect(screen.getByText('/həˈloʊ/')).toBeInTheDocument();
      expect(screen.getByText('exclamation')).toBeInTheDocument();
      expect(screen.getByText('used as a greeting')).toBeInTheDocument();
    });
  });

  it('should display error message when search fails', async () => {
    mockedApi.lookupWord.mockRejectedValue(new Error('Word not found'));

    renderSearchPage();

    const input = screen.getByPlaceholderText('Enter a word...');
    const button = screen.getByRole('button', { name: /search/i });

    fireEvent.change(input, { target: { value: 'nonexistent' } });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText('Word not found')).toBeInTheDocument();
    });
  });

  it('should show "Add to Wordbook" button when word is not in wordbook', async () => {
    const mockEntry = {
      entry: {
        word: 'hello',
        phonetics: [],
        meanings: [{ partOfSpeech: 'noun', definitions: [{ definition: 'test' }] }],
      },
      in_wordbook: false,
    };
    mockedApi.lookupWord.mockResolvedValue(mockEntry);

    renderSearchPage();

    fireEvent.change(screen.getByPlaceholderText('Enter a word...'), { target: { value: 'hello' } });
    fireEvent.click(screen.getByRole('button', { name: /search/i }));

    await waitFor(() => {
      expect(screen.getByText('+ Add to Wordbook')).toBeInTheDocument();
    });
  });

  it('should show "In Wordbook" button when word is in wordbook', async () => {
    const mockEntry = {
      entry: {
        word: 'hello',
        phonetics: [],
        meanings: [{ partOfSpeech: 'noun', definitions: [{ definition: 'test' }] }],
      },
      in_wordbook: true,
    };
    mockedApi.lookupWord.mockResolvedValue(mockEntry);

    renderSearchPage();

    fireEvent.change(screen.getByPlaceholderText('Enter a word...'), { target: { value: 'hello' } });
    fireEvent.click(screen.getByRole('button', { name: /search/i }));

    await waitFor(() => {
      expect(screen.getByText(/in wordbook/i)).toBeInTheDocument();
    });
  });

  it('should trigger search on Enter key', async () => {
    mockedApi.lookupWord.mockResolvedValue({
      entry: { word: 'test', phonetics: [], meanings: [] },
      in_wordbook: false,
    });

    renderSearchPage();

    const input = screen.getByPlaceholderText('Enter a word...');
    fireEvent.change(input, { target: { value: 'test' } });
    fireEvent.keyDown(input, { key: 'Enter', code: 'Enter' });

    await waitFor(() => {
      expect(mockedApi.lookupWord).toHaveBeenCalledWith('test');
    });
  });
});
