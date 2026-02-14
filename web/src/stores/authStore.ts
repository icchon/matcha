import { create } from 'zustand';
import type { AuthProvider } from '@/types';
import { clearTokens, setTokens } from '@/api/client';

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

  // In-memory tokens do not survive page reload.
  // initialize is kept for API compatibility but is a no-op.
  initialize: () => {
    // No-op: tokens are in-memory only and lost on reload.
    // User must re-login after page refresh.
  },
}));
