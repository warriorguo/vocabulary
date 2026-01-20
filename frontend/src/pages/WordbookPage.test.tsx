import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { WordbookPage } from './WordbookPage';
import { api } from '../services/api';

vi.mock('../services/api', () => ({
  api: {
    getWordbook: vi.fn(),
    removeFromWordbook: vi.fn(),
  },
}));

const mockedApi = vi.mocked(api);

const renderWordbookPage = () => {
  return render(
    <BrowserRouter>
      <WordbookPage />
    </BrowserRouter>
  );
};

describe('WordbookPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render page title', async () => {
    mockedApi.getWordbook.mockResolvedValue([]);

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByText('My Wordbook')).toBeInTheDocument();
    });
  });

  it('should render navigation link to search', async () => {
    mockedApi.getWordbook.mockResolvedValue([]);

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByRole('link', { name: /search/i })).toBeInTheDocument();
    });
  });

  it('should display empty state when wordbook is empty', async () => {
    mockedApi.getWordbook.mockResolvedValue([]);

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByText('Your wordbook is empty.')).toBeInTheDocument();
      expect(screen.getByText(/search for words/i)).toBeInTheDocument();
    });
  });

  it('should display loading state initially', () => {
    mockedApi.getWordbook.mockImplementation(() => new Promise(() => {})); // Never resolves

    renderWordbookPage();

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('should display wordbook entries', async () => {
    const mockEntries = [
      {
        id: 1,
        user_id: 'default',
        word: 'hello',
        short_definition: 'a greeting',
        created_at: '2024-01-15T10:00:00Z',
      },
      {
        id: 2,
        user_id: 'default',
        word: 'world',
        short_definition: 'the earth',
        created_at: '2024-01-14T10:00:00Z',
      },
    ];
    mockedApi.getWordbook.mockResolvedValue(mockEntries);

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByText('hello')).toBeInTheDocument();
      expect(screen.getByText('a greeting')).toBeInTheDocument();
      expect(screen.getByText('world')).toBeInTheDocument();
      expect(screen.getByText('the earth')).toBeInTheDocument();
    });
  });

  it('should call removeFromWordbook when delete button is clicked', async () => {
    const mockEntries = [
      {
        id: 1,
        user_id: 'default',
        word: 'hello',
        short_definition: 'a greeting',
        created_at: '2024-01-15T10:00:00Z',
      },
    ];
    mockedApi.getWordbook.mockResolvedValue(mockEntries);
    mockedApi.removeFromWordbook.mockResolvedValue(undefined);

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByText('hello')).toBeInTheDocument();
    });

    const deleteButton = screen.getByLabelText('Delete hello');
    fireEvent.click(deleteButton);

    await waitFor(() => {
      expect(mockedApi.removeFromWordbook).toHaveBeenCalledWith('hello');
    });
  });

  it('should remove entry from list after successful deletion', async () => {
    const mockEntries = [
      {
        id: 1,
        user_id: 'default',
        word: 'hello',
        short_definition: 'a greeting',
        created_at: '2024-01-15T10:00:00Z',
      },
    ];
    mockedApi.getWordbook.mockResolvedValue(mockEntries);
    mockedApi.removeFromWordbook.mockResolvedValue(undefined);

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByText('hello')).toBeInTheDocument();
    });

    const deleteButton = screen.getByLabelText('Delete hello');
    fireEvent.click(deleteButton);

    await waitFor(() => {
      expect(screen.queryByText('hello')).not.toBeInTheDocument();
    });
  });

  it('should display error message when loading fails', async () => {
    mockedApi.getWordbook.mockRejectedValue(new Error('Network error'));

    renderWordbookPage();

    await waitFor(() => {
      expect(screen.getByText('Failed to load wordbook')).toBeInTheDocument();
    });
  });

  it('should format dates correctly', async () => {
    const mockEntries = [
      {
        id: 1,
        user_id: 'default',
        word: 'test',
        short_definition: 'a test',
        created_at: '2024-01-15T10:00:00Z',
      },
    ];
    mockedApi.getWordbook.mockResolvedValue(mockEntries);

    renderWordbookPage();

    await waitFor(() => {
      // The date format is "Jan 15, 2024" based on the formatDate function
      expect(screen.getByText(/jan.*15.*2024/i)).toBeInTheDocument();
    });
  });
});
