import { useState, useCallback } from 'react';
import { toast } from 'sonner';
import * as settingsApi from '@/api/settings';
import type { Block } from '@/types';

interface UseBlockListReturn {
  readonly blocks: readonly Block[];
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly fetchBlockList: () => Promise<void>;
  readonly unblock: (userId: string) => Promise<void>;
}

export function useBlockList(): UseBlockListReturn {
  const [blocks, setBlocks] = useState<readonly Block[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBlockList = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await settingsApi.getBlockList();
      setBlocks(result);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load block list';
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const unblock = useCallback(async (userId: string) => {
    setError(null);
    try {
      await settingsApi.unblockUser(userId);
      setBlocks((prev) => prev.filter((b) => b.blockedId !== userId));
      toast.success('User unblocked.');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to unblock user';
      setError(message);
      toast.error(message);
    }
  }, []);

  return { blocks, isLoading, error, fetchBlockList, unblock };
}
