import { useEffect, useState, useCallback, type FC } from 'react';
import { toast } from 'sonner';
import { Spinner } from '@/components/ui/Spinner';
import { ProfileCard } from '@/features/users/components/ProfileCard';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

type TabValue = 'liked-by-me' | 'who-liked-me';

const LikesPage: FC = () => {
  const [activeTab, setActiveTab] = useState<TabValue>('liked-by-me');
  const [likedByMeProfiles, setLikedByMeProfiles] = useState<readonly UserProfileDetail[]>([]);
  const [whoLikedMeProfiles, setWhoLikedMeProfiles] = useState<readonly UserProfileDetail[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchProfiles = useCallback(async () => {
    setIsLoading(true);
    try {
      const [likedByMe, whoLikedMe] = await Promise.all([
        usersApi.getLikedUsers(),
        usersApi.getWhoLikedMe(),
      ]);

      const likedIds = likedByMe.map((l) => l.likedId);
      const whoLikedIds = whoLikedMe.map((l) => l.likerId);
      const allIds = [...new Set([...likedIds, ...whoLikedIds])];

      const profiles = await Promise.all(
        allIds.map((id) => usersApi.getUserProfile(id).catch(() => null)),
      );

      const profileMap = new Map<string, UserProfileDetail>();
      for (const p of profiles) {
        if (p) profileMap.set(p.userId, p);
      }

      setLikedByMeProfiles(likedIds.map((id) => profileMap.get(id)).filter(Boolean) as UserProfileDetail[]);
      setWhoLikedMeProfiles(whoLikedIds.map((id) => profileMap.get(id)).filter(Boolean) as UserProfileDetail[]);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load likes';
      toast.error(message);
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

  const currentProfiles = activeTab === 'liked-by-me' ? likedByMeProfiles : whoLikedMeProfiles;

  return (
    <div className="mx-auto max-w-4xl px-4 py-6">
      <h1 className="mb-6 text-2xl font-bold">Likes</h1>

      <div className="mb-6 flex gap-2" role="tablist">
        <button
          role="tab"
          aria-selected={activeTab === 'liked-by-me'}
          className={`rounded-lg px-4 py-2 text-sm font-medium ${
            activeTab === 'liked-by-me'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
          onClick={() => setActiveTab('liked-by-me')}
        >
          Liked by me
        </button>
        <button
          role="tab"
          aria-selected={activeTab === 'who-liked-me'}
          className={`rounded-lg px-4 py-2 text-sm font-medium ${
            activeTab === 'who-liked-me'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
          onClick={() => setActiveTab('who-liked-me')}
        >
          Who liked me
        </button>
      </div>

      {currentProfiles.length === 0 ? (
        <p className="py-8 text-center text-gray-500">No likes yet</p>
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

export { LikesPage };
