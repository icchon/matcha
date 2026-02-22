import type { FC } from 'react';
import { useParams } from 'react-router-dom';
import { Spinner } from '@/components/ui/Spinner';
import { Badge } from '@/components/ui/Badge';
import { OnlineIndicator } from '@/features/users/components/OnlineIndicator';
import { ActionButtons } from '@/features/users/components/ActionButtons';
import { useUserProfile } from '@/features/users/hooks/useUserProfile';
import { useProfileActions } from '@/features/users/hooks/useProfileActions';
import { useAuthStore } from '@/stores/authStore';

const UserProfilePage: FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const currentUserId = useAuthStore((s) => s.userId);
  const { profile, isLoading, error } = useUserProfile(userId);
  const actions = useProfileActions(userId);
  const isOwnProfile = profile !== null && profile.userId === currentUserId;

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

      {!isOwnProfile && (
        <ActionButtons
          isLiked={actions.isLiked}
          isBlocked={actions.isBlocked}
          isLoading={actions.actionLoading}
          onLike={actions.handleLike}
          onUnlike={actions.handleUnlike}
          onBlock={actions.handleBlock}
          onUnblock={actions.handleUnblock}
          onReport={actions.handleReport}
        />
      )}
    </div>
  );
};

export { UserProfilePage };
