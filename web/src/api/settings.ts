import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Block, MessageResponse } from '@/types';

export interface ChangePasswordRequest {
  readonly currentPassword: string;
  readonly newPassword: string;
}

export async function deleteAccount(): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.DELETE_ME);
}

// Reuses PASSWORD_RESET endpoint — backend differentiates by presence of currentPassword vs token
export async function changePassword(
  params: ChangePasswordRequest,
): Promise<MessageResponse> {
  return apiClient.post<MessageResponse>(API_PATHS.AUTH.PASSWORD_RESET, params);
}

export async function getBlockList(): Promise<Block[]> {
  return apiClient.get<Block[]>(API_PATHS.USERS.MY_BLOCKS);
}

// [MOCK] BE-08 #25: endpoint not yet implemented
export async function unblockUser(userId: string): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.UNBLOCK(userId));
}
