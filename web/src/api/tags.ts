import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Tag, UserTag } from '@/types';

export async function getTags(): Promise<Tag[]> {
  return apiClient.get<Tag[]>(API_PATHS.TAGS);
}

export async function addTag(tagId: number): Promise<UserTag> {
  return apiClient.post<UserTag>(API_PATHS.USERS.MY_TAGS, { tagId });
}

export async function removeTag(tagId: number): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.DELETE_TAG(tagId));
}
