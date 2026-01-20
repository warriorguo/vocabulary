import axios from 'axios';
import type { LookupResponse, WordbookResponse, WordbookEntry, AddWordRequest } from '../types';

const API_BASE = '/api';

const client = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const api = {
  async lookupWord(word: string): Promise<LookupResponse> {
    const response = await client.get<LookupResponse>('/dict', {
      params: { word },
    });
    return response.data;
  },

  async getWordbook(): Promise<WordbookEntry[]> {
    const response = await client.get<WordbookResponse>('/wordbook');
    return response.data.entries;
  },

  async addToWordbook(request: AddWordRequest): Promise<WordbookEntry> {
    const response = await client.post<{ entry: WordbookEntry }>('/wordbook', request);
    return response.data.entry;
  },

  async removeFromWordbook(word: string): Promise<void> {
    await client.delete(`/wordbook/${encodeURIComponent(word)}`);
  },
};
