import { useCallback } from 'react';
import { useAsyncAction } from '@/hooks/useAsyncAction';
import * as settingsApi from '@/api/settings';
import type { ChangePasswordRequest } from '@/api/settings';

interface UseChangePasswordReturn {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly changePassword: (params: ChangePasswordRequest) => Promise<boolean>;
}

export function useChangePassword(): UseChangePasswordReturn {
  const { isLoading, error, execute } = useAsyncAction(settingsApi.changePassword, {
    successMessage: 'Password changed successfully!',
    fallbackError: 'Failed to change password',
  });

  const changePassword = useCallback(
    async (params: ChangePasswordRequest): Promise<boolean> => {
      const result = await execute(params);
      return result !== undefined;
    },
    [execute],
  );

  return { isLoading, error, changePassword };
}
