import type { WsMessage, MessageHandler } from './types';

const handlers = new Map<string, Set<MessageHandler>>();

export function registerHandler(type: string, handler: MessageHandler): void {
  const existing = handlers.get(type);
  if (existing) {
    const next = new Set(existing);
    next.add(handler);
    handlers.set(type, next);
  } else {
    handlers.set(type, new Set([handler]));
  }
}

export function unregisterHandler(type: string, handler: MessageHandler): void {
  const existing = handlers.get(type);
  if (!existing) return;
  const next = new Set(existing);
  next.delete(handler);
  if (next.size === 0) {
    handlers.delete(type);
  } else {
    handlers.set(type, next);
  }
}

export function dispatchMessage(event: MessageEvent): void {
  try {
    const message = JSON.parse(event.data as string) as WsMessage;
    if (typeof message.type !== 'string') return;
    const typeHandlers = handlers.get(message.type);
    if (typeHandlers) {
      for (const handler of typeHandlers) {
        handler(message.payload);
      }
    }
  } catch {
    // Ignore malformed messages
  }
}

export function clearAllHandlers(): void {
  handlers.clear();
}

export function getHandlerCount(type: string): number {
  return handlers.get(type)?.size ?? 0;
}
