import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useProfileStore } from '@/stores/profileStore';
import * as profileApi from '@/api/profile';
import { ApiClientError } from '@/api/client';
import type { UserProfile } from '@/types';

vi.mock('@/api/profile');

const mockGetMyProfile = vi.mocked(profileApi.getMyProfile);
const mockCreateProfile = vi.mocked(profileApi.createProfile);
const mockUpdateProfile = vi.mocked(profileApi.updateProfile);

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

beforeEach(() => {
  vi.resetAllMocks();
  useProfileStore.setState({
    profile: null,
    isLoading: false,
    error: null,
  });
});

describe('profileStore initial state', () => {
  it('has correct initial state', () => {
    const state = useProfileStore.getState();

    expect(
      state.profile,
      'Initial profile should be null. Check initialState in profileStore.',
    ).toBeNull();
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

describe('fetchProfile', () => {
  it('fetches profile, updates state, and returns true', async () => {
    mockGetMyProfile.mockResolvedValue(sampleProfile);

    const result = await useProfileStore.getState().fetchProfile();

    expect(
      result,
      'fetchProfile should return true when a profile is found.',
    ).toBe(true);
    const state = useProfileStore.getState();
    expect(
      state.profile?.userId,
      'fetchProfile should set profile data from API response.',
    ).toBe('user-123');
    expect(
      state.isLoading,
      'isLoading should be false after fetch completes.',
    ).toBe(false);
    expect(
      state.error,
      'error should be null after successful fetch.',
    ).toBeNull();
  });

  it('sets error on failure and returns false', async () => {
    mockGetMyProfile.mockRejectedValue(new Error('Network error'));

    const result = await useProfileStore.getState().fetchProfile();

    expect(
      result,
      'fetchProfile should return false on non-404/405 errors.',
    ).toBe(false);
    const state = useProfileStore.getState();
    expect(
      state.profile,
      'profile should remain null when fetchProfile fails.',
    ).toBeNull();
    expect(
      state.error,
      'error should be set to the error message on failure.',
    ).toBe('Network error');
    expect(
      state.isLoading,
      'isLoading should be false after fetch fails.',
    ).toBe(false);
  });

  it('maps 5xx errors to generic fallback message', async () => {
    const serverError = Object.assign(new Error('Internal Server Error'), { status: 500 });
    mockGetMyProfile.mockRejectedValue(serverError);

    await useProfileStore.getState().fetchProfile();

    const state = useProfileStore.getState();
    expect(
      state.error,
      '5xx errors should show generic fallback, not raw server message. Check toUserFacingMessage.',
    ).toBe('Failed to fetch profile');
  });

  it('treats 404 as "no profile yet" (profile=null, error=null)', async () => {
    mockGetMyProfile.mockRejectedValue(
      new ApiClientError(404, { error: 'Not Found' }),
    );

    await useProfileStore.getState().fetchProfile();

    const state = useProfileStore.getState();
    expect(
      state.profile,
      '404 should result in profile=null (no profile exists yet).',
    ).toBeNull();
    expect(
      state.error,
      '404 should NOT set an error — it simply means the profile has not been created.',
    ).toBeNull();
    expect(state.isLoading, 'isLoading should be false after 404.').toBe(false);
  });

  it('treats 405 as "no profile yet" (BE endpoint not implemented)', async () => {
    mockGetMyProfile.mockRejectedValue(
      new ApiClientError(405, { error: 'Method Not Allowed' }),
    );

    await useProfileStore.getState().fetchProfile();

    const state = useProfileStore.getState();
    expect(
      state.profile,
      '405 should result in profile=null (BE GET not implemented yet).',
    ).toBeNull();
    expect(
      state.error,
      '405 should NOT set an error — it is a temporary BE limitation.',
    ).toBeNull();
    expect(state.isLoading, 'isLoading should be false after 405.').toBe(false);
  });
});

describe('createProfile', () => {
  it('creates a new profile and returns true', async () => {
    mockCreateProfile.mockResolvedValue(sampleProfile);

    const result = await useProfileStore.getState().createProfile({
      firstName: 'John',
      lastName: 'Doe',
      username: 'johndoe',
      gender: 'male',
      sexualPreference: 'heterosexual',
      birthday: '1995-06-15',
      biography: 'Hello world',
    });

    expect(
      result,
      'createProfile should return true on success so the caller can show a toast.',
    ).toBe(true);
    const state = useProfileStore.getState();
    expect(
      mockCreateProfile,
      'createProfile should call profileApi.createProfile.',
    ).toHaveBeenCalled();
    expect(
      state.profile?.firstName,
      'Profile should be updated after successful create.',
    ).toBe('John');
  });

  it('sets error on create failure and returns false', async () => {
    mockCreateProfile.mockRejectedValue(new Error('Create failed'));

    const result = await useProfileStore.getState().createProfile({
      firstName: 'John',
      lastName: 'Doe',
      username: 'johndoe',
      gender: 'male',
      sexualPreference: 'heterosexual',
      birthday: '1995-06-15',
      biography: 'Hello world',
    });

    expect(
      result,
      'createProfile should return false on failure.',
    ).toBe(false);
    const state = useProfileStore.getState();
    expect(
      state.error,
      'error should contain the failure message.',
    ).toBe('Create failed');
  });
});

describe('updateProfile', () => {
  it('updates an existing profile and returns true', async () => {
    const updatedProfile = { ...sampleProfile, biography: 'Updated bio' };
    mockUpdateProfile.mockResolvedValue(updatedProfile);

    const result = await useProfileStore.getState().updateProfile({ biography: 'Updated bio' });

    expect(
      result,
      'updateProfile should return true on success so the caller can show a toast.',
    ).toBe(true);
    const state = useProfileStore.getState();
    expect(
      mockUpdateProfile,
      'updateProfile should call profileApi.updateProfile.',
    ).toHaveBeenCalled();
    expect(
      state.profile?.biography,
      'Profile biography should be updated after successful save.',
    ).toBe('Updated bio');
  });

  it('sets error on update failure and returns false', async () => {
    mockUpdateProfile.mockRejectedValue(new Error('Save failed'));

    const result = await useProfileStore.getState().updateProfile({ biography: 'New bio' });

    expect(
      result,
      'updateProfile should return false on failure.',
    ).toBe(false);
    const state = useProfileStore.getState();
    expect(
      state.error,
      'error should contain the failure message.',
    ).toBe('Save failed');
  });
});
