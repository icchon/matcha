import { getAccessToken } from '@/api/client';
import { dispatchMessage } from './messageRouter';
import type { ConnectionStatus } from './types';

const INITIAL_BACKOFF_MS = 1000;
const MAX_BACKOFF_MS = 30000;
const BACKOFF_MULTIPLIER = 2;

let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let currentBackoff = INITIAL_BACKOFF_MS;
let intentionalDisconnect = false;
let currentUrl: string | null = null;

type StatusSetter = (state: { readonly connectionStatus: ConnectionStatus; readonly error?: string | null }) => void;

function clearReconnectTimer(): void {
  if (reconnectTimer !== null) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
}

// Browser WebSocket API does not support custom headers.
// Token is passed as a query param. Server-side (nginx) must:
// 1. Strip token param from access logs to prevent credential leakage
// 2. Convert query param to Authorization header before forwarding to WS gateway
// Future: use short-lived WS tickets instead of the access token itself
function buildWsUrl(baseUrl: string): string {
  const token = getAccessToken();
  if (!token) return baseUrl;

  const separator = baseUrl.includes('?') ? '&' : '?';
  return `${baseUrl}${separator}token=${encodeURIComponent(token)}`;
}

function scheduleReconnect(set: StatusSetter): void {
  clearReconnectTimer();
  set({ connectionStatus: 'reconnecting' });

  reconnectTimer = setTimeout(() => {
    reconnectTimer = null;
    if (currentUrl) {
      createConnection(currentUrl, set);
    }
  }, currentBackoff);

  currentBackoff = Math.min(currentBackoff * BACKOFF_MULTIPLIER, MAX_BACKOFF_MS);
}

function createConnection(url: string, set: StatusSetter): void {
  const wsUrl = buildWsUrl(url);
  const socket = new WebSocket(wsUrl);

  socket.onopen = () => {
    currentBackoff = INITIAL_BACKOFF_MS;
    set({ connectionStatus: 'connected', error: null });
  };

  socket.onclose = (event: CloseEvent) => {
    if (intentionalDisconnect || event.code === 1000) {
      set({ connectionStatus: 'disconnected' });
      return;
    }
    scheduleReconnect(set);
  };

  socket.onerror = () => {
    set({ connectionStatus: ws ? 'connected' : 'connecting', error: 'WebSocket connection error' });
  };

  socket.onmessage = dispatchMessage;

  ws = socket;
}

export function connect(url: string, set: StatusSetter): void {
  if (ws && ws.readyState !== WebSocket.CLOSED) {
    ws.close();
  }
  clearReconnectTimer();

  intentionalDisconnect = false;
  currentUrl = url;
  set({ connectionStatus: 'connecting', error: null });

  createConnection(url, set);
}

export function disconnect(set: StatusSetter): void {
  intentionalDisconnect = true;
  clearReconnectTimer();
  currentBackoff = INITIAL_BACKOFF_MS;
  currentUrl = null;

  if (ws) {
    ws.close();
    ws = null;
  }

  set({ connectionStatus: 'disconnected' });
}

export function send(message: { readonly type: string; readonly payload: unknown }): void {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(message));
  }
}
