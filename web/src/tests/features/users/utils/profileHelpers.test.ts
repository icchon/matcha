import { describe, it, expect } from 'vitest';
import { calculateAge, getProfilePicUrl } from '@/features/users/utils/profileHelpers';
import type { UserProfileDetail } from '@/types';

const baseProfile: UserProfileDetail = {
  userId: '00000000-0000-0000-0000-000000000001',
  firstName: 'Alice',
  lastName: 'Smith',
  username: 'alice',
  gender: 'female',
  sexualPreference: 'bisexual',
  birthday: '1995-06-15',
  occupation: 'Designer',
  biography: 'Hello',
  locationName: 'Tokyo',
  fameRating: 75,
  pictures: [],
  tags: [],
  isOnline: true,
  lastConnection: null,
  distance: 3.5,
};

describe('getProfilePicUrl', () => {
  it('returns profile pic URL when isProfilePic is true', () => {
    const profile: UserProfileDetail = {
      ...baseProfile,
      pictures: [
        { id: 1, userId: 'u1', url: '/img/other.jpg', isProfilePic: false, createdAt: '2024-01-01' },
        { id: 2, userId: 'u1', url: '/img/main.jpg', isProfilePic: true, createdAt: '2024-01-02' },
      ],
    };

    expect(
      getProfilePicUrl(profile),
      'Should return the URL of the picture with isProfilePic=true.',
    ).toBe('/img/main.jpg');
  });

  it('returns first picture URL when no profile pic is set', () => {
    const profile: UserProfileDetail = {
      ...baseProfile,
      pictures: [
        { id: 1, userId: 'u1', url: '/img/first.jpg', isProfilePic: false, createdAt: '2024-01-01' },
      ],
    };

    expect(
      getProfilePicUrl(profile),
      'Should fall back to first picture when no isProfilePic=true exists.',
    ).toBe('/img/first.jpg');
  });

  it('returns null when no pictures exist', () => {
    expect(
      getProfilePicUrl(baseProfile),
      'Should return null when profile has no pictures.',
    ).toBeNull();
  });
});

describe('calculateAge', () => {
  it('returns null for null birthday', () => {
    expect(
      calculateAge(null),
      'Should return null when birthday is null.',
    ).toBeNull();
  });

  it('returns null for invalid date string', () => {
    expect(
      calculateAge('not-a-date'),
      'Should return null for invalid date string.',
    ).toBeNull();
  });

  it('calculates age correctly', () => {
    const now = new Date();
    const pastYear = now.getFullYear() - 30;
    const birthday = `${pastYear}-01-01`;

    const age = calculateAge(birthday);
    expect(
      age,
      'Should calculate age based on year difference. Check birthday-before-today logic.',
    ).toBeGreaterThanOrEqual(29);
    expect(age).toBeLessThanOrEqual(30);
  });
});
