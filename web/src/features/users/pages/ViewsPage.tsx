import { useEffect, useState, useCallback, type FC } from 'react';
import { Spinner } from '@/components/ui/Spinner';
import { ProfileCard } from '@/features/users/components/ProfileCard';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

type TabValue = 'viewed-by-me' | 'who-viewed-me';

const ViewsPage: FC = () => {
  const [activeTab, setActiveTab] = useState<TabValue>('viewed-by-me');
  const [viewedByMeProfiles, setViewedByMeProfiles] = useState<readonly UserProfileDetail[]>([]);
  const [whoViewedMeProfiles, setWhoViewedMeProfiles] = useState<readonly UserProfileDetail[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchProfiles = useCallback(async () => {
    setIsLoading(true);
    try {
      const [viewedByMe, whoViewedMe] = await Promise.all([
        usersApi.getViewedUsers(),
        usersApi.getWhoViewedMe(),
      ]);

      const viewedIds = viewedByMe.map((v) => v.viewedId);
      const viewerIds = whoViewedMe.map((v) => v.viewerId);
      const allIds = [...new Set([...viewedIds, ...viewerIds])];

      const profiles = await Promise.all(
        allIds.map((id) => usersApi.getUserProfile(id).catch(() => null)),
      );

      const profileMap = new Map<string, UserProfileDetail>();
      for (const p of profiles) {
        if (p) profileMap.set(p.userId, p);
      }

      setViewedByMeProfiles(viewedIds.map((id) => profileMap.get(id)).filter(Boolean) as UserProfileDetail[]);
      setWhoViewedMeProfiles(viewerIds.map((id) => profileMap.get(id)).filter(Boolean) as UserProfileDetail[]);
    } catch {
      // Non-critical: empty lists shown on error
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchProfiles();
  }, [fetchProfiles]);

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Spinner size="lg" />
      </div>
    );
  }

  const currentProfiles = activeTab === 'viewed-by-me' ? viewedByMeProfiles : whoViewedMeProfiles;

  return (
    <div className="mx-auto max-w-4xl px-4 py-6">
      <h1 className="mb-6 text-2xl font-bold">Profile Views</h1>

      <div className="mb-6 flex gap-2" role="tablist">
        <button
          role="tab"
          aria-selected={activeTab === 'viewed-by-me'}
          className={`rounded-lg px-4 py-2 text-sm font-medium ${
            activeTab === 'viewed-by-me'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
          onClick={() => setActiveTab('viewed-by-me')}
        >
          Profiles I viewed
        </button>
        <button
          role="tab"
          aria-selected={activeTab === 'who-viewed-me'}
          className={`rounded-lg px-4 py-2 text-sm font-medium ${
            activeTab === 'who-viewed-me'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
          onClick={() => setActiveTab('who-viewed-me')}
        >
          Who viewed me
        </button>
      </div>

      {currentProfiles.length === 0 ? (
        <p className="py-8 text-center text-gray-500">No views yet</p>
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

export { ViewsPage };
