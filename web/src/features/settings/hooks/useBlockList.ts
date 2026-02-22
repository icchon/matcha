import { useState, useCallback } from 'react';
import { toast } from 'sonner';
import * as settingsApi from '@/api/settings';
import { ApiClientError } from '@/api/client';
import type { Block } from '@/types';

function mapBlockListErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof ApiClientError) {
    switch (err.status) {
      case 401:
        return 'Session expired. Please log in again.';
      case 404:
        return 'User not found.';
      default:
        return fallback;
    }
  }
  return fallback;
}

interface UseBlockListReturn {
  readonly blocks: readonly Block[];
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly unblockingId: string | null;
  readonly fetchBlockList: () => Promise<void>;
  readonly unblock: (userId: string) => Promise<void>;
}

export function useBlockList(): UseBlockListReturn {
  const [blocks, setBlocks] = useState<readonly Block[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [unblockingId, setUnblockingId] = useState<string | null>(null);

  const fetchBlockList = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await settingsApi.getBlockList();
      setBlocks(result);
    } catch (err) {
      const message = mapBlockListErrorMessage(err, 'Failed to load block list');
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const unblock = useCallback(async (userId: string) => {
    setError(null);
    setUnblockingId(userId);
    try {
      await settingsApi.unblockUser(userId);
      setBlocks((prev) => prev.filter((b) => b.blockedId !== userId));
      toast.success('User unblocked.');
    } catch (err) {
      const message = mapBlockListErrorMessage(err, 'Failed to unblock user');
      setError(message);
      toast.error(message);
    } finally {
      setUnblockingId(null);
    }
  }, []);

  return { blocks, isLoading, error, unblockingId, fetchBlockList, unblock };
}
