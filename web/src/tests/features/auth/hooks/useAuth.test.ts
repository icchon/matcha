import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useAuth } from '@/features/auth/hooks/useAuth';
import * as authApi from '@/api/auth';
import { useAuthStore } from '@/stores/authStore';
import { clearTokens } from '@/api/client';
import type { LoginResponse, SignupResponse, MessageResponse } from '@/types';

vi.mock('@/api/auth');
vi.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
}));
vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
  },
}));

const mockNavigate = vi.fn();
const { toast } = await import('sonner');

const loginResponse: LoginResponse = {
  userId: 'user-123',
  isVerified: true,
  authMethod: 'local',
  accessToken: 'access-token',
  refreshToken: 'refresh-token',
};

beforeEach(() => {
  vi.resetAllMocks();
  clearTokens();
  useAuthStore.setState({
    userId: null,
    isAuthenticated: false,
    isVerified: false,
    authMethod: null,
  });
});

describe('useAuth', () => {
  describe('login', () => {
    it('calls login API, updates store, and navigates to home', async () => {
      vi.mocked(authApi.login).mockResolvedValue(loginResponse);

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.login({ email: 'user@example.com', password: 'password123' });
      });

      expect(authApi.login).toHaveBeenCalledWith({
        email: 'user@example.com',
        password: 'password123',
      });
      expect(
        useAuthStore.getState().isAuthenticated,
        'After successful login, store should be authenticated.',
      ).toBe(true);
      expect(useAuthStore.getState().userId).toBe('user-123');
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });

    it('sets error on login failure', async () => {
      vi.mocked(authApi.login).mockRejectedValue(new Error('Invalid credentials'));

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.login({ email: 'user@example.com', password: 'wrong' });
      });

      expect(
        result.current.error,
        'Error should be set on login failure.',
      ).toBe('Invalid credentials');
      expect(toast.error).toHaveBeenCalled();
    });
  });

  describe('signup', () => {
    it('calls signup API, shows success toast, and navigates to login', async () => {
      const signupResponse: SignupResponse = { message: 'Check email' };
      vi.mocked(authApi.signup).mockResolvedValue(signupResponse);

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.signup({ email: 'new@example.com', password: 'password123' });
      });

      expect(authApi.signup).toHaveBeenCalledWith({
        email: 'new@example.com',
        password: 'password123',
      });
      expect(toast.success).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith('/login');
    });
  });

  describe('logout', () => {
    it('calls logout API, clears store, and navigates to login', async () => {
      vi.mocked(authApi.logout).mockResolvedValue(undefined);
      useAuthStore.getState().login({
        ...loginResponse,
      });

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.logout();
      });

      expect(authApi.logout).toHaveBeenCalled();
      expect(useAuthStore.getState().isAuthenticated).toBe(false);
      expect(mockNavigate).toHaveBeenCalledWith('/login');
    });
  });

  describe('verifyEmail', () => {
    it('calls verifyEmail API and shows success toast', async () => {
      vi.mocked(authApi.verifyEmail).mockResolvedValue(undefined);

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.verifyEmail('token-123');
      });

      expect(authApi.verifyEmail).toHaveBeenCalledWith('token-123');
      expect(toast.success).toHaveBeenCalled();
    });

    it('sets error on verification failure', async () => {
      vi.mocked(authApi.verifyEmail).mockRejectedValue(new Error('Invalid token'));

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.verifyEmail('bad-token');
      });

      expect(result.current.error).toBe('Invalid token');
      expect(toast.error).toHaveBeenCalled();
    });
  });

  describe('forgotPassword', () => {
    it('calls forgotPassword API and shows success toast', async () => {
      const response: MessageResponse = { message: 'Check email' };
      vi.mocked(authApi.forgotPassword).mockResolvedValue(response);

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.forgotPassword({ email: 'user@example.com' });
      });

      expect(authApi.forgotPassword).toHaveBeenCalledWith({ email: 'user@example.com' });
      expect(toast.success).toHaveBeenCalled();
    });
  });

  describe('resetPassword', () => {
    it('calls resetPassword API and navigates to login', async () => {
      vi.mocked(authApi.resetPassword).mockResolvedValue(undefined);

      const { result } = renderHook(() => useAuth());

      await act(async () => {
        await result.current.resetPassword({ token: 'tok', password: 'newpass123' });
      });

      expect(authApi.resetPassword).toHaveBeenCalledWith({
        token: 'tok',
        password: 'newpass123',
      });
      expect(toast.success).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith('/login');
    });
  });

  describe('isLoading state', () => {
    it('is true while API call is in progress', async () => {
      let resolveLogin: (v: LoginResponse) => void;
      vi.mocked(authApi.login).mockReturnValue(
        new Promise((resolve) => {
          resolveLogin = resolve;
        }),
      );

      const { result } = renderHook(() => useAuth());
      expect(result.current.isLoading).toBe(false);

      let loginPromise: Promise<void>;
      act(() => {
        loginPromise = result.current.login({
          email: 'user@example.com',
          password: 'password123',
        });
      });

      expect(
        result.current.isLoading,
        'isLoading should be true while the API call is pending.',
      ).toBe(true);

      await act(async () => {
        resolveLogin!(loginResponse);
        await loginPromise!;
      });

      expect(result.current.isLoading).toBe(false);
    });
  });
});
