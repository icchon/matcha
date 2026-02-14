import { create } from 'zustand';
import { getAccessToken } from '@/api/client';

export interface WsMessage {
  readonly type: string;
  readonly payload: unknown;
}

type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';
type MessageHandler = (payload: unknown) => void;

interface WsState {
  readonly connectionStatus: ConnectionStatus;
  readonly error: string | null;
}

interface WsActions {
  readonly connect: (url: string) => void;
  readonly disconnect: () => void;
  readonly send: (message: WsMessage) => void;
  readonly registerHandler: (type: string, handler: MessageHandler) => void;
  readonly unregisterHandler: (type: string) => void;
}

type WsStore = WsState & WsActions;

const INITIAL_BACKOFF_MS = 1000;
const MAX_BACKOFF_MS = 30000;
const BACKOFF_MULTIPLIER = 2;

// Module-level mutable refs (not in Zustand state to avoid serialization issues)
let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let currentBackoff = INITIAL_BACKOFF_MS;
let intentionalDisconnect = false;
let currentUrl: string | null = null;
const handlers = new Map<string, MessageHandler>();

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

function handleMessage(event: MessageEvent): void {
  try {
    const message = JSON.parse(event.data as string) as WsMessage;
    if (typeof message.type !== 'string') return;
    const handler = handlers.get(message.type);
    if (handler) {
      handler(message.payload);
    }
  } catch {
    // Ignore malformed messages
  }
}

function scheduleReconnect(set: (state: Partial<WsState>) => void): void {
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

function createConnection(
  url: string,
  set: (state: Partial<WsState>) => void,
): void {
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
    set({ error: 'WebSocket connection error' });
  };

  socket.onmessage = handleMessage;

  ws = socket;
}

export const useWsStore = create<WsStore>()((set) => ({
  connectionStatus: 'disconnected',
  error: null,

  connect: (url: string) => {
    // Disconnect existing connection first
    if (ws && ws.readyState !== WebSocket.CLOSED) {
      ws.close();
    }
    clearReconnectTimer();

    intentionalDisconnect = false;
    currentUrl = url;
    set({ connectionStatus: 'connecting', error: null });

    createConnection(url, set);
  },

  disconnect: () => {
    intentionalDisconnect = true;
    clearReconnectTimer();
    currentBackoff = INITIAL_BACKOFF_MS;
    currentUrl = null;
    handlers.clear();

    if (ws) {
      ws.close();
      ws = null;
    }

    set({ connectionStatus: 'disconnected' });
  },

  send: (message: WsMessage) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    }
  },

  registerHandler: (type: string, handler: MessageHandler) => {
    handlers.set(type, handler);
  },

  unregisterHandler: (type: string) => {
    handlers.delete(type);
  },
}));
