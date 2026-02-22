import { useState, useCallback } from 'react';
import { toast } from 'sonner';
import * as usersApi from '@/api/users';

interface ApiError {
  readonly status?: number;
  readonly message?: string;
}

function getErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof Error) {
    const apiErr = err as Error & ApiError;
    if (typeof apiErr.status === 'number' && apiErr.status >= 500) {
      return 'Something went wrong. Please try again later.';
    }
    return apiErr.message;
  }
  return fallback;
}

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
      toast.error(getErrorMessage(err, 'Failed to like user'));
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
      toast.error(getErrorMessage(err, 'Failed to unlike user'));
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
      toast.error(getErrorMessage(err, 'Failed to block user'));
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
      toast.error(getErrorMessage(err, 'Failed to unblock user'));
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
      toast.error(getErrorMessage(err, 'Failed to report user'));
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
