import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAsyncAction } from '@/hooks/useAsyncAction';
import * as settingsApi from '@/api/settings';
import type { ChangePasswordRequest } from '@/api/settings';
import { useAuthStore } from '@/stores/authStore';

interface UseChangePasswordReturn {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly changePassword: (params: ChangePasswordRequest) => Promise<boolean>;
}

export function useChangePassword(): UseChangePasswordReturn {
  const navigate = useNavigate();
  const logout = useAuthStore((s) => s.logout);
  const { isLoading, error, execute } = useAsyncAction(settingsApi.changePassword, {
    successMessage: 'Password changed successfully! Please log in again.',
    fallbackError: 'Failed to change password',
  });

  const changePassword = useCallback(
    async (params: ChangePasswordRequest): Promise<boolean> => {
      const result = await execute(params);
      if (result !== undefined) {
        logout();
        navigate('/login');
        return true;
      }
      return false;
    },
    [execute, logout, navigate],
  );

  return { isLoading, error, changePassword };
}
