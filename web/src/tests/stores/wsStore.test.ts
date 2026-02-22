import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { setTokens, clearTokens } from '@/api/client';
import { clearAllHandlers } from '@/features/websocket/messageRouter';
import { useWsStore } from '@/stores/wsStore';
import type { WsMessage } from '@/stores/wsStore';

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

  simulateMessage(data: unknown): void {
    this.onmessage?.({ data: JSON.stringify(data) } as MessageEvent);
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
});

afterEach(() => {
  (globalThis as unknown as Record<string, unknown>).WebSocket = OriginalWebSocket;
});

describe('wsStore', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    clearAllHandlers();
    useWsStore.getState().disconnect();
    useWsStore.setState({
      connectionStatus: 'disconnected',
      error: null,
    });
    mockInstances = [];
    // Set a default token so buildWsUrl does not throw
    setTokens('test-access-token', 'test-refresh-token');
  });

  afterEach(() => {
    clearTokens();
    vi.useRealTimers();
  });

  // --- Connection lifecycle ---

  describe('connect', () => {
    it('sets status to connecting then connected on open', () => {
      useWsStore.getState().connect('ws://test.example/ws');

      expect(
        useWsStore.getState().connectionStatus,
        'After connect() call, status should be "connecting". Check that connect sets status before creating WebSocket.',
      ).toBe('connecting');

      latestMock().simulateOpen();

      expect(
        useWsStore.getState().connectionStatus,
        'After WebSocket open event, status should be "connected". Check onopen handler.',
      ).toBe('connected');
    });

    it('creates WebSocket with token as query param', () => {
      setTokens('test-token-123', 'refresh-token');

      useWsStore.getState().connect('ws://test.example/ws');
      const ws = latestMock();

      expect(
        ws.url,
        'WebSocket URL should include token query param for auth. Check buildWsUrl() URL construction.',
      ).toContain('token=test-token-123');
    });

    it('clears previous error on new connect', () => {
      useWsStore.setState({ error: 'previous error' });

      useWsStore.getState().connect('ws://test.example/ws');

      expect(
        useWsStore.getState().error,
        'Error should be cleared when starting a new connection. Check connect() clears error.',
      ).toBeNull();
    });

    it('disconnects existing connection before creating new one', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      const firstWs = latestMock();
      firstWs.simulateOpen();

      useWsStore.getState().connect('ws://test.example/ws2');

      expect(
        firstWs.close,
        'Previous WebSocket should be closed when connecting again. Check connect() calls disconnect first.',
      ).toHaveBeenCalled();
    });
  });

  describe('disconnect', () => {
    it('closes WebSocket and sets status to disconnected', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      const ws = latestMock();
      ws.simulateOpen();

      useWsStore.getState().disconnect();

      expect(
        ws.close,
        'disconnect() should call WebSocket.close(). Check disconnect implementation.',
      ).toHaveBeenCalled();
      expect(
        useWsStore.getState().connectionStatus,
        'After disconnect(), status should be "disconnected".',
      ).toBe('disconnected');
    });

    it('is safe to call when not connected', () => {
      expect(() => useWsStore.getState().disconnect()).not.toThrow();
    });

    it('does NOT clear handlers (handlers persist across reconnections)', () => {
      const handler = vi.fn();
      useWsStore.getState().registerHandler('chat.message', handler);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      useWsStore.getState().disconnect();

      // Re-connect and send message — handler SHOULD still fire
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();
      latestMock().simulateMessage({ type: 'chat.message', payload: { id: '1', senderId: 'a', receiverId: 'b', content: 'hello', timestamp: '2026-01-01T00:00:00Z' } });

      expect(
        handler,
        'Handlers should persist after disconnect(). Only clearAllHandlers() removes them.',
      ).toHaveBeenCalledWith({ id: '1', senderId: 'a', receiverId: 'b', content: 'hello', timestamp: '2026-01-01T00:00:00Z' });

      // Cleanup
      clearAllHandlers();
    });
  });

  // --- Auto-reconnect ---

  describe('auto-reconnect', () => {
    it('reconnects with exponential backoff on unclean close', () => {
      // Stub Math.random to remove jitter for deterministic timing
      vi.spyOn(Math, 'random').mockReturnValue(0);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      expect(
        useWsStore.getState().connectionStatus,
        'Should be connected before testing reconnect.',
      ).toBe('connected');

      latestMock().simulateClose(1006, false);

      expect(
        useWsStore.getState().connectionStatus,
        'After unclean close, status should be "reconnecting". Check onclose handler.',
      ).toBe('reconnecting');

      vi.advanceTimersByTime(1000);
      expect(
        mockInstances.length,
        'After 1s, should attempt first reconnect (2 total WS instances). Check backoff logic.',
      ).toBe(2);

      latestMock().simulateClose(1006, false);
      vi.advanceTimersByTime(2000);
      expect(
        mockInstances.length,
        'After 2s backoff, should attempt second reconnect (3 total). Check exponential backoff.',
      ).toBe(3);

      latestMock().simulateClose(1006, false);
      vi.advanceTimersByTime(4000);
      expect(
        mockInstances.length,
        'After 4s backoff, should attempt third reconnect (4 total). Check exponential backoff.',
      ).toBe(4);

      vi.spyOn(Math, 'random').mockRestore();
    });

    it('caps backoff at 30 seconds', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();
      latestMock().simulateClose(1006, false);

      for (let i = 0; i < 6; i++) {
        vi.advanceTimersByTime(30000);
        latestMock().simulateClose(1006, false);
      }

      const instancesBefore = mockInstances.length;
      vi.advanceTimersByTime(29999);
      expect(
        mockInstances.length,
        'Should NOT reconnect before 30s cap. Check max backoff constant.',
      ).toBe(instancesBefore);

      vi.advanceTimersByTime(1);
      expect(
        mockInstances.length,
        'Should reconnect at exactly 30s cap. Check max backoff constant.',
      ).toBe(instancesBefore + 1);

      vi.spyOn(Math, 'random').mockRestore();
    });

    it('does not reconnect on clean close (code 1000)', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      latestMock().simulateClose(1000, true);

      expect(
        useWsStore.getState().connectionStatus,
        'After clean close (1000), status should be "disconnected", not "reconnecting".',
      ).toBe('disconnected');

      vi.advanceTimersByTime(5000);
      expect(
        mockInstances.length,
        'No reconnect should happen after clean close. Check wasClean / code handling.',
      ).toBe(1);
    });

    it('resets backoff on successful reconnect', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();
      latestMock().simulateClose(1006, false);

      vi.advanceTimersByTime(1000);
      latestMock().simulateOpen();

      expect(
        useWsStore.getState().connectionStatus,
        'After successful reconnect, status should be "connected".',
      ).toBe('connected');

      latestMock().simulateClose(1006, false);
      const instancesBefore = mockInstances.length;

      vi.advanceTimersByTime(1000);
      expect(
        mockInstances.length,
        'After successful reconnect + disconnect, backoff should reset to 1s. Check backoff reset in onopen.',
      ).toBe(instancesBefore + 1);

      vi.spyOn(Math, 'random').mockRestore();
    });

    it('does not reconnect after manual disconnect', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      useWsStore.getState().disconnect();

      vi.advanceTimersByTime(5000);
      expect(
        mockInstances.length,
        'After manual disconnect(), no reconnect should occur. Check intentional disconnect flag.',
      ).toBe(1);
    });

    it('stops reconnecting after MAX_RECONNECT_ATTEMPTS', () => {
      vi.spyOn(Math, 'random').mockReturnValue(0);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      // Exhaust all 10 reconnect attempts
      for (let i = 0; i < 10; i++) {
        latestMock().simulateClose(1006, false);
        vi.advanceTimersByTime(30000);
      }

      // 11th close should not schedule a reconnect
      latestMock().simulateClose(1006, false);
      const instancesBefore = mockInstances.length;
      vi.advanceTimersByTime(60000);

      expect(
        mockInstances.length,
        'Should stop reconnecting after MAX_RECONNECT_ATTEMPTS (10). Check reconnect counter.',
      ).toBe(instancesBefore);

      expect(
        useWsStore.getState().connectionStatus,
        'After exhausting reconnect attempts, status should be "disconnected".',
      ).toBe('disconnected');

      vi.spyOn(Math, 'random').mockRestore();
    });

    it('does not reconnect on terminal close codes (4401, 4403)', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      latestMock().simulateClose(4401, false);

      expect(
        useWsStore.getState().connectionStatus,
        'After terminal close code 4401, status should be "disconnected".',
      ).toBe('disconnected');

      vi.advanceTimersByTime(5000);
      expect(
        mockInstances.length,
        'No reconnect should happen after terminal close code. Check TERMINAL_CLOSE_CODES.',
      ).toBe(1);
    });
  });

  // --- Error handling ---

  describe('error handling', () => {
    it('sets error state on WebSocket error', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateError();

      expect(
        useWsStore.getState().error,
        'After WebSocket error event, error state should be set. Check onerror handler.',
      ).toBeTruthy();
    });
  });

  // --- Message send ---

  describe('send', () => {
    it('sends JSON-stringified message when connected', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      const ws = latestMock();
      ws.simulateOpen();

      const message: WsMessage = { type: 'chat.message', payload: { id: '1', senderId: 'a', receiverId: 'b', content: 'hello', timestamp: '' } };
      useWsStore.getState().send(message);

      expect(
        ws.send,
        'send() should call WebSocket.send with JSON string. Check send implementation.',
      ).toHaveBeenCalledWith(JSON.stringify(message));
    });

    it('does not send when not connected', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      const ws = latestMock();

      const message: WsMessage = { type: 'chat.message', payload: { id: '1', senderId: 'a', receiverId: 'b', content: '', timestamp: '' } };
      useWsStore.getState().send(message);

      expect(
        ws.send,
        'send() should not call WebSocket.send when not in OPEN state.',
      ).not.toHaveBeenCalled();
    });
  });

  // --- Message routing ---

  describe('message routing', () => {
    it('routes incoming messages to registered handlers by type', () => {
      const chatHandler = vi.fn();
      useWsStore.getState().registerHandler('chat.message', chatHandler);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      const payload = { id: '1', senderId: 'a', receiverId: 'b', content: 'hello', timestamp: '2026-01-01T00:00:00Z' };
      latestMock().simulateMessage({ type: 'chat.message', payload });

      expect(
        chatHandler,
        'Registered handler for "chat.message" should be called with payload. Check message routing in onmessage.',
      ).toHaveBeenCalledWith(payload);
    });

    it('supports multiple handlers for the SAME message type', () => {
      const handler1 = vi.fn();
      const handler2 = vi.fn();

      useWsStore.getState().registerHandler('chat.message', handler1);
      useWsStore.getState().registerHandler('chat.message', handler2);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      const payload = { id: '1', senderId: 'a', receiverId: 'b', content: 'hi', timestamp: '2026-01-01T00:00:00Z' };
      latestMock().simulateMessage({ type: 'chat.message', payload });

      expect(
        handler1,
        'First handler for "chat.message" should be called. Multiple handlers per type must be supported.',
      ).toHaveBeenCalledWith(payload);
      expect(
        handler2,
        'Second handler for "chat.message" should also be called. Handlers must not overwrite each other.',
      ).toHaveBeenCalledWith(payload);
    });

    it('supports multiple handlers for different message types', () => {
      const chatHandler = vi.fn();
      const notifHandler = vi.fn();

      useWsStore.getState().registerHandler('chat.message', chatHandler);
      useWsStore.getState().registerHandler('notification', notifHandler);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      const notifPayload = { id: 'n1', type: 'like', message: 'Someone liked you', timestamp: '2026-01-01T00:00:00Z', read: false };
      latestMock().simulateMessage({ type: 'notification', payload: notifPayload });

      expect(
        chatHandler,
        'chat.message handler should NOT be called for a "notification" message.',
      ).not.toHaveBeenCalled();
      expect(
        notifHandler,
        'notification handler should be called with correct payload.',
      ).toHaveBeenCalledWith(notifPayload);
    });

    it('unregisterHandler removes only the specified handler', () => {
      const handler1 = vi.fn();
      const handler2 = vi.fn();
      useWsStore.getState().registerHandler('chat.message', handler1);
      useWsStore.getState().registerHandler('chat.message', handler2);
      useWsStore.getState().unregisterHandler('chat.message', handler1);

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      const payload = { id: '1', senderId: 'a', receiverId: 'b', content: 'test', timestamp: '2026-01-01T00:00:00Z' };
      latestMock().simulateMessage({ type: 'chat.message', payload });

      expect(
        handler1,
        'After unregisterHandler(handler1), handler1 should not be called.',
      ).not.toHaveBeenCalled();
      expect(
        handler2,
        'handler2 should still be called after only handler1 was unregistered.',
      ).toHaveBeenCalled();
    });

    it('ignores messages with unknown types gracefully', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      expect(() => {
        latestMock().simulateMessage({ type: 'unknown_type', payload: {} });
      }).not.toThrow();
    });

    it('ignores malformed (non-JSON) messages gracefully', () => {
      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();

      const ws = latestMock();
      expect(() => {
        ws.onmessage?.({ data: 'not-json' } as MessageEvent);
      }).not.toThrow();
    });

    it('clearAllHandlers removes all handlers', () => {
      const handler = vi.fn();
      useWsStore.getState().registerHandler('chat.message', handler);

      useWsStore.getState().clearAllHandlers();

      useWsStore.getState().connect('ws://test.example/ws');
      latestMock().simulateOpen();
      latestMock().simulateMessage({ type: 'chat.message', payload: {} });

      expect(
        handler,
        'After clearAllHandlers(), no handlers should be called.',
      ).not.toHaveBeenCalled();
    });
  });
});
