import { useEffect, useState, useCallback } from 'react';
import { toast } from 'sonner';
import * as usersApi from '@/api/users';
import { getErrorMessage } from '@/features/users/utils/errorHelpers';
import type { UserProfileDetail } from '@/types';

interface TabbedListConfig<T> {
  readonly fetchMyList: () => Promise<readonly T[]>;
  readonly fetchTheirList: () => Promise<readonly T[]>;
  readonly extractMyIds: (items: readonly T[]) => readonly string[];
  readonly extractTheirIds: (items: readonly T[]) => readonly string[];
  readonly errorMessage: string;
}

interface UseTabbedProfileListResult<TTab extends string> {
  readonly activeTab: TTab;
  readonly setActiveTab: (tab: TTab) => void;
  readonly myProfiles: readonly UserProfileDetail[];
  readonly theirProfiles: readonly UserProfileDetail[];
  readonly isLoading: boolean;
}

export function useTabbedProfileList<T, TTab extends string>(
  config: TabbedListConfig<T>,
  defaultTab: TTab,
): UseTabbedProfileListResult<TTab> {
  const [activeTab, setActiveTab] = useState<TTab>(defaultTab);
  const [myProfiles, setMyProfiles] = useState<readonly UserProfileDetail[]>([]);
  const [theirProfiles, setTheirProfiles] = useState<readonly UserProfileDetail[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchProfiles = useCallback(async () => {
    setIsLoading(true);
    try {
      const [myList, theirList] = await Promise.all([
        config.fetchMyList(),
        config.fetchTheirList(),
      ]);

      const myIds = config.extractMyIds(myList);
      const theirIds = config.extractTheirIds(theirList);
      const allIds = [...new Set([...myIds, ...theirIds])];

      // TODO(BE-XX): Replace N+1 getUserProfile calls with a bulk-fetch endpoint
      // (e.g., GET /users/profiles?ids=id1,id2,...) to avoid per-user requests
      const profiles = await Promise.all(
        allIds.map((id) => usersApi.getUserProfile(id).catch(() => null)),
      );

      const profileMap = new Map<string, UserProfileDetail>();
      for (const p of profiles) {
        if (p) profileMap.set(p.userId, p);
      }

      setMyProfiles(
        myIds.map((id) => profileMap.get(id)).filter(Boolean) as UserProfileDetail[],
      );
      setTheirProfiles(
        theirIds.map((id) => profileMap.get(id)).filter(Boolean) as UserProfileDetail[],
      );
    } catch (err) {
      toast.error(getErrorMessage(err, config.errorMessage));
    } finally {
      setIsLoading(false);
    }
  }, [config]);

  useEffect(() => {
    fetchProfiles();
  }, [fetchProfiles]);

  return { activeTab, setActiveTab, myProfiles, theirProfiles, isLoading };
}
