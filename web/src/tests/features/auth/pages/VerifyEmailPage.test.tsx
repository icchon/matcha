import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { VerifyEmailPage } from '@/features/auth/pages/VerifyEmailPage';

const mockVerifyEmail = vi.fn();
const mockSendVerificationEmail = vi.fn();

vi.mock('@/features/auth/hooks/useAuth', () => ({
  useAuth: () => ({
    verifyEmail: mockVerifyEmail,
    sendVerificationEmail: mockSendVerificationEmail,
    isLoading: false,
    error: null,
  }),
}));

vi.mock('sonner', () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

beforeEach(() => {
  vi.clearAllMocks();
});

function renderWithRoute(token: string) {
  return render(
    <MemoryRouter initialEntries={[`/verify/${token}`]}>
      <Routes>
        <Route path="/verify/:token" element={<VerifyEmailPage />} />
      </Routes>
    </MemoryRouter>,
  );
}

describe('VerifyEmailPage', () => {
  it('calls verifyEmail with token from URL params on mount', () => {
    renderWithRoute('abc-token-123');

    expect(
      mockVerifyEmail,
      'VerifyEmailPage should call verifyEmail with the token from the URL on mount.',
    ).toHaveBeenCalledWith('abc-token-123');
  });

  it('renders a heading', () => {
    renderWithRoute('abc-token-123');

    expect(screen.getByRole('heading')).toBeInTheDocument();
  });

  it('shows a resend verification form', () => {
    renderWithRoute('abc-token-123');

    expect(
      screen.getByLabelText(/email/i),
      'VerifyEmailPage should have an email input for resending verification.',
    ).toBeInTheDocument();
  });
});
