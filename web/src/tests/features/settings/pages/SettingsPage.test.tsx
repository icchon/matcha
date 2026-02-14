import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { SettingsPage } from '@/features/settings/pages/SettingsPage';

vi.mock('@/features/settings/components/ChangePasswordForm', () => ({
  ChangePasswordForm: () => <div data-testid="change-password-form">ChangePasswordForm</div>,
}));

vi.mock('@/features/settings/components/BlockList', () => ({
  BlockList: () => <div data-testid="block-list">BlockList</div>,
}));

vi.mock('@/features/settings/components/DeleteAccountSection', () => ({
  DeleteAccountSection: () => <div data-testid="delete-account-section">DeleteAccountSection</div>,
}));

describe('SettingsPage', () => {
  it('renders the Settings heading', () => {
    render(<SettingsPage />);

    expect(
      screen.getByRole('heading', { name: /settings/i }),
      'Should render a "Settings" heading.',
    ).toBeInTheDocument();
  });

  it('renders Change Password section', () => {
    render(<SettingsPage />);

    expect(
      screen.getByText(/change password/i),
      'Should render the "Change Password" section heading.',
    ).toBeInTheDocument();
    expect(screen.getByTestId('change-password-form')).toBeInTheDocument();
  });

  it('renders Blocked Users section', () => {
    render(<SettingsPage />);

    expect(
      screen.getByText(/blocked users/i),
      'Should render the "Blocked Users" section heading.',
    ).toBeInTheDocument();
    expect(screen.getByTestId('block-list')).toBeInTheDocument();
  });

  it('renders Delete Account section at the bottom', () => {
    render(<SettingsPage />);

    expect(screen.getByTestId('delete-account-section')).toBeInTheDocument();
  });

  it('renders sections in correct order: password, blocks, delete', () => {
    render(<SettingsPage />);

    const sections = screen.getAllByTestId(/.+/);
    const ids = sections.map((el) => el.getAttribute('data-testid'));

    const passwordIdx = ids.indexOf('change-password-form');
    const blockIdx = ids.indexOf('block-list');
    const deleteIdx = ids.indexOf('delete-account-section');

    expect(
      passwordIdx,
      'Change Password should appear before Block List. Check section order.',
    ).toBeLessThan(blockIdx);
    expect(
      blockIdx,
      'Block List should appear before Delete Account. Check section order.',
    ).toBeLessThan(deleteIdx);
  });
});
