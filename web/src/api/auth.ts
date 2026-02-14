import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type {
  LoginRequest,
  LoginResponse,
  SignupRequest,
  SignupResponse,
  PasswordForgotRequest,
  PasswordResetRequest,
  MessageResponse,
} from '@/types';
import type { RawLoginResponse } from '@/types/raw';
import type { AuthProvider } from '@/types';

const VALID_AUTH_PROVIDERS: readonly string[] = [
  'local',
  'google',
  'facebook',
  'apple',
  'github',
  'x',
] satisfies readonly AuthProvider[];

function validateAuthProvider(value: string): AuthProvider {
  if (!VALID_AUTH_PROVIDERS.includes(value)) {
    throw new Error(`Invalid auth_method from server: ${value}`);
  }
  return value as AuthProvider;
}

function mapLoginResponse(raw: RawLoginResponse): LoginResponse {
  return {
    userId: raw.user_id,
    isVerified: raw.is_verified,
    authMethod: validateAuthProvider(raw.auth_method),
    accessToken: raw.access_token,
    refreshToken: raw.refresh_token,
  };
}

export async function login(params: LoginRequest): Promise<LoginResponse> {
  const raw = await apiClient.post<RawLoginResponse>(API_PATHS.AUTH.LOGIN, params);
  return mapLoginResponse(raw);
}

export async function signup(params: SignupRequest): Promise<SignupResponse> {
  return apiClient.post<SignupResponse>(API_PATHS.AUTH.SIGNUP, params);
}

export async function logout(): Promise<void> {
  await apiClient.post<undefined>(API_PATHS.AUTH.LOGOUT);
}

export async function verifyEmail(token: string): Promise<void> {
  await apiClient.get<undefined>(API_PATHS.AUTH.VERIFY_EMAIL(token));
}

export async function sendVerificationEmail(params: {
  email: string;
}): Promise<MessageResponse> {
  return apiClient.post<MessageResponse>(API_PATHS.AUTH.SEND_VERIFICATION, params);
}

export async function forgotPassword(params: PasswordForgotRequest): Promise<MessageResponse> {
  return apiClient.post<MessageResponse>(API_PATHS.AUTH.PASSWORD_FORGOT, params);
}

export async function resetPassword(params: PasswordResetRequest): Promise<void> {
  await apiClient.post<undefined>(API_PATHS.AUTH.PASSWORD_RESET, params);
}

// [MOCK] OAuth requires provider registration (client IDs, redirect URIs).
// Returns null until OAuth is configured.
export async function oauthLogin(_provider: string): Promise<LoginResponse | null> {
  return null;
}
