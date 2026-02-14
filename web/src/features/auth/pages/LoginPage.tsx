import type { FC } from 'react';
import { LoginForm } from '@/features/auth/components/LoginForm';
import { OAuthButtons } from '@/features/auth/components/OAuthButtons';

const LoginPage: FC = () => {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md space-y-6">
        <h1 className="text-center text-2xl font-bold text-gray-900">Log In</h1>
        <LoginForm />
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-300" />
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="bg-gray-50 px-2 text-gray-500">or</span>
          </div>
        </div>
        <OAuthButtons />
      </div>
    </div>
  );
};

export { LoginPage };
