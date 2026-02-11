export type NotificationType = 'like' | 'view' | 'match' | 'unlike' | 'message';

export interface Notification {
  readonly id: number;
  readonly recipientId: string;
  readonly senderId: string | null;
  readonly type: NotificationType;
  readonly isRead: boolean | null;
  readonly createdAt: string;
}
