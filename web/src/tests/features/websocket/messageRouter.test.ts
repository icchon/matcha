import { describe, it, expect, beforeEach, vi } from 'vitest';
import {
  registerHandler,
  unregisterHandler,
  dispatchMessage,
  clearAllHandlers,
  getHandlerCount,
} from '@/features/websocket/messageRouter';

describe('messageRouter', () => {
  beforeEach(() => {
    clearAllHandlers();
  });

  it('dispatches message to registered handler', () => {
    const handler = vi.fn();
    registerHandler('chat.message', handler);

    dispatchMessage({ data: JSON.stringify({ type: 'chat.message', payload: { text: 'hi' } }) } as MessageEvent);

    expect(
      handler,
      'Handler should be called with the message payload.',
    ).toHaveBeenCalledWith({ text: 'hi' });
  });

  it('supports multiple handlers for the same type', () => {
    const h1 = vi.fn();
    const h2 = vi.fn();
    registerHandler('chat.message', h1);
    registerHandler('chat.message', h2);

    dispatchMessage({ data: JSON.stringify({ type: 'chat.message', payload: 'x' }) } as MessageEvent);

    expect(h1, 'First handler should be called.').toHaveBeenCalledWith('x');
    expect(h2, 'Second handler should also be called.').toHaveBeenCalledWith('x');
    expect(getHandlerCount('chat.message'), 'Should have 2 handlers registered.').toBe(2);
  });

  it('unregisterHandler removes only the specified handler', () => {
    const h1 = vi.fn();
    const h2 = vi.fn();
    registerHandler('chat.message', h1);
    registerHandler('chat.message', h2);

    unregisterHandler('chat.message', h1);

    dispatchMessage({ data: JSON.stringify({ type: 'chat.message', payload: null }) } as MessageEvent);

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
    registerHandler('a', vi.fn());
    registerHandler('b', vi.fn());

    clearAllHandlers();

    expect(getHandlerCount('a'), 'After clear, type a should have 0 handlers.').toBe(0);
    expect(getHandlerCount('b'), 'After clear, type b should have 0 handlers.').toBe(0);
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

    dispatchMessage({ data: JSON.stringify({ payload: {} }) } as MessageEvent);

    expect(handler, 'Handler should not be called for message without type.').not.toHaveBeenCalled();
  });

  it('ignores messages with unregistered type', () => {
    const handler = vi.fn();
    registerHandler('chat.message', handler);

    dispatchMessage({ data: JSON.stringify({ type: 'unknown', payload: {} }) } as MessageEvent);

    expect(handler, 'Handler should not be called for unregistered type.').not.toHaveBeenCalled();
  });
});
