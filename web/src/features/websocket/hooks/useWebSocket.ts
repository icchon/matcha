import { useEffect } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { useWsStore } from '@/stores/wsStore';

function isLocalhost(hostname: string): boolean {
  return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '[::1]';
}

export function getWsUrl(): string {
  const { protocol, host, hostname } = window.location;
  const wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:';

  if (wsProtocol === 'ws:' && !isLocalhost(hostname)) {
    throw new Error('Insecure WebSocket (ws://) is not allowed for non-localhost hosts');
  }

  return `${wsProtocol}//${host}/ws`;
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
    }

    return () => {
      disconnect();
    };
  }, [isAuthenticated, connect, disconnect]);

  return { connectionStatus, error };
}
