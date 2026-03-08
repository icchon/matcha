import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useTagStore } from '@/stores/tagStore';
import * as tagsApi from '@/api/tags';
import type { Tag, UserTag } from '@/types';

vi.mock('@/api/tags');

const mockGetTags = vi.mocked(tagsApi.getTags);
const mockAddTag = vi.mocked(tagsApi.addTag);
const mockRemoveTag = vi.mocked(tagsApi.removeTag);

const sampleTags: Tag[] = [
  { id: 1, name: 'hiking' },
  { id: 2, name: 'cooking' },
];

beforeEach(() => {
  vi.resetAllMocks();
  useTagStore.setState({
    tags: [],
    allTags: [],
    isLoading: false,
    error: null,
  });
});

describe('tagStore initial state', () => {
  it('has correct initial state', () => {
    const state = useTagStore.getState();

    expect(
      state.tags,
      'Initial tags (user tags) should be empty array.',
    ).toEqual([]);
    expect(
      state.allTags,
      'Initial allTags should be empty array.',
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

describe('fetchTags', () => {
  it('fetches all available tags', async () => {
    mockGetTags.mockResolvedValue(sampleTags);

    await useTagStore.getState().fetchTags();

    const state = useTagStore.getState();
    expect(
      state.allTags,
      'allTags should be populated with fetched tags.',
    ).toHaveLength(2);
    expect(
      state.allTags[0].name,
      'First tag should be hiking.',
    ).toBe('hiking');
  });
});

describe('addTag', () => {
  it('adds a tag to user tags', async () => {
    useTagStore.setState({ allTags: sampleTags });
    const userTag: UserTag = { userId: 'user-123', tagId: 1 };
    mockAddTag.mockResolvedValue(userTag);

    await useTagStore.getState().addTag(1);

    const state = useTagStore.getState();
    expect(
      state.tags,
      'tags array should contain the newly added tag.',
    ).toHaveLength(1);
    expect(
      state.tags[0].id,
      'The added tag should match the tag from allTags with the given tagId.',
    ).toBe(1);
  });

  it('sets error on addTag failure', async () => {
    mockAddTag.mockRejectedValue(new Error('Add failed'));

    await useTagStore.getState().addTag(1);

    expect(
      useTagStore.getState().error,
      'error should be set when addTag fails.',
    ).toBe('Add failed');
  });

  it('maps 5xx addTag errors to generic fallback message', async () => {
    const serverError = Object.assign(new Error('Internal Server Error'), { status: 503 });
    mockAddTag.mockRejectedValue(serverError);

    await useTagStore.getState().addTag(1);

    expect(
      useTagStore.getState().error,
      '5xx errors should show generic fallback, not raw server message. Check toUserFacingMessage.',
    ).toBe('Failed to add tag');
  });
});

describe('removeTag', () => {
  it('removes a tag from user tags', async () => {
    useTagStore.setState({ tags: [sampleTags[0]] });
    mockRemoveTag.mockResolvedValue(undefined);

    await useTagStore.getState().removeTag(1);

    const state = useTagStore.getState();
    expect(
      state.tags,
      'tags array should be empty after removing the only tag.',
    ).toHaveLength(0);
  });

  it('sets error on removeTag failure', async () => {
    useTagStore.setState({ tags: [sampleTags[0]] });
    mockRemoveTag.mockRejectedValue(new Error('Remove failed'));

    await useTagStore.getState().removeTag(1);

    expect(
      useTagStore.getState().error,
      'error should be set when removeTag fails.',
    ).toBe('Remove failed');
  });
});
