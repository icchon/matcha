import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { setTokens, clearTokens } from '@/api/client';
import { buildWsUrl, connect, disconnect, send, resetConnectionState } from '@/features/websocket/connection';
import type { ConnectionStatus } from '@/features/websocket/types';

// --- MockWebSocket ---
class MockWebSocket {
  static readonly CONNECTING = 0;
  static readonly OPEN = 1;
  static readonly CLOSING = 2;
  static readonly CLOSED = 3;

  readonly CONNECTING = 0;
  readonly OPEN = 1;
  readonly CLOSING = 2;
  readonly CLOSED = 3;

  url: string;
  readyState: number = MockWebSocket.CONNECTING;
  onopen: ((ev: Event) => void) | null = null;
  onclose: ((ev: CloseEvent) => void) | null = null;
  onmessage: ((ev: MessageEvent) => void) | null = null;
  onerror: ((ev: Event) => void) | null = null;
  protocol = '';
  extensions = '';
  bufferedAmount = 0;
  binaryType: BinaryType = 'blob';

  close = vi.fn(() => {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.({ code: 1000, reason: '', wasClean: true } as CloseEvent);
  });

  send = vi.fn();

  simulateOpen(): void {
    this.readyState = MockWebSocket.OPEN;
    this.onopen?.({} as Event);
  }

  simulateError(): void {
    this.onerror?.({} as Event);
  }

  simulateClose(code = 1006, wasClean = false): void {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.({ code, reason: '', wasClean } as CloseEvent);
  }

  addEventListener = vi.fn();
  removeEventListener = vi.fn();
  dispatchEvent = vi.fn(() => true);

  constructor(url: string | URL, _protocols?: string | string[]) {
    this.url = typeof url === 'string' ? url : url.toString();
    mockInstances.push(this);
  }
}

let mockInstances: MockWebSocket[] = [];

function latestMock(): MockWebSocket {
  const instance = mockInstances[mockInstances.length - 1];
  if (!instance) throw new Error('No MockWebSocket instance created');
  return instance;
}

const OriginalWebSocket = globalThis.WebSocket;

beforeEach(() => {
  mockInstances = [];
  (globalThis as unknown as Record<string, unknown>).WebSocket = MockWebSocket as unknown as typeof WebSocket;
  resetConnectionState();
});

afterEach(() => {
  (globalThis as unknown as Record<string, unknown>).WebSocket = OriginalWebSocket;
});

describe('buildWsUrl', () => {
  afterEach(() => {
    clearTokens();
  });

  it('throws when no token is available', () => {
    clearTokens();

    expect(
      () => buildWsUrl('ws://test.example/ws'),
      'buildWsUrl should throw when getAccessToken() returns null. This prevents unauthenticated WS connections.',
    ).toThrow('Cannot connect to WebSocket without an access token');
  });

  it('appends token as query param when token exists', () => {
    setTokens('my-token', 'refresh');

    const url = buildWsUrl('ws://test.example/ws');

    expect(
      url,
      'buildWsUrl should append token= query param to the base URL.',
    ).toBe('ws://test.example/ws?token=my-token');
  });

  it('uses & separator when base URL already has query params', () => {
    setTokens('my-token', 'refresh');

    const url = buildWsUrl('ws://test.example/ws?foo=bar');

    expect(
      url,
      'buildWsUrl should use & when URL already contains ?.',
    ).toBe('ws://test.example/ws?foo=bar&token=my-token');
  });

  it('encodes special characters in the token', () => {
    setTokens('token with spaces&special=chars', 'refresh');

    const url = buildWsUrl('ws://test.example/ws');

    expect(
      url,
      'Token should be URI-encoded in the query param.',
    ).toContain('token=token%20with%20spaces%26special%3Dchars');
  });
});

describe('send', () => {
  const mockSet = vi.fn();

  beforeEach(() => {
    vi.useFakeTimers();
    mockSet.mockClear();
    setTokens('test-token', 'refresh');
  });

  afterEach(() => {
    clearTokens();
    vi.useRealTimers();
  });

  it('sends message when WebSocket is OPEN', () => {
    connect('ws://test.example/ws', mockSet);
    const ws = latestMock();
    ws.simulateOpen();

    send({ type: 'chat.message', payload: { text: 'hello' } });

    expect(
      ws.send,
      'send() should call WebSocket.send when readyState is OPEN.',
    ).toHaveBeenCalledWith(JSON.stringify({ type: 'chat.message', payload: { text: 'hello' } }));
  });

  it('does not send when WebSocket is CONNECTING', () => {
    connect('ws://test.example/ws', mockSet);
    const ws = latestMock();
    // readyState is still CONNECTING

    send({ type: 'chat.message', payload: {} });

    expect(
      ws.send,
      'send() should not call WebSocket.send when readyState is CONNECTING.',
    ).not.toHaveBeenCalled();
  });

  it('does not send when WebSocket is CLOSED', () => {
    connect('ws://test.example/ws', mockSet);
    const ws = latestMock();
    ws.simulateOpen();
    disconnect(mockSet);

    send({ type: 'chat.message', payload: {} });

    expect(
      ws.send,
      'send() should not call WebSocket.send when WebSocket has been closed.',
    ).not.toHaveBeenCalled();
  });

  it('does not send when no connection exists', () => {
    // No connect() called, so ws is null
    expect(
      () => send({ type: 'chat.message', payload: {} }),
      'send() should not throw when no WebSocket connection exists.',
    ).not.toThrow();
  });
});

describe('onerror', () => {
  const mockSet = vi.fn();

  beforeEach(() => {
    vi.useFakeTimers();
    mockSet.mockClear();
    setTokens('test-token', 'refresh');
  });

  afterEach(() => {
    clearTokens();
    vi.useRealTimers();
  });

  it('reports "connecting" when socket is not open', () => {
    connect('ws://test.example/ws', mockSet);
    mockSet.mockClear();

    latestMock().simulateError();

    expect(
      mockSet,
      'onerror should set status to "connecting" when socket.readyState is not OPEN. Verifies fix for stale ws variable.',
    ).toHaveBeenCalledWith({ connectionStatus: 'connecting', error: 'WebSocket connection error' });
  });

  it('reports "connected" when socket is open', () => {
    connect('ws://test.example/ws', mockSet);
    latestMock().simulateOpen();
    mockSet.mockClear();

    // Manually keep readyState as OPEN but fire error
    const ws = latestMock();
    ws.readyState = MockWebSocket.OPEN;
    ws.onerror?.({} as Event);

    expect(
      mockSet,
      'onerror should set status to "connected" when socket.readyState is OPEN.',
    ).toHaveBeenCalledWith({ connectionStatus: 'connected', error: 'WebSocket connection error' });
  });
});
