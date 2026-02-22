import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useBlockList } from '@/features/settings/hooks/useBlockList';
import * as settingsApi from '@/api/settings';
import type { Block } from '@/types';

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
    ).toBe('Failed to unblock user');
    expect(
      result.current.blocks,
      'blocks should remain unchanged on unblock failure.',
    ).toEqual(blocks);
  });

  it('tracks unblockingId during unblock operation', async () => {
    vi.mocked(settingsApi.getBlockList).mockResolvedValue(blocks);
    let resolveUnblock: () => void;
    vi.mocked(settingsApi.unblockUser).mockReturnValue(
      new Promise((resolve) => { resolveUnblock = resolve; }),
    );

    const { result } = renderHook(() => useBlockList());

    await act(async () => {
      await result.current.fetchBlockList();
    });

    expect(result.current.unblockingId).toBeNull();

    let promise: Promise<void>;
    act(() => {
      promise = result.current.unblock('user-2');
    });

    expect(
      result.current.unblockingId,
      'unblockingId should be set to the user being unblocked.',
    ).toBe('user-2');

    await act(async () => {
      resolveUnblock!();
      await promise!;
    });

    expect(result.current.unblockingId).toBeNull();
  });

  it('maps ApiClientError status to user-friendly messages', async () => {
    const { ApiClientError: MockApiClientError } = await import('@/api/client');
    vi.mocked(settingsApi.getBlockList).mockRejectedValue(
      new MockApiClientError(401, { error: 'unauthorized' }),
    );

    const { result } = renderHook(() => useBlockList());

    await act(async () => {
      await result.current.fetchBlockList();
    });

    expect(
      result.current.error,
      'Should show user-friendly message for 401 error.',
    ).toBe('Session expired. Please log in again.');
  });
});
