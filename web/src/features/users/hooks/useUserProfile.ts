import { useEffect, useState } from 'react';
import * as usersApi from '@/api/users';
import type { UserProfileDetail } from '@/types';

interface UseUserProfileResult {
  readonly profile: UserProfileDetail | null;
  readonly isLoading: boolean;
  readonly error: string | null;
}

export function useUserProfile(userId: string | undefined): UseUserProfileResult {
  const [profile, setProfile] = useState<UserProfileDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!userId) {
      setIsLoading(false);
      return;
    }

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

  return { profile, isLoading, error };
}
