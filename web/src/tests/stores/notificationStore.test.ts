import { describe, it, expect, beforeEach } from 'vitest';
import { useNotificationStore } from '@/stores/notificationStore';
import type { Notification } from '@/stores/notificationStore';

describe('notificationStore', () => {
  beforeEach(() => {
    useNotificationStore.setState({
      notifications: [],
      unreadCount: 0,
    });
  });

  it('has correct initial state', () => {
    const state = useNotificationStore.getState();

    expect(
      state.notifications.length,
      'Initial notifications should be an empty array. Check initialState.',
    ).toBe(0);
    expect(
      state.unreadCount,
      'Initial unreadCount should be 0. Check initialState.',
    ).toBe(0);
  });

  it('onNotification adds notification to the list', () => {
    const notification: Notification = {
      id: 'notif-1',
      type: 'like',
      message: 'Someone liked you!',
      timestamp: '2026-01-01T00:00:00Z',
      read: false,
    };

    useNotificationStore.getState().onNotification(notification);

    const state = useNotificationStore.getState();
    expect(
      state.notifications.length,
      'After onNotification, notifications array should have 1 item.',
    ).toBe(1);
    expect(
      state.notifications[0]?.message,
      'Stored notification message should match input.',
    ).toBe('Someone liked you!');
    expect(
      state.unreadCount,
      'Unread count should increment to 1 for unread notification.',
    ).toBe(1);
  });

  it('maintains immutability when adding notifications', () => {
    const notif1: Notification = {
      id: 'notif-1',
      type: 'like',
      message: 'First',
      timestamp: '2026-01-01T00:00:00Z',
      read: false,
    };
    const notif2: Notification = {
      id: 'notif-2',
      type: 'view',
      message: 'Second',
      timestamp: '2026-01-01T00:01:00Z',
      read: false,
    };

    useNotificationStore.getState().onNotification(notif1);
    const arrayAfterFirst = useNotificationStore.getState().notifications;

    useNotificationStore.getState().onNotification(notif2);
    const arrayAfterSecond = useNotificationStore.getState().notifications;

    expect(
      arrayAfterFirst,
      'Notifications array should be a new reference after adding (immutability).',
    ).not.toBe(arrayAfterSecond);
  });

  it('does not increment unreadCount for already-read notifications', () => {
    const notification: Notification = {
      id: 'notif-1',
      type: 'like',
      message: 'Already read',
      timestamp: '2026-01-01T00:00:00Z',
      read: true,
    };

    useNotificationStore.getState().onNotification(notification);

    expect(
      useNotificationStore.getState().unreadCount,
      'Unread count should NOT increment for read notifications.',
    ).toBe(0);
  });
});
