import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAuthStore } from '@/stores/authStore';
import { STORAGE_KEYS } from '@/lib/constants';

function createFakeJwt(payload: Record<string, unknown>): string {
  const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const body = btoa(JSON.stringify(payload));
  const signature = 'fake-signature';
  return `${header}.${body}.${signature}`;
}

describe('authStore', () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({
      userId: null,
      isAuthenticated: false,
      isVerified: false,
      authMethod: null,
    });
  });

  it('has correct initial state', () => {
    const state = useAuthStore.getState();

    expect(
      state.userId,
      'Initial userId should be null. Check initialState in authStore.',
    ).toBeNull();
    expect(
      state.isAuthenticated,
      'Initial isAuthenticated should be false. Check initialState in authStore.',
    ).toBe(false);
    expect(
      state.isVerified,
      'Initial isVerified should be false. Check initialState in authStore.',
    ).toBe(false);
    expect(
      state.authMethod,
      'Initial authMethod should be null. Check initialState in authStore.',
    ).toBeNull();
  });

  it('login sets authenticated state and stores tokens', () => {
    useAuthStore.getState().login({
      userId: 'user-123',
      isVerified: true,
      authMethod: 'local',
      accessToken: 'access-abc',
      refreshToken: 'refresh-xyz',
    });

    const state = useAuthStore.getState();
    expect(
      state.userId,
      'After login, userId should match the provided value.',
    ).toBe('user-123');
    expect(
      state.isAuthenticated,
      'After login, isAuthenticated should be true.',
    ).toBe(true);
    expect(
      state.isVerified,
      'After login, isVerified should reflect the provided value.',
    ).toBe(true);
    expect(
      state.authMethod,
      'After login, authMethod should match the provided value.',
    ).toBe('local');
    expect(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN),
      'login should store access token in localStorage.',
    ).toBe('access-abc');
    expect(
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN),
      'login should store refresh token in localStorage.',
    ).toBe('refresh-xyz');
  });

  it('logout resets state and clears tokens', () => {
    useAuthStore.getState().login({
      userId: 'user-123',
      isVerified: true,
      authMethod: 'local',
      accessToken: 'access-abc',
      refreshToken: 'refresh-xyz',
    });

    useAuthStore.getState().logout();

    const state = useAuthStore.getState();
    expect(
      state.userId,
      'After logout, userId should be reset to null.',
    ).toBeNull();
    expect(
      state.isAuthenticated,
      'After logout, isAuthenticated should be reset to false.',
    ).toBe(false);
    expect(
      state.isVerified,
      'After logout, isVerified should be reset to false.',
    ).toBe(false);
    expect(
      state.authMethod,
      'After logout, authMethod should be reset to null.',
    ).toBeNull();
    expect(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN),
      'After logout, access token should be removed from localStorage.',
    ).toBeNull();
    expect(
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN),
      'After logout, refresh token should be removed from localStorage.',
    ).toBeNull();
  });

  it('initialize hydrates state from a valid JWT', () => {
    const futureExp = Math.floor(Date.now() / 1000) + 3600;
    const token = createFakeJwt({
      sub: 'user-456',
      is_verified: true,
      auth_method: 'google',
      exp: futureExp,
    });
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, token);

    useAuthStore.getState().initialize();

    const state = useAuthStore.getState();
    expect(
      state.isAuthenticated,
      'initialize should set isAuthenticated to true when a valid, non-expired JWT exists.',
    ).toBe(true);
    expect(
      state.userId,
      'initialize should extract userId from the JWT "sub" claim.',
    ).toBe('user-456');
    expect(
      state.isVerified,
      'initialize should extract isVerified from the JWT "is_verified" claim.',
    ).toBe(true);
    expect(
      state.authMethod,
      'initialize should extract authMethod from the JWT "auth_method" claim.',
    ).toBe('google');
  });

  it('initialize clears tokens and stays unauthenticated for expired JWT', () => {
    const pastExp = Math.floor(Date.now() / 1000) - 3600;
    const token = createFakeJwt({ sub: 'user-789', exp: pastExp });
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, token);

    useAuthStore.getState().initialize();

    expect(
      useAuthStore.getState().isAuthenticated,
      'initialize should not authenticate when the JWT is expired (exp < now).',
    ).toBe(false);
    expect(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN),
      'initialize should clear the expired token from localStorage.',
    ).toBeNull();
  });

  it('initialize clears tokens for malformed token', () => {
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, 'not-a-jwt');

    useAuthStore.getState().initialize();

    expect(
      useAuthStore.getState().isAuthenticated,
      'initialize should not authenticate when the token is malformed (not a valid JWT).',
    ).toBe(false);
    expect(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN),
      'initialize should clear the malformed token from localStorage.',
    ).toBeNull();
  });

  it('initialize does not set isAuthenticated when no token', () => {
    useAuthStore.getState().initialize();

    expect(
      useAuthStore.getState().isAuthenticated,
      'initialize should leave isAuthenticated as false when no access token in localStorage.',
    ).toBe(false);
  });

  it('initialize handles JWT without exp claim gracefully', () => {
    const token = createFakeJwt({ sub: 'user-no-exp', is_verified: false });
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, token);

    useAuthStore.getState().initialize();

    const state = useAuthStore.getState();
    expect(
      state.isAuthenticated,
      'initialize should authenticate when JWT has no exp claim (no expiry to check).',
    ).toBe(true);
    expect(
      state.userId,
      'initialize should still extract sub from a JWT without exp.',
    ).toBe('user-no-exp');
  });

  it('initialize uses current time for expiry check', () => {
    const futureExp = Math.floor(Date.now() / 1000) + 10;
    const token = createFakeJwt({ sub: 'user-time', exp: futureExp });
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, token);

    vi.useFakeTimers();
    vi.setSystemTime(new Date((futureExp + 1) * 1000));

    useAuthStore.getState().initialize();

    expect(
      useAuthStore.getState().isAuthenticated,
      'initialize should reject a token whose exp is in the past relative to Date.now().',
    ).toBe(false);

    vi.useRealTimers();
  });
});
