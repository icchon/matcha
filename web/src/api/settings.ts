import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Block, MessageResponse } from '@/types';

export interface ChangePasswordRequest {
  readonly currentPassword: string;
  readonly newPassword: string;
}

export interface DeleteAccountRequest {
  readonly currentPassword: string;
}

export async function deleteAccount(params: DeleteAccountRequest): Promise<void> {
  // Uses POST instead of DELETE to include re-authentication body
  await apiClient.post(API_PATHS.USERS.DELETE_ACCOUNT, params);
}

// Uses dedicated CHANGE_PASSWORD path (currently shares endpoint with PASSWORD_RESET;
// backend differentiates by presence of currentPassword vs token)
export async function changePassword(
  params: ChangePasswordRequest,
): Promise<MessageResponse> {
  return apiClient.post<MessageResponse>(API_PATHS.AUTH.CHANGE_PASSWORD, params);
}

export async function getBlockList(): Promise<Block[]> {
  return apiClient.get<Block[]>(API_PATHS.USERS.MY_BLOCKS);
}

// [MOCK] BE-08 #25: endpoint not yet implemented
export async function unblockUser(userId: string): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.UNBLOCK(userId));
}
