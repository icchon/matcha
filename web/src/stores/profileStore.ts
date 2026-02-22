import { create } from 'zustand';
import type { UserProfile } from '@/types';
import type { CreateProfileRequest } from '@/api/profile';
import * as profileApi from '@/api/profile';
import { toUserFacingMessage } from '@/lib/errorUtils';

interface ProfileState {
  readonly profile: UserProfile | null;
  readonly isLoading: boolean;
  readonly error: string | null;
}

interface ProfileActions {
  readonly fetchProfile: () => Promise<void>;
  readonly createProfile: (params: CreateProfileRequest) => Promise<void>;
  readonly updateProfile: (params: Parameters<typeof profileApi.updateProfile>[0]) => Promise<void>;
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
      const message = toUserFacingMessage(err, 'Failed to fetch profile');
      set({ error: message, isLoading: false });
    }
  },

  createProfile: async (params) => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileApi.createProfile(params);
      set({ profile, isLoading: false });
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to create profile');
      set({ error: message, isLoading: false });
    }
  },

  updateProfile: async (params) => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileApi.updateProfile(params);
      set({ profile, isLoading: false });
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to update profile');
      set({ error: message, isLoading: false });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
