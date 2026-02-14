import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useChangePassword } from '@/features/settings/hooks/useChangePassword';
import * as settingsApi from '@/api/settings';
import type { MessageResponse } from '@/types';

vi.mock('@/api/settings');
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
  it('calls changePassword API, shows success toast, and returns true', async () => {
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
    expect(toast.success).toHaveBeenCalledWith('Password changed successfully!');
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(
      success,
      'changePassword should return true on success for conditional form reset.',
    ).toBe(true);
  });

  it('sets error on changePassword failure and returns false', async () => {
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
});
