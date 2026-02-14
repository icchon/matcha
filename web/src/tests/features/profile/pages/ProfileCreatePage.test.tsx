import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { ProfileCreatePage } from '@/features/profile/pages/ProfileCreatePage';
import { useProfileStore } from '@/stores/profileStore';

vi.mock('@/stores/profileStore');

const mockSaveProfile = vi.fn();
const mockFetchTags = vi.fn();
const mockUploadPicture = vi.fn();
const mockDeletePicture = vi.fn();
const mockAddTag = vi.fn();
const mockRemoveTag = vi.fn();

function setupMockStore(overrides: Record<string, unknown> = {}) {
  vi.mocked(useProfileStore).mockImplementation((selector: unknown) => {
    const state = {
      profile: null,
      pictures: [],
      tags: [],
      allTags: [{ id: 1, name: 'hiking' }, { id: 2, name: 'cooking' }],
      isLoading: false,
      error: null,
      saveProfile: mockSaveProfile,
      fetchTags: mockFetchTags,
      uploadPicture: mockUploadPicture,
      deletePicture: mockDeletePicture,
      addTag: mockAddTag,
      removeTag: mockRemoveTag,
      clearError: vi.fn(),
      fetchProfile: vi.fn(),
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
      <ProfileCreatePage />
    </MemoryRouter>,
  );
}

describe('ProfileCreatePage', () => {
  it('renders page title', () => {
    renderPage();

    expect(
      screen.getByText(/create.*profile/i),
      'ProfileCreatePage should display a "Create Profile" heading.',
    ).toBeInTheDocument();
  });

  it('renders profile form', () => {
    renderPage();

    expect(
      screen.getByLabelText(/first name/i),
      'ProfileCreatePage should render the ProfileForm with a First Name field.',
    ).toBeInTheDocument();
  });

  it('renders photo uploader', () => {
    renderPage();

    expect(
      screen.getByText(/photos/i),
      'ProfileCreatePage should render the PhotoUploader section.',
    ).toBeInTheDocument();
  });

  it('renders tag manager', () => {
    renderPage();

    expect(
      screen.getByText(/interest tags/i),
      'ProfileCreatePage should render the TagManager section.',
    ).toBeInTheDocument();
  });

  it('calls saveProfile when form is submitted', async () => {
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
        mockSaveProfile,
        'Submitting the form should call saveProfile from the store.',
      ).toHaveBeenCalled();
    });
  });

  it('displays error message when error is set', () => {
    setupMockStore({ error: 'Something went wrong' });
    renderPage();

    expect(
      screen.getByText('Something went wrong'),
      'Error message should be displayed when store has an error.',
    ).toBeInTheDocument();
  });
});
