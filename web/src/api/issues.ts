import type { DropletIssue, AddIssueRequest, ResolveIssueRequest } from './types';
import { getAuthHeaders } from '../hooks/useAuth';

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

export async function createIssue(dropletId: string, body: AddIssueRequest): Promise<DropletIssue> {
  return apiFetch<DropletIssue>(`/api/droplets/${encodeURIComponent(dropletId)}/issues`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function resolveIssue(issueId: string, body: ResolveIssueRequest): Promise<void> {
  return apiFetch<void>(`/api/issues/${encodeURIComponent(issueId)}/resolve`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function rejectIssue(issueId: string, body: ResolveIssueRequest): Promise<void> {
  return apiFetch<void>(`/api/issues/${encodeURIComponent(issueId)}/reject`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function fetchIssues(
  dropletId: string,
  filters?: { open?: boolean; flagged_by?: string },
): Promise<DropletIssue[]> {
  const params = new URLSearchParams();
  if (filters?.open !== undefined) params.set('open', String(filters.open));
  if (filters?.flagged_by) params.set('flagged_by', filters.flagged_by);
  const qs = params.toString();
  const url = qs
    ? `/api/droplets/${encodeURIComponent(dropletId)}/issues?${qs}`
    : `/api/droplets/${encodeURIComponent(dropletId)}/issues`;
  return apiFetch<DropletIssue[]>(url);
}