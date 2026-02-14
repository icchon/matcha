import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import {
  apiClient,
  ApiClientError,
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
} from '@/api/client';
import { STORAGE_KEYS } from '@/lib/constants';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

function jsonResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  } as Response;
}

function noContentResponse(): Response {
  return {
    ok: true,
    status: 204,
    json: () => Promise.reject(new Error('No content')),
  } as Response;
}

describe('Token management', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('setTokens stores both tokens in localStorage', () => {
    setTokens('access-123', 'refresh-456');

    expect(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN),
      'Access token should be stored under STORAGE_KEYS.ACCESS_TOKEN. Check setTokens uses the correct key.',
    ).toBe('access-123');
    expect(
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN),
      'Refresh token should be stored under STORAGE_KEYS.REFRESH_TOKEN. Check setTokens uses the correct key.',
    ).toBe('refresh-456');
  });

  it('getAccessToken returns stored token', () => {
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, 'my-token');

    expect(
      getAccessToken(),
      'getAccessToken should read from STORAGE_KEYS.ACCESS_TOKEN in localStorage.',
    ).toBe('my-token');
  });

  it('getAccessToken returns null when no token stored', () => {
    expect(
      getAccessToken(),
      'getAccessToken should return null when localStorage has no access token.',
    ).toBeNull();
  });

  it('getRefreshToken returns stored token', () => {
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, 'refresh-abc');

    expect(
      getRefreshToken(),
      'getRefreshToken should read from STORAGE_KEYS.REFRESH_TOKEN in localStorage.',
    ).toBe('refresh-abc');
  });

  it('clearTokens removes both tokens from localStorage', () => {
    setTokens('a', 'b');
    clearTokens();

    expect(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN),
      'clearTokens should remove access token from localStorage.',
    ).toBeNull();
    expect(
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN),
      'clearTokens should remove refresh token from localStorage.',
    ).toBeNull();
  });
});

describe('apiClient', () => {
  beforeEach(() => {
    localStorage.clear();
    mockFetch.mockReset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('GET request sends Authorization header when token exists', async () => {
    setTokens('bearer-token', 'refresh');
    mockFetch.mockResolvedValue(jsonResponse({ data: 'ok' }));

    await apiClient.get('/test');

    const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(url, 'GET should prepend /api/v1 to the path.').toBe('/api/v1/test');
    expect(
      (options.headers as Record<string, string>)['Authorization'],
      'Authorization header should be "Bearer <token>" when access token is stored.',
    ).toBe('Bearer bearer-token');
  });

  it('GET request omits Authorization header when no token', async () => {
    mockFetch.mockResolvedValue(jsonResponse({ data: 'ok' }));

    await apiClient.get('/test');

    const [, options] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(
      (options.headers as Record<string, string>)['Authorization'],
      'Authorization header should be absent when no token is stored.',
    ).toBeUndefined();
  });

  it('POST sends JSON body', async () => {
    mockFetch.mockResolvedValue(jsonResponse({ id: 1 }, 201));

    const result = await apiClient.post<{ id: number }>('/items', { name: 'test' });

    const [, options] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(options.method, 'POST request should use POST method.').toBe('POST');
    expect(
      JSON.parse(options.body as string),
      'POST body should be JSON-serialized request payload.',
    ).toEqual({ name: 'test' });
    expect(result.id, 'POST should return parsed JSON response body.').toBe(1);
  });

  it('PUT sends JSON body', async () => {
    mockFetch.mockResolvedValue(jsonResponse({ updated: true }));

    const result = await apiClient.put<{ updated: boolean }>('/items/1', { name: 'updated' });

    const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(url).toBe('/api/v1/items/1');
    expect(options.method, 'PUT request should use PUT method.').toBe('PUT');
    expect(result.updated).toBe(true);
  });

  it('DELETE request sends correct method', async () => {
    mockFetch.mockResolvedValue(noContentResponse());

    await apiClient.delete('/items/1');

    const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(url).toBe('/api/v1/items/1');
    expect(options.method, 'DELETE request should use DELETE method.').toBe('DELETE');
  });

  it('handles 204 No Content response', async () => {
    mockFetch.mockResolvedValue(noContentResponse());

    const result = await apiClient.delete('/items/1');

    expect(
      result,
      'A 204 response should resolve to undefined since there is no body.',
    ).toBeUndefined();
  });

  it('throws ApiClientError on 401 response', async () => {
    mockFetch.mockResolvedValue(jsonResponse({ error: 'Unauthorized' }, 401));

    try {
      await apiClient.get('/protected');
      expect.fail('Should have thrown ApiClientError for 401 response.');
    } catch (err) {
      expect(err).toBeInstanceOf(ApiClientError);
      const apiErr = err as ApiClientError;
      expect(
        apiErr.status,
        'ApiClientError.status should match HTTP status code 401.',
      ).toBe(401);
      expect(
        apiErr.body.error,
        'ApiClientError.body.error should contain the error message from the response.',
      ).toBe('Unauthorized');
    }
  });

  it('throws ApiClientError on 404 response', async () => {
    mockFetch.mockResolvedValue(jsonResponse({ error: 'Not found' }, 404));

    try {
      await apiClient.get('/missing');
      expect.fail('Should have thrown ApiClientError for 404 response.');
    } catch (err) {
      expect(err).toBeInstanceOf(ApiClientError);
      const apiErr = err as ApiClientError;
      expect(apiErr.status, 'ApiClientError.status should be 404 for not-found responses.').toBe(
        404,
      );
    }
  });

  it('upload sends FormData without Content-Type header', async () => {
    setTokens('bearer-token', 'refresh');
    mockFetch.mockResolvedValue(jsonResponse({ url: '/images/pic.jpg' }, 201));

    const formData = new FormData();
    formData.append('file', new Blob(['data']), 'photo.jpg');

    const result = await apiClient.upload<{ url: string }>('/me/profile/pictures', formData);

    const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(url, 'upload should prepend /api/v1 to the path.').toBe(
      '/api/v1/me/profile/pictures',
    );
    expect(options.method, 'upload should use POST method.').toBe('POST');
    expect(
      options.body,
      'upload should send FormData as body directly (not JSON-serialized).',
    ).toBe(formData);
    expect(
      (options.headers as Record<string, string>)['Content-Type'],
      'upload should NOT set Content-Type header so browser sets multipart/form-data boundary automatically.',
    ).toBeUndefined();
    expect(
      (options.headers as Record<string, string>)['Authorization'],
      'upload should include Authorization header when token exists.',
    ).toBe('Bearer bearer-token');
    expect(result.url).toBe('/images/pic.jpg');
  });

  it('handles non-JSON error response gracefully', async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 500,
      json: () => Promise.reject(new Error('Invalid JSON')),
    } as unknown as Response);

    try {
      await apiClient.get('/broken');
      expect.fail('Should have thrown ApiClientError for 500 response.');
    } catch (err) {
      expect(err).toBeInstanceOf(ApiClientError);
      const apiErr = err as ApiClientError;
      expect(apiErr.status).toBe(500);
      expect(
        apiErr.body.error,
        'When response body is not valid JSON, error message should include the status code.',
      ).toContain('500');
    }
  });
});
