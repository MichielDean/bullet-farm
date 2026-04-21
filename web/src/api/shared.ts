import { getAuthHeaders, clearStoredKey } from '../hooks/useAuth';

export { clearStoredKey } from '../hooks/useAuth';

export async function apiFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  const authHeaders = getAuthHeaders();
  Object.assign(headers, authHeaders);
  const resp = await fetch(url, { ...options, headers: { ...headers, ...options?.headers } });
  if (resp.status === 401) {
    clearStoredKey();
    window.dispatchEvent(new CustomEvent('cistern:auth-expired'));
    throw new Error('API error 401: Authentication required');
  }
  if (!resp.ok) {
    const body = await resp.text().catch(() => resp.statusText);
    throw new Error(`API error ${resp.status}: ${body}`);
  }
  if (resp.status === 204) return undefined as T;
  return resp.json();
}