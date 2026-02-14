import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LoginForm } from '@/features/auth/components/LoginForm';

const mockLogin = vi.fn();

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useAuth: () => ({
    login: mockLogin,
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

describe('LoginForm', () => {
  it('renders email and password fields', () => {
    render(<LoginForm />);

    expect(
      screen.getByLabelText(/email/i),
      'LoginForm should render an email input with a label.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/password/i),
      'LoginForm should render a password input with a label.',
    ).toBeInTheDocument();
  });

  it('renders submit button', () => {
    render(<LoginForm />);

    expect(
      screen.getByRole('button', { name: /log in/i }),
      'LoginForm should have a submit button labeled "Log In".',
    ).toBeInTheDocument();
  });

  it('calls login with form data on valid submission', async () => {
    const user = userEvent.setup();
    render(<LoginForm />);

    await user.type(screen.getByLabelText(/email/i), 'user@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /log in/i }));

    expect(
      mockLogin,
      'LoginForm should call useAuth().login with email and password on submit.',
    ).toHaveBeenCalledWith({
      email: 'user@example.com',
      password: 'password123',
    });
  });

  it('shows validation error for invalid email', async () => {
    const user = userEvent.setup();
    render(<LoginForm />);

    await user.type(screen.getByLabelText(/email/i), 'bad');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /log in/i }));

    expect(
      await screen.findByText(/valid email/i),
      'Should show email validation error for invalid email format.',
    ).toBeInTheDocument();
    expect(mockLogin).not.toHaveBeenCalled();
  });

  it('shows validation error for short password', async () => {
    const user = userEvent.setup();
    render(<LoginForm />);

    await user.type(screen.getByLabelText(/email/i), 'user@example.com');
    await user.type(screen.getByLabelText(/password/i), 'short');
    await user.click(screen.getByRole('button', { name: /log in/i }));

    expect(
      await screen.findByText(/8 characters/i),
      'Should show password validation error when password is too short.',
    ).toBeInTheDocument();
    expect(mockLogin).not.toHaveBeenCalled();
  });

  it('has a link to signup page', () => {
    render(<LoginForm />);

    expect(
      screen.getByRole('link', { name: /sign up/i }),
      'LoginForm should have a link to the signup page.',
    ).toHaveAttribute('href', '/signup');
  });

  it('has a link to forgot password page', () => {
    render(<LoginForm />);

    expect(
      screen.getByRole('link', { name: /forgot/i }),
      'LoginForm should have a link to the forgot password page.',
    ).toHaveAttribute('href', '/forgot-password');
  });
});
