import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { UserProfilePage } from '@/features/users/pages/UserProfilePage';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

vi.mock('@/api/users');

const mockProfile: UserProfileDetail = {
  userId: '00000000-0000-0000-0000-000000000001',
  firstName: 'Alice',
  lastName: 'Smith',
  username: 'alice',
  gender: 'female',
  sexualPreference: 'bisexual',
  birthday: '1995-06-15',
  occupation: 'Designer',
  biography: 'Love art and coffee',
  locationName: 'Tokyo',
  fameRating: 75,
  pictures: [
    { id: 1, userId: '00000000-0000-0000-0000-000000000001', url: '/images/alice.jpg', isProfilePic: true, createdAt: '2024-01-01' },
    { id: 2, userId: '00000000-0000-0000-0000-000000000001', url: '/images/alice2.jpg', isProfilePic: false, createdAt: '2024-01-02' },
  ],
  tags: [
    { id: 1, name: 'art' },
    { id: 2, name: 'coffee' },
  ],
  isOnline: true,
  lastConnection: null,
  distance: 3.5,
};

function renderPage(userId = '00000000-0000-0000-0000-000000000001') {
  return render(
    <MemoryRouter initialEntries={[`/users/${userId}`]}>
      <Routes>
        <Route path="/users/:userId" element={<UserProfilePage />} />
      </Routes>
    </MemoryRouter>,
  );
}

beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(usersApi.getUserProfile).mockResolvedValue(mockProfile);
  vi.mocked(usersApi.likeUser).mockResolvedValue({ matched: false });
  vi.mocked(usersApi.unlikeUser).mockResolvedValue(undefined);
  vi.mocked(usersApi.blockUser).mockResolvedValue(undefined);
  vi.mocked(usersApi.unblockUser).mockResolvedValue(undefined);
  vi.mocked(usersApi.reportUser).mockResolvedValue({ message: 'Report submitted' });
});

describe('UserProfilePage', () => {
  it('shows loading spinner while fetching profile', () => {
    vi.mocked(usersApi.getUserProfile).mockReturnValue(new Promise(() => {}));
    renderPage();

    expect(
      screen.getByRole('status'),
      'Should show a loading spinner while the profile is being fetched.',
    ).toBeInTheDocument();
  });

  it('renders profile details after loading', async () => {
    renderPage();

    await waitFor(() => {
      expect(
        screen.getByText('Alice'),
        'Should display user first name after loading.',
      ).toBeInTheDocument();
    });

    expect(screen.getByText(/Designer/)).toBeInTheDocument();
    expect(screen.getByText(/Love art and coffee/)).toBeInTheDocument();
    expect(screen.getByText(/Tokyo/)).toBeInTheDocument();
    expect(screen.getByText('art')).toBeInTheDocument();
    expect(screen.getByText('coffee')).toBeInTheDocument();
  });

  it('renders all photos', async () => {
    renderPage();

    await waitFor(() => {
      expect(screen.getByText('Alice')).toBeInTheDocument();
    });

    const images = screen.getAllByRole('img');
    expect(
      images.length,
      'Should render all profile pictures. Check photo gallery rendering.',
    ).toBe(2);
  });

  it('renders online indicator', async () => {
    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('online-indicator')).toBeInTheDocument();
    });
  });

  it('renders action buttons', async () => {
    renderPage();

    await waitFor(() => {
      expect(screen.getByText('Alice')).toBeInTheDocument();
    });

    expect(screen.getByRole('button', { name: /like/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /block/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /report/i })).toBeInTheDocument();
  });

  it('calls likeUser API when like button is clicked', async () => {
    const user = userEvent.setup();
    renderPage();

    await waitFor(() => {
      expect(screen.getByText('Alice')).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: /like/i }));

    expect(
      usersApi.likeUser,
      'Should call likeUser API with the user ID.',
    ).toHaveBeenCalledWith('00000000-0000-0000-0000-000000000001');
  });

  it('shows error message when profile fetch fails', async () => {
    vi.mocked(usersApi.getUserProfile).mockRejectedValue(new Error('Not found'));
    renderPage();

    await waitFor(() => {
      expect(
        screen.getByText(/not found/i),
        'Should display error message when profile fetch fails.',
      ).toBeInTheDocument();
    });
  });

  it('shows fame rating', async () => {
    renderPage();

    await waitFor(() => {
      expect(screen.getByText('Alice')).toBeInTheDocument();
    });

    expect(screen.getByText(/75/)).toBeInTheDocument();
  });
});
