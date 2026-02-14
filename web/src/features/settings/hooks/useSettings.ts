import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import * as settingsApi from '@/api/settings';
import { clearTokens } from '@/api/client';
import type { ChangePasswordRequest } from '@/api/settings';
import type { Block } from '@/types';

interface UseChangePasswordReturn {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly changePassword: (params: ChangePasswordRequest) => Promise<boolean>;
}

export function useChangePassword(): UseChangePasswordReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const changePassword = useCallback(async (params: ChangePasswordRequest): Promise<boolean> => {
    setIsLoading(true);
    setError(null);
    try {
      await settingsApi.changePassword(params);
      toast.success('Password changed successfully!');
      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to change password';
      setError(message);
      toast.error(message);
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return { isLoading, error, changePassword };
}

interface UseDeleteAccountReturn {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly deleteAccount: () => Promise<void>;
}

export function useDeleteAccount(): UseDeleteAccountReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const deleteAccount = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      await settingsApi.deleteAccount();
      clearTokens();
      toast.success('Account deleted successfully.');
      navigate('/login');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to delete account';
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, [navigate]);

  return { isLoading, error, deleteAccount };
}

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
