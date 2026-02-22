import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  getUserProfile,
  likeUser,
  unlikeUser,
  blockUser,
  unblockUser,
  reportUser,
  getLikedUsers,
  getWhoLikedMe,
  getViewedUsers,
  getWhoViewedMe,
} from '@/api/users';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { RawUserProfileResponse } from '@/types/raw';

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const mockGet = vi.mocked(apiClient.get);
const mockPost = vi.mocked(apiClient.post);
const mockDelete = vi.mocked(apiClient.delete);

const UUID_1 = '00000000-0000-0000-0000-000000000001';
const UUID_2 = '00000000-0000-0000-0000-000000000002';
const UUID_3 = '00000000-0000-0000-0000-000000000003';
const UUID_4 = '00000000-0000-0000-0000-000000000004';
const UUID_5 = '00000000-0000-0000-0000-000000000005';
const UUID_6 = '00000000-0000-0000-0000-000000000006';

const rawProfile: RawUserProfileResponse = {
  user_id: UUID_1,
  first_name: 'John',
  last_name: 'Doe',
  username: 'johndoe',
  gender: 'male',
  sexual_preference: 'heterosexual',
  birthday: '1990-01-01',
  occupation: 'Engineer',
  biography: 'Hello world',
  location_name: 'Paris',
  fame_rating: 42,
  pictures: [
    { id: 1, user_id: UUID_1, url: '/images/1.jpg', is_profile_pic: true, created_at: '2024-01-01' },
  ],
  tags: [{ id: 1, name: 'coding' }],
  is_online: true,
  last_connection: '2024-01-01T12:00:00Z',
  distance: 5.2,
};

beforeEach(() => {
  vi.resetAllMocks();
});

describe('getUserProfile', () => {
  it('calls GET /users/:userId/profile and maps snake_case to camelCase', async () => {
    mockGet.mockResolvedValue(rawProfile);

    const result = await getUserProfile(UUID_1);

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.PROFILE.GET(UUID_1));
    expect(
      result.userId,
      'getUserProfile should map user_id to userId. Check mapUserProfileResponse.',
    ).toBe(UUID_1);
    expect(result.firstName).toBe('John');
    expect(result.lastName).toBe('Doe');
    expect(result.username).toBe('johndoe');
    expect(result.gender).toBe('male');
    expect(result.sexualPreference).toBe('heterosexual');
    expect(result.birthday).toBe('1990-01-01');
    expect(result.occupation).toBe('Engineer');
    expect(result.biography).toBe('Hello world');
    expect(result.locationName).toBe('Paris');
    expect(result.fameRating).toBe(42);
    expect(result.isOnline).toBe(true);
    expect(result.lastConnection).toBe('2024-01-01T12:00:00Z');
    expect(result.distance).toBe(5.2);
  });

  it('maps pictures from snake_case to camelCase', async () => {
    mockGet.mockResolvedValue(rawProfile);

    const result = await getUserProfile(UUID_1);

    expect(
      result.pictures,
      'Pictures should be mapped with camelCase keys. Check mapPicture.',
    ).toHaveLength(1);
    expect(result.pictures[0].userId).toBe(UUID_1);
    expect(result.pictures[0].url).toBe('/images/1.jpg');
    expect(result.pictures[0].isProfilePic).toBe(true);
    expect(result.pictures[0].createdAt).toBe('2024-01-01');
  });

  it('maps tags correctly', async () => {
    mockGet.mockResolvedValue(rawProfile);

    const result = await getUserProfile(UUID_1);

    expect(result.tags).toHaveLength(1);
    expect(result.tags[0]).toEqual({ id: 1, name: 'coding' });
  });

  it('rejects invalid userId format', async () => {
    await expect(
      getUserProfile('invalid-id'),
    ).rejects.toThrow('Invalid user ID format');
  });

  it('strips http: picture URLs for security (only allows https: and same-origin)', async () => {
    const profileWithHttpPic = {
      ...rawProfile,
      pictures: [
        { id: 1, user_id: UUID_1, url: 'http://evil.com/pic.jpg', is_profile_pic: true, created_at: '2024-01-01' },
      ],
    };
    mockGet.mockResolvedValue(profileWithHttpPic);

    const result = await getUserProfile(UUID_1);

    expect(
      result.pictures[0].url,
      'http: URLs from external origins should be stripped. Only https: and same-origin allowed.',
    ).toBe('');
  });

  it('allows https: picture URLs', async () => {
    const profileWithHttpsPic = {
      ...rawProfile,
      pictures: [
        { id: 1, user_id: UUID_1, url: 'https://cdn.example.com/pic.jpg', is_profile_pic: true, created_at: '2024-01-01' },
      ],
    };
    mockGet.mockResolvedValue(profileWithHttpsPic);

    const result = await getUserProfile(UUID_1);

    expect(
      result.pictures[0].url,
      'https: URLs should be allowed.',
    ).toBe('https://cdn.example.com/pic.jpg');
  });
});

describe('likeUser', () => {
  it('calls POST /users/:userId/like and returns match status', async () => {
    mockPost.mockResolvedValue({ matched: true });

    const result = await likeUser(UUID_2);

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.USERS.LIKE(UUID_2));
    expect(
      result.matched,
      'likeUser should return { matched: boolean }. Check the POST response.',
    ).toBe(true);
  });
});

describe('unlikeUser', () => {
  it('calls DELETE /users/:userId/like', async () => {
    mockDelete.mockResolvedValue(undefined);

    await unlikeUser(UUID_2);

    expect(
      mockDelete,
      'unlikeUser should call DELETE on the like endpoint. Check API_PATHS.USERS.UNLIKE.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.UNLIKE(UUID_2));
  });
});

describe('blockUser', () => {
  it('calls POST /users/:userId/block', async () => {
    mockPost.mockResolvedValue(undefined);

    await blockUser(UUID_3);

    expect(
      mockPost,
      'blockUser should call POST on the block endpoint.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.BLOCK(UUID_3));
  });
});

describe('unblockUser', () => {
  it('[MOCK] calls DELETE /users/:userId/block and returns success stub', async () => {
    mockDelete.mockResolvedValue(undefined);

    await unblockUser(UUID_3);

    expect(
      mockDelete,
      '[MOCK] unblockUser should call DELETE on the block endpoint. BE-08 #25: endpoint not yet implemented.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.UNBLOCK(UUID_3));
  });
});

describe('reportUser', () => {
  it('[MOCK] throws "Report feature is not yet available" error', async () => {
    await expect(
      reportUser(UUID_4, 'spam'),
    ).rejects.toThrow('Report feature is not yet available');
  });
});

describe('getLikedUsers', () => {
  it('calls GET /me/likes and returns like list', async () => {
    const likes = [
      { liker_id: 'me', liked_id: UUID_2, created_at: '2024-01-01' },
    ];
    mockGet.mockResolvedValue(likes);

    const result = await getLikedUsers();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.USERS.MY_LIKES);
    expect(
      result,
      'getLikedUsers should map snake_case to camelCase. Check mapLike.',
    ).toHaveLength(1);
    expect(result[0].likerId).toBe('me');
    expect(result[0].likedId).toBe(UUID_2);
    expect(result[0].createdAt).toBe('2024-01-01');
  });
});

describe('getWhoLikedMe', () => {
  it('calls GET /me/profile/likes and returns like list', async () => {
    const likes = [
      { liker_id: UUID_5, liked_id: 'me', created_at: '2024-02-01' },
    ];
    mockGet.mockResolvedValue(likes);

    const result = await getWhoLikedMe();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.PROFILE.WHO_LIKED_ME);
    expect(result).toHaveLength(1);
    expect(result[0].likerId).toBe(UUID_5);
  });
});

describe('getViewedUsers', () => {
  it('calls GET /me/views and returns view list', async () => {
    const views = [
      { viewer_id: 'me', viewed_id: UUID_2, view_time: '2024-01-01T10:00:00Z' },
    ];
    mockGet.mockResolvedValue(views);

    const result = await getViewedUsers();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.USERS.MY_VIEWS);
    expect(
      result,
      'getViewedUsers should map snake_case to camelCase. Check mapView.',
    ).toHaveLength(1);
    expect(result[0].viewerId).toBe('me');
    expect(result[0].viewedId).toBe(UUID_2);
    expect(result[0].viewTime).toBe('2024-01-01T10:00:00Z');
  });
});

describe('getWhoViewedMe', () => {
  it('calls GET /me/profile/views and returns view list', async () => {
    const views = [
      { viewer_id: UUID_6, viewed_id: 'me', view_time: '2024-02-01T10:00:00Z' },
    ];
    mockGet.mockResolvedValue(views);

    const result = await getWhoViewedMe();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.PROFILE.WHO_VIEWED_ME);
    expect(result).toHaveLength(1);
    expect(result[0].viewerId).toBe(UUID_6);
  });
});
