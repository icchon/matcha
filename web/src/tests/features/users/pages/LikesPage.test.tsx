import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { LikesPage } from '@/features/users/pages/LikesPage';
import * as usersApi from '@/api/users';
import type { Like, UserProfileDetail } from '@/types';

vi.mock('@/api/users');

const mockLikedByMe: readonly Like[] = [
  { likerId: 'me', likedId: '00000000-0000-0000-0000-000000000002', createdAt: '2024-01-01' },
];

const mockLikedMe: readonly Like[] = [
  { likerId: '00000000-0000-0000-0000-000000000003', likedId: 'me', createdAt: '2024-01-02' },
];

const mockProfile: UserProfileDetail = {
  userId: '00000000-0000-0000-0000-000000000002',
  firstName: 'Bob',
  lastName: 'Jones',
  username: 'bob',
  gender: 'male',
  sexualPreference: 'heterosexual',
  birthday: '1992-03-20',
  occupation: 'Dev',
  biography: 'Hello',
  locationName: 'Paris',
  fameRating: 50,
  pictures: [{ id: 1, userId: '00000000-0000-0000-0000-000000000002', url: '/images/bob.jpg', isProfilePic: true, createdAt: '2024-01-01' }],
  tags: [],
  isOnline: false,
  lastConnection: '2024-01-15T10:00:00Z',
};

const mockProfile2: UserProfileDetail = {
  ...mockProfile,
  userId: '00000000-0000-0000-0000-000000000003',
  firstName: 'Carol',
};

beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(usersApi.getLikedUsers).mockResolvedValue(mockLikedByMe);
  vi.mocked(usersApi.getWhoLikedMe).mockResolvedValue(mockLikedMe);
  vi.mocked(usersApi.getUserProfile).mockImplementation(async (userId: string) => {
    if (userId === '00000000-0000-0000-0000-000000000002') return mockProfile;
    if (userId === '00000000-0000-0000-0000-000000000003') return mockProfile2;
    throw new Error('Not found');
  });
});

describe('LikesPage', () => {
  it('shows loading spinner initially', () => {
    vi.mocked(usersApi.getLikedUsers).mockReturnValue(new Promise(() => {}));
    vi.mocked(usersApi.getWhoLikedMe).mockReturnValue(new Promise(() => {}));
    render(<MemoryRouter><LikesPage /></MemoryRouter>);

    expect(
      screen.getByRole('status'),
      'Should show loading spinner while fetching likes data.',
    ).toBeInTheDocument();
  });

  it('renders tabs for "Liked by me" and "Who liked me"', async () => {
    render(<MemoryRouter><LikesPage /></MemoryRouter>);

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /liked by me/i }),
        'Should render "Liked by me" tab.',
      ).toBeInTheDocument();
    });
    expect(screen.getByRole('tab', { name: /who liked me/i })).toBeInTheDocument();
  });

  it('shows users I liked in the default tab', async () => {
    render(<MemoryRouter><LikesPage /></MemoryRouter>);

    await waitFor(() => {
      expect(
        screen.getByText('Bob'),
        'Should display profiles of users I liked.',
      ).toBeInTheDocument();
    });
  });

  it('switches to "Who liked me" tab and shows those users', async () => {
    const user = userEvent.setup();
    render(<MemoryRouter><LikesPage /></MemoryRouter>);

    await waitFor(() => {
      expect(screen.getByRole('tab', { name: /who liked me/i })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('tab', { name: /who liked me/i }));

    await waitFor(() => {
      expect(
        screen.getByText('Carol'),
        'Should display profiles of users who liked me after switching tab.',
      ).toBeInTheDocument();
    });
  });

  it('shows empty message when no likes', async () => {
    vi.mocked(usersApi.getLikedUsers).mockResolvedValue([]);
    render(<MemoryRouter><LikesPage /></MemoryRouter>);

    await waitFor(() => {
      expect(
        screen.getByText(/no likes yet/i),
        'Should display empty state message when no likes exist.',
      ).toBeInTheDocument();
    });
  });
});
