import { describe, it, expect, vi, beforeEach } from 'vitest';
import { usePictureStore } from '@/stores/pictureStore';
import * as picturesApi from '@/api/pictures';
import type { Picture } from '@/types';

vi.mock('@/api/pictures');

const mockUploadPicture = vi.mocked(picturesApi.uploadPicture);
const mockDeletePicture = vi.mocked(picturesApi.deletePicture);

const samplePicture: Picture = {
  id: 1,
  userId: 'user-123',
  url: '/images/pic1.jpg',
  isProfilePic: true,
  createdAt: '2025-01-01T00:00:00Z',
};

beforeEach(() => {
  vi.resetAllMocks();
  usePictureStore.setState({
    pictures: [],
    isLoading: false,
    error: null,
  });
});

describe('pictureStore initial state', () => {
  it('has correct initial state', () => {
    const state = usePictureStore.getState();

    expect(
      state.pictures,
      'Initial pictures should be empty array.',
    ).toEqual([]);
    expect(
      state.isLoading,
      'Initial isLoading should be false.',
    ).toBe(false);
    expect(
      state.error,
      'Initial error should be null.',
    ).toBeNull();
  });
});

describe('uploadPicture', () => {
  it('adds uploaded picture to pictures array', async () => {
    mockUploadPicture.mockResolvedValue(samplePicture);
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' });

    await usePictureStore.getState().uploadPicture(file);

    const state = usePictureStore.getState();
    expect(
      state.pictures,
      'pictures array should contain the newly uploaded picture.',
    ).toHaveLength(1);
    expect(
      state.pictures[0].id,
      'The uploaded picture should match the API response.',
    ).toBe(1);
  });

  it('rejects upload when already at max 5 pictures', async () => {
    const fivePictures: Picture[] = Array.from({ length: 5 }, (_, i) => ({
      id: i + 1,
      userId: 'user-123',
      url: `/images/pic${i + 1}.jpg`,
      isProfilePic: i === 0,
      createdAt: '2025-01-01T00:00:00Z',
    }));
    usePictureStore.setState({ pictures: fivePictures });
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' });

    await usePictureStore.getState().uploadPicture(file);

    const state = usePictureStore.getState();
    expect(
      mockUploadPicture,
      'uploadPicture should NOT call API when at max pictures (5).',
    ).not.toHaveBeenCalled();
    expect(
      state.error,
      'error should indicate max pictures reached.',
    ).toBe('Maximum 5 pictures allowed');
    expect(
      state.pictures,
      'pictures array should remain unchanged.',
    ).toHaveLength(5);
  });

  it('sets error on upload failure', async () => {
    mockUploadPicture.mockRejectedValue(new Error('Upload failed'));
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' });

    await usePictureStore.getState().uploadPicture(file);

    expect(
      usePictureStore.getState().error,
      'error should be set when upload fails.',
    ).toBe('Upload failed');
  });
});

describe('deletePicture', () => {
  it('removes picture from pictures array', async () => {
    usePictureStore.setState({ pictures: [samplePicture] });
    mockDeletePicture.mockResolvedValue(undefined);

    await usePictureStore.getState().deletePicture(1);

    const state = usePictureStore.getState();
    expect(
      state.pictures,
      'pictures array should be empty after deleting the only picture.',
    ).toHaveLength(0);
  });

  it('sets error on delete failure', async () => {
    usePictureStore.setState({ pictures: [samplePicture] });
    mockDeletePicture.mockRejectedValue(new Error('Delete failed'));

    await usePictureStore.getState().deletePicture(1);

    expect(
      usePictureStore.getState().error,
      'error should be set when delete fails.',
    ).toBe('Delete failed');
  });
});
