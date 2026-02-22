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
      newPassword: 'Newpass1!',
      confirmPassword: 'Newpass1!',
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
      newPassword: 'Newpass1!',
      confirmPassword: 'Newpass1!',
    });

    expect(
      result.success,
      'currentPassword shorter than 8 chars should fail. Check min length validation.',
    ).toBe(false);
  });

  it('rejects when newPassword is too short', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'Short1!',
      confirmPassword: 'Short1!',
    });

    expect(
      result.success,
      'newPassword shorter than 8 chars should fail. Check min length validation.',
    ).toBe(false);
  });

  it('rejects when newPassword lacks uppercase letter', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'newpass1!',
      confirmPassword: 'newpass1!',
    });

    expect(
      result.success,
      'newPassword without uppercase letter should fail. Check .regex(/[A-Z]/) validation.',
    ).toBe(false);
  });

  it('rejects when newPassword lacks a number', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'Newpasss!',
      confirmPassword: 'Newpasss!',
    });

    expect(
      result.success,
      'newPassword without a number should fail. Check .regex(/[0-9]/) validation.',
    ).toBe(false);
  });

  it('rejects when newPassword lacks a special character', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'Newpass12',
      confirmPassword: 'Newpass12',
    });

    expect(
      result.success,
      'newPassword without special character should fail. Check .regex(/[^a-zA-Z0-9]/) validation.',
    ).toBe(false);
  });

  it('rejects when newPassword is same as currentPassword', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'Samepass1!',
      newPassword: 'Samepass1!',
      confirmPassword: 'Samepass1!',
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
      newPassword: 'Newpass1!',
      confirmPassword: 'Different1!',
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

  it('rejects when confirmPassword is too short', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpass123',
      newPassword: 'Newpass1!',
      confirmPassword: 'short',
    });

    expect(
      result.success,
      'confirmPassword shorter than 8 chars should fail. Check min(8) validation.',
    ).toBe(false);
  });
});

describe('deleteAccountSchema', () => {
  it('accepts when confirmText is exactly "DELETE" and currentPassword is provided', () => {
    const input: DeleteAccountFormData = { confirmText: 'DELETE', currentPassword: 'mypass123' };
    const result = deleteAccountSchema.safeParse(input);

    expect(
      result.success,
      'confirmText "DELETE" with currentPassword should pass validation.',
    ).toBe(true);
  });

  it('rejects when confirmText is not "DELETE"', () => {
    const result = deleteAccountSchema.safeParse({ confirmText: 'delete', currentPassword: 'mypass123' });

    expect(
      result.success,
      'confirmText must be exactly "DELETE" (case-sensitive). Check refine().',
    ).toBe(false);
  });

  it('rejects when confirmText is empty', () => {
    const result = deleteAccountSchema.safeParse({ confirmText: '', currentPassword: 'mypass123' });

    expect(
      result.success,
      'Empty confirmText should fail. Check min(1) validation.',
    ).toBe(false);
  });

  it('rejects when currentPassword is empty', () => {
    const result = deleteAccountSchema.safeParse({ confirmText: 'DELETE', currentPassword: '' });

    expect(
      result.success,
      'Empty currentPassword should fail. Check min(1) validation on currentPassword.',
    ).toBe(false);
  });
});
