import { describe, it, expect, beforeEach, vi } from 'vitest';
import {
  registerHandler,
  unregisterHandler,
  dispatchMessage,
  clearAllHandlers,
  getHandlerCount,
} from '@/features/websocket/messageRouter';

const validChatPayload = {
  id: 'msg-1',
  senderId: 'user-a',
  receiverId: 'user-b',
  content: 'hi',
  timestamp: '2026-01-01T00:00:00Z',
};

const validNotificationPayload = {
  id: 'notif-1',
  type: 'like',
  message: 'Someone liked you',
  timestamp: '2026-01-01T00:00:00Z',
  read: false,
};

function makeEvent(data: unknown): MessageEvent {
  return { data: JSON.stringify(data) } as MessageEvent;
}

describe('messageRouter', () => {
  beforeEach(() => {
    clearAllHandlers();
  });

  it('dispatches message to registered handler', () => {
    const handler = vi.fn();
    registerHandler('chat.message', handler);

    dispatchMessage(makeEvent({ type: 'chat.message', payload: validChatPayload }));

    expect(
      handler,
      'Handler should be called with the message payload.',
    ).toHaveBeenCalledWith(validChatPayload);
  });

  it('supports multiple handlers for the same type', () => {
    const h1 = vi.fn();
    const h2 = vi.fn();
    registerHandler('chat.message', h1);
    registerHandler('chat.message', h2);

    dispatchMessage(makeEvent({ type: 'chat.message', payload: validChatPayload }));

    expect(h1, 'First handler should be called.').toHaveBeenCalledWith(validChatPayload);
    expect(h2, 'Second handler should also be called.').toHaveBeenCalledWith(validChatPayload);
    expect(getHandlerCount('chat.message'), 'Should have 2 handlers registered.').toBe(2);
  });

  it('unregisterHandler removes only the specified handler', () => {
    const h1 = vi.fn();
    const h2 = vi.fn();
    registerHandler('chat.message', h1);
    registerHandler('chat.message', h2);

    unregisterHandler('chat.message', h1);

    dispatchMessage(makeEvent({ type: 'chat.message', payload: validChatPayload }));

    expect(h1, 'Unregistered handler should not be called.').not.toHaveBeenCalled();
    expect(h2, 'Remaining handler should still be called.').toHaveBeenCalled();
    expect(getHandlerCount('chat.message'), 'Should have 1 handler remaining.').toBe(1);
  });

  it('unregisterHandler cleans up empty Set', () => {
    const h = vi.fn();
    registerHandler('chat.message', h);
    unregisterHandler('chat.message', h);

    expect(getHandlerCount('chat.message'), 'Empty type should report 0 handlers.').toBe(0);
  });

  it('unregisterHandler is safe for unknown type', () => {
    expect(() => unregisterHandler('nonexistent', vi.fn())).not.toThrow();
  });

  it('clearAllHandlers removes everything', () => {
    registerHandler('chat.message', vi.fn());
    registerHandler('notification', vi.fn());

    clearAllHandlers();

    expect(getHandlerCount('chat.message'), 'After clear, chat.message should have 0 handlers.').toBe(0);
    expect(getHandlerCount('notification'), 'After clear, notification should have 0 handlers.').toBe(0);
  });

  it('ignores malformed JSON gracefully', () => {
    registerHandler('chat.message', vi.fn());

    expect(() => {
      dispatchMessage({ data: 'not-json' } as MessageEvent);
    }).not.toThrow();
  });

  it('ignores messages without a string type field', () => {
    const handler = vi.fn();
    registerHandler('chat.message', handler);

    dispatchMessage(makeEvent({ payload: {} }));

    expect(handler, 'Handler should not be called for message without type.').not.toHaveBeenCalled();
  });

  it('ignores messages with disallowed type', () => {
    const handler = vi.fn();
    registerHandler('chat.message', handler);

    dispatchMessage(makeEvent({ type: 'unknown', payload: {} }));

    expect(handler, 'Handler should not be called for disallowed type. z.enum should reject it.').not.toHaveBeenCalled();
  });

  it('drops chat.message with invalid payload', () => {
    const handler = vi.fn();
    registerHandler('chat.message', handler);

    // Missing required fields (senderId, receiverId, etc.)
    dispatchMessage(makeEvent({ type: 'chat.message', payload: { text: 'hi' } }));

    expect(
      handler,
      'Handler should NOT be called when chat.message payload fails Zod validation. Payload must have id, senderId, receiverId, content, timestamp.',
    ).not.toHaveBeenCalled();
  });

  it('drops chat.ack with invalid payload', () => {
    const handler = vi.fn();
    registerHandler('chat.ack', handler);

    dispatchMessage(makeEvent({ type: 'chat.ack', payload: { wrong: 'shape' } }));

    expect(
      handler,
      'Handler should NOT be called when chat.ack payload fails Zod validation. Payload must have messageId, status.',
    ).not.toHaveBeenCalled();
  });

  it('dispatches valid chat.ack payload', () => {
    const handler = vi.fn();
    registerHandler('chat.ack', handler);
    const ackPayload = { messageId: 'msg-1', status: 'delivered' };

    dispatchMessage(makeEvent({ type: 'chat.ack', payload: ackPayload }));

    expect(handler, 'Handler should be called for valid chat.ack payload.').toHaveBeenCalledWith(ackPayload);
  });

  it('drops chat.read with invalid payload', () => {
    const handler = vi.fn();
    registerHandler('chat.read', handler);

    dispatchMessage(makeEvent({ type: 'chat.read', payload: null }));

    expect(
      handler,
      'Handler should NOT be called when chat.read payload fails Zod validation. Payload must have conversationId, readAt.',
    ).not.toHaveBeenCalled();
  });

  it('dispatches valid chat.read payload', () => {
    const handler = vi.fn();
    registerHandler('chat.read', handler);
    const readPayload = { conversationId: 'user-a', readAt: '2026-01-01T00:00:00Z' };

    dispatchMessage(makeEvent({ type: 'chat.read', payload: readPayload }));

    expect(handler, 'Handler should be called for valid chat.read payload.').toHaveBeenCalledWith(readPayload);
  });

  it('drops notification with invalid payload', () => {
    const handler = vi.fn();
    registerHandler('notification', handler);

    dispatchMessage(makeEvent({ type: 'notification', payload: { id: 'n1' } }));

    expect(
      handler,
      'Handler should NOT be called when notification payload fails Zod validation. Payload must have id, type, message, timestamp, read.',
    ).not.toHaveBeenCalled();
  });

  it('dispatches valid notification payload', () => {
    const handler = vi.fn();
    registerHandler('notification', handler);

    dispatchMessage(makeEvent({ type: 'notification', payload: validNotificationPayload }));

    expect(handler, 'Handler should be called for valid notification payload.').toHaveBeenCalledWith(validNotificationPayload);
  });
});
