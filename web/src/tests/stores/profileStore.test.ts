import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useProfileStore } from '@/stores/profileStore';
import * as profileApi from '@/api/profile';
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
  it('fetches profile and updates state', async () => {
    mockGetMyProfile.mockResolvedValue(sampleProfile);

    await useProfileStore.getState().fetchProfile();

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

  it('sets error on failure', async () => {
    mockGetMyProfile.mockRejectedValue(new Error('Network error'));

    await useProfileStore.getState().fetchProfile();

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
});

describe('createProfile', () => {
  it('creates a new profile', async () => {
    mockCreateProfile.mockResolvedValue(sampleProfile);

    await useProfileStore.getState().createProfile({
      firstName: 'John',
      lastName: 'Doe',
      username: 'johndoe',
      gender: 'male',
      sexualPreference: 'heterosexual',
      birthday: '1995-06-15',
      biography: 'Hello world',
    });

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

  it('sets error on create failure', async () => {
    mockCreateProfile.mockRejectedValue(new Error('Create failed'));

    await useProfileStore.getState().createProfile({
      firstName: 'John',
      lastName: 'Doe',
      username: 'johndoe',
      gender: 'male',
      sexualPreference: 'heterosexual',
      birthday: '1995-06-15',
      biography: 'Hello world',
    });

    const state = useProfileStore.getState();
    expect(
      state.error,
      'error should contain the failure message.',
    ).toBe('Create failed');
  });
});

describe('updateProfile', () => {
  it('updates an existing profile', async () => {
    const updatedProfile = { ...sampleProfile, biography: 'Updated bio' };
    mockUpdateProfile.mockResolvedValue(updatedProfile);

    await useProfileStore.getState().updateProfile({ biography: 'Updated bio' });

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

  it('sets error on update failure', async () => {
    mockUpdateProfile.mockRejectedValue(new Error('Save failed'));

    await useProfileStore.getState().updateProfile({ biography: 'New bio' });

    const state = useProfileStore.getState();
    expect(
      state.error,
      'error should contain the failure message.',
    ).toBe('Save failed');
  });
});
