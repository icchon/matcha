import { useEffect } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { useWsStore } from '@/stores/wsStore';

function getWsUrl(): string {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${protocol}//${window.location.host}/ws`;
}

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

  useEffect(() => {
    if (isAuthenticated) {
      connect(getWsUrl());
    } else {
      disconnect();
    }

    return () => {
      disconnect();
    };
  }, [isAuthenticated, connect, disconnect]);

  return { connectionStatus, error };
}
