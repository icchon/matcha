import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { DeleteAccountSection } from '@/features/settings/components/DeleteAccountSection';

const mockDeleteAccount = vi.fn();

vi.mock('@/features/settings/hooks/useDeleteAccount', () => ({
  useDeleteAccount: () => ({
    isLoading: false,
    error: null,
    deleteAccount: mockDeleteAccount,
  }),
}));

beforeEach(() => {
  vi.resetAllMocks();
});

describe('DeleteAccountSection', () => {
  it('renders danger zone heading and delete button', () => {
    render(<DeleteAccountSection />);

    expect(
      screen.getByText(/danger zone/i),
      'Should render a "Danger Zone" heading.',
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: /delete account/i }),
      'Should render a "Delete Account" button.',
    ).toBeInTheDocument();
  });

  it('opens confirmation modal when delete button is clicked', async () => {
    const user = userEvent.setup();
    render(<DeleteAccountSection />);

    await user.click(screen.getByRole('button', { name: /delete account/i }));

    expect(
      screen.getByRole('dialog'),
      'Should open a confirmation modal dialog.',
    ).toBeInTheDocument();
    expect(
      screen.getByText(/type DELETE to confirm/i),
      'Modal should instruct user to type DELETE.',
    ).toBeInTheDocument();
  });

  it('does not call deleteAccount when confirm text is wrong', async () => {
    const user = userEvent.setup();
    render(<DeleteAccountSection />);

    await user.click(screen.getByRole('button', { name: /delete account/i }));
    await user.type(screen.getByLabelText(/confirm/i), 'wrong');
    await user.click(screen.getByRole('button', { name: /confirm delete/i }));

    await waitFor(() => {
      expect(
        mockDeleteAccount,
        'deleteAccount should not be called when confirmText is not "DELETE".',
      ).not.toHaveBeenCalled();
    });
  });

  it('calls deleteAccount when "DELETE" is typed and confirmed', async () => {
    mockDeleteAccount.mockResolvedValue(undefined);
    const user = userEvent.setup();
    render(<DeleteAccountSection />);

    await user.click(screen.getByRole('button', { name: /delete account/i }));
    await user.type(screen.getByLabelText(/confirm/i), 'DELETE');
    await user.click(screen.getByRole('button', { name: /confirm delete/i }));

    await waitFor(() => {
      expect(
        mockDeleteAccount,
        'deleteAccount should be called after typing "DELETE" and confirming.',
      ).toHaveBeenCalled();
    });
  });

  it('closes modal when cancel is clicked', async () => {
    const user = userEvent.setup();
    render(<DeleteAccountSection />);

    await user.click(screen.getByRole('button', { name: /delete account/i }));
    expect(screen.getByRole('dialog')).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: /cancel/i }));

    await waitFor(() => {
      expect(
        screen.queryByRole('dialog'),
        'Modal should close when Cancel is clicked.',
      ).not.toBeInTheDocument();
    });
  });
});
