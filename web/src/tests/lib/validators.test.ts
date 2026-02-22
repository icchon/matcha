import { describe, it, expect } from 'vitest';
import {
  loginSchema,
  signupSchema,
  forgotPasswordSchema,
  resetPasswordSchema,
  sendVerificationSchema,
  profileSchema,
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

describe('profileSchema', () => {
  const validProfile = {
    firstName: 'John',
    lastName: 'Doe',
    username: 'johndoe_123',
    gender: 'male' as const,
    sexualPreference: 'heterosexual' as const,
    birthday: '2000-01-01',
    biography: 'Hello world',
  };

  it('accepts a valid profile', () => {
    const result = profileSchema.safeParse(validProfile);
    expect(result.success, 'Valid profile data should pass validation.').toBe(true);
  });

  describe('username', () => {
    it('accepts letters, numbers, underscores, and hyphens', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'user_name-123' });
      expect(result.success, 'Username with allowed chars should pass.').toBe(true);
    });

    it('rejects special characters', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'user@name!' });
      expect(
        result.success,
        'Username with special chars should fail. Regex: /^[a-zA-Z0-9_-]+$/',
      ).toBe(false);
    });

    it('rejects spaces', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'user name' });
      expect(result.success, 'Username with spaces should fail.').toBe(false);
    });

    it('rejects username shorter than 3 characters', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'ab' });
      expect(result.success, 'Username under 3 chars should fail. Min length is 3.').toBe(false);
    });

    it('accepts username at exactly 3 characters', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'abc' });
      expect(result.success, 'Username at min boundary (3) should pass.').toBe(true);
    });

    it('accepts username at exactly 30 characters', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'a'.repeat(30) });
      expect(result.success, 'Username at max boundary (30) should pass.').toBe(true);
    });

    it('rejects username longer than 30 characters', () => {
      const result = profileSchema.safeParse({ ...validProfile, username: 'a'.repeat(31) });
      expect(result.success, 'Username over 30 chars should fail. Max length is 30.').toBe(false);
    });
  });

  describe('birthday age validation', () => {
    it('rejects a user under 18 years old', () => {
      const today = new Date();
      const under18 = new Date(today.getFullYear() - 17, today.getMonth(), today.getDate());
      const result = profileSchema.safeParse({
        ...validProfile,
        birthday: under18.toISOString().split('T')[0],
      });
      expect(
        result.success,
        'Birthday indicating age < 18 should fail. Minimum age is 18.',
      ).toBe(false);
    });

    it('accepts a user exactly 18 years old', () => {
      const today = new Date();
      const exactly18 = new Date(today.getFullYear() - 18, today.getMonth(), today.getDate());
      const result = profileSchema.safeParse({
        ...validProfile,
        birthday: exactly18.toISOString().split('T')[0],
      });
      expect(
        result.success,
        'Birthday indicating exactly 18 years should pass.',
      ).toBe(true);
    });

    it('accepts a user over 18 years old', () => {
      const result = profileSchema.safeParse({
        ...validProfile,
        birthday: '1990-05-15',
      });
      expect(result.success, 'Birthday indicating age > 18 should pass.').toBe(true);
    });

    it('rejects an invalid date string', () => {
      const result = profileSchema.safeParse({
        ...validProfile,
        birthday: 'not-a-date',
      });
      expect(
        result.success,
        'Invalid date string should fail the age refine check.',
      ).toBe(false);
    });
  });

  describe('occupation', () => {
    it('accepts when occupation is omitted', () => {
      const { occupation: _, ...withoutOccupation } = validProfile;
      const result = profileSchema.safeParse(withoutOccupation);
      expect(result.success, 'Occupation is optional and should pass when omitted.').toBe(true);
    });

    it('accepts empty string for occupation', () => {
      const result = profileSchema.safeParse({ ...validProfile, occupation: '' });
      expect(result.success, 'Empty string occupation should pass.').toBe(true);
    });

    it('accepts a valid occupation string', () => {
      const result = profileSchema.safeParse({ ...validProfile, occupation: 'Engineer' });
      expect(result.success, 'Valid occupation string should pass.').toBe(true);
    });

    it('rejects occupation longer than 100 characters', () => {
      const result = profileSchema.safeParse({ ...validProfile, occupation: 'x'.repeat(101) });
      expect(
        result.success,
        'Occupation over 100 chars should fail. Max length is 100.',
      ).toBe(false);
    });
  });
});
