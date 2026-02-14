import type { ApiError } from '@/types';

const BASE_URL = '/api/v1';

export class ApiClientError extends Error {
  readonly status: number;
  readonly body: ApiError;

  constructor(status: number, body: ApiError) {
    super(body.error);
    this.name = 'ApiClientError';
    this.status = status;
    this.body = body;
  }
}

// In-memory token storage (H-1: avoids localStorage XSS risk)
// Trade-off: tokens do not survive page reload — user must re-login
let accessToken: string | null = null;
let refreshToken: string | null = null;

export function getAccessToken(): string | null {
  return accessToken;
}

export function getRefreshToken(): string | null {
  return refreshToken;
}

export function setTokens(access: string, refresh: string): void {
  accessToken = access;
  refreshToken = refresh;
}

export function clearTokens(): void {
  accessToken = null;
  refreshToken = null;
}

function buildHeaders(): Record<string, string> {
  const token = getAccessToken();
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
}

function buildAuthHeader(): Record<string, string> {
  const token = getAccessToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function request<T>(path: string, options: RequestInit): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, options);

  if (!response.ok) {
    // H-2: 401 interceptor — clear tokens on unauthorized
    // [MOCK] POST /auth/token/refresh not yet routed in BE
    // Future: attempt refresh before logout
    if (response.status === 401) {
      clearTokens();
    }

    const body: ApiError = await response.json().catch(() => ({
      error: `Request failed with status ${response.status}`,
    }));
    throw new ApiClientError(response.status, body);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}

export const apiClient = {
  async get<T>(path: string): Promise<T> {
    return request<T>(path, {
      method: 'GET',
      headers: buildHeaders(),
    });
  },

  async post<T>(path: string, body?: unknown): Promise<T> {
    return request<T>(path, {
      method: 'POST',
      headers: buildHeaders(),
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  },

  async put<T>(path: string, body?: unknown): Promise<T> {
    return request<T>(path, {
      method: 'PUT',
      headers: buildHeaders(),
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  },

  async upload<T>(path: string, formData: FormData): Promise<T> {
    return request<T>(path, {
      method: 'POST',
      headers: buildAuthHeader(),
      body: formData,
    });
  },

  async delete(path: string): Promise<void> {
    await request<undefined>(path, {
      method: 'DELETE',
      headers: buildHeaders(),
    });
  },
};
