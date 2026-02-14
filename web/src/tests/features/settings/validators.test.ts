import { describe, it, expect } from 'vitest';
import {
  changePasswordSchema,
  deleteAccountSchema,
  type ChangePasswordFormData,
  type DeleteAccountFormData,
} from '@/features/settings/validators';

describe('changePasswordSchema', () => {
  it('accepts valid input with matching passwords', () => {
    const input: ChangePasswordFormData = {
      currentPassword: 'oldpass123',
      newPassword: 'newpass123',
      confirmPassword: 'newpass123',
    };
    const result = changePasswordSchema.safeParse(input);

    expect(
      result.success,
      'Valid input with matching passwords should pass validation. Check schema refinement.',
    ).toBe(true);
  });

  it('rejects when currentPassword is too short', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'short',
      newPassword: 'newpass123',
      confirmPassword: 'newpass123',
    });

    expect(
      result.success,
      'currentPassword shorter than 8 chars should fail. Check min length validation.',
    ).toBe(false);
  });

  it('rejects when newPassword is too short', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'short',
      confirmPassword: 'short',
    });

    expect(
      result.success,
      'newPassword shorter than 8 chars should fail. Check min length validation.',
    ).toBe(false);
  });

  it('rejects when newPassword is same as currentPassword', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'samepass123',
      newPassword: 'samepass123',
      confirmPassword: 'samepass123',
    });

    expect(
      result.success,
      'New password must differ from current password. Check refine().',
    ).toBe(false);
    if (!result.success) {
      const newPwError = result.error.issues.find(
        (i) => i.path.includes('newPassword'),
      );
      expect(
        newPwError?.message,
        'Error should say new password must be different.',
      ).toBe('New password must be different from current password');
    }
  });

  it('rejects when confirmPassword does not match newPassword', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'newpass123',
      confirmPassword: 'different123',
    });

    expect(
      result.success,
      'Mismatched newPassword and confirmPassword should fail. Check refine().',
    ).toBe(false);
    if (!result.success) {
      const confirmError = result.error.issues.find(
        (i) => i.path.includes('confirmPassword'),
      );
      expect(
        confirmError?.message,
        'Error message should indicate passwords do not match.',
      ).toBe('Passwords do not match');
    }
  });

  it('rejects when confirmPassword is empty', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'newpass123',
      confirmPassword: '',
    });

    expect(
      result.success,
      'Empty confirmPassword should fail. Check min(1) validation.',
    ).toBe(false);
  });
});

describe('deleteAccountSchema', () => {
  it('accepts when confirmText is exactly "DELETE"', () => {
    const input: DeleteAccountFormData = { confirmText: 'DELETE' };
    const result = deleteAccountSchema.safeParse(input);

    expect(
      result.success,
      'confirmText "DELETE" should pass validation. Check literal or refine().',
    ).toBe(true);
  });

  it('rejects when confirmText is not "DELETE"', () => {
    const result = deleteAccountSchema.safeParse({ confirmText: 'delete' });

    expect(
      result.success,
      'confirmText must be exactly "DELETE" (case-sensitive). Check refine().',
    ).toBe(false);
  });

  it('rejects when confirmText is empty', () => {
    const result = deleteAccountSchema.safeParse({ confirmText: '' });

    expect(
      result.success,
      'Empty confirmText should fail. Check min(1) validation.',
    ).toBe(false);
  });
});
