import { describe, it, expect } from 'vitest';
import { getErrorMessage } from '@/features/users/utils/errorHelpers';

describe('getErrorMessage', () => {
  it('returns generic message for 5xx errors', () => {
    const err = Object.assign(new Error('Internal Server Error'), { status: 500 });

    expect(
      getErrorMessage(err, 'fallback'),
      'Should return generic message for status >= 500 to avoid leaking server details.',
    ).toBe('Something went wrong. Please try again later.');
  });

  it('returns authorization message for 401 errors', () => {
    const err = Object.assign(new Error('Unauthorized'), { status: 401 });

    expect(
      getErrorMessage(err, 'fallback'),
      'Should return authorization message for 401 status.',
    ).toBe('You are not authorized to perform this action.');
  });

  it('returns authorization message for 403 errors', () => {
    const err = Object.assign(new Error('Forbidden'), { status: 403 });

    expect(
      getErrorMessage(err, 'fallback'),
      'Should return authorization message for 403 status.',
    ).toBe('You are not authorized to perform this action.');
  });

  it('returns error message for other Error instances', () => {
    const err = new Error('Something specific');

    expect(
      getErrorMessage(err, 'fallback'),
      'Should return the Error message when no special status handling applies.',
    ).toBe('Something specific');
  });

  it('returns error message for 4xx errors (not 401/403)', () => {
    const err = Object.assign(new Error('Not Found'), { status: 404 });

    expect(
      getErrorMessage(err, 'fallback'),
      'Should return API error message for 4xx errors that are not 401/403.',
    ).toBe('Not Found');
  });

  it('returns fallback for non-Error values', () => {
    expect(
      getErrorMessage('string error', 'fallback'),
      'Should return fallback for non-Error values.',
    ).toBe('fallback');
  });

  it('returns fallback for null/undefined', () => {
    expect(getErrorMessage(null, 'fallback')).toBe('fallback');
    expect(getErrorMessage(undefined, 'fallback')).toBe('fallback');
  });
});
