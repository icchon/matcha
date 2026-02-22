interface ApiError {
  readonly status?: number;
  readonly message?: string;
}

export function getErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof Error) {
    const apiErr = err as Error & ApiError;
    if (typeof apiErr.status === 'number') {
      if (apiErr.status >= 500) return 'Something went wrong. Please try again later.';
      if (apiErr.status === 401 || apiErr.status === 403) return 'You are not authorized to perform this action.';
    }
    return apiErr.message;
  }
  return fallback;
}
