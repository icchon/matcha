import { useState, useCallback } from 'react';
import { toast } from 'sonner';
import * as usersApi from '@/api/users';

interface UseProfileActionsResult {
  readonly isLiked: boolean;
  readonly isBlocked: boolean;
  readonly actionLoading: boolean;
  readonly handleLike: () => Promise<void>;
  readonly handleUnlike: () => Promise<void>;
  readonly handleBlock: () => Promise<void>;
  readonly handleUnblock: () => Promise<void>;
  readonly handleReport: () => Promise<void>;
}

export function useProfileActions(userId: string | undefined): UseProfileActionsResult {
  // TODO(FE-XX): Derive initial isLiked/isBlocked from API response or separate endpoint
  const [isLiked, setIsLiked] = useState(false);
  const [isBlocked, setIsBlocked] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

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

  return {
    isLiked,
    isBlocked,
    actionLoading,
    handleLike,
    handleUnlike,
    handleBlock,
    handleUnblock,
    handleReport,
  };
}
