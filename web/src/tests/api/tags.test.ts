import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getTags, addTag, removeTag } from '@/api/tags';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Tag, UserTag } from '@/types';

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
const mockDelete = vi.mocked(apiClient.delete);

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
