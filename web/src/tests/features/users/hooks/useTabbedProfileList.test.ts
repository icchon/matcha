import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useTabbedProfileList } from '@/features/users/hooks/useTabbedProfileList';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

vi.mock('@/api/users');
vi.mock('sonner', () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

const mockProfile1: UserProfileDetail = {
  userId: '00000000-0000-0000-0000-000000000002',
  firstName: 'Bob',
  lastName: 'Jones',
  username: 'bob',
  gender: 'male',
  sexualPreference: 'heterosexual',
  birthday: '1992-03-20',
  occupation: 'Dev',
  biography: 'Hello',
  locationName: 'Paris',
  fameRating: 50,
  pictures: [],
  tags: [],
  isOnline: false,
  lastConnection: '2024-01-15T10:00:00Z',
};

const mockProfile2: UserProfileDetail = {
  ...mockProfile1,
  userId: '00000000-0000-0000-0000-000000000003',
  firstName: 'Carol',
};

interface MockItem {
  readonly myId: string;
  readonly theirId: string;
}

const defaultConfig = {
  fetchMyList: vi.fn<() => Promise<readonly MockItem[]>>().mockResolvedValue([
    { myId: '00000000-0000-0000-0000-000000000002', theirId: '' },
  ]),
  fetchTheirList: vi.fn<() => Promise<readonly MockItem[]>>().mockResolvedValue([
    { myId: '', theirId: '00000000-0000-0000-0000-000000000003' },
  ]),
  extractMyIds: (items: readonly MockItem[]) => items.map((i) => i.myId).filter(Boolean),
  extractTheirIds: (items: readonly MockItem[]) => items.map((i) => i.theirId).filter(Boolean),
  errorMessage: 'Failed to load',
};

beforeEach(() => {
  vi.resetAllMocks();
  defaultConfig.fetchMyList.mockResolvedValue([
    { myId: '00000000-0000-0000-0000-000000000002', theirId: '' },
  ]);
  defaultConfig.fetchTheirList.mockResolvedValue([
    { myId: '', theirId: '00000000-0000-0000-0000-000000000003' },
  ]);
  vi.mocked(usersApi.getUserProfile).mockImplementation(async (userId: string) => {
    if (userId === '00000000-0000-0000-0000-000000000002') return mockProfile1;
    if (userId === '00000000-0000-0000-0000-000000000003') return mockProfile2;
    throw new Error('Not found');
  });
});

describe('useTabbedProfileList', () => {
  it('starts in loading state', () => {
    defaultConfig.fetchMyList.mockReturnValue(new Promise(() => {}));
    defaultConfig.fetchTheirList.mockReturnValue(new Promise(() => {}));

    const { result } = renderHook(() =>
      useTabbedProfileList<MockItem, string>(defaultConfig, 'tab1'),
    );

    expect(
      result.current.isLoading,
      'Should be loading initially while fetching lists.',
    ).toBe(true);
  });

  it('resolves profiles into myProfiles and theirProfiles', async () => {
    const { result } = renderHook(() =>
      useTabbedProfileList<MockItem, string>(defaultConfig, 'tab1'),
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(
      result.current.myProfiles.length,
      'Should have 1 profile in myProfiles after fetch.',
    ).toBe(1);
    expect(result.current.myProfiles[0].firstName).toBe('Bob');
    expect(result.current.theirProfiles.length).toBe(1);
    expect(result.current.theirProfiles[0].firstName).toBe('Carol');
  });

  it('shows generic error for 5xx errors instead of leaking server details', async () => {
    const serverError = Object.assign(new Error('Internal Server Error'), { status: 500 });
    defaultConfig.fetchMyList.mockRejectedValue(serverError);
    const { toast } = await import('sonner');
    const { result } = renderHook(() =>
      useTabbedProfileList<MockItem, string>(defaultConfig, 'tab1'),
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(
      toast.error,
      'Should show generic message for 5xx errors via shared getErrorMessage. Check useTabbedProfileList error handling.',
    ).toHaveBeenCalledWith('Something went wrong. Please try again later.');
  });

  it('shows authorization error for 401/403 errors', async () => {
    const authError = Object.assign(new Error('Unauthorized'), { status: 401 });
    defaultConfig.fetchMyList.mockRejectedValue(authError);
    const { toast } = await import('sonner');
    const { result } = renderHook(() =>
      useTabbedProfileList<MockItem, string>(defaultConfig, 'tab1'),
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(
      toast.error,
      'Should show authorization error for 401 via shared getErrorMessage.',
    ).toHaveBeenCalledWith('You are not authorized to perform this action.');
  });

  it('allows switching active tab', async () => {
    const { result } = renderHook(() =>
      useTabbedProfileList<MockItem, string>(defaultConfig, 'tab1'),
    );

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.activeTab).toBe('tab1');

    act(() => {
      result.current.setActiveTab('tab2');
    });

    expect(
      result.current.activeTab,
      'Active tab should change when setActiveTab is called.',
    ).toBe('tab2');
  });
});
