import { create } from 'zustand';
import type { Tag } from '@/types';
import * as tagsApi from '@/api/tags';
import { toUserFacingMessage } from '@/lib/errorUtils';

interface TagState {
  readonly tags: readonly Tag[];
  readonly allTags: readonly Tag[];
  readonly isLoading: boolean;
  readonly error: string | null;
}

interface TagActions {
  readonly fetchTags: () => Promise<void>;
  readonly addTag: (tagId: number) => Promise<void>;
  readonly removeTag: (tagId: number) => Promise<void>;
  readonly clearError: () => void;
}

type TagStore = TagState & TagActions;

const initialState: TagState = {
  tags: [],
  allTags: [],
  isLoading: false,
  error: null,
};

export const useTagStore = create<TagStore>()((set, get) => ({
  ...initialState,

  fetchTags: async () => {
    set({ isLoading: true, error: null });
    try {
      const allTags = await tagsApi.getTags();
      set({ allTags, isLoading: false });
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to fetch tags');
      set({ error: message, isLoading: false });
    }
  },

  addTag: async (tagId: number) => {
    set({ isLoading: true, error: null });
    try {
      await tagsApi.addTag(tagId);
      const { allTags, tags } = get();
      const tag = allTags.find((t) => t.id === tagId);
      if (tag) {
        set({ tags: [...tags, tag], isLoading: false });
      } else {
        set({ isLoading: false });
      }
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to add tag');
      set({ error: message, isLoading: false });
    }
  },

  removeTag: async (tagId: number) => {
    set({ isLoading: true, error: null });
    try {
      await tagsApi.removeTag(tagId);
      set({
        tags: get().tags.filter((t) => t.id !== tagId),
        isLoading: false,
      });
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to remove tag');
      set({ error: message, isLoading: false });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
