import { describe, it, expect, beforeEach } from 'vitest';
import { useAuthStore } from '@/stores/authStore';
import { getAccessToken, getRefreshToken, clearTokens } from '@/api/client';

describe('authStore', () => {
  beforeEach(() => {
    clearTokens();
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

  it('login sets authenticated state and stores tokens in memory', () => {
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
      getAccessToken(),
      'login should store access token in memory via setTokens.',
    ).toBe('access-abc');
    expect(
      getRefreshToken(),
      'login should store refresh token in memory via setTokens.',
    ).toBe('refresh-xyz');
  });

  it('login does NOT persist tokens to localStorage', () => {
    useAuthStore.getState().login({
      userId: 'user-123',
      isVerified: true,
      authMethod: 'local',
      accessToken: 'access-abc',
      refreshToken: 'refresh-xyz',
    });

    expect(
      localStorage.getItem('matcha_access_token'),
      'Tokens must NOT be persisted to localStorage (H-1 XSS fix). Use in-memory only.',
    ).toBeNull();
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
      getAccessToken(),
      'After logout, access token should be cleared from memory.',
    ).toBeNull();
    expect(
      getRefreshToken(),
      'After logout, refresh token should be cleared from memory.',
    ).toBeNull();
  });

  it('initialize is a no-op (in-memory tokens do not survive reload)', () => {
    useAuthStore.getState().initialize();

    expect(
      useAuthStore.getState().isAuthenticated,
      'initialize should leave isAuthenticated as false since in-memory tokens do not persist across page loads.',
    ).toBe(false);
  });
});
