import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { OAuthButtons } from '@/features/auth/components/OAuthButtons';

vi.mock('sonner', () => ({
  toast: {
    info: vi.fn(),
  },
}));

const { toast } = await import('sonner');

describe('OAuthButtons', () => {
  it('renders Google and GitHub buttons', () => {
    render(<OAuthButtons />);

    expect(
      screen.getByRole('button', { name: /google/i }),
      'Should render a Google OAuth button.',
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: /github/i }),
      'Should render a GitHub OAuth button.',
    ).toBeInTheDocument();
  });

  it('shows "coming soon" toast when Google button is clicked', async () => {
    const user = userEvent.setup();
    render(<OAuthButtons />);

    await user.click(screen.getByRole('button', { name: /google/i }));

    expect(
      toast.info,
      '[MOCK] Google OAuth should show a "coming soon" toast.',
    ).toHaveBeenCalled();
  });

  it('shows "coming soon" toast when GitHub button is clicked', async () => {
    const user = userEvent.setup();
    render(<OAuthButtons />);

    await user.click(screen.getByRole('button', { name: /github/i }));

    expect(
      toast.info,
      '[MOCK] GitHub OAuth should show a "coming soon" toast.',
    ).toHaveBeenCalled();
  });
});
