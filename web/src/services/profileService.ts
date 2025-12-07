import apiClient from '../api';

export interface UserProfile {
    user_id: string;
    first_name: string;
    last_name: string;
    username: string;
    gender: string;
    sexual_preference: string;
    birthday: string;
    occupation: string;
    biography: string;
    location_name: string;
}

export interface Picture {
    id: number;
    url: string;
    is_profile_pic?: boolean;
}

export const getMyProfile = () => {
    return apiClient.get<UserProfile>('/me/profile');
}

export const updateMyProfile = (profile: UserProfile) => {
    return apiClient.put('/me/profile', profile);
}

export const getUserProfile = (userId: string) => {
    return apiClient.get<UserProfile>(`/users/${userId}/profile`);
}

export const getMyPictures = () => {
    return apiClient.get<Picture[]>('/me/profile/pictures');
}

export const getUserPictures = (userId: string) => {
    return apiClient.get<Picture[]>(`/users/${userId}/pictures`);
}

export const uploadPicture = (formData: FormData) => {
    return apiClient.post('/me/profile/pictures', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    });
}

export const setProfilePic = (pictureId: number) => {
    return apiClient.put(`/me/profile/pictures/${pictureId}/status`, { is_profile_pic: true });
}

export const deletePicture = (pictureId: number) => {
    return apiClient.delete(`/me/profile/pictures/${pictureId}`);
}