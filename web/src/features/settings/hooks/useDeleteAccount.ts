import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import * as settingsApi from '@/api/settings';
import { useAuthStore } from '@/stores/authStore';
import { ApiClientError } from '@/api/client';
import type { DeleteAccountRequest } from '@/api/settings';

function mapDeleteErrorMessage(err: unknown): string {
  if (err instanceof ApiClientError) {
    switch (err.status) {
      case 401:
        return 'Incorrect password. Please try again.';
      case 403:
        return 'You do not have permission to perform this action.';
      default:
        return 'Failed to delete account. Please try again later.';
    }
  }
  return 'Failed to delete account';
}

interface UseDeleteAccountReturn {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly deleteAccount: (params: DeleteAccountRequest) => Promise<void>;
}

export function useDeleteAccount(): UseDeleteAccountReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const logout = useAuthStore((s) => s.logout);

  const deleteAccount = useCallback(async (params: DeleteAccountRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      await settingsApi.deleteAccount(params);
      logout();
      toast.success('Account deleted successfully.');
      navigate('/login');
    } catch (err) {
      const message = mapDeleteErrorMessage(err);
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, [navigate, logout]);

  return { isLoading, error, deleteAccount };
}
