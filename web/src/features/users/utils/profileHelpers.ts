import type { UserProfileDetail } from '@/types';

export function getProfilePicUrl(profile: UserProfileDetail): string | null {
  const profilePic = profile.pictures.find((p) => p.isProfilePic);
  return profilePic?.url ?? profile.pictures[0]?.url ?? null;
}

export function calculateAge(birthday: string | null): number | null {
  if (!birthday) return null;
  const birth = new Date(birthday);
  if (isNaN(birth.getTime())) return null;
  const now = new Date();
  let age = now.getFullYear() - birth.getFullYear();
  const monthDiff = now.getMonth() - birth.getMonth();
  if (monthDiff < 0 || (monthDiff === 0 && now.getDate() < birth.getDate())) {
    age -= 1;
  }
  return age;
}
