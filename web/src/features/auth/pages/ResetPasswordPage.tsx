import type { FC } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { resetPasswordSchema, type ResetPasswordFormData } from '@/lib/validators';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

const ResetPasswordPage: FC = () => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const { resetPassword, isLoading } = useAuth();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
  });

  if (!token) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
        <div className="w-full max-w-md space-y-6">
          <h1 className="text-center text-2xl font-bold text-gray-900">Invalid Reset Link</h1>
          <p className="text-center text-red-600">
            This password reset link is invalid or has expired.
          </p>
          <p className="text-center">
            <Link to="/forgot-password" className="text-blue-600 hover:underline">
              Request a new reset link
            </Link>
          </p>
        </div>
      </div>
    );
  }

  const onSubmit = (data: ResetPasswordFormData) => {
    resetPassword({ token, password: data.password });
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md space-y-6">
        <h1 className="text-center text-2xl font-bold text-gray-900">Reset Password</h1>
        <form onSubmit={handleSubmit(onSubmit)} noValidate className="flex flex-col gap-4">
          <Input
            label="New Password"
            type="password"
            autoComplete="new-password"
            error={errors.password?.message}
            {...register('password')}
          />
          <Input
            label="Confirm Password"
            type="password"
            autoComplete="new-password"
            error={errors.passwordConfirm?.message}
            {...register('passwordConfirm')}
          />
          <Button type="submit" loading={isLoading}>
            Reset Password
          </Button>
        </form>
      </div>
    </div>
  );
};

export { ResetPasswordPage };
