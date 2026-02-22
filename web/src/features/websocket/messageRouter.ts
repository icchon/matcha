import { z } from 'zod';
import type { MessageHandler } from './types';

const chatMessagePayloadSchema = z.object({
  id: z.string(),
  senderId: z.string(),
  receiverId: z.string(),
  content: z.string(),
  timestamp: z.string(),
});

const chatAckPayloadSchema = z.object({
  messageId: z.string(),
  status: z.string(),
});

const chatReadPayloadSchema = z.object({
  conversationId: z.string(),
  readAt: z.string(),
});

const notificationPayloadSchema = z.object({
  id: z.string(),
  type: z.string(),
  message: z.string(),
  timestamp: z.string(),
  read: z.boolean(),
});

const ALLOWED_TYPES = ['chat.message', 'chat.ack', 'chat.read', 'notification'] as const;

const payloadSchemas: Record<string, z.ZodSchema> = {
  'chat.message': chatMessagePayloadSchema,
  'chat.ack': chatAckPayloadSchema,
  'chat.read': chatReadPayloadSchema,
  'notification': notificationPayloadSchema,
};

const wsMessageSchema = z.object({
  type: z.enum(ALLOWED_TYPES),
  payload: z.unknown(),
});

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
    const raw: unknown = JSON.parse(event.data as string);
    const parsed = wsMessageSchema.safeParse(raw);
    if (!parsed.success) return;

    const message = parsed.data;

    const payloadSchema = payloadSchemas[message.type];
    if (payloadSchema) {
      const payloadResult = payloadSchema.safeParse(message.payload);
      if (!payloadResult.success) return;
    }

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
