import { describe, it, expect } from 'vitest';
import {
  loginSchema,
  signupSchema,
  forgotPasswordSchema,
  resetPasswordSchema,
  sendVerificationSchema,
} from '@/lib/validators';

describe('loginSchema', () => {
  it('accepts valid email and password', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: 'password123' });
    expect(result.success, 'Valid login credentials should pass validation.').toBe(true);
  });

  it('rejects invalid email', () => {
    const result = loginSchema.safeParse({ email: 'not-an-email', password: 'password123' });
    expect(result.success, 'Invalid email format should fail validation.').toBe(false);
  });

  it('rejects empty email', () => {
    const result = loginSchema.safeParse({ email: '', password: 'password123' });
    expect(result.success, 'Empty email should fail validation.').toBe(false);
  });

  it('rejects password shorter than 8 characters', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: 'short' });
    expect(
      result.success,
      'Password under 8 chars should fail. BE MIN_PASSWORD_LENGTH = 8.',
    ).toBe(false);
  });

  it('rejects empty password', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: '' });
    expect(result.success, 'Empty password should fail validation.').toBe(false);
  });
});

describe('signupSchema', () => {
  const validData = {
    email: 'user@example.com',
    password: 'password123',
    passwordConfirm: 'password123',
  };

  it('accepts valid signup data', () => {
    const result = signupSchema.safeParse(validData);
    expect(result.success, 'Valid signup data should pass validation.').toBe(true);
  });

  it('rejects invalid email', () => {
    const result = signupSchema.safeParse({ ...validData, email: 'bad' });
    expect(result.success, 'Invalid email should fail validation.').toBe(false);
  });

  it('rejects short password', () => {
    const result = signupSchema.safeParse({
      ...validData,
      password: 'short',
      passwordConfirm: 'short',
    });
    expect(result.success, 'Password under 8 chars should fail.').toBe(false);
  });

  it('rejects mismatched passwords', () => {
    const result = signupSchema.safeParse({
      ...validData,
      passwordConfirm: 'different123',
    });
    expect(
      result.success,
      'Password and passwordConfirm must match. Use .refine() to check.',
    ).toBe(false);
  });

  it('rejects empty passwordConfirm', () => {
    const result = signupSchema.safeParse({
      ...validData,
      passwordConfirm: '',
    });
    expect(result.success, 'Empty passwordConfirm should fail.').toBe(false);
  });
});

describe('forgotPasswordSchema', () => {
  it('accepts valid email', () => {
    const result = forgotPasswordSchema.safeParse({ email: 'user@example.com' });
    expect(result.success, 'Valid email should pass.').toBe(true);
  });

  it('rejects invalid email', () => {
    const result = forgotPasswordSchema.safeParse({ email: 'bad' });
    expect(result.success, 'Invalid email should fail.').toBe(false);
  });

  it('rejects empty email', () => {
    const result = forgotPasswordSchema.safeParse({ email: '' });
    expect(result.success, 'Empty email should fail.').toBe(false);
  });
});

describe('resetPasswordSchema', () => {
  it('accepts valid password pair', () => {
    const result = resetPasswordSchema.safeParse({
      password: 'newpass123',
      passwordConfirm: 'newpass123',
    });
    expect(result.success, 'Matching passwords >= 8 chars should pass.').toBe(true);
  });

  it('rejects short password', () => {
    const result = resetPasswordSchema.safeParse({
      password: 'short',
      passwordConfirm: 'short',
    });
    expect(result.success, 'Password under 8 chars should fail.').toBe(false);
  });

  it('rejects mismatched passwords', () => {
    const result = resetPasswordSchema.safeParse({
      password: 'newpass123',
      passwordConfirm: 'different1',
    });
    expect(result.success, 'Mismatched passwords should fail.').toBe(false);
  });
});

describe('sendVerificationSchema', () => {
  it('accepts valid email', () => {
    const result = sendVerificationSchema.safeParse({ email: 'user@example.com' });
    expect(result.success, 'Valid email should pass.').toBe(true);
  });

  it('rejects invalid email', () => {
    const result = sendVerificationSchema.safeParse({ email: 'bad' });
    expect(result.success, 'Invalid email should fail.').toBe(false);
  });
});
