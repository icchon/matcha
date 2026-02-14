import { useMemo, type FC } from 'react';
import * as usersApi from '@/api/users';
import { useTabbedProfileList } from '@/features/users/hooks/useTabbedProfileList';
import { TabbedProfileList } from '@/features/users/components/TabbedProfileList';
import type { View } from '@/types';
import type { TabConfig } from '@/features/users/components/TabbedProfileList';

const TABS: readonly [TabConfig, TabConfig] = [
  { value: 'viewed-by-me', label: 'Profiles I viewed' },
  { value: 'who-viewed-me', label: 'Who viewed me' },
] as const;

const ViewsPage: FC = () => {
  const config = useMemo(() => ({
    fetchMyList: usersApi.getViewedUsers,
    fetchTheirList: usersApi.getWhoViewedMe,
    extractMyIds: (items: readonly View[]) => items.map((v) => v.viewedId),
    extractTheirIds: (items: readonly View[]) => items.map((v) => v.viewerId),
    errorMessage: 'Failed to load views',
  }), []);

  const { activeTab, setActiveTab, myProfiles, theirProfiles, isLoading } =
    useTabbedProfileList<View, string>(config, TABS[0].value);

  return (
    <TabbedProfileList
      title="Profile Views"
      tabs={TABS}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      myProfiles={myProfiles}
      theirProfiles={theirProfiles}
      isLoading={isLoading}
      emptyMessage="No views yet"
    />
  );
};

export { ViewsPage };
