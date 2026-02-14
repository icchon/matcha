import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SignupForm } from '@/features/auth/components/SignupForm';

const mockSignup = vi.fn();

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useAuth: () => ({
    signup: mockSignup,
    isLoading: false,
    error: null,
  }),
}));

vi.mock('react-router-dom', () => ({
  useNavigate: () => vi.fn(),
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => <a href={to}>{children}</a>,
}));

beforeEach(() => {
  vi.clearAllMocks();
});

describe('SignupForm', () => {
  it('renders email, password, and confirm password fields', () => {
    render(<SignupForm />);

    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
  });

  it('renders submit button', () => {
    render(<SignupForm />);

    expect(
      screen.getByRole('button', { name: /sign up/i }),
      'SignupForm should have a submit button labeled "Sign Up".',
    ).toBeInTheDocument();
  });

  it('calls signup with email and password on valid submission', async () => {
    const user = userEvent.setup();
    render(<SignupForm />);

    await user.type(screen.getByLabelText(/email/i), 'new@example.com');
    await user.type(screen.getByLabelText(/^password$/i), 'password123');
    await user.type(screen.getByLabelText(/confirm password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign up/i }));

    expect(
      mockSignup,
      'SignupForm should call signup with email and password (not passwordConfirm).',
    ).toHaveBeenCalledWith({
      email: 'new@example.com',
      password: 'password123',
    });
  });

  it('shows error when passwords do not match', async () => {
    const user = userEvent.setup();
    render(<SignupForm />);

    await user.type(screen.getByLabelText(/email/i), 'new@example.com');
    await user.type(screen.getByLabelText(/^password$/i), 'password123');
    await user.type(screen.getByLabelText(/confirm password/i), 'different99');
    await user.click(screen.getByRole('button', { name: /sign up/i }));

    expect(
      await screen.findByText(/do not match/i),
      'Should show password mismatch error.',
    ).toBeInTheDocument();
    expect(mockSignup).not.toHaveBeenCalled();
  });

  it('shows error for short password', async () => {
    const user = userEvent.setup();
    render(<SignupForm />);

    await user.type(screen.getByLabelText(/email/i), 'new@example.com');
    await user.type(screen.getByLabelText(/^password$/i), 'short');
    await user.type(screen.getByLabelText(/confirm password/i), 'short');
    await user.click(screen.getByRole('button', { name: /sign up/i }));

    expect(
      await screen.findByText(/8 characters/i),
      'Should show password length error.',
    ).toBeInTheDocument();
    expect(mockSignup).not.toHaveBeenCalled();
  });

  it('has a link to login page', () => {
    render(<SignupForm />);

    expect(
      screen.getByRole('link', { name: /log in/i }),
      'SignupForm should have a link to the login page.',
    ).toHaveAttribute('href', '/login');
  });
});
