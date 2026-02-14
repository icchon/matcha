import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { UserProfile } from '@/types';
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
