import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { ProfileCard } from '@/features/users/components/ProfileCard';
import type { UserProfileDetail } from '@/types';

const mockProfile: UserProfileDetail = {
  userId: 'user-1',
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
    { id: 1, userId: 'user-1', url: '/images/alice.jpg', isProfilePic: true, createdAt: '2024-01-01' },
  ],
  tags: [
    { id: 1, name: 'art' },
    { id: 2, name: 'coffee' },
  ],
  isOnline: true,
  lastConnection: null,
  distance: 3.5,
};

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('ProfileCard', () => {
  it('renders user name and age', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    expect(
      screen.getByText('Alice'),
      'Should display first name. Check ProfileCard rendering.',
    ).toBeInTheDocument();
  });

  it('renders location', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    expect(
      screen.getByText(/Tokyo/),
      'Should display locationName. Check ProfileCard rendering.',
    ).toBeInTheDocument();
  });

  it('renders online status indicator', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    expect(
      screen.getByTestId('online-indicator'),
      'Should render OnlineIndicator component.',
    ).toBeInTheDocument();
  });

  it('renders tags as badges', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    expect(screen.getByText('art')).toBeInTheDocument();
    expect(screen.getByText('coffee')).toBeInTheDocument();
  });

  it('renders fame rating', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    expect(
      screen.getByText(/75/),
      'Should display fameRating. Check ProfileCard rendering.',
    ).toBeInTheDocument();
  });

  it('renders profile picture', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    const img = screen.getByRole('img');
    expect(
      img.getAttribute('src'),
      'Should render the profile picture URL.',
    ).toBe('/images/alice.jpg');
  });

  it('renders link to user profile page', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    const link = screen.getByRole('link', { name: /view profile/i });
    expect(
      link.getAttribute('href'),
      'Should link to /users/:userId.',
    ).toBe('/users/user-1');
  });

  it('calls onLike when like button is clicked', async () => {
    const onLike = vi.fn();
    const user = userEvent.setup();
    renderWithRouter(<ProfileCard profile={mockProfile} onLike={onLike} />);

    const likeButton = screen.getByRole('button', { name: /like/i });
    await user.click(likeButton);

    expect(
      onLike,
      'onLike callback should be called with userId when like button clicked.',
    ).toHaveBeenCalledWith('user-1');
  });

  it('does not render like button when onLike is not provided', () => {
    renderWithRouter(<ProfileCard profile={mockProfile} />);

    expect(screen.queryByRole('button', { name: /like/i })).not.toBeInTheDocument();
  });

  it('renders placeholder when no profile picture', () => {
    const profileNoPic: UserProfileDetail = {
      ...mockProfile,
      pictures: [],
    };
    renderWithRouter(<ProfileCard profile={profileNoPic} />);

    expect(
      screen.getByTestId('avatar-placeholder'),
      'Should render a placeholder when no pictures are available.',
    ).toBeInTheDocument();
  });
});
