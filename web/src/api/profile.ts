import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { UserProfile, Picture, Tag, UserTag } from '@/types';
import type { Gender, SexualPreference } from '@/types';

export interface CreateProfileRequest {
  readonly firstName: string;
  readonly lastName: string;
  readonly username: string;
  readonly gender: Gender;
  readonly sexualPreference: SexualPreference;
  readonly birthday: string;
  readonly biography: string;
  readonly occupation?: string;
}

export interface UpdateProfileRequest {
  readonly firstName?: string;
  readonly lastName?: string;
  readonly username?: string;
  readonly gender?: Gender;
  readonly sexualPreference?: SexualPreference;
  readonly birthday?: string;
  readonly biography?: string;
  readonly occupation?: string;
  readonly locationName?: string;
}

export async function createProfile(params: CreateProfileRequest): Promise<UserProfile> {
  return apiClient.post<UserProfile>(API_PATHS.PROFILE.CREATE, params);
}

export async function updateProfile(params: UpdateProfileRequest): Promise<UserProfile> {
  return apiClient.put<UserProfile>(API_PATHS.PROFILE.UPDATE, params);
}

export async function getMyProfile(): Promise<UserProfile> {
  return apiClient.get<UserProfile>(API_PATHS.PROFILE.CREATE);
}

export async function uploadPicture(file: File): Promise<Picture> {
  const formData = new FormData();
  formData.append('file', file);
  return apiClient.upload<Picture>(API_PATHS.PROFILE.PICTURES, formData);
}

export async function deletePicture(pictureId: number): Promise<void> {
  await apiClient.delete(API_PATHS.PROFILE.DELETE_PICTURE(pictureId));
}

export async function getTags(): Promise<Tag[]> {
  return apiClient.get<Tag[]>(API_PATHS.TAGS);
}

export async function addTag(tagId: number): Promise<UserTag> {
  return apiClient.post<UserTag>(API_PATHS.USERS.MY_TAGS, { tagId });
}

export async function removeTag(tagId: number): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.DELETE_TAG(tagId));
}
