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

export const useNotificationStore = create<NotificationStore>()((set) => ({
  notifications: [],
  unreadCount: 0,

  onNotification: (notification: Notification) => {
    set((state) => ({
      notifications: [...state.notifications, notification],
      unreadCount: notification.read
        ? state.unreadCount
        : state.unreadCount + 1,
    }));
  },
}));
