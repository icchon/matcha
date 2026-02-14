import { useEffect, useState, useCallback, type FC } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from 'sonner';
import { Spinner } from '@/components/ui/Spinner';
import { Badge } from '@/components/ui/Badge';
import { OnlineIndicator } from '@/features/users/components/OnlineIndicator';
import { ActionButtons } from '@/features/users/components/ActionButtons';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

const UserProfilePage: FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const [profile, setProfile] = useState<UserProfileDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLiked, setIsLiked] = useState(false);
  const [isBlocked, setIsBlocked] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    if (!userId) return;

    let cancelled = false;
    setIsLoading(true);
    setError(null);

    usersApi.getUserProfile(userId).then((data) => {
      if (!cancelled) {
        setProfile(data);
        setIsLoading(false);
      }
    }).catch((err) => {
      if (!cancelled) {
        const message = err instanceof Error ? err.message : 'Failed to load profile';
        setError(message);
        setIsLoading(false);
      }
    });

    return () => { cancelled = true; };
  }, [userId]);

  const handleLike = useCallback(async () => {
    if (!userId) return;
    setActionLoading(true);
    try {
      const result = await usersApi.likeUser(userId);
      setIsLiked(true);
      if (result.matched) {
        toast.success("It's a match!");
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to like user';
      toast.error(message);
    } finally {
      setActionLoading(false);
    }
  }, [userId]);

  const handleUnlike = useCallback(async () => {
    if (!userId) return;
    setActionLoading(true);
    try {
      await usersApi.unlikeUser(userId);
      setIsLiked(false);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to unlike user';
      toast.error(message);
    } finally {
      setActionLoading(false);
    }
  }, [userId]);

  const handleBlock = useCallback(async () => {
    if (!userId) return;
    setActionLoading(true);
    try {
      await usersApi.blockUser(userId);
      setIsBlocked(true);
      toast.success('User blocked');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to block user';
      toast.error(message);
    } finally {
      setActionLoading(false);
    }
  }, [userId]);

  const handleUnblock = useCallback(async () => {
    if (!userId) return;
    setActionLoading(true);
    try {
      await usersApi.unblockUser(userId);
      setIsBlocked(false);
      toast.success('User unblocked');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to unblock user';
      toast.error(message);
    } finally {
      setActionLoading(false);
    }
  }, [userId]);

  const handleReport = useCallback(async () => {
    if (!userId) return;
    setActionLoading(true);
    try {
      await usersApi.reportUser(userId, 'inappropriate');
      toast.success('Report submitted');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to report user';
      toast.error(message);
    } finally {
      setActionLoading(false);
    }
  }, [userId]);

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="py-12 text-center text-red-600">
        <p>{error}</p>
      </div>
    );
  }

  if (!profile) return null;

  return (
    <div className="mx-auto max-w-2xl space-y-6 px-4 py-6">
      <div className="space-y-4">
        {profile.pictures.length > 0 && (
          <div className="grid grid-cols-2 gap-2">
            {profile.pictures.map((pic) => (
              <img
                key={pic.id}
                src={pic.url}
                alt={`${profile.firstName ?? 'User'}'s photo`}
                className="h-64 w-full rounded-lg object-cover"
              />
            ))}
          </div>
        )}

        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold">
            {profile.firstName}
            {profile.lastName && <span className="ml-1">{profile.lastName}</span>}
          </h1>
          <OnlineIndicator isOnline={profile.isOnline} lastConnection={profile.lastConnection} />
        </div>

        {profile.occupation && (
          <p className="text-gray-600">{profile.occupation}</p>
        )}

        {profile.locationName && (
          <p className="text-sm text-gray-500">{profile.locationName}</p>
        )}

        {profile.fameRating !== null && (
          <p className="text-sm text-gray-500">
            Fame: <span className="font-medium">{profile.fameRating}</span>
          </p>
        )}

        {profile.biography && (
          <div>
            <h2 className="mb-1 font-semibold">About</h2>
            <p className="text-gray-700">{profile.biography}</p>
          </div>
        )}

        {profile.tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {profile.tags.map((tag) => (
              <Badge key={tag.id}>{tag.name}</Badge>
            ))}
          </div>
        )}
      </div>

      <ActionButtons
        userId={profile.userId}
        isLiked={isLiked}
        isBlocked={isBlocked}
        isLoading={actionLoading}
        onLike={handleLike}
        onUnlike={handleUnlike}
        onBlock={handleBlock}
        onUnblock={handleUnblock}
        onReport={handleReport}
      />
    </div>
  );
};

export { UserProfilePage };
