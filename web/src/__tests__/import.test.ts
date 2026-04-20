import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { importIssue } from '../api/import';
import type { Droplet } from '../api/types';

const mockDroplet: Droplet = {
  id: 'mr-abc1',
  repo: 'myrepo',
  title: 'Jira Issue Title',
  description: 'Jira issue description',
  priority: 2,
  complexity: 2,
  status: 'open',
  assignee: '',
  current_cataractae: '',
  created_at: '2026-04-20T00:00:00Z',
  updated_at: '2026-04-20T00:00:00Z',
};

describe('import API', () => {
  beforeEach(() => { localStorage.clear(); });
  afterEach(() => { vi.restoreAllMocks(); });

  it('sends POST to /api/import', async () => {
    vi.spyOn(window, 'fetch').mockResolvedValue({
      ok: true,
      status: 201,
      json: () => Promise.resolve(mockDroplet),
      text: () => Promise.resolve(JSON.stringify(mockDroplet)),
    } as Response);

    const result = await importIssue({
      provider: 'jira',
      key: 'PROJ-123',
      repo: 'myrepo',
      complexity: 2,
      priority: 1,
    });
    expect(result.id).toBe('mr-abc1');
    expect(result.title).toBe('Jira Issue Title');
  });
});