import type { FC } from 'react';
import { Link } from 'react-router-dom';
import { Card } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { OnlineIndicator } from '@/features/users/components/OnlineIndicator';
import { getProfilePicUrl, calculateAge } from '@/features/users/utils/profileHelpers';
import type { UserProfileDetail } from '@/types';

interface ProfileCardProps {
  readonly profile: UserProfileDetail;
  readonly onLike?: (userId: string) => void;
}

const ProfileCard: FC<ProfileCardProps> = ({ profile, onLike }) => {
  const picUrl = getProfilePicUrl(profile);
  const age = calculateAge(profile.birthday);

  return (
    <Card className="flex flex-col gap-3">
      {picUrl ? (
        <img
          src={picUrl}
          alt={`${profile.firstName ?? 'User'}'s photo`}
          className="h-48 w-full rounded-md object-cover"
        />
      ) : (
        <div
          data-testid="avatar-placeholder"
          className="flex h-48 w-full items-center justify-center rounded-md bg-gray-200 text-gray-400"
        >
          No photo
        </div>
      )}

      <div className="flex items-center justify-between">
        <div>
          <span className="text-lg font-semibold">
            {profile.firstName}
            {age !== null && <span className="ml-1 text-gray-500">{age}</span>}
          </span>
        </div>
        <OnlineIndicator isOnline={profile.isOnline} lastConnection={profile.lastConnection} />
      </div>

      {profile.locationName && (
        <p className="text-sm text-gray-500">{profile.locationName}</p>
      )}

      {profile.fameRating !== null && (
        <p className="text-sm text-gray-500">
          Fame: <span className="font-medium">{profile.fameRating}</span>
        </p>
      )}

      {profile.tags.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {profile.tags.map((tag) => (
            <Badge key={tag.id}>{tag.name}</Badge>
          ))}
        </div>
      )}

      <div className="flex items-center gap-2">
        <Link
          to={`/users/${profile.userId}`}
          className="text-sm text-blue-600 hover:underline"
        >
          View Profile
        </Link>
        {onLike && (
          <Button
            type="button"
            onClick={() => onLike(profile.userId)}
            className="ml-auto"
          >
            Like
          </Button>
        )}
      </div>
    </Card>
  );
};

export { ProfileCard };
export type { ProfileCardProps };
