/** Raw API response types matching backend snake_case JSON field names */

export interface RawLoginResponse {
  readonly user_id: string;
  readonly is_verified: boolean;
  readonly auth_method: string;
  readonly access_token: string;
  readonly refresh_token: string;
}
