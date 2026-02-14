import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useChangePassword, useDeleteAccount, useBlockList } from '@/features/settings/hooks/useSettings';
import * as settingsApi from '@/api/settings';
import type { Block, MessageResponse } from '@/types';

vi.mock('@/api/settings');
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

describe('useChangePassword', () => {
  it('calls changePassword API and shows success toast', async () => {
    const response: MessageResponse = { message: 'Password changed' };
    vi.mocked(settingsApi.changePassword).mockResolvedValue(response);

    const { result } = renderHook(() => useChangePassword());

    await act(async () => {
      await result.current.changePassword({
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
  });

  it('sets error on changePassword failure', async () => {
    vi.mocked(settingsApi.changePassword).mockRejectedValue(new Error('Wrong password'));

    const { result } = renderHook(() => useChangePassword());

    await act(async () => {
      await result.current.changePassword({
        currentPassword: 'wrong',
        newPassword: 'newpass123',
      });
    });

    expect(
      result.current.error,
      'Error should be set when changePassword fails. Check error extraction.',
    ).toBe('Wrong password');
    expect(toast.error).toHaveBeenCalled();
  });

  it('tracks isLoading state during API call', async () => {
    let resolveApi: (v: MessageResponse) => void;
    vi.mocked(settingsApi.changePassword).mockReturnValue(
      new Promise((resolve) => { resolveApi = resolve; }),
    );

    const { result } = renderHook(() => useChangePassword());
    expect(result.current.isLoading).toBe(false);

    let promise: Promise<void>;
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

describe('useDeleteAccount', () => {
  it('calls deleteAccount API, shows toast, and navigates to login', async () => {
    vi.mocked(settingsApi.deleteAccount).mockResolvedValue(undefined);

    const { result } = renderHook(() => useDeleteAccount());

    await act(async () => {
      await result.current.deleteAccount();
    });

    expect(settingsApi.deleteAccount).toHaveBeenCalled();
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

describe('useBlockList', () => {
  const blocks: Block[] = [
    { blockerId: 'me', blockedId: 'user-2' },
    { blockerId: 'me', blockedId: 'user-3' },
  ];

  it('fetches block list on load', async () => {
    vi.mocked(settingsApi.getBlockList).mockResolvedValue(blocks);

    const { result } = renderHook(() => useBlockList());

    await act(async () => {
      await result.current.fetchBlockList();
    });

    expect(settingsApi.getBlockList).toHaveBeenCalled();
    expect(
      result.current.blocks,
      'blocks should contain the fetched block list.',
    ).toEqual(blocks);
  });

  it('calls unblockUser API and removes user from list', async () => {
    vi.mocked(settingsApi.getBlockList).mockResolvedValue(blocks);
    vi.mocked(settingsApi.unblockUser).mockResolvedValue(undefined);

    const { result } = renderHook(() => useBlockList());

    await act(async () => {
      await result.current.fetchBlockList();
    });

    await act(async () => {
      await result.current.unblock('user-2');
    });

    expect(settingsApi.unblockUser).toHaveBeenCalledWith('user-2');
    expect(toast.success).toHaveBeenCalledWith('User unblocked.');
    expect(
      result.current.blocks,
      'After unblocking user-2, blocks should only contain user-3.',
    ).toEqual([{ blockerId: 'me', blockedId: 'user-3' }]);
  });

  it('sets error on unblock failure', async () => {
    vi.mocked(settingsApi.getBlockList).mockResolvedValue(blocks);
    vi.mocked(settingsApi.unblockUser).mockRejectedValue(new Error('Failed'));

    const { result } = renderHook(() => useBlockList());

    await act(async () => {
      await result.current.fetchBlockList();
    });

    await act(async () => {
      await result.current.unblock('user-2');
    });

    expect(
      result.current.error,
      'Error should be set on unblock failure.',
    ).toBe('Failed');
    expect(
      result.current.blocks,
      'blocks should remain unchanged on unblock failure.',
    ).toEqual(blocks);
  });
});
