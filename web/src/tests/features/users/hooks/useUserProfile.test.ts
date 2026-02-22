import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
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
