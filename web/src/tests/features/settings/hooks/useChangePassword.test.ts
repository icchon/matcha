import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useChangePassword } from '@/features/settings/hooks/useChangePassword';
import * as settingsApi from '@/api/settings';
import { ApiClientError } from '@/api/client';
import type { MessageResponse } from '@/types';

const mockNavigate = vi.fn();
const mockLogout = vi.fn();

vi.mock('@/api/settings');
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

describe('useChangePassword', () => {
  it('calls changePassword API, shows success toast, logs out, and navigates to login', async () => {
    const response: MessageResponse = { message: 'Password changed' };
    vi.mocked(settingsApi.changePassword).mockResolvedValue(response);

    const { result } = renderHook(() => useChangePassword());

    let success: boolean = false;
    await act(async () => {
      success = await result.current.changePassword({
        currentPassword: 'oldpass123',
        newPassword: 'newpass123',
      });
    });

    expect(
      settingsApi.changePassword,
      'changePassword should call the API with currentPassword and newPassword.',
    ).toHaveBeenCalledWith({
      currentPassword: 'oldpass123',
      newPassword: 'newpass123',
    });
    expect(toast.success).toHaveBeenCalledWith('Password changed successfully! Please log in again.');
    expect(
      mockLogout,
      'logout should be called to clear session after password change.',
    ).toHaveBeenCalled();
    expect(
      mockNavigate,
      'Should navigate to /login after password change.',
    ).toHaveBeenCalledWith('/login');
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(
      success,
      'changePassword should return true on success for conditional form reset.',
    ).toBe(true);
  });

  it('sets error on changePassword failure and returns false without navigating', async () => {
    vi.mocked(settingsApi.changePassword).mockRejectedValue(new Error('Wrong password'));

    const { result } = renderHook(() => useChangePassword());

    let success: boolean = true;
    await act(async () => {
      success = await result.current.changePassword({
        currentPassword: 'wrong',
        newPassword: 'newpass123',
      });
    });

    expect(
      result.current.error,
      'Error should be set when changePassword fails. Check error extraction.',
    ).toBe('Wrong password');
    expect(toast.error).toHaveBeenCalled();
    expect(
      mockLogout,
      'logout should not be called on failure.',
    ).not.toHaveBeenCalled();
    expect(
      mockNavigate,
      'Should not navigate on failure.',
    ).not.toHaveBeenCalled();
    expect(
      success,
      'changePassword should return false on failure so form is not reset.',
    ).toBe(false);
  });

  it('tracks isLoading state during API call', async () => {
    let resolveApi: (v: MessageResponse) => void;
    vi.mocked(settingsApi.changePassword).mockReturnValue(
      new Promise((resolve) => { resolveApi = resolve; }),
    );

    const { result } = renderHook(() => useChangePassword());
    expect(result.current.isLoading).toBe(false);

    let promise: Promise<boolean>;
    act(() => {
      promise = result.current.changePassword({
        currentPassword: 'old123456',
        newPassword: 'new123456',
      });
    });

    expect(
      result.current.isLoading,
      'isLoading should be true while API call is in progress.',
    ).toBe(true);

    await act(async () => {
      resolveApi!({ message: 'ok' });
      await promise!;
    });

    expect(result.current.isLoading).toBe(false);
  });

  it('maps 401 ApiClientError to "Incorrect current password." message', async () => {
    vi.mocked(settingsApi.changePassword).mockRejectedValue(
      new ApiClientError(401, { error: 'Unauthorized' }),
    );

    const { result } = renderHook(() => useChangePassword());

    await act(async () => {
      await result.current.changePassword({
        currentPassword: 'wrong',
        newPassword: 'Newpass1!',
      });
    });

    expect(
      result.current.error,
      'A 401 status should be mapped to a user-friendly "Incorrect current password." message.',
    ).toBe('Incorrect current password.');
    expect(toast.error).toHaveBeenCalledWith('Incorrect current password.');
  });

  it('maps 422 ApiClientError to "Invalid password format." message', async () => {
    vi.mocked(settingsApi.changePassword).mockRejectedValue(
      new ApiClientError(422, { error: 'Validation failed' }),
    );

    const { result } = renderHook(() => useChangePassword());

    await act(async () => {
      await result.current.changePassword({
        currentPassword: 'oldpass123',
        newPassword: 'weak',
      });
    });

    expect(
      result.current.error,
      'A 422 status should be mapped to a user-friendly "Invalid password format." message.',
    ).toBe('Invalid password format.');
    expect(toast.error).toHaveBeenCalledWith('Invalid password format.');
  });

  it('re-throws non-mapped ApiClientError with original message', async () => {
    vi.mocked(settingsApi.changePassword).mockRejectedValue(
      new ApiClientError(500, { error: 'Internal server error' }),
    );

    const { result } = renderHook(() => useChangePassword());

    await act(async () => {
      await result.current.changePassword({
        currentPassword: 'oldpass123',
        newPassword: 'Newpass1!',
      });
    });

    expect(
      result.current.error,
      'Non-mapped ApiClientError should fall through with original error message.',
    ).toBe('Internal server error');
  });
});
