import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ChangePasswordForm } from '@/features/settings/components/ChangePasswordForm';

const mockChangePassword = vi.fn();

vi.mock('@/features/settings/hooks/useChangePassword', () => ({
  useChangePassword: () => ({
    isLoading: false,
    error: null,
    changePassword: mockChangePassword,
  }),
}));

beforeEach(() => {
  vi.resetAllMocks();
});

describe('ChangePasswordForm', () => {
  it('renders current password, new password, and confirm password fields', () => {
    render(<ChangePasswordForm />);

    expect(
      screen.getByLabelText(/current password/i),
      'Should render a Current Password input field.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/new password/i),
      'Should render a New Password input field.',
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText(/confirm password/i),
      'Should render a Confirm Password input field.',
    ).toBeInTheDocument();
  });

  it('renders a submit button', () => {
    render(<ChangePasswordForm />);

    expect(
      screen.getByRole('button', { name: /change password/i }),
      'Should render a "Change Password" submit button.',
    ).toBeInTheDocument();
  });

  it('shows validation errors for empty fields on submit', async () => {
    const user = userEvent.setup();
    render(<ChangePasswordForm />);

    await user.click(screen.getByRole('button', { name: /change password/i }));

    await waitFor(() => {
      expect(
        screen.getAllByRole('alert').length,
        'Should show validation errors when fields are empty. Check Zod schema integration.',
      ).toBeGreaterThan(0);
    });
  });

  it('shows mismatch error when passwords do not match', async () => {
    const user = userEvent.setup();
    render(<ChangePasswordForm />);

    await user.type(screen.getByLabelText(/current password/i), 'oldpass123');
    await user.type(screen.getByLabelText(/new password/i), 'Newpass1!');
    await user.type(screen.getByLabelText(/confirm password/i), 'Different1!');
    await user.click(screen.getByRole('button', { name: /change password/i }));

    await waitFor(() => {
      expect(
        screen.getByText(/passwords do not match/i),
        'Should show "Passwords do not match" error. Check refine() in schema.',
      ).toBeInTheDocument();
    });
  });

  it('calls changePassword with correct params on valid submit', async () => {
    mockChangePassword.mockResolvedValue(undefined);
    const user = userEvent.setup();
    render(<ChangePasswordForm />);

    await user.type(screen.getByLabelText(/current password/i), 'oldpass123');
    await user.type(screen.getByLabelText(/new password/i), 'Newpass1!');
    await user.type(screen.getByLabelText(/confirm password/i), 'Newpass1!');
    await user.click(screen.getByRole('button', { name: /change password/i }));

    await waitFor(() => {
      expect(
        mockChangePassword,
        'Should call changePassword hook with currentPassword and newPassword.',
      ).toHaveBeenCalledWith({
        currentPassword: 'oldpass123',
        newPassword: 'Newpass1!',
      });
    });
  });
});
