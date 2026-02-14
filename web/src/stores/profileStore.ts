import { create } from 'zustand';
import type { UserProfile, Picture, Tag } from '@/types';
import type { CreateProfileRequest, UpdateProfileRequest } from '@/api/profile';
import * as profileApi from '@/api/profile';

const MAX_PICTURES = 5;

interface ProfileState {
  readonly profile: UserProfile | null;
  readonly pictures: readonly Picture[];
  readonly tags: readonly Tag[];
  readonly allTags: readonly Tag[];
  readonly isLoading: boolean;
  readonly error: string | null;
}

interface ProfileActions {
  readonly fetchProfile: () => Promise<void>;
  readonly saveProfile: (params: CreateProfileRequest | UpdateProfileRequest) => Promise<void>;
  readonly uploadPicture: (file: File) => Promise<void>;
  readonly deletePicture: (pictureId: number) => Promise<void>;
  readonly fetchTags: () => Promise<void>;
  readonly addTag: (tagId: number) => Promise<void>;
  readonly removeTag: (tagId: number) => Promise<void>;
  readonly clearError: () => void;
}

type ProfileStore = ProfileState & ProfileActions;

const initialState: ProfileState = {
  profile: null,
  pictures: [],
  tags: [],
  allTags: [],
  isLoading: false,
  error: null,
};

export const useProfileStore = create<ProfileStore>()((set, get) => ({
  ...initialState,

  fetchProfile: async () => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileApi.getMyProfile();
      set({ profile, isLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch profile';
      set({ error: message, isLoading: false });
    }
  },

  saveProfile: async (params) => {
    set({ isLoading: true, error: null });
    try {
      const { profile: existing } = get();
      const profile = existing
        ? await profileApi.updateProfile(params)
        : await profileApi.createProfile(params as CreateProfileRequest);
      set({ profile, isLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to save profile';
      set({ error: message, isLoading: false });
    }
  },

  uploadPicture: async (file: File) => {
    const { pictures } = get();
    if (pictures.length >= MAX_PICTURES) {
      set({ error: 'Maximum 5 pictures allowed' });
      return;
    }
    set({ isLoading: true, error: null });
    try {
      const picture = await profileApi.uploadPicture(file);
      set({ pictures: [...get().pictures, picture], isLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to upload picture';
      set({ error: message, isLoading: false });
    }
  },

  deletePicture: async (pictureId: number) => {
    set({ isLoading: true, error: null });
    try {
      await profileApi.deletePicture(pictureId);
      set({
        pictures: get().pictures.filter((p) => p.id !== pictureId),
        isLoading: false,
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to delete picture';
      set({ error: message, isLoading: false });
    }
  },

  fetchTags: async () => {
    set({ isLoading: true, error: null });
    try {
      const allTags = await profileApi.getTags();
      set({ allTags, isLoading: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch tags';
      set({ error: message, isLoading: false });
    }
  },

  addTag: async (tagId: number) => {
    set({ isLoading: true, error: null });
    try {
      await profileApi.addTag(tagId);
      const { allTags, tags } = get();
      const tag = allTags.find((t) => t.id === tagId);
      if (tag) {
        set({ tags: [...tags, tag], isLoading: false });
      } else {
        set({ isLoading: false });
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to add tag';
      set({ error: message, isLoading: false });
    }
  },

  removeTag: async (tagId: number) => {
    set({ isLoading: true, error: null });
    try {
      await profileApi.removeTag(tagId);
      set({
        tags: get().tags.filter((t) => t.id !== tagId),
        isLoading: false,
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to remove tag';
      set({ error: message, isLoading: false });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
