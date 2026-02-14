import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useWebSocket } from '@/features/websocket/hooks/useWebSocket';
import { useAuthStore } from '@/stores/authStore';
import { useWsStore } from '@/stores/wsStore';

// Mock wsStore actions
const mockConnect = vi.fn();
const mockDisconnect = vi.fn();

vi.mock('@/stores/wsStore', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/stores/wsStore')>();
  return {
    ...actual,
    useWsStore: vi.fn((selector?: (state: ReturnType<typeof actual.useWsStore.getState>) => unknown) => {
      const state = {
        connectionStatus: 'disconnected' as const,
        error: null,
        connect: mockConnect,
        disconnect: mockDisconnect,
        send: vi.fn(),
        registerHandler: vi.fn(),
        unregisterHandler: vi.fn(),
      };
      return selector ? selector(state) : state;
    }),
  };
});

describe('useWebSocket', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.setState({
      userId: null,
      isAuthenticated: false,
      isVerified: false,
      authMethod: null,
    });
  });

  it('returns connection status from wsStore', () => {
    const { result } = renderHook(() => useWebSocket());

    expect(
      result.current.connectionStatus,
      'useWebSocket should expose connectionStatus from wsStore.',
    ).toBe('disconnected');
  });

  it('connects when user becomes authenticated', () => {
    const { rerender } = renderHook(() => useWebSocket());

    act(() => {
      useAuthStore.setState({ isAuthenticated: true });
    });
    rerender();

    // The hook should observe auth state change and connect
    // Since we're mocking at the module level, we verify connect was called
    expect(
      mockConnect,
      'useWebSocket should call connect() when isAuthenticated becomes true.',
    ).toHaveBeenCalled();
  });

  it('disconnects on unmount', () => {
    act(() => {
      useAuthStore.setState({ isAuthenticated: true });
    });

    const { unmount } = renderHook(() => useWebSocket());
    unmount();

    expect(
      mockDisconnect,
      'useWebSocket should call disconnect() on unmount for cleanup.',
    ).toHaveBeenCalled();
  });

  it('disconnects when user logs out', () => {
    act(() => {
      useAuthStore.setState({ isAuthenticated: true });
    });

    const { rerender } = renderHook(() => useWebSocket());

    act(() => {
      useAuthStore.setState({ isAuthenticated: false });
    });
    rerender();

    expect(
      mockDisconnect,
      'useWebSocket should call disconnect() when isAuthenticated becomes false (logout).',
    ).toHaveBeenCalled();
  });
});
