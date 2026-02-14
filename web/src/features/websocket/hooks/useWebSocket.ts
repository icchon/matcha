import { useEffect, useRef } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { useWsStore } from '@/stores/wsStore';

const WS_URL = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`;

interface UseWebSocketReturn {
  readonly connectionStatus: string;
  readonly error: string | null;
}

export function useWebSocket(): UseWebSocketReturn {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const connectionStatus = useWsStore((s) => s.connectionStatus);
  const error = useWsStore((s) => s.error);
  const connect = useWsStore((s) => s.connect);
  const disconnect = useWsStore((s) => s.disconnect);

  const prevAuthRef = useRef(isAuthenticated);

  useEffect(() => {
    if (isAuthenticated && !prevAuthRef.current) {
      connect(WS_URL);
    } else if (!isAuthenticated && prevAuthRef.current) {
      disconnect();
    }
    prevAuthRef.current = isAuthenticated;
  }, [isAuthenticated, connect, disconnect]);

  // Also connect on mount if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      connect(WS_URL);
    }

    return () => {
      disconnect();
    };
    // Only run on mount/unmount
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { connectionStatus, error };
}
