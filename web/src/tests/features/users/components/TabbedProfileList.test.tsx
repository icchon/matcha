import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { TabbedProfileList } from '@/features/users/components/TabbedProfileList';
import type { UserProfileDetail } from '@/types';
import type { TabConfig } from '@/features/users/components/TabbedProfileList';

const TABS: readonly [TabConfig, TabConfig] = [
  { value: 'tab-a', label: 'Tab A' },
  { value: 'tab-b', label: 'Tab B' },
];

const mockProfile: UserProfileDetail = {
  userId: '00000000-0000-0000-0000-000000000001',
  firstName: 'Alice',
  lastName: 'Smith',
  username: 'alice',
  gender: 'female',
  sexualPreference: 'bisexual',
  birthday: '1995-06-15',
  occupation: 'Designer',
  biography: 'Hello',
  locationName: 'Tokyo',
  fameRating: 75,
  pictures: [{ id: 1, userId: '00000000-0000-0000-0000-000000000001', url: '/img/a.jpg', isProfilePic: true, createdAt: '2024-01-01' }],
  tags: [],
  isOnline: true,
  lastConnection: null,
  distance: 3.5,
};

describe('TabbedProfileList', () => {
  it('renders loading spinner when isLoading is true', () => {
    render(
      <MemoryRouter>
        <TabbedProfileList
          title="Test"
          tabs={TABS}
          activeTab="tab-a"
          onTabChange={() => {}}
          myProfiles={[]}
          theirProfiles={[]}
          isLoading={true}
          emptyMessage="No items"
        />
      </MemoryRouter>,
    );

    expect(
      screen.getByRole('status'),
      'Should show spinner when isLoading=true.',
    ).toBeInTheDocument();
  });

  it('renders title and tabs', () => {
    render(
      <MemoryRouter>
        <TabbedProfileList
          title="My Title"
          tabs={TABS}
          activeTab="tab-a"
          onTabChange={() => {}}
          myProfiles={[mockProfile]}
          theirProfiles={[]}
          isLoading={false}
          emptyMessage="No items"
        />
      </MemoryRouter>,
    );

    expect(screen.getByText('My Title')).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: 'Tab A' })).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: 'Tab B' })).toBeInTheDocument();
  });

  it('renders profiles for active first tab', () => {
    render(
      <MemoryRouter>
        <TabbedProfileList
          title="Test"
          tabs={TABS}
          activeTab="tab-a"
          onTabChange={() => {}}
          myProfiles={[mockProfile]}
          theirProfiles={[]}
          isLoading={false}
          emptyMessage="No items"
        />
      </MemoryRouter>,
    );

    expect(
      screen.getByText('Alice'),
      'Should render profiles from myProfiles when first tab is active.',
    ).toBeInTheDocument();
  });

  it('renders empty message when no profiles', () => {
    render(
      <MemoryRouter>
        <TabbedProfileList
          title="Test"
          tabs={TABS}
          activeTab="tab-a"
          onTabChange={() => {}}
          myProfiles={[]}
          theirProfiles={[]}
          isLoading={false}
          emptyMessage="Nothing here"
        />
      </MemoryRouter>,
    );

    expect(
      screen.getByText('Nothing here'),
      'Should display emptyMessage when currentProfiles is empty.',
    ).toBeInTheDocument();
  });

  it('calls onTabChange when a tab is clicked', async () => {
    const onTabChange = vi.fn();
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <TabbedProfileList
          title="Test"
          tabs={TABS}
          activeTab="tab-a"
          onTabChange={onTabChange}
          myProfiles={[]}
          theirProfiles={[]}
          isLoading={false}
          emptyMessage="No items"
        />
      </MemoryRouter>,
    );

    await user.click(screen.getByRole('tab', { name: 'Tab B' }));

    expect(
      onTabChange,
      'onTabChange should be called with the clicked tab value.',
    ).toHaveBeenCalledWith('tab-b');
  });
});
