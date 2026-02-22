import { create } from 'zustand';
import type { NotificationPayload } from '@/features/websocket/types';

export type Notification = NotificationPayload;

interface NotificationState {
  readonly notifications: readonly Notification[];
  readonly unreadCount: number;
}

interface NotificationActions {
  readonly onNotification: (notification: Notification) => void;
}

type NotificationStore = NotificationState & NotificationActions;

const MAX_NOTIFICATIONS = 100;

export const useNotificationStore = create<NotificationStore>()((set) => ({
  notifications: [],
  unreadCount: 0,

  onNotification: (notification: Notification) => {
    set((state) => {
      const updated = [...state.notifications, notification];
      const capped = updated.length > MAX_NOTIFICATIONS
        ? updated.slice(updated.length - MAX_NOTIFICATIONS)
        : updated;

      return {
        notifications: capped,
        unreadCount: notification.read
          ? state.unreadCount
          : state.unreadCount + 1,
      };
    });
  },
}));
