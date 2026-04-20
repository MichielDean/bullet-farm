import { getAuthHeaders } from '../hooks/useAuth';
import type { Droplet, ImportRequest } from './types';

async function apiFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  const authHeaders = getAuthHeaders();
  Object.assign(headers, authHeaders);
  const resp = await fetch(url, { ...options, headers: { ...headers, ...options?.headers } });
  if (!resp.ok) {
    const body = await resp.text().catch(() => resp.statusText);
    throw new Error(`API error ${resp.status}: ${body}`);
  }
  if (resp.status === 204) return undefined as T;
  return resp.json();
}

export interface ImportPreview {
  key: string;
  title: string;
  description: string;
  priority: number;
  labels: string[];
  source_url: string;
}

export async function fetchImportPreview(provider: string, key: string): Promise<ImportPreview> {
  const params = new URLSearchParams({ provider, key });
  return apiFetch<ImportPreview>(`/api/import/preview?${params}`);
}

export async function importIssue(req: ImportRequest): Promise<Droplet> {
  return apiFetch<Droplet>('/api/import', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}