import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Picture } from '@/types';

export async function uploadPicture(file: File): Promise<Picture> {
  const formData = new FormData();
  formData.append('file', file);
  return apiClient.upload<Picture>(API_PATHS.PROFILE.PICTURES, formData);
}

export async function deletePicture(pictureId: number): Promise<void> {
  await apiClient.delete(API_PATHS.PROFILE.DELETE_PICTURE(pictureId));
}
