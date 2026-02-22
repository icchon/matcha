import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useChatStore } from '@/stores/chatStore';
import type { ChatMessage, ChatConversation } from '@/stores/chatStore';

describe('chatStore', () => {
  beforeEach(() => {
    useChatStore.setState({
      conversations: new Map<string, ChatConversation>(),
      unreadCount: 0,
    });
  });

  it('has correct initial state', () => {
    const state = useChatStore.getState();

    expect(
      state.conversations.size,
      'Initial conversations should be an empty Map. Check initialState.',
    ).toBe(0);
    expect(
      state.unreadCount,
      'Initial unreadCount should be 0. Check initialState.',
    ).toBe(0);
  });

  it('onMessage is a placeholder that stores the message', () => {
    const message: ChatMessage = {
      id: 'msg-1',
      senderId: 'user-a',
      receiverId: 'user-b',
      content: 'Hello!',
      timestamp: '2026-01-01T00:00:00Z',
    };

    useChatStore.getState().onMessage(message);

    const state = useChatStore.getState();
    const conv = state.conversations.get('user-a');

    expect(
      conv,
      'After onMessage, a conversation entry should exist for the senderId.',
    ).toBeDefined();
    expect(
      conv?.messages.length,
      'Conversation should contain exactly 1 message after onMessage.',
    ).toBe(1);
    expect(
      conv?.messages[0]?.content,
      'Stored message content should match input.',
    ).toBe('Hello!');
  });

  it('onAck is callable (placeholder)', () => {
    expect(
      () => useChatStore.getState().onAck({ messageId: 'msg-1', status: 'delivered' }),
      'onAck should be callable without throwing.',
    ).not.toThrow();
  });

  it('onRead is callable (placeholder)', () => {
    expect(
      () => useChatStore.getState().onRead({ conversationId: 'user-a', readAt: '2026-01-01T00:00:00Z' }),
      'onRead should be callable without throwing.',
    ).not.toThrow();
  });

  it('caps messages at 200 per conversation, dropping oldest', () => {
    const senderId = 'user-a';
    for (let i = 0; i < 201; i++) {
      const message: ChatMessage = {
        id: `msg-${i}`,
        senderId,
        receiverId: 'user-b',
        content: `Message ${i}`,
        timestamp: `2026-01-01T00:${String(i).padStart(2, '0')}:00Z`,
      };
      useChatStore.getState().onMessage(message);
    }

    const conv = useChatStore.getState().conversations.get(senderId);
    expect(
      conv?.messages.length,
      'After adding 201 messages, only 200 should remain (MAX_MESSAGES_PER_CONVERSATION=200). Check the slicing logic in onMessage.',
    ).toBe(200);
    expect(
      conv?.messages[0]?.id,
      'The oldest message (msg-0) should have been dropped. The first remaining message should be msg-1. Check slice offset in onMessage.',
    ).toBe('msg-1');
  });

  it('maintains immutability when adding messages', () => {
    const message1: ChatMessage = {
      id: 'msg-1',
      senderId: 'user-a',
      receiverId: 'user-b',
      content: 'First',
      timestamp: '2026-01-01T00:00:00Z',
    };
    const message2: ChatMessage = {
      id: 'msg-2',
      senderId: 'user-a',
      receiverId: 'user-b',
      content: 'Second',
      timestamp: '2026-01-01T00:01:00Z',
    };

    useChatStore.getState().onMessage(message1);
    const stateAfterFirst = useChatStore.getState().conversations;

    useChatStore.getState().onMessage(message2);
    const stateAfterSecond = useChatStore.getState().conversations;

    expect(
      stateAfterFirst,
      'Conversations map should be a new reference after adding a message (immutability).',
    ).not.toBe(stateAfterSecond);
  });
});
