import { z } from 'zod';

export const changePasswordSchema = z
  .object({
    currentPassword: z.string().min(8, 'Password must be at least 8 characters'),
    newPassword: z
      .string()
      .min(8, 'Password must be at least 8 characters')
      .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
      .regex(/[0-9]/, 'Password must contain at least one number')
      .regex(/[^a-zA-Z0-9]/, 'Password must contain at least one special character'),
    confirmPassword: z.string().min(8, 'Password must be at least 8 characters'),
  })
  .refine((data) => data.newPassword !== data.currentPassword, {
    message: 'New password must be different from current password',
    path: ['newPassword'],
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

export type ChangePasswordFormData = z.infer<typeof changePasswordSchema>;

export const deleteAccountSchema = z
  .object({
    currentPassword: z.string().min(1, 'Please enter your current password'),
    confirmText: z.string().min(1, 'Please type DELETE to confirm'),
  })
  .refine((data) => data.confirmText === 'DELETE', {
    message: 'Please type DELETE to confirm',
    path: ['confirmText'],
  });

export type DeleteAccountFormData = z.infer<typeof deleteAccountSchema>;
