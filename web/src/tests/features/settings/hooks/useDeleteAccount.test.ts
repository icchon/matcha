import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useDeleteAccount } from '@/features/settings/hooks/useDeleteAccount';
import * as settingsApi from '@/api/settings';
import { clearTokens } from '@/api/client';

vi.mock('@/api/settings');
vi.mock('@/api/client', () => ({
  clearTokens: vi.fn(),
}));
vi.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
}));
vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

const mockNavigate = vi.fn();
const { toast } = await import('sonner');

beforeEach(() => {
  vi.resetAllMocks();
});

describe('useDeleteAccount', () => {
  it('calls deleteAccount API, clears tokens, shows toast, and navigates to login', async () => {
    vi.mocked(settingsApi.deleteAccount).mockResolvedValue(undefined);

    const { result } = renderHook(() => useDeleteAccount());

    await act(async () => {
      await result.current.deleteAccount();
    });

    expect(settingsApi.deleteAccount).toHaveBeenCalled();
    expect(
      clearTokens,
      'clearTokens should be called after account deletion to clear in-memory auth state.',
    ).toHaveBeenCalled();
    expect(toast.success).toHaveBeenCalledWith('Account deleted successfully.');
    expect(mockNavigate).toHaveBeenCalledWith('/login');
  });

  it('sets error on deleteAccount failure', async () => {
    vi.mocked(settingsApi.deleteAccount).mockRejectedValue(new Error('Server error'));

    const { result } = renderHook(() => useDeleteAccount());

    await act(async () => {
      await result.current.deleteAccount();
    });

    expect(
      result.current.error,
      'Error should be set when deleteAccount fails.',
    ).toBe('Server error');
    expect(toast.error).toHaveBeenCalled();
  });
});
