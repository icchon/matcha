import { useMemo, type FC } from 'react';
import * as usersApi from '@/api/users';
import { useTabbedProfileList } from '@/features/users/hooks/useTabbedProfileList';
import { TabbedProfileList } from '@/features/users/components/TabbedProfileList';
import type { Like } from '@/types';
import type { TabConfig } from '@/features/users/components/TabbedProfileList';

const TABS: readonly [TabConfig, TabConfig] = [
  { value: 'liked-by-me', label: 'Liked by me' },
  { value: 'who-liked-me', label: 'Who liked me' },
] as const;

const LikesPage: FC = () => {
  const config = useMemo(() => ({
    fetchMyList: usersApi.getLikedUsers,
    fetchTheirList: usersApi.getWhoLikedMe,
    extractMyIds: (items: readonly Like[]) => items.map((l) => l.likedId),
    extractTheirIds: (items: readonly Like[]) => items.map((l) => l.likerId),
    errorMessage: 'Failed to load likes',
  }), []);

  const { activeTab, setActiveTab, myProfiles, theirProfiles, isLoading } =
    useTabbedProfileList<Like, string>(config, TABS[0].value);

  return (
    <TabbedProfileList
      title="Likes"
      tabs={TABS}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      myProfiles={myProfiles}
      theirProfiles={theirProfiles}
      isLoading={isLoading}
      emptyMessage="No likes yet"
    />
  );
};

export { LikesPage };
