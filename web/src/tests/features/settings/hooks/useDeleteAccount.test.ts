import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useDeleteAccount } from '@/features/settings/hooks/useDeleteAccount';
import * as settingsApi from '@/api/settings';
import { ApiClientError } from '@/api/client';

vi.mock('@/api/settings');
vi.mock('@/api/client', () => ({
  ApiClientError: class extends Error {
    readonly status: number;
    readonly body: { error: string };
    constructor(status: number, body: { error: string }) {
      super(body.error);
      this.name = 'ApiClientError';
      this.status = status;
      this.body = body;
    }
  },
}));

const mockNavigate = vi.fn();
const mockLogout = vi.fn();

vi.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
}));
vi.mock('@/stores/authStore', () => ({
  useAuthStore: (selector: (s: { logout: () => void }) => unknown) =>
    selector({ logout: mockLogout }),
}));
vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

const { toast } = await import('sonner');

beforeEach(() => {
  vi.resetAllMocks();
});

describe('useDeleteAccount', () => {
  it('calls deleteAccount API with currentPassword, logs out, shows toast, and navigates to login', async () => {
    vi.mocked(settingsApi.deleteAccount).mockResolvedValue(undefined);

    const { result } = renderHook(() => useDeleteAccount());

    await act(async () => {
      await result.current.deleteAccount({ currentPassword: 'mypass123' });
    });

    expect(
      settingsApi.deleteAccount,
      'deleteAccount should call API with currentPassword.',
    ).toHaveBeenCalledWith({ currentPassword: 'mypass123' });
    expect(
      mockLogout,
      'logout should be called after account deletion to clear auth state.',
    ).toHaveBeenCalled();
    expect(toast.success).toHaveBeenCalledWith('Account deleted successfully.');
    expect(mockNavigate).toHaveBeenCalledWith('/login');
  });

  it('maps 401 errors to user-friendly message', async () => {
    const { ApiClientError: MockApiClientError } = await import('@/api/client');
    vi.mocked(settingsApi.deleteAccount).mockRejectedValue(
      new MockApiClientError(401, { error: 'unauthorized' }),
    );

    const { result } = renderHook(() => useDeleteAccount());

    await act(async () => {
      await result.current.deleteAccount({ currentPassword: 'wrong' });
    });

    expect(
      result.current.error,
      'Should show user-friendly message for 401 status.',
    ).toBe('Incorrect password. Please try again.');
    expect(toast.error).toHaveBeenCalledWith('Incorrect password. Please try again.');
  });

  it('sets generic error for non-ApiClientError failures', async () => {
    vi.mocked(settingsApi.deleteAccount).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useDeleteAccount());

    await act(async () => {
      await result.current.deleteAccount({ currentPassword: 'mypass123' });
    });

    expect(
      result.current.error,
      'Should set a generic error for non-API errors.',
    ).toBe('Failed to delete account');
    expect(toast.error).toHaveBeenCalled();
  });
});
