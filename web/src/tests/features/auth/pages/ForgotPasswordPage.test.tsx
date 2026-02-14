import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { ForgotPasswordPage } from '@/features/auth/pages/ForgotPasswordPage';

const mockForgotPassword = vi.fn();

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useAuth: () => ({
    forgotPassword: mockForgotPassword,
    isLoading: false,
    error: null,
  }),
}));

describe('ForgotPasswordPage', () => {
  it('renders email input and submit button', () => {
    render(
      <MemoryRouter>
        <ForgotPasswordPage />
      </MemoryRouter>,
    );
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /send/i })).toBeInTheDocument();
  });

  it('calls forgotPassword with email on submit', async () => {
    const user = userEvent.setup();
    render(
      <MemoryRouter>
        <ForgotPasswordPage />
      </MemoryRouter>,
    );

    await user.type(screen.getByLabelText(/email/i), 'user@example.com');
    await user.click(screen.getByRole('button', { name: /send/i }));

    expect(mockForgotPassword).toHaveBeenCalledWith({ email: 'user@example.com' });
  });

  it('renders a heading', () => {
    render(
      <MemoryRouter>
        <ForgotPasswordPage />
      </MemoryRouter>,
    );
    expect(screen.getByRole('heading', { name: /forgot password/i })).toBeInTheDocument();
  });
});
