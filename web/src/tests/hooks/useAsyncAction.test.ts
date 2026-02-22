import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useAsyncAction } from '@/hooks/useAsyncAction';

vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

const { toast } = await import('sonner');

beforeEach(() => {
  vi.resetAllMocks();
});

describe('useAsyncAction', () => {
  it('executes action and returns result on success', async () => {
    const action = vi.fn().mockResolvedValue('result');

    const { result } = renderHook(() => useAsyncAction(action));

    let returned: string | undefined;
    await act(async () => {
      returned = await result.current.execute();
    });

    expect(returned, 'execute should return the action result on success.').toBe('result');
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('shows success toast when successMessage is provided', async () => {
    const action = vi.fn().mockResolvedValue('ok');

    const { result } = renderHook(() =>
      useAsyncAction(action, { successMessage: 'Done!' }),
    );

    await act(async () => {
      await result.current.execute();
    });

    expect(toast.success).toHaveBeenCalledWith('Done!');
  });

  it('sets error and shows error toast on failure', async () => {
    const action = vi.fn().mockRejectedValue(new Error('Something failed'));

    const { result } = renderHook(() => useAsyncAction(action));

    const returned = await act(async () => {
      return await result.current.execute();
    });

    expect(
      result.current.error,
      'Error should be set from the thrown Error message.',
    ).toBe('Something failed');
    expect(toast.error).toHaveBeenCalledWith('Something failed');
    expect(returned, 'execute should return undefined on failure.').toBeUndefined();
  });

  it('uses fallbackError for non-Error exceptions', async () => {
    const action = vi.fn().mockRejectedValue('string error');

    const { result } = renderHook(() =>
      useAsyncAction(action, { fallbackError: 'Custom fallback' }),
    );

    await act(async () => {
      await result.current.execute();
    });

    expect(
      result.current.error,
      'Should use fallbackError when thrown value is not an Error instance.',
    ).toBe('Custom fallback');
  });

  it('tracks isLoading state during execution', async () => {
    let resolveAction: (v: string) => void;
    const action = vi.fn().mockReturnValue(
      new Promise((resolve) => { resolveAction = resolve; }),
    );

    const { result } = renderHook(() => useAsyncAction(action));
    expect(result.current.isLoading).toBe(false);

    let promise: Promise<string | undefined>;
    act(() => {
      promise = result.current.execute();
    });

    expect(
      result.current.isLoading,
      'isLoading should be true while action is in progress.',
    ).toBe(true);

    await act(async () => {
      resolveAction!('done');
      await promise!;
    });

    expect(result.current.isLoading).toBe(false);
  });

  it('clears error with clearError', async () => {
    const action = vi.fn().mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useAsyncAction(action));

    await act(async () => {
      await result.current.execute();
    });

    expect(result.current.error).toBe('fail');

    act(() => {
      result.current.clearError();
    });

    expect(result.current.error).toBeNull();
  });
});
