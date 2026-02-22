import { useState, useCallback, useRef } from 'react';
import { toast } from 'sonner';

interface UseAsyncActionOptions {
  readonly successMessage?: string;
  readonly fallbackError?: string;
}

interface UseAsyncActionReturn<T extends unknown[], R> {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly execute: (...args: T) => Promise<R | undefined>;
  readonly clearError: () => void;
}

export function useAsyncAction<T extends unknown[], R>(
  action: (...args: T) => Promise<R>,
  options: UseAsyncActionOptions = {},
): UseAsyncActionReturn<T, R> {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const optionsRef = useRef(options);
  optionsRef.current = options;

  const clearError = useCallback(() => setError(null), []);

  const execute = useCallback(
    async (...args: T): Promise<R | undefined> => {
      setIsLoading(true);
      setError(null);
      try {
        const result = await action(...args);
        if (optionsRef.current.successMessage) {
          toast.success(optionsRef.current.successMessage);
        }
        return result;
      } catch (err) {
        const message =
          err instanceof Error ? err.message : (optionsRef.current.fallbackError ?? 'An error occurred');
        setError(message);
        toast.error(message);
        return undefined;
      } finally {
        setIsLoading(false);
      }
    },
    [action],
  );

  return { isLoading, error, execute, clearError };
}
