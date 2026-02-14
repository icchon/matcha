import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useProfileStore } from '@/stores/profileStore';
import * as profileApi from '@/api/profile';
import type { UserProfile, Picture, Tag, UserTag } from '@/types';

vi.mock('@/api/profile');

const mockGetMyProfile = vi.mocked(profileApi.getMyProfile);
const mockCreateProfile = vi.mocked(profileApi.createProfile);
const mockUpdateProfile = vi.mocked(profileApi.updateProfile);
const mockUploadPicture = vi.mocked(profileApi.uploadPicture);
const mockDeletePicture = vi.mocked(profileApi.deletePicture);
const mockGetTags = vi.mocked(profileApi.getTags);
const mockAddTag = vi.mocked(profileApi.addTag);
const mockRemoveTag = vi.mocked(profileApi.removeTag);

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

const samplePicture: Picture = {
  id: 1,
  userId: 'user-123',
  url: '/images/pic1.jpg',
  isProfilePic: true,
  createdAt: '2025-01-01T00:00:00Z',
};

const sampleTags: Tag[] = [
  { id: 1, name: 'hiking' },
  { id: 2, name: 'cooking' },
];

beforeEach(() => {
  vi.resetAllMocks();
  useProfileStore.setState({
    profile: null,
    pictures: [],
    tags: [],
    allTags: [],
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
      state.pictures,
      'Initial pictures should be empty array.',
    ).toEqual([]);
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

describe('saveProfile (create)', () => {
  it('creates a new profile when no profile exists', async () => {
    mockCreateProfile.mockResolvedValue(sampleProfile);

    await useProfileStore.getState().saveProfile({
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
      'saveProfile should call createProfile when no existing profile.',
    ).toHaveBeenCalled();
    expect(
      state.profile?.firstName,
      'Profile should be updated after successful create.',
    ).toBe('John');
  });
});

describe('saveProfile (update)', () => {
  it('updates an existing profile', async () => {
    useProfileStore.setState({ profile: sampleProfile });
    const updatedProfile = { ...sampleProfile, biography: 'Updated bio' };
    mockUpdateProfile.mockResolvedValue(updatedProfile);

    await useProfileStore.getState().saveProfile({ biography: 'Updated bio' });

    const state = useProfileStore.getState();
    expect(
      mockUpdateProfile,
      'saveProfile should call updateProfile when profile already exists.',
    ).toHaveBeenCalled();
    expect(
      state.profile?.biography,
      'Profile biography should be updated after successful save.',
    ).toBe('Updated bio');
  });

  it('sets error on save failure', async () => {
    mockUpdateProfile.mockRejectedValue(new Error('Save failed'));
    useProfileStore.setState({ profile: sampleProfile });

    await useProfileStore.getState().saveProfile({ biography: 'New bio' });

    const state = useProfileStore.getState();
    expect(
      state.error,
      'error should contain the failure message.',
    ).toBe('Save failed');
  });
});

describe('uploadPicture', () => {
  it('adds uploaded picture to pictures array', async () => {
    mockUploadPicture.mockResolvedValue(samplePicture);
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' });

    await useProfileStore.getState().uploadPicture(file);

    const state = useProfileStore.getState();
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
    useProfileStore.setState({ pictures: fivePictures });
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' });

    await useProfileStore.getState().uploadPicture(file);

    const state = useProfileStore.getState();
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

    await useProfileStore.getState().uploadPicture(file);

    expect(
      useProfileStore.getState().error,
      'error should be set when upload fails.',
    ).toBe('Upload failed');
  });
});

describe('deletePicture', () => {
  it('removes picture from pictures array', async () => {
    useProfileStore.setState({ pictures: [samplePicture] });
    mockDeletePicture.mockResolvedValue(undefined);

    await useProfileStore.getState().deletePicture(1);

    const state = useProfileStore.getState();
    expect(
      state.pictures,
      'pictures array should be empty after deleting the only picture.',
    ).toHaveLength(0);
  });

  it('sets error on delete failure', async () => {
    useProfileStore.setState({ pictures: [samplePicture] });
    mockDeletePicture.mockRejectedValue(new Error('Delete failed'));

    await useProfileStore.getState().deletePicture(1);

    expect(
      useProfileStore.getState().error,
      'error should be set when delete fails.',
    ).toBe('Delete failed');
  });
});

describe('fetchTags', () => {
  it('fetches all available tags', async () => {
    mockGetTags.mockResolvedValue(sampleTags);

    await useProfileStore.getState().fetchTags();

    const state = useProfileStore.getState();
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
    useProfileStore.setState({ allTags: sampleTags });
    const userTag: UserTag = { userId: 'user-123', tagId: 1 };
    mockAddTag.mockResolvedValue(userTag);

    await useProfileStore.getState().addTag(1);

    const state = useProfileStore.getState();
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

    await useProfileStore.getState().addTag(1);

    expect(
      useProfileStore.getState().error,
      'error should be set when addTag fails.',
    ).toBe('Add failed');
  });
});

describe('removeTag', () => {
  it('removes a tag from user tags', async () => {
    useProfileStore.setState({ tags: [sampleTags[0]] });
    mockRemoveTag.mockResolvedValue(undefined);

    await useProfileStore.getState().removeTag(1);

    const state = useProfileStore.getState();
    expect(
      state.tags,
      'tags array should be empty after removing the only tag.',
    ).toHaveLength(0);
  });

  it('sets error on removeTag failure', async () => {
    useProfileStore.setState({ tags: [sampleTags[0]] });
    mockRemoveTag.mockRejectedValue(new Error('Remove failed'));

    await useProfileStore.getState().removeTag(1);

    expect(
      useProfileStore.getState().error,
      'error should be set when removeTag fails.',
    ).toBe('Remove failed');
  });
});
