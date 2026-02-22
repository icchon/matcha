/** Raw API response types matching backend snake_case JSON field names */

export interface RawLoginResponse {
  readonly user_id: string;
  readonly is_verified: boolean;
  readonly auth_method: string;
  readonly access_token: string;
  readonly refresh_token: string;
}

export interface RawPicture {
  readonly id: number;
  readonly user_id: string;
  readonly url: string;
  readonly is_profile_pic: boolean | null;
  readonly created_at: string;
}

export interface RawUserProfileResponse {
  readonly user_id: string;
  readonly first_name: string | null;
  readonly last_name: string | null;
  readonly username: string | null;
  readonly gender: string | null;
  readonly sexual_preference: string | null;
  readonly birthday: string | null;
  readonly occupation: string | null;
  readonly biography: string | null;
  readonly location_name: string | null;
  readonly fame_rating: number | null;
  readonly pictures: readonly RawPicture[];
  readonly tags: readonly { readonly id: number; readonly name: string }[];
  readonly is_online: boolean;
  readonly last_connection: string | null;
  readonly distance?: number | null;
}

export interface RawLike {
  readonly liker_id: string;
  readonly liked_id: string;
  readonly created_at: string;
}

export interface RawView {
  readonly viewer_id: string;
  readonly viewed_id: string;
  readonly view_time: string;
}
