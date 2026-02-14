import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { ViewsPage } from '@/features/users/pages/ViewsPage';
import * as usersApi from '@/api/users';
import type { View, UserProfileDetail } from '@/types';

vi.mock('@/api/users');

const mockViewedByMe: readonly View[] = [
  { viewerId: 'me', viewedId: 'user-2', viewTime: '2024-01-01T10:00:00Z' },
];

const mockViewedMe: readonly View[] = [
  { viewerId: 'user-3', viewedId: 'me', viewTime: '2024-01-02T10:00:00Z' },
];

const mockProfile: UserProfileDetail = {
  userId: 'user-2',
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
  pictures: [{ id: 1, userId: 'user-2', url: '/images/bob.jpg', isProfilePic: true, createdAt: '2024-01-01' }],
  tags: [],
  isOnline: false,
  lastConnection: '2024-01-15T10:00:00Z',
};

const mockProfile2: UserProfileDetail = {
  ...mockProfile,
  userId: 'user-3',
  firstName: 'Carol',
};

beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(usersApi.getViewedUsers).mockResolvedValue(mockViewedByMe);
  vi.mocked(usersApi.getWhoViewedMe).mockResolvedValue(mockViewedMe);
  vi.mocked(usersApi.getUserProfile).mockImplementation(async (userId: string) => {
    if (userId === 'user-2') return mockProfile;
    if (userId === 'user-3') return mockProfile2;
    throw new Error('Not found');
  });
});

describe('ViewsPage', () => {
  it('shows loading spinner initially', () => {
    vi.mocked(usersApi.getViewedUsers).mockReturnValue(new Promise(() => {}));
    vi.mocked(usersApi.getWhoViewedMe).mockReturnValue(new Promise(() => {}));
    render(<MemoryRouter><ViewsPage /></MemoryRouter>);

    expect(
      screen.getByRole('status'),
      'Should show loading spinner while fetching views data.',
    ).toBeInTheDocument();
  });

  it('renders tabs for "Profiles I viewed" and "Who viewed me"', async () => {
    render(<MemoryRouter><ViewsPage /></MemoryRouter>);

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /profiles i viewed/i }),
        'Should render "Profiles I viewed" tab.',
      ).toBeInTheDocument();
    });
    expect(screen.getByRole('tab', { name: /who viewed me/i })).toBeInTheDocument();
  });

  it('shows profiles I viewed in the default tab', async () => {
    render(<MemoryRouter><ViewsPage /></MemoryRouter>);

    await waitFor(() => {
      expect(
        screen.getByText('Bob'),
        'Should display profiles I viewed.',
      ).toBeInTheDocument();
    });
  });

  it('switches to "Who viewed me" tab and shows those users', async () => {
    const user = userEvent.setup();
    render(<MemoryRouter><ViewsPage /></MemoryRouter>);

    await waitFor(() => {
      expect(screen.getByRole('tab', { name: /who viewed me/i })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('tab', { name: /who viewed me/i }));

    await waitFor(() => {
      expect(
        screen.getByText('Carol'),
        'Should display profiles of users who viewed me after switching tab.',
      ).toBeInTheDocument();
    });
  });

  it('shows empty message when no views', async () => {
    vi.mocked(usersApi.getViewedUsers).mockResolvedValue([]);
    render(<MemoryRouter><ViewsPage /></MemoryRouter>);

    await waitFor(() => {
      expect(
        screen.getByText(/no views yet/i),
        'Should display empty state message when no views exist.',
      ).toBeInTheDocument();
    });
  });
});
