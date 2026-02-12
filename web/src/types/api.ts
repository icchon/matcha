export interface LoginRequest {
  readonly email: string;
  readonly password: string;
}

import type { AuthProvider } from './user.ts';

export interface LoginResponse {
  readonly userId: string;
  readonly isVerified: boolean;
  readonly authMethod: AuthProvider;
  readonly accessToken: string;
  readonly refreshToken: string;
}

export interface SignupRequest {
  readonly email: string;
  readonly password: string;
}

export interface SignupResponse {
  readonly message: string;
}

export interface OAuthLoginRequest {
  readonly code: string;
  readonly codeVerifier: string;
}

export interface PasswordForgotRequest {
  readonly email: string;
}

export interface PasswordResetRequest {
  readonly token: string;
  readonly password: string;
}

export interface ApiError {
  readonly error: string;
}

export interface MessageResponse {
  readonly message: string;
}

export interface LikeUserResponse {
  readonly matched: boolean;
}

export interface Like {
  readonly likerId: string;
  readonly likedId: string;
  readonly createdAt: string;
}

export interface View {
  readonly viewerId: string;
  readonly viewedId: string;
  readonly viewTime: string;
}

export interface Block {
  readonly blockerId: string;
  readonly blockedId: string;
}
