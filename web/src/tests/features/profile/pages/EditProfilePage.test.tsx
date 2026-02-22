import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { EditProfilePage } from '@/features/profile/pages/EditProfilePage';
import { useProfileStore } from '@/stores/profileStore';
import { usePictureStore } from '@/stores/pictureStore';
import { useTagStore } from '@/stores/tagStore';
import type { UserProfile } from '@/types';

vi.mock('@/stores/profileStore');
vi.mock('@/stores/pictureStore');
vi.mock('@/stores/tagStore');

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

const mockUpdateProfile = vi.fn();
const mockFetchProfile = vi.fn();
const mockFetchTags = vi.fn();
const mockUploadPicture = vi.fn();
const mockDeletePicture = vi.fn();
const mockAddTag = vi.fn();
const mockRemoveTag = vi.fn();

function setupMockStores(overrides: Record<string, unknown> = {}) {
  const profileState = {
    profile: sampleProfile,
    isLoading: false,
    error: null,
    updateProfile: mockUpdateProfile,
    fetchProfile: mockFetchProfile,
    clearError: vi.fn(),
    ...('profile' in overrides || 'isLoading' in overrides || 'error' in overrides
      ? overrides
      : {}),
  };

  const pictureState = {
    pictures: [],
    isLoading: false,
    error: null,
    uploadPicture: mockUploadPicture,
    deletePicture: mockDeletePicture,
    clearError: vi.fn(),
  };

  const tagState = {
    tags: [{ id: 1, name: 'hiking' }],
    allTags: [{ id: 1, name: 'hiking' }, { id: 2, name: 'cooking' }],
    isLoading: false,
    error: null,
    fetchTags: mockFetchTags,
    addTag: mockAddTag,
    removeTag: mockRemoveTag,
    clearError: vi.fn(),
  };

  vi.mocked(useProfileStore).mockImplementation((selector: unknown) => {
    if (typeof selector === 'function') {
      return (selector as (s: typeof profileState) => unknown)(profileState);
    }
    return profileState;
  });

  vi.mocked(usePictureStore).mockImplementation((selector: unknown) => {
    if (typeof selector === 'function') {
      return (selector as (s: typeof pictureState) => unknown)(pictureState);
    }
    return pictureState;
  });

  vi.mocked(useTagStore).mockImplementation((selector: unknown) => {
    if (typeof selector === 'function') {
      return (selector as (s: typeof tagState) => unknown)(tagState);
    }
    return tagState;
  });
}

beforeEach(() => {
  vi.resetAllMocks();
  setupMockStores();
});

function renderPage() {
  return render(
    <MemoryRouter>
      <EditProfilePage />
    </MemoryRouter>,
  );
}

describe('EditProfilePage', () => {
  it('renders page title for editing', () => {
    renderPage();

    expect(
      screen.getByText(/edit.*profile/i),
      'EditProfilePage should display an "Edit Profile" heading.',
    ).toBeInTheDocument();
  });

  it('pre-fills form with existing profile data', () => {
    renderPage();

    expect(
      (screen.getByLabelText(/first name/i) as HTMLInputElement).value,
      'First name should be pre-filled from existing profile.',
    ).toBe('John');
    expect(
      (screen.getByLabelText(/biography/i) as HTMLTextAreaElement).value,
      'Biography should be pre-filled from existing profile.',
    ).toBe('Hello world');
  });

  it('renders existing user tags', () => {
    renderPage();

    expect(
      screen.getByText('hiking'),
      'EditProfilePage should show existing user tags.',
    ).toBeInTheDocument();
  });

  it('shows loading spinner when isLoading', () => {
    setupMockStores({ isLoading: true });
    renderPage();

    expect(
      screen.getByRole('status'),
      'Should show loading spinner when profile is loading.',
    ).toBeInTheDocument();
  });

  it('displays error message when error is set', () => {
    setupMockStores({ error: 'Failed to load profile' });
    renderPage();

    expect(
      screen.getByText('Failed to load profile'),
      'Error message should be displayed when store has an error.',
    ).toBeInTheDocument();
  });
});
