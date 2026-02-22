import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useUserProfile } from '@/features/users/hooks/useUserProfile';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

vi.mock('@/api/users');

const mockProfile: UserProfileDetail = {
  userId: '00000000-0000-0000-0000-000000000001',
  firstName: 'Alice',
  lastName: 'Smith',
  username: 'alice',
  gender: 'female',
  sexualPreference: 'bisexual',
  birthday: '1995-06-15',
  occupation: 'Designer',
  biography: 'Love art',
  locationName: 'Tokyo',
  fameRating: 75,
  pictures: [],
  tags: [],
  isOnline: true,
  lastConnection: null,
  distance: 3.5,
};

beforeEach(() => {
  vi.resetAllMocks();
});

describe('useUserProfile', () => {
  it('returns loading state initially', () => {
    vi.mocked(usersApi.getUserProfile).mockReturnValue(new Promise(() => {}));
    const { result } = renderHook(() => useUserProfile('00000000-0000-0000-0000-000000000001'));

    expect(
      result.current.isLoading,
      'Should be loading while fetching profile.',
    ).toBe(true);
    expect(result.current.profile).toBeNull();
    expect(result.current.error).toBeNull();
  });

  it('returns profile after successful fetch', async () => {
    vi.mocked(usersApi.getUserProfile).mockResolvedValue(mockProfile);
    const { result } = renderHook(() => useUserProfile('00000000-0000-0000-0000-000000000001'));

    await waitFor(() => {
      expect(
        result.current.isLoading,
        'Should stop loading after fetch resolves.',
      ).toBe(false);
    });

    expect(result.current.profile).toEqual(mockProfile);
    expect(result.current.error).toBeNull();
  });

  it('returns error after failed fetch', async () => {
    vi.mocked(usersApi.getUserProfile).mockRejectedValue(new Error('Not found'));
    const { result } = renderHook(() => useUserProfile('00000000-0000-0000-0000-000000000001'));

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(
      result.current.error,
      'Should capture error message from rejected promise.',
    ).toBe('Not found');
    expect(result.current.profile).toBeNull();
  });

  it('resets profile to null when userId changes to prevent stale data flash', async () => {
    const secondProfile: UserProfileDetail = {
      ...mockProfile,
      userId: '00000000-0000-0000-0000-000000000002',
      firstName: 'Bob',
    };

    vi.mocked(usersApi.getUserProfile).mockResolvedValue(mockProfile);
    const { result, rerender } = renderHook(
      ({ userId }: { userId: string }) => useUserProfile(userId),
      { initialProps: { userId: '00000000-0000-0000-0000-000000000001' } },
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });
    expect(result.current.profile?.firstName).toBe('Alice');

    // Simulate slow second fetch so we can observe the null reset
    let resolveSecond: (value: UserProfileDetail) => void;
    vi.mocked(usersApi.getUserProfile).mockReturnValue(
      new Promise((resolve) => { resolveSecond = resolve; }),
    );

    rerender({ userId: '00000000-0000-0000-0000-000000000002' });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(true);
    });
    expect(
      result.current.profile,
      'Profile should be null during userId change to prevent stale data flash. Check setProfile(null) in useEffect.',
    ).toBeNull();

    // Resolve the second fetch
    await act(async () => {
      resolveSecond!(secondProfile);
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });
    expect(result.current.profile?.firstName).toBe('Bob');
  });

  it('does not fetch when userId is undefined and sets isLoading to false', async () => {
    const { result } = renderHook(() => useUserProfile(undefined));

    await waitFor(() => {
      expect(
        result.current.isLoading,
        'isLoading should be false when userId is undefined to avoid stuck loading state.',
      ).toBe(false);
    });

    expect(
      usersApi.getUserProfile,
      'Should not call API when userId is undefined.',
    ).not.toHaveBeenCalled();
  });
});
