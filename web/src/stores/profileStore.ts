import { create } from 'zustand';
import type { UserProfile } from '@/types';
import type { CreateProfileRequest, UpdateProfileRequest } from '@/api/profile';
import * as profileApi from '@/api/profile';

interface ProfileState {
  readonly profile: UserProfile | null;
  readonly isLoading: boolean;
  readonly error: string | null;
}

interface ProfileActions {
  readonly fetchProfile: () => Promise<void>;
  readonly saveProfile: (params: CreateProfileRequest | UpdateProfileRequest) => Promise<void>;
  readonly clearError: () => void;
}

type ProfileStore = ProfileState & ProfileActions;

const initialState: ProfileState = {
  profile: null,
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

  clearError: () => {
    set({ error: null });
  },
}));
