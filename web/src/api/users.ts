import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { LikeUserResponse, Like, View, MessageResponse, Tag, Picture } from '@/types';
import type { UserProfileDetail } from '@/types';
import type { RawUserProfileResponse, RawLike, RawView, RawPicture } from '@/types/raw';

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

function validateUserId(userId: string): string {
  if (!UUID_REGEX.test(userId)) {
    throw new Error('Invalid user ID format');
  }
  return userId;
}

function isSafeImageUrl(url: string): boolean {
  try {
    const parsed = new URL(url, window.location.origin);
    return parsed.protocol === 'https:' || parsed.protocol === 'http:' || parsed.pathname.startsWith('/');
  } catch {
    return false;
  }
}

function mapPicture(raw: RawPicture): Picture {
  return {
    id: raw.id,
    userId: raw.user_id,
    url: isSafeImageUrl(raw.url) ? raw.url : '',
    isProfilePic: raw.is_profile_pic,
    createdAt: raw.created_at,
  };
}

function mapUserProfileResponse(raw: RawUserProfileResponse): UserProfileDetail {
  return {
    userId: raw.user_id,
    firstName: raw.first_name,
    lastName: raw.last_name,
    username: raw.username,
    gender: raw.gender as UserProfileDetail['gender'],
    sexualPreference: raw.sexual_preference as UserProfileDetail['sexualPreference'],
    birthday: raw.birthday,
    occupation: raw.occupation,
    biography: raw.biography,
    locationName: raw.location_name,
    fameRating: raw.fame_rating,
    pictures: raw.pictures.map(mapPicture),
    tags: raw.tags.map((t): Tag => ({ id: t.id, name: t.name })),
    isOnline: raw.is_online,
    lastConnection: raw.last_connection,
    distance: raw.distance,
  };
}

function mapLike(raw: RawLike): Like {
  return {
    likerId: raw.liker_id,
    likedId: raw.liked_id,
    createdAt: raw.created_at,
  };
}

function mapView(raw: RawView): View {
  return {
    viewerId: raw.viewer_id,
    viewedId: raw.viewed_id,
    viewTime: raw.view_time,
  };
}

export async function getUserProfile(userId: string): Promise<UserProfileDetail> {
  const raw = await apiClient.get<RawUserProfileResponse>(API_PATHS.PROFILE.GET(validateUserId(userId)));
  return mapUserProfileResponse(raw);
}

export async function likeUser(userId: string): Promise<LikeUserResponse> {
  return apiClient.post<LikeUserResponse>(API_PATHS.USERS.LIKE(validateUserId(userId)));
}

export async function unlikeUser(userId: string): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.UNLIKE(validateUserId(userId)));
}

export async function blockUser(userId: string): Promise<void> {
  await apiClient.post<undefined>(API_PATHS.USERS.BLOCK(validateUserId(userId)));
}

// [MOCK] BE-08 #25: endpoint not yet implemented
export async function unblockUser(userId: string): Promise<void> {
  await apiClient.delete(API_PATHS.USERS.UNBLOCK(validateUserId(userId)));
}

// [MOCK] BE-08 #25: endpoint not yet implemented
// TODO(FE-XX): Add report reason selection UI before BE-08 #25 is implemented
export async function reportUser(_userId: string, _reason: string): Promise<MessageResponse> {
  return { message: 'Report submitted' };
}

export async function getLikedUsers(): Promise<readonly Like[]> {
  const raw = await apiClient.get<readonly RawLike[]>(API_PATHS.USERS.MY_LIKES);
  return raw.map(mapLike);
}

export async function getWhoLikedMe(): Promise<readonly Like[]> {
  const raw = await apiClient.get<readonly RawLike[]>(API_PATHS.PROFILE.WHO_LIKED_ME);
  return raw.map(mapLike);
}

export async function getViewedUsers(): Promise<readonly View[]> {
  const raw = await apiClient.get<readonly RawView[]>(API_PATHS.USERS.MY_VIEWS);
  return raw.map(mapView);
}

export async function getWhoViewedMe(): Promise<readonly View[]> {
  const raw = await apiClient.get<readonly RawView[]>(API_PATHS.PROFILE.WHO_VIEWED_ME);
  return raw.map(mapView);
}
