import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { ProfilePage } from '@/features/profile/pages/ProfilePage';
import { useProfileStore } from '@/stores/profileStore';
import { usePictureStore } from '@/stores/pictureStore';
import { useTagStore } from '@/stores/tagStore';
import type { UserProfile } from '@/types';

vi.mock('@/stores/profileStore');
vi.mock('@/stores/pictureStore');
vi.mock('@/stores/tagStore');
vi.mock('sonner', () => ({ toast: { success: vi.fn() } }));

import { toast } from 'sonner';

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

const mockCreateProfile = vi.fn().mockResolvedValue(true);
const mockUpdateProfile = vi.fn().mockResolvedValue(true);
const mockFetchProfile = vi.fn().mockResolvedValue(true);
const mockFetchTags = vi.fn();
const mockUploadPicture = vi.fn();
const mockDeletePicture = vi.fn();
const mockAddTag = vi.fn();
const mockRemoveTag = vi.fn();

interface StoreOverrides {
  readonly profile?: UserProfile | null;
  readonly isLoading?: boolean;
  readonly error?: string | null;
}

function setupMockStores(overrides: StoreOverrides = {}) {
  const profileState = {
    profile: null as UserProfile | null,
    isLoading: false,
    error: null as string | null,
    createProfile: mockCreateProfile,
    updateProfile: mockUpdateProfile,
    fetchProfile: mockFetchProfile,
    clearError: vi.fn(),
    ...overrides,
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
    tags: [],
    allTags: [{ id: 1, name: 'hiking' }, { id: 2, name: 'cooking' }],
    isLoading: false,
    error: null,
    fetchTags: mockFetchTags,
    addTag: mockAddTag,
    removeTag: mockRemoveTag,
    clearError: vi.fn(),
  };

  const mockedStore = vi.mocked(useProfileStore);
  mockedStore.mockImplementation((selector: unknown) => {
    if (typeof selector === 'function') {
      return (selector as (s: typeof profileState) => unknown)(profileState);
    }
    return profileState;
  });
  mockedStore.getState = vi.fn().mockReturnValue(profileState) as typeof mockedStore.getState;

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
  mockCreateProfile.mockResolvedValue(true);
  mockUpdateProfile.mockResolvedValue(true);
  mockFetchProfile.mockResolvedValue(true);
  setupMockStores();
});

function renderPage() {
  return render(
    <MemoryRouter>
      <ProfilePage />
    </MemoryRouter>,
  );
}

describe('ProfilePage', () => {
  it('shows loading spinner during initial fetch', () => {
    setupMockStores({ isLoading: true });
    renderPage();

    expect(
      screen.getByRole('status'),
      'Should show a loading spinner when isLoading=true and no profile/error yet.',
    ).toBeInTheDocument();
  });

  it('shows "Create Profile" heading when profile is null and no error (new user)', () => {
    setupMockStores({ profile: null, error: null });
    renderPage();

    expect(
      screen.getByText(/create.*profile/i),
      'When profile=null and error=null, should show "Create Profile" heading for new users.',
    ).toBeInTheDocument();
  });

  it('shows "Edit Profile" heading and pre-fills form when profile exists', () => {
    setupMockStores({ profile: sampleProfile });
    renderPage();

    expect(
      screen.getByText(/edit.*profile/i),
      'When profile exists, should show "Edit Profile" heading.',
    ).toBeInTheDocument();
    expect(
      (screen.getByLabelText(/first name/i) as HTMLInputElement).value,
      'First name should be pre-filled from existing profile data.',
    ).toBe('John');
    expect(
      (screen.getByLabelText(/biography/i) as HTMLTextAreaElement).value,
      'Biography should be pre-filled from existing profile data.',
    ).toBe('Hello world');
  });

  it('shows error alert and hides form when there is a server error', () => {
    setupMockStores({ error: 'Something went wrong' });
    renderPage();

    expect(
      screen.getByRole('alert'),
      'Should show an error alert when error is set.',
    ).toBeInTheDocument();
    expect(
      screen.getByText('Something went wrong'),
      'Error message text should be visible.',
    ).toBeInTheDocument();
    expect(
      screen.queryByLabelText(/first name/i),
      'Form should be hidden when there is a server error.',
    ).not.toBeInTheDocument();
  });

  it('calls createProfile when submitting in create mode', async () => {
    setupMockStores({ profile: null, error: null });
    const user = userEvent.setup();
    renderPage();

    await user.type(screen.getByLabelText(/first name/i), 'John');
    await user.type(screen.getByLabelText(/last name/i), 'Doe');
    await user.type(screen.getByLabelText(/username/i), 'johndoe');
    await user.selectOptions(screen.getByLabelText(/gender/i), 'male');
    await user.selectOptions(screen.getByLabelText(/sexual preference/i), 'heterosexual');
    await user.type(screen.getByLabelText(/birthday/i), '1995-06-15');
    await user.type(screen.getByLabelText(/biography/i), 'Hello world');
    await user.click(screen.getByRole('button', { name: /save/i }));

    await waitFor(() => {
      expect(
        mockCreateProfile,
        'Submitting in create mode should call createProfile.',
      ).toHaveBeenCalled();
    });
  });

  it('calls updateProfile when submitting in edit mode', async () => {
    setupMockStores({ profile: sampleProfile });
    const user = userEvent.setup();
    renderPage();

    await user.clear(screen.getByLabelText(/biography/i));
    await user.type(screen.getByLabelText(/biography/i), 'Updated bio');
    await user.click(screen.getByRole('button', { name: /save/i }));

    await waitFor(() => {
      expect(
        mockUpdateProfile,
        'Submitting in edit mode should call updateProfile.',
      ).toHaveBeenCalled();
    });
  });

  it('calls fetchProfile and fetchTags on mount', () => {
    renderPage();

    expect(
      mockFetchProfile,
      'ProfilePage should call fetchProfile on mount to determine create/edit mode.',
    ).toHaveBeenCalled();
    expect(
      mockFetchTags,
      'ProfilePage should call fetchTags on mount to load available tags.',
    ).toHaveBeenCalled();
  });

  it('shows success toast after creating profile', async () => {
    setupMockStores({ profile: null, error: null });
    const user = userEvent.setup();
    renderPage();

    await user.type(screen.getByLabelText(/first name/i), 'John');
    await user.type(screen.getByLabelText(/last name/i), 'Doe');
    await user.type(screen.getByLabelText(/username/i), 'johndoe');
    await user.selectOptions(screen.getByLabelText(/gender/i), 'male');
    await user.selectOptions(screen.getByLabelText(/sexual preference/i), 'heterosexual');
    await user.type(screen.getByLabelText(/birthday/i), '1995-06-15');
    await user.type(screen.getByLabelText(/biography/i), 'Hello world');
    await user.click(screen.getByRole('button', { name: /save/i }));

    await waitFor(() => {
      expect(
        toast.success,
        'Should show "Profile created!" toast after successful create.',
      ).toHaveBeenCalledWith('Profile created!');
    });
  });

  it('shows success toast after updating profile', async () => {
    setupMockStores({ profile: sampleProfile });
    const user = userEvent.setup();
    renderPage();

    await user.clear(screen.getByLabelText(/biography/i));
    await user.type(screen.getByLabelText(/biography/i), 'Updated bio');
    await user.click(screen.getByRole('button', { name: /save/i }));

    await waitFor(() => {
      expect(
        toast.success,
        'Should show "Profile updated!" toast after successful update.',
      ).toHaveBeenCalledWith('Profile updated!');
    });
  });
});
