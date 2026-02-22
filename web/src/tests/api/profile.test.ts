import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  createProfile,
  updateProfile,
  getMyProfile,
} from '@/api/profile';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { UserProfile } from '@/types';

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    upload: vi.fn(),
  },
}));

const mockGet = vi.mocked(apiClient.get);
const mockPost = vi.mocked(apiClient.post);
const mockPut = vi.mocked(apiClient.put);

const sampleProfile: UserProfile = {
  userId: 'user-123',
  firstName: 'John',
  lastName: 'Doe',
  username: 'johndoe',
  gender: 'male',
  sexualPreference: 'heterosexual',
  birthday: '1995-06-15',
  occupation: 'Developer',
  biography: 'Hello world',
  locationName: 'Tokyo',
  fameRating: 42,
};

beforeEach(() => {
  vi.resetAllMocks();
});

describe('createProfile', () => {
  it('calls POST /me/profile/ with profile data', async () => {
    const payload = {
      firstName: 'John',
      lastName: 'Doe',
      username: 'johndoe',
      gender: 'male' as const,
      sexualPreference: 'heterosexual' as const,
      birthday: '1995-06-15',
      biography: 'Hello world',
    };
    mockPost.mockResolvedValue(sampleProfile);

    const result = await createProfile(payload);

    expect(
      mockPost,
      'createProfile should call POST on PROFILE.CREATE path. Check API_PATHS.PROFILE.CREATE.',
    ).toHaveBeenCalledWith(API_PATHS.PROFILE.CREATE, payload);
    expect(
      result.userId,
      'createProfile should return the created profile. Check return type.',
    ).toBe('user-123');
  });
});

describe('updateProfile', () => {
  it('calls PUT /me/profile/ with updated fields', async () => {
    const payload = { biography: 'Updated bio' };
    mockPut.mockResolvedValue({ ...sampleProfile, biography: 'Updated bio' });

    const result = await updateProfile(payload);

    expect(
      mockPut,
      'updateProfile should call PUT on PROFILE.UPDATE path.',
    ).toHaveBeenCalledWith(API_PATHS.PROFILE.UPDATE, payload);
    expect(
      result.biography,
      'updateProfile should return updated profile with new biography.',
    ).toBe('Updated bio');
  });
});

describe('getMyProfile', () => {
  it('calls GET /me/profile/', async () => {
    mockGet.mockResolvedValue(sampleProfile);

    const result = await getMyProfile();

    expect(
      mockGet,
      'getMyProfile should call GET on PROFILE.MY_PROFILE path. Check API_PATHS.PROFILE.MY_PROFILE.',
    ).toHaveBeenCalledWith(API_PATHS.PROFILE.MY_PROFILE);
    expect(
      result.userId,
      'getMyProfile should return the user profile.',
    ).toBe('user-123');
  });
});
