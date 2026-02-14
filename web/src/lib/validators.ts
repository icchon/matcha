import { z } from 'zod';

export const loginSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
});

export type LoginFormData = z.infer<typeof loginSchema>;

export const signupSchema = z
  .object({
    email: z.string().email('Please enter a valid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    passwordConfirm: z.string().min(1, 'Please confirm your password'),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: 'Passwords do not match',
    path: ['passwordConfirm'],
  });

export type SignupFormData = z.infer<typeof signupSchema>;

export const forgotPasswordSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
});

export type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

export const resetPasswordSchema = z
  .object({
    password: z.string().min(8, 'Password must be at least 8 characters'),
    passwordConfirm: z.string().min(1, 'Please confirm your password'),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: 'Passwords do not match',
    path: ['passwordConfirm'],
  });

export type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

export const sendVerificationSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
});

export type SendVerificationFormData = z.infer<typeof sendVerificationSchema>;

export const profileSchema = z.object({
  firstName: z.string().min(1, 'First name is required').max(50, 'First name must be 50 characters or less'),
  lastName: z.string().min(1, 'Last name is required').max(50, 'Last name must be 50 characters or less'),
  username: z.string().min(3, 'Username must be at least 3 characters').max(30, 'Username must be 30 characters or less'),
  gender: z.enum(['male', 'female', 'other'], { required_error: 'Gender is required' }),
  sexualPreference: z.enum(['heterosexual', 'homosexual', 'bisexual'], { required_error: 'Sexual preference is required' }),
  birthday: z.string().min(1, 'Birthday is required'),
  biography: z.string().min(1, 'Biography is required').max(500, 'Biography must be 500 characters or less'),
  occupation: z.string().max(100, 'Occupation must be 100 characters or less').optional().or(z.literal('')),
});

export type ProfileFormData = z.infer<typeof profileSchema>;
