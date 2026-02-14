import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { SignupPage } from '@/features/auth/pages/SignupPage';

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useAuth: () => ({
    signup: vi.fn(),
    isLoading: false,
    error: null,
  }),
}));

vi.mock('sonner', () => ({
  toast: { info: vi.fn() },
}));

describe('SignupPage', () => {
  it('renders the signup form', () => {
    render(
      <MemoryRouter>
        <SignupPage />
      </MemoryRouter>,
    );
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
  });

  it('renders a page heading', () => {
    render(
      <MemoryRouter>
        <SignupPage />
      </MemoryRouter>,
    );
    expect(
      screen.getByRole('heading', { name: /sign up/i }),
      'SignupPage should have a heading.',
    ).toBeInTheDocument();
  });

  it('renders OAuth buttons', () => {
    render(
      <MemoryRouter>
        <SignupPage />
      </MemoryRouter>,
    );
    expect(screen.getByRole('button', { name: /google/i })).toBeInTheDocument();
  });
});
