export interface Message {
  readonly id: number;
  readonly senderId: string;
  readonly recipientId: string;
  readonly content: string;
  readonly sentAt: string;
  readonly isRead: boolean | null;
}

export interface ChatPartner {
  readonly userId: string;
  readonly username: string | null;
  readonly profilePicUrl: string | null;
}
