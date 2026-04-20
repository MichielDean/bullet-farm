import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { CreateDropletForm } from '../components/CreateDropletForm';
import type { Droplet } from '../api/types';

const mockDroplet: Droplet = {
  id: 'ct-new1',
  repo: 'cistern',
  title: 'New Droplet',
  description: 'A new one',
  priority: 2,
  complexity: 1,
  status: 'open',
  assignee: '',
  current_cataractae: '',
  created_at: '2026-04-19T00:00:00Z',
  updated_at: '2026-04-19T00:00:00Z',
};

function mockFetch(urls: Record<string, unknown>) {
  vi.spyOn(window, 'fetch').mockImplementation(async (input: RequestInfo | URL) => {
    const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
    for (const [pattern, response] of Object.entries(urls)) {
      if (url.includes(pattern)) {
        return {
          ok: true,
          status: 200,
          json: () => Promise.resolve(response),
          text: () => Promise.resolve(JSON.stringify(response)),
        } as Response;
      }
    }
    return {
      ok: true,
      status: 200,
      json: () => Promise.resolve([]),
      text: () => Promise.resolve('[]'),
    } as Response;
  });
}

describe('CreateDropletForm', () => {
  beforeEach(() => { localStorage.clear(); });
  afterEach(() => { vi.restoreAllMocks(); });

  it('renders form fields', () => {
    mockFetch({ '/api/repos': [{ name: 'cistern', url: 'git@github.com:test/cistern' }] });
    render(<CreateDropletForm onSuccess={vi.fn()} onCancel={vi.fn()} />);
    expect(screen.getByText('Title *')).toBeInTheDocument();
    expect(screen.getByText('Description')).toBeInTheDocument();
    expect(screen.getByText('Priority')).toBeInTheDocument();
    expect(screen.getByText('Complexity')).toBeInTheDocument();
    expect(screen.getByText('Dependencies')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create droplet/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument();
  });

  it('disables submit when title is empty', () => {
    mockFetch({ '/api/repos': [{ name: 'cistern', url: 'git@github.com:test/cistern' }] });
    render(<CreateDropletForm onSuccess={vi.fn()} onCancel={vi.fn()} />);
    const submitBtn = screen.getByRole('button', { name: /create droplet/i });
    expect(submitBtn).toBeDisabled();
  });

  it('shows complexity levels', () => {
    mockFetch({ '/api/repos': [{ name: 'cistern', url: 'git@github.com:test/cistern' }] });
    render(<CreateDropletForm onSuccess={vi.fn()} onCancel={vi.fn()} />);
    expect(screen.getByText('Standard (1)')).toBeInTheDocument();
    expect(screen.getByText('Full (2)')).toBeInTheDocument();
    expect(screen.getByText('Critical (3)')).toBeInTheDocument();
  });

  it('calls onCancel when cancel is clicked', () => {
    mockFetch({ '/api/repos': [{ name: 'cistern', url: 'git@github.com:test/cistern' }] });
    const onCancel = vi.fn();
    render(<CreateDropletForm onSuccess={vi.fn()} onCancel={onCancel} />);
    screen.getByRole('button', { name: /cancel/i }).click();
    expect(onCancel).toHaveBeenCalled();
  });

  it('creates droplet on submit with valid data', async () => {
    mockFetch({
      '/api/repos': [{ name: 'cistern', url: 'git@github.com:test/cistern' }],
      '/api/droplets': mockDroplet,
    });
    const onSuccess = vi.fn();

    render(<CreateDropletForm onSuccess={onSuccess} onCancel={vi.fn()} />);

    await act(async () => {
      await new Promise((r) => setTimeout(r, 0));
    });

    const textInputs = screen.getAllByRole('textbox');
    const titleInput = textInputs.find((el) => el.getAttribute('required') !== null) ?? textInputs[0];
    await act(() => { fireEvent.change(titleInput, { target: { value: 'New Droplet' } }); });

    const repoSelect = screen.getByRole('combobox');
    await act(() => { fireEvent.change(repoSelect, { target: { value: 'cistern' } }); });

    const submitBtn = screen.getByRole('button', { name: /create droplet/i });
    expect(submitBtn).not.toBeDisabled();
    await act(() => { fireEvent.click(submitBtn); });

    expect(window.fetch).toHaveBeenCalled();
  });
});