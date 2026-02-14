import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import * as authApi from '@/api/auth';
import { useAuthStore } from '@/stores/authStore';
import type { LoginRequest, PasswordForgotRequest, PasswordResetRequest, SignupRequest } from '@/types';

interface UseAuthReturn {
  readonly isLoading: boolean;
  readonly error: string | null;
  readonly login: (params: LoginRequest) => Promise<void>;
  readonly signup: (params: SignupRequest) => Promise<void>;
  readonly logout: () => Promise<void>;
  readonly verifyEmail: (token: string) => Promise<void>;
  readonly sendVerificationEmail: (params: { email: string }) => Promise<void>;
  readonly forgotPassword: (params: PasswordForgotRequest) => Promise<void>;
  readonly resetPassword: (params: PasswordResetRequest) => Promise<void>;
}

export function useAuth(): UseAuthReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const storeLogin = useAuthStore((s) => s.login);
  const storeLogout = useAuthStore((s) => s.logout);

  const login = useCallback(
    async (params: LoginRequest) => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await authApi.login(params);
        storeLogin(response);
        navigate('/');
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Login failed';
        setError(message);
        toast.error(message);
      } finally {
        setIsLoading(false);
      }
    },
    [storeLogin, navigate],
  );

  const signup = useCallback(
    async (params: SignupRequest) => {
      setIsLoading(true);
      setError(null);
      try {
        await authApi.signup(params);
        toast.success('Account created. Please check your email to verify.');
        navigate('/login');
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Signup failed';
        setError(message);
        toast.error(message);
      } finally {
        setIsLoading(false);
      }
    },
    [navigate],
  );

  const logout = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      await authApi.logout();
    } catch {
      // Logout API failure is non-critical — still clear local state
    } finally {
      storeLogout();
      setIsLoading(false);
      navigate('/login');
    }
  }, [storeLogout, navigate]);

  const verifyEmail = useCallback(async (token: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await authApi.verifyEmail(token);
      toast.success('Email verified successfully!');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Verification failed';
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const sendVerificationEmail = useCallback(async (params: { email: string }) => {
    setIsLoading(true);
    setError(null);
    try {
      await authApi.sendVerificationEmail(params);
      toast.success('Verification email sent. Please check your inbox.');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to send verification email';
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const forgotPassword = useCallback(async (params: PasswordForgotRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      await authApi.forgotPassword(params);
      toast.success('Password reset email sent. Please check your inbox.');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to send reset email';
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const resetPassword = useCallback(
    async (params: PasswordResetRequest) => {
      setIsLoading(true);
      setError(null);
      try {
        await authApi.resetPassword(params);
        toast.success('Password reset successfully!');
        navigate('/login');
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Password reset failed';
        setError(message);
        toast.error(message);
      } finally {
        setIsLoading(false);
      }
    },
    [navigate],
  );

  return {
    isLoading,
    error,
    login,
    signup,
    logout,
    verifyEmail,
    sendVerificationEmail,
    forgotPassword,
    resetPassword,
  };
}
