import { create } from 'zustand';
import type { UserProfile } from '@/types';
import type { CreateProfileRequest } from '@/api/profile';
import * as profileApi from '@/api/profile';
import { ApiClientError } from '@/api/client';
import { toUserFacingMessage } from '@/lib/errorUtils';

interface ProfileState {
  readonly profile: UserProfile | null;
  readonly isLoading: boolean;
  readonly error: string | null;
}

interface ProfileActions {
  /** Returns true if a profile was found, false if none exists or an error occurred.
   *  Check `error` state to distinguish "no profile" from "fetch failure". */
  readonly fetchProfile: () => Promise<boolean>;
  readonly createProfile: (params: CreateProfileRequest) => Promise<boolean>;
  readonly updateProfile: (params: Parameters<typeof profileApi.updateProfile>[0]) => Promise<boolean>;
  readonly clearError: () => void;
}

type ProfileStore = ProfileState & ProfileActions;

const initialState: ProfileState = {
  profile: null,
  isLoading: false,
  error: null,
};

export const useProfileStore = create<ProfileStore>()((set) => ({
  ...initialState,

  fetchProfile: async () => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileApi.getMyProfile();
      set({ profile, isLoading: false });
      return true;
    } catch (err) {
      // 404: profile not created yet
      if (err instanceof ApiClientError && err.status === 404) {
        set({ profile: null, error: null, isLoading: false });
        return false;
      }
      // TODO(FE-04): Remove 405 fallback once BE implements GET /me/profile/
      if (err instanceof ApiClientError && err.status === 405) {
        if (import.meta.env.DEV) {
          console.warn('[profileStore] GET /me/profile/ returned 405 — endpoint not yet implemented');
        }
        set({ profile: null, error: null, isLoading: false });
        return false;
      }
      const message = toUserFacingMessage(err, 'Failed to fetch profile');
      set({ error: message, isLoading: false });
      return false;
    }
  },

  createProfile: async (params) => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileApi.createProfile(params);
      set({ profile, isLoading: false });
      return true;
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to create profile');
      set({ error: message, isLoading: false });
      return false;
    }
  },

  updateProfile: async (params) => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileApi.updateProfile(params);
      set({ profile, isLoading: false });
      return true;
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to update profile');
      set({ error: message, isLoading: false });
      return false;
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
