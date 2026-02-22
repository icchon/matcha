export type Gender = 'male' | 'female' | 'other';

export type SexualPreference = 'heterosexual' | 'homosexual' | 'bisexual';

export type AuthProvider = 'local' | 'google' | 'facebook' | 'apple' | 'github' | 'x';

export interface User {
  readonly id: string;
  readonly createdAt: string;
  readonly lastConnection: string | null;
}

export interface UserProfile {
  readonly userId: string;
  readonly firstName: string | null;
  readonly lastName: string | null;
  readonly username: string | null;
  readonly gender: Gender | null;
  readonly sexualPreference: SexualPreference | null;
  readonly birthday: string | null;
  readonly occupation: string | null;
  readonly biography: string | null;
  readonly locationName: string | null;
  readonly fameRating: number | null;
  readonly distance?: number | null;
}

export interface UserData {
  readonly userId: string;
  readonly latitude: number | null;
  readonly longitude: number | null;
  readonly internalScore: number | null;
}

export interface Tag {
  readonly id: number;
  readonly name: string;
}

export interface UserTag {
  readonly userId: string;
  readonly tagId: number;
}

export interface Picture {
  readonly id: number;
  readonly userId: string;
  readonly url: string;
  readonly isProfilePic: boolean | null;
  readonly createdAt: string;
}

export interface UserProfileDetail {
  readonly userId: string;
  readonly firstName: string | null;
  readonly lastName: string | null;
  readonly username: string | null;
  readonly gender: Gender | null;
  readonly sexualPreference: SexualPreference | null;
  readonly birthday: string | null;
  readonly occupation: string | null;
  readonly biography: string | null;
  readonly locationName: string | null;
  readonly fameRating: number | null;
  readonly pictures: readonly Picture[];
  readonly tags: readonly Tag[];
  readonly isOnline: boolean;
  readonly lastConnection: string | null;
  readonly distance?: number | null;
}
