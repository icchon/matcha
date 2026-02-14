import { z } from 'zod';

export const changePasswordSchema = z
  .object({
    currentPassword: z.string().min(8, 'Password must be at least 8 characters'),
    newPassword: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string().min(1, 'Please confirm your password'),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

export type ChangePasswordFormData = z.infer<typeof changePasswordSchema>;

export const deleteAccountSchema = z
  .object({
    confirmText: z.string().min(1, 'Please type DELETE to confirm'),
  })
  .refine((data) => data.confirmText === 'DELETE', {
    message: 'Please type DELETE to confirm',
    path: ['confirmText'],
  });

export type DeleteAccountFormData = z.infer<typeof deleteAccountSchema>;
