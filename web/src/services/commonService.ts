import apiClient from '../api';
import { UserProfile } from './profileService'; // Assuming UserProfile is shared

export interface Tag {
  id: number;
  name: string;
}

export interface Notification {
  id: number;
  recipient_id: string;
  sender_id?: string; // Optional because sql.NullString
  type: string;
  is_read?: boolean; // Optional because sql.NullBool
  created_at: string;
}

export const getAllProfiles = () => {
    return apiClient.get<UserProfile[]>('/profiles');
};

export const getRecommendedProfiles = () => {
    return apiClient.get<UserProfile[]>('/profiles/recommends');
};

export const getAllTags = () => {
    return apiClient.get<Tag[]>('/tags');
};

export const getMyTags = () => {
    return apiClient.get<Tag[]>('/me/tags');
};

export const addMyTag = (tagId: number) => {
    return apiClient.post('/me/tags', { tag_id: tagId });
};

export const deleteMyTag = (tagId: number) => {
    return apiClient.delete(`/me/tags/${tagId}`);
};

export const getMyNotifications = () => {
    return apiClient.get<Notification[]>('/me/notifications');
};
