/**
 * Maps raw errors to user-facing messages.
 * 5xx server errors are replaced with a generic fallback to avoid leaking internal details.
 */
export function toUserFacingMessage(err: unknown, fallback: string): string {
  if (err instanceof Error) {
    const apiErr = err as Error & { status?: number };
    if (typeof apiErr.status === 'number' && apiErr.status >= 500) {
      return fallback;
    }
    return apiErr.message;
  }
  return fallback;
}
