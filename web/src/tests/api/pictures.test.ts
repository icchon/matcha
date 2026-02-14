import { describe, it, expect, vi, beforeEach } from 'vitest';
import { uploadPicture, deletePicture } from '@/api/pictures';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Picture } from '@/types';

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    upload: vi.fn(),
  },
}));

const mockDelete = vi.mocked(apiClient.delete);
const mockUpload = vi.mocked(apiClient.upload);

const samplePicture: Picture = {
  id: 1,
  userId: 'user-123',
  url: '/images/pic1.jpg',
  isProfilePic: true,
  createdAt: '2025-01-01T00:00:00Z',
};

beforeEach(() => {
  vi.resetAllMocks();
});

describe('uploadPicture', () => {
  it('calls upload on PROFILE.PICTURES with FormData', async () => {
    const file = new File(['image-data'], 'photo.jpg', { type: 'image/jpeg' });
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
