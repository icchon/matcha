import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { ResetPasswordPage } from '@/features/auth/pages/ResetPasswordPage';

const mockResetPassword = vi.fn();

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useAuth: () => ({
    resetPassword: mockResetPassword,
    isLoading: false,
    error: null,
  }),
}));

describe('ResetPasswordPage', () => {
  it('renders password fields and submit button', () => {
    render(
      <MemoryRouter initialEntries={['/reset-password?token=abc']}>
        <ResetPasswordPage />
      </MemoryRouter>,
    );
    expect(screen.getByLabelText(/^new password$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/confirm/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /reset/i })).toBeInTheDocument();
  });

  it('calls resetPassword with token from query param and new password', async () => {
    const user = userEvent.setup();
    render(
      <MemoryRouter initialEntries={['/reset-password?token=abc-token']}>
        <ResetPasswordPage />
      </MemoryRouter>,
    );

    await user.type(screen.getByLabelText(/^new password$/i), 'newpass123');
    await user.type(screen.getByLabelText(/confirm/i), 'newpass123');
    await user.click(screen.getByRole('button', { name: /reset/i }));

    expect(
      mockResetPassword,
      'ResetPasswordPage should extract token from ?token= query param.',
    ).toHaveBeenCalledWith({ token: 'abc-token', password: 'newpass123' });
  });

  it('renders a heading', () => {
    render(
      <MemoryRouter initialEntries={['/reset-password?token=abc']}>
        <ResetPasswordPage />
      </MemoryRouter>,
    );
    expect(screen.getByRole('heading', { name: /reset password/i })).toBeInTheDocument();
  });

  it('shows invalid link message when token is missing', () => {
    render(
      <MemoryRouter initialEntries={['/reset-password']}>
        <ResetPasswordPage />
      </MemoryRouter>,
    );
    expect(
      screen.getByText(/invalid or has expired/i),
      'Should show error when no token query param is present.',
    ).toBeInTheDocument();
    expect(
      screen.getByRole('link', { name: /new reset link/i }),
      'Should link to forgot-password page.',
    ).toHaveAttribute('href', '/forgot-password');
  });
});
