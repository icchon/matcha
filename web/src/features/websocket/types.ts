// All WS message payload types are defined here to avoid circular dependencies.
// Domain stores (chatStore, notificationStore) import from this file.

export interface ChatMessagePayload {
  readonly id: string;
  readonly senderId: string;
  readonly receiverId: string;
  readonly content: string;
  readonly timestamp: string;
}

export interface AckPayload {
  readonly messageId: string;
  readonly status: string;
}

export interface ReadPayload {
  readonly conversationId: string;
  readonly readAt: string;
}

export interface NotificationPayload {
  readonly id: string;
  readonly type: string;
  readonly message: string;
  readonly timestamp: string;
  readonly read: boolean;
}

export type WsMessage =
  | { readonly type: 'chat.message'; readonly payload: ChatMessagePayload }
  | { readonly type: 'chat.ack'; readonly payload: AckPayload }
  | { readonly type: 'chat.read'; readonly payload: ReadPayload }
  | { readonly type: 'notification'; readonly payload: NotificationPayload };

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

export type MessageHandler = (payload: WsMessage['payload']) => void;
