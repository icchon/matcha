import apiClient from '../api';
import { UserProfile } from './profileService';

export interface Message {
  id: number;
  sender_id: string;
  recipient_id: string;
  content: string;
  sent_at: string;
}

export interface ChatOverview {
  other_user: UserProfile;
  last_message: Message | null;
}

export const getChatsForUser = () => {
    return apiClient.get<ChatOverview[]>('/me/chats');
};

export const getChatMessages = (userId: string, limit: number = 20, offset: number = 0) => {
    return apiClient.get<Message[]>(`/chats/${userId}/messages`, {
        params: { limit, offset }
    });
};
