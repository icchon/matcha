import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  createProfile,
  updateProfile,
  getMyProfile,
  uploadPicture,
  deletePicture,
  getTags,
  addTag,
  removeTag,
} from '@/api/profile';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { UserProfile, Picture, Tag, UserTag } from '@/types';

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
const mockDelete = vi.mocked(apiClient.delete);
const mockUpload = vi.mocked(apiClient.upload);

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

const samplePicture: Picture = {
  id: 1,
  userId: 'user-123',
  url: '/images/pic1.jpg',
  isProfilePic: true,
  createdAt: '2025-01-01T00:00:00Z',
};

const sampleTag: Tag = {
  id: 1,
  name: 'hiking',
};

const sampleUserTag: UserTag = {
  userId: 'user-123',
  tagId: 1,
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
      'getMyProfile should call GET on PROFILE.CREATE path (same endpoint for GET own profile).',
    ).toHaveBeenCalledWith(API_PATHS.PROFILE.CREATE);
    expect(
      result.userId,
      'getMyProfile should return the user profile.',
    ).toBe('user-123');
  });
});

describe('uploadPicture', () => {
  it('calls upload on PROFILE.PICTURES with FormData', async () => {
    const file = new File(['image-data'], 'photo.jpg', { type: 'image/jpeg' });
    const formData = new FormData();
    formData.append('file', file);
    mockUpload.mockResolvedValue(samplePicture);

    const result = await uploadPicture(file);

    expect(
      mockUpload,
      'uploadPicture should call apiClient.upload with FormData containing the file.',
    ).toHaveBeenCalledWith(API_PATHS.PROFILE.PICTURES, expect.any(FormData));
    expect(
      result.id,
      'uploadPicture should return the created picture object.',
    ).toBe(1);
  });
});

describe('deletePicture', () => {
  it('calls DELETE /me/profile/pictures/:id', async () => {
    mockDelete.mockResolvedValue(undefined);

    await deletePicture(5);

    expect(
      mockDelete,
      'deletePicture should call DELETE with the picture ID in the path.',
    ).toHaveBeenCalledWith(API_PATHS.PROFILE.DELETE_PICTURE(5));
  });
});

describe('getTags', () => {
  it('calls GET /tags/ and returns tag list', async () => {
    const tags: Tag[] = [sampleTag, { id: 2, name: 'cooking' }];
    mockGet.mockResolvedValue(tags);

    const result = await getTags();

    expect(
      mockGet,
      'getTags should call GET on TAGS path.',
    ).toHaveBeenCalledWith(API_PATHS.TAGS);
    expect(
      result,
      'getTags should return an array of tags. Check that the response is passed through.',
    ).toHaveLength(2);
  });
});

describe('addTag', () => {
  it('calls POST /me/tags/ with tagId', async () => {
    mockPost.mockResolvedValue(sampleUserTag);

    const result = await addTag(1);

    expect(
      mockPost,
      'addTag should call POST on USERS.MY_TAGS with the tag ID.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.MY_TAGS, { tagId: 1 });
    expect(
      result.tagId,
      'addTag should return the created UserTag.',
    ).toBe(1);
  });
});

describe('removeTag', () => {
  it('calls DELETE /me/tags/:tagId', async () => {
    mockDelete.mockResolvedValue(undefined);

    await removeTag(3);

    expect(
      mockDelete,
      'removeTag should call DELETE with the tag ID in the path.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.DELETE_TAG(3));
  });
});
