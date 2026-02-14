import type { FC } from 'react';
import { Spinner } from '@/components/ui/Spinner';
import { ProfileCard } from '@/features/users/components/ProfileCard';
import type { UserProfileDetail } from '@/types';

interface TabConfig {
  readonly value: string;
  readonly label: string;
}

interface TabbedProfileListProps {
  readonly title: string;
  readonly tabs: readonly [TabConfig, TabConfig];
  readonly activeTab: string;
  readonly onTabChange: (tab: string) => void;
  readonly myProfiles: readonly UserProfileDetail[];
  readonly theirProfiles: readonly UserProfileDetail[];
  readonly isLoading: boolean;
  readonly emptyMessage: string;
}

const TabbedProfileList: FC<TabbedProfileListProps> = ({
  title,
  tabs,
  activeTab,
  onTabChange,
  myProfiles,
  theirProfiles,
  isLoading,
  emptyMessage,
}) => {
  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Spinner size="lg" />
      </div>
    );
  }

  const currentProfiles = activeTab === tabs[0].value ? myProfiles : theirProfiles;

  return (
    <div className="mx-auto max-w-4xl px-4 py-6">
      <h1 className="mb-6 text-2xl font-bold">{title}</h1>

      <div className="mb-6 flex gap-2" role="tablist">
        {tabs.map((tab) => (
          <button
            key={tab.value}
            role="tab"
            aria-selected={activeTab === tab.value}
            className={`rounded-lg px-4 py-2 text-sm font-medium ${
              activeTab === tab.value
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
            onClick={() => onTabChange(tab.value)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {currentProfiles.length === 0 ? (
        <p className="py-8 text-center text-gray-500">{emptyMessage}</p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {currentProfiles.map((profile) => (
            <ProfileCard key={profile.userId} profile={profile} />
          ))}
        </div>
      )}
    </div>
  );
};

export { TabbedProfileList };
export type { TabbedProfileListProps, TabConfig };
