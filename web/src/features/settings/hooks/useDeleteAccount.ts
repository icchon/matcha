import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import * as settingsApi from '@/api/settings';
import { clearTokens } from '@/api/client';

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
