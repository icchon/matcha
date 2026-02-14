import { create } from 'zustand';

export interface Notification {
  readonly id: string;
  readonly type: string;
  readonly message: string;
  readonly timestamp: string;
  readonly read: boolean;
}

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
