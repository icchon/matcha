import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { EditProfilePage } from '@/features/profile/pages/EditProfilePage';
import { useProfileStore } from '@/stores/profileStore';
import type { UserProfile } from '@/types';

vi.mock('@/stores/profileStore');

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

const mockSaveProfile = vi.fn();
const mockFetchProfile = vi.fn();
const mockFetchTags = vi.fn();
const mockUploadPicture = vi.fn();
const mockDeletePicture = vi.fn();
const mockAddTag = vi.fn();
const mockRemoveTag = vi.fn();

function setupMockStore(overrides: Record<string, unknown> = {}) {
  vi.mocked(useProfileStore).mockImplementation((selector: unknown) => {
    const state = {
      profile: sampleProfile,
      pictures: [],
      tags: [{ id: 1, name: 'hiking' }],
      allTags: [{ id: 1, name: 'hiking' }, { id: 2, name: 'cooking' }],
      isLoading: false,
      error: null,
      saveProfile: mockSaveProfile,
      fetchProfile: mockFetchProfile,
      fetchTags: mockFetchTags,
      uploadPicture: mockUploadPicture,
      deletePicture: mockDeletePicture,
      addTag: mockAddTag,
      removeTag: mockRemoveTag,
      clearError: vi.fn(),
      ...overrides,
    };
    if (typeof selector === 'function') {
      return (selector as (s: typeof state) => unknown)(state);
    }
    return state;
  });
}

beforeEach(() => {
  vi.resetAllMocks();
  setupMockStore();
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
    setupMockStore({ isLoading: true });
    renderPage();

    expect(
      screen.getByRole('status'),
      'Should show loading spinner when profile is loading.',
    ).toBeInTheDocument();
  });

  it('displays error message when error is set', () => {
    setupMockStore({ error: 'Failed to load profile' });
    renderPage();

    expect(
      screen.getByText('Failed to load profile'),
      'Error message should be displayed when store has an error.',
    ).toBeInTheDocument();
  });
});
