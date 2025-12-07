import apiClient from '../api';

export const likeUser = (userId: string) => {
    return apiClient.post(`/users/${userId}/like`);
};

export const unlikeUser = (userId: string) => {
    return apiClient.delete(`/users/${userId}/like`);
};

export const blockUser = (userId: string) => {
    return apiClient.post(`/users/${userId}/block`);
};
