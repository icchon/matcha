import { type FC, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { sendVerificationSchema, type SendVerificationFormData } from '@/lib/validators';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

const VerifyEmailPage: FC = () => {
  const { token } = useParams<{ token: string }>();
  const { verifyEmail, sendVerificationEmail, isLoading, error } = useAuth();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SendVerificationFormData>({
    resolver: zodResolver(sendVerificationSchema),
  });

  useEffect(() => {
    if (token) {
      void verifyEmail(token);
    }
  }, [token, verifyEmail]);

  const onResend = (data: SendVerificationFormData) => {
    sendVerificationEmail({ email: data.email });
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md space-y-6">
        <h1 className="text-center text-2xl font-bold text-gray-900">Email Verification</h1>
        {error ? (
          <p className="text-center text-red-600">{error}</p>
        ) : (
          <p className="text-center text-gray-600">Verifying your email...</p>
        )}
        <div className="border-t border-gray-300 pt-4">
          <p className="mb-4 text-center text-sm text-gray-600">
            Didn&apos;t receive the email? Resend verification below.
          </p>
          <form onSubmit={handleSubmit(onResend)} noValidate className="flex flex-col gap-4">
            <Input
              label="Email"
              type="email"
              autoComplete="email"
              error={errors.email?.message}
              {...register('email')}
            />
            <Button type="submit" loading={isLoading}>
              Resend Verification Email
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
};

export { VerifyEmailPage };
