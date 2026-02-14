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

const rawProfile: RawUserProfileResponse = {
  user_id: 'user-1',
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
    { id: 1, user_id: 'user-1', url: '/images/1.jpg', is_profile_pic: true, created_at: '2024-01-01' },
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

    const result = await getUserProfile('user-1');

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.PROFILE.GET('user-1'));
    expect(
      result.userId,
      'getUserProfile should map user_id to userId. Check mapUserProfileResponse.',
    ).toBe('user-1');
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

    const result = await getUserProfile('user-1');

    expect(
      result.pictures,
      'Pictures should be mapped with camelCase keys. Check mapPicture.',
    ).toHaveLength(1);
    expect(result.pictures[0].userId).toBe('user-1');
    expect(result.pictures[0].url).toBe('/images/1.jpg');
    expect(result.pictures[0].isProfilePic).toBe(true);
    expect(result.pictures[0].createdAt).toBe('2024-01-01');
  });

  it('maps tags correctly', async () => {
    mockGet.mockResolvedValue(rawProfile);

    const result = await getUserProfile('user-1');

    expect(result.tags).toHaveLength(1);
    expect(result.tags[0]).toEqual({ id: 1, name: 'coding' });
  });
});

describe('likeUser', () => {
  it('calls POST /users/:userId/like and returns match status', async () => {
    mockPost.mockResolvedValue({ matched: true });

    const result = await likeUser('user-2');

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.USERS.LIKE('user-2'));
    expect(
      result.matched,
      'likeUser should return { matched: boolean }. Check the POST response.',
    ).toBe(true);
  });
});

describe('unlikeUser', () => {
  it('calls DELETE /users/:userId/like', async () => {
    mockDelete.mockResolvedValue(undefined);

    await unlikeUser('user-2');

    expect(
      mockDelete,
      'unlikeUser should call DELETE on the like endpoint. Check API_PATHS.USERS.UNLIKE.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.UNLIKE('user-2'));
  });
});

describe('blockUser', () => {
  it('calls POST /users/:userId/block', async () => {
    mockPost.mockResolvedValue(undefined);

    await blockUser('user-3');

    expect(
      mockPost,
      'blockUser should call POST on the block endpoint.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.BLOCK('user-3'));
  });
});

describe('unblockUser', () => {
  it('[MOCK] calls DELETE /users/:userId/block and returns success stub', async () => {
    mockDelete.mockResolvedValue(undefined);

    await unblockUser('user-3');

    expect(
      mockDelete,
      '[MOCK] unblockUser should call DELETE on the block endpoint. BE-08 #25: endpoint not yet implemented.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.UNBLOCK('user-3'));
  });
});

describe('reportUser', () => {
  it('[MOCK] returns success stub', async () => {
    const result = await reportUser('user-4', 'spam');

    expect(
      result,
      '[MOCK] reportUser should return { message } stub. BE-08 #25: endpoint not yet implemented.',
    ).toEqual({ message: 'Report submitted' });
  });
});

describe('getLikedUsers', () => {
  it('calls GET /me/likes and returns like list', async () => {
    const likes = [
      { liker_id: 'me', liked_id: 'user-2', created_at: '2024-01-01' },
    ];
    mockGet.mockResolvedValue(likes);

    const result = await getLikedUsers();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.USERS.MY_LIKES);
    expect(
      result,
      'getLikedUsers should map snake_case to camelCase. Check mapLike.',
    ).toHaveLength(1);
    expect(result[0].likerId).toBe('me');
    expect(result[0].likedId).toBe('user-2');
    expect(result[0].createdAt).toBe('2024-01-01');
  });
});

describe('getWhoLikedMe', () => {
  it('calls GET /me/profile/likes and returns like list', async () => {
    const likes = [
      { liker_id: 'user-5', liked_id: 'me', created_at: '2024-02-01' },
    ];
    mockGet.mockResolvedValue(likes);

    const result = await getWhoLikedMe();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.PROFILE.WHO_LIKED_ME);
    expect(result).toHaveLength(1);
    expect(result[0].likerId).toBe('user-5');
  });
});

describe('getViewedUsers', () => {
  it('calls GET /me/views and returns view list', async () => {
    const views = [
      { viewer_id: 'me', viewed_id: 'user-2', view_time: '2024-01-01T10:00:00Z' },
    ];
    mockGet.mockResolvedValue(views);

    const result = await getViewedUsers();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.USERS.MY_VIEWS);
    expect(
      result,
      'getViewedUsers should map snake_case to camelCase. Check mapView.',
    ).toHaveLength(1);
    expect(result[0].viewerId).toBe('me');
    expect(result[0].viewedId).toBe('user-2');
    expect(result[0].viewTime).toBe('2024-01-01T10:00:00Z');
  });
});

describe('getWhoViewedMe', () => {
  it('calls GET /me/profile/views and returns view list', async () => {
    const views = [
      { viewer_id: 'user-6', viewed_id: 'me', view_time: '2024-02-01T10:00:00Z' },
    ];
    mockGet.mockResolvedValue(views);

    const result = await getWhoViewedMe();

    expect(mockGet).toHaveBeenCalledWith(API_PATHS.PROFILE.WHO_VIEWED_ME);
    expect(result).toHaveLength(1);
    expect(result[0].viewerId).toBe('user-6');
  });
});
