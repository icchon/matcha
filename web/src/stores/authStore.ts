import { create } from 'zustand';
import type { AuthProvider } from '@/types';
import { clearTokens, getAccessToken, setTokens } from '@/api/client';

interface AuthState {
  readonly userId: string | null;
  readonly isAuthenticated: boolean;
  readonly isVerified: boolean;
  readonly authMethod: AuthProvider | null;
}

interface AuthActions {
  readonly login: (params: {
    userId: string;
    isVerified: boolean;
    authMethod: AuthProvider;
    accessToken: string;
    refreshToken: string;
  }) => void;
  readonly logout: () => void;
  readonly initialize: () => void;
}

type AuthStore = AuthState & AuthActions;

const initialState: AuthState = {
  userId: null,
  isAuthenticated: false,
  isVerified: false,
  authMethod: null,
};

function decodeTokenPayload(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return null;
    const payload = JSON.parse(atob(parts[1]));
    return payload;
  } catch {
    return null;
  }
}

export const useAuthStore = create<AuthStore>()((set) => ({
  ...initialState,

  login: ({ userId, isVerified, authMethod, accessToken, refreshToken }) => {
    setTokens(accessToken, refreshToken);
    set({
      userId,
      isAuthenticated: true,
      isVerified,
      authMethod,
    });
  },

  logout: () => {
    clearTokens();
    set({ ...initialState });
  },

  initialize: () => {
    const token = getAccessToken();
    if (!token) return;

    const payload = decodeTokenPayload(token);
    if (!payload) {
      clearTokens();
      return;
    }

    const exp = payload.exp as number | undefined;
    if (exp && exp * 1000 < Date.now()) {
      clearTokens();
      return;
    }

    set({
      userId: (payload.sub as string) ?? null,
      isAuthenticated: true,
      isVerified: (payload.is_verified as boolean) ?? false,
      authMethod: (payload.auth_method as AuthProvider) ?? null,
    });
  },
}));
