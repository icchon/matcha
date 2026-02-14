import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  login,
  signup,
  logout,
  verifyEmail,
  sendVerificationEmail,
  forgotPassword,
  resetPassword,
  oauthLogin,
} from '@/api/auth';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { RawLoginResponse } from '@/types/raw';

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const mockGet = vi.mocked(apiClient.get);
const mockPost = vi.mocked(apiClient.post);

const rawLoginResponse: RawLoginResponse = {
  user_id: 'uuid-123',
  is_verified: true,
  auth_method: 'local',
  access_token: 'access-token',
  refresh_token: 'refresh-token',
};

beforeEach(() => {
  vi.resetAllMocks();
});

describe('login', () => {
  it('calls POST /auth/login and maps snake_case response to camelCase', async () => {
    mockPost.mockResolvedValue(rawLoginResponse);

    const result = await login({ email: 'user@example.com', password: 'password123' });

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.AUTH.LOGIN, {
      email: 'user@example.com',
      password: 'password123',
    });
    expect(
      result.userId,
      'login should map user_id to userId. Check mapLoginResponse.',
    ).toBe('uuid-123');
    expect(result.isVerified).toBe(true);
    expect(result.authMethod).toBe('local');
    expect(result.accessToken).toBe('access-token');
    expect(result.refreshToken).toBe('refresh-token');
  });
});

describe('signup', () => {
  it('calls POST /auth/signup with email and password', async () => {
    mockPost.mockResolvedValue({ message: 'Check email' });

    const result = await signup({ email: 'new@example.com', password: 'password123' });

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.AUTH.SIGNUP, {
      email: 'new@example.com',
      password: 'password123',
    });
    expect(result.message).toBe('Check email');
  });
});

describe('logout', () => {
  it('calls POST /auth/logout', async () => {
    mockPost.mockResolvedValue(undefined);

    await logout();

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.AUTH.LOGOUT);
  });
});

describe('verifyEmail', () => {
  it('calls GET /auth/verify/{token}', async () => {
    mockGet.mockResolvedValue(undefined);

    await verifyEmail('abc-token-123');

    expect(
      mockGet,
      'verifyEmail should call GET with the token in the URL path.',
    ).toHaveBeenCalledWith(API_PATHS.AUTH.VERIFY_EMAIL('abc-token-123'));
  });
});

describe('sendVerificationEmail', () => {
  it('calls POST /auth/verify/mail with email', async () => {
    mockPost.mockResolvedValue({ message: 'Check email' });

    const result = await sendVerificationEmail({ email: 'user@example.com' });

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.AUTH.SEND_VERIFICATION, {
      email: 'user@example.com',
    });
    expect(result.message).toBe('Check email');
  });
});

describe('forgotPassword', () => {
  it('calls POST /auth/password/forgot with email', async () => {
    mockPost.mockResolvedValue({ message: 'Check email' });

    const result = await forgotPassword({ email: 'user@example.com' });

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.AUTH.PASSWORD_FORGOT, {
      email: 'user@example.com',
    });
    expect(result.message).toBe('Check email');
  });
});

describe('resetPassword', () => {
  it('calls POST /auth/password/reset with token and password', async () => {
    mockPost.mockResolvedValue(undefined);

    await resetPassword({ token: 'reset-token', password: 'newpass123' });

    expect(mockPost).toHaveBeenCalledWith(API_PATHS.AUTH.PASSWORD_RESET, {
      token: 'reset-token',
      password: 'newpass123',
    });
  });
});

describe('oauthLogin', () => {
  it('shows mock toast and returns null (OAuth not yet available)', async () => {
    const result = await oauthLogin('google');

    expect(
      result,
      '[MOCK] oauthLogin should return null until OAuth provider registration is complete.',
    ).toBeNull();
  });
});
