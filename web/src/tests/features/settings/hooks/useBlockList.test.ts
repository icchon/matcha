import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useBlockList } from '@/features/settings/hooks/useBlockList';
import * as settingsApi from '@/api/settings';
import type { Block } from '@/types';

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
