import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useProfileActions } from '@/features/users/hooks/useProfileActions';
import * as usersApi from '@/api/users';

vi.mock('@/api/users');
vi.mock('sonner', () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

const USER_ID = '00000000-0000-0000-0000-000000000001';

beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(usersApi.likeUser).mockResolvedValue({ matched: false });
  vi.mocked(usersApi.unlikeUser).mockResolvedValue(undefined);
  vi.mocked(usersApi.blockUser).mockResolvedValue(undefined);
  vi.mocked(usersApi.unblockUser).mockResolvedValue(undefined);
  vi.mocked(usersApi.reportUser).mockResolvedValue({ message: 'Report submitted' });
});

describe('useProfileActions', () => {
  it('starts with isLiked=false and isBlocked=false', () => {
    const { result } = renderHook(() => useProfileActions(USER_ID));

    expect(result.current.isLiked, 'Initial isLiked should be false.').toBe(false);
    expect(result.current.isBlocked, 'Initial isBlocked should be false.').toBe(false);
    expect(result.current.actionLoading).toBe(false);
  });

  it('sets isLiked=true after successful like', async () => {
    const { result } = renderHook(() => useProfileActions(USER_ID));

    await act(async () => {
      await result.current.handleLike();
    });

    expect(
      result.current.isLiked,
      'isLiked should be true after successful handleLike call.',
    ).toBe(true);
    expect(usersApi.likeUser).toHaveBeenCalledWith(USER_ID);
  });

  it('sets isLiked=false after successful unlike', async () => {
    const { result } = renderHook(() => useProfileActions(USER_ID));

    await act(async () => {
      await result.current.handleLike();
    });
    await act(async () => {
      await result.current.handleUnlike();
    });

    expect(
      result.current.isLiked,
      'isLiked should be false after handleUnlike.',
    ).toBe(false);
  });

  it('sets isBlocked=true after successful block', async () => {
    const { result } = renderHook(() => useProfileActions(USER_ID));

    await act(async () => {
      await result.current.handleBlock();
    });

    expect(
      result.current.isBlocked,
      'isBlocked should be true after handleBlock.',
    ).toBe(true);
  });

  it('sets isBlocked=false after successful unblock', async () => {
    const { result } = renderHook(() => useProfileActions(USER_ID));

    await act(async () => {
      await result.current.handleBlock();
    });
    await act(async () => {
      await result.current.handleUnblock();
    });

    expect(result.current.isBlocked).toBe(false);
  });

  it('calls reportUser API on handleReport', async () => {
    const { result } = renderHook(() => useProfileActions(USER_ID));

    await act(async () => {
      await result.current.handleReport();
    });

    expect(
      usersApi.reportUser,
      'Should call reportUser with userId and reason.',
    ).toHaveBeenCalledWith(USER_ID, 'inappropriate');
  });

  it('does not call API when userId is undefined', async () => {
    const { result } = renderHook(() => useProfileActions(undefined));

    await act(async () => {
      await result.current.handleLike();
    });

    expect(usersApi.likeUser, 'Should not call API when userId is undefined.').not.toHaveBeenCalled();
  });
});
