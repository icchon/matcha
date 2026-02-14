import { create } from 'zustand';
import type { ChatMessagePayload, AckPayload, ReadPayload } from '@/features/websocket/types';

export type ChatMessage = ChatMessagePayload;

export interface ChatConversation {
  readonly messages: readonly ChatMessage[];
  readonly lastMessageAt: string;
}

interface ChatState {
  readonly conversations: Map<string, ChatConversation>;
  readonly unreadCount: number;
}

interface ChatActions {
  readonly onMessage: (message: ChatMessage) => void;
  readonly onAck: (payload: AckPayload) => void;
  readonly onRead: (payload: ReadPayload) => void;
}

type ChatStore = ChatState & ChatActions;

export const useChatStore = create<ChatStore>()((set) => ({
  conversations: new Map<string, ChatConversation>(),
  unreadCount: 0,

  onMessage: (message: ChatMessage) => {
    set((state) => {
      const newConversations = new Map(state.conversations);
      const existing = newConversations.get(message.senderId);
      const updatedConversation: ChatConversation = {
        messages: [...(existing?.messages ?? []), message],
        lastMessageAt: message.timestamp,
      };
      newConversations.set(message.senderId, updatedConversation);

      return { conversations: newConversations };
    });
  },

  // Placeholder — will be implemented in chat feature
  onAck: (_payload: AckPayload) => {
    // No-op placeholder
  },

  // Placeholder — will be implemented in chat feature
  onRead: (_payload: ReadPayload) => {
    // No-op placeholder
  },
}));
