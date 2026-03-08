import { describe, it, expect } from 'vitest';
import { toUserFacingMessage } from '@/lib/errorUtils';

describe('toUserFacingMessage', () => {
  it('returns error message for client errors (4xx)', () => {
    const err = Object.assign(new Error('Not Found'), { status: 404 });

    expect(
      toUserFacingMessage(err, 'Something went wrong'),
      'Client errors (4xx) should expose the original message, not the fallback.',
    ).toBe('Not Found');
  });

  it('returns fallback for 5xx server errors', () => {
    const err = Object.assign(new Error('Internal Server Error'), { status: 500 });

    expect(
      toUserFacingMessage(err, 'Something went wrong'),
      '5xx errors should return the fallback to avoid leaking server details.',
    ).toBe('Something went wrong');
  });

  it('returns fallback for 502 Bad Gateway', () => {
    const err = Object.assign(new Error('Bad Gateway'), { status: 502 });

    expect(
      toUserFacingMessage(err, 'Service unavailable'),
      '502 is a 5xx error and should return fallback.',
    ).toBe('Service unavailable');
  });

  it('returns error message when no status is present', () => {
    const err = new Error('Network timeout');

    expect(
      toUserFacingMessage(err, 'Something went wrong'),
      'Errors without status should expose the original message.',
    ).toBe('Network timeout');
  });

  it('returns fallback for non-Error values', () => {
    expect(
      toUserFacingMessage('string error', 'Something went wrong'),
      'Non-Error values (string, number, etc.) should return fallback.',
    ).toBe('Something went wrong');
  });

  it('returns fallback for null/undefined', () => {
    expect(
      toUserFacingMessage(null, 'Fallback'),
      'null should return fallback.',
    ).toBe('Fallback');

    expect(
      toUserFacingMessage(undefined, 'Fallback'),
      'undefined should return fallback.',
    ).toBe('Fallback');
  });
});
