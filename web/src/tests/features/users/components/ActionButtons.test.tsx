import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ActionButtons } from '@/features/users/components/ActionButtons';

describe('ActionButtons', () => {
  const defaultProps = {
    userId: 'user-1',
    isLiked: false,
    isBlocked: false,
    onLike: vi.fn(),
    onUnlike: vi.fn(),
    onBlock: vi.fn(),
    onUnblock: vi.fn(),
    onReport: vi.fn(),
  };

  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders Like button when not liked', () => {
    render(<ActionButtons {...defaultProps} />);

    expect(
      screen.getByRole('button', { name: /like/i }),
      'Should show Like button when isLiked is false.',
    ).toBeInTheDocument();
  });

  it('renders Unlike button when already liked', () => {
    render(<ActionButtons {...defaultProps} isLiked={true} />);

    expect(
      screen.getByRole('button', { name: /unlike/i }),
      'Should show Unlike button when isLiked is true.',
    ).toBeInTheDocument();
  });

  it('calls onLike when Like button is clicked', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} />);

    await user.click(screen.getByRole('button', { name: /like/i }));

    expect(
      defaultProps.onLike,
      'onLike should be called when Like button is clicked.',
    ).toHaveBeenCalled();
  });

  it('calls onUnlike when Unlike button is clicked', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} isLiked={true} />);

    await user.click(screen.getByRole('button', { name: /unlike/i }));

    expect(
      defaultProps.onUnlike,
      'onUnlike should be called when Unlike button is clicked.',
    ).toHaveBeenCalled();
  });

  it('renders Block button when not blocked', () => {
    render(<ActionButtons {...defaultProps} />);

    expect(
      screen.getByRole('button', { name: /block/i }),
      'Should show Block button when isBlocked is false.',
    ).toBeInTheDocument();
  });

  it('renders Unblock button when blocked', () => {
    render(<ActionButtons {...defaultProps} isBlocked={true} />);

    expect(
      screen.getByRole('button', { name: /unblock/i }),
      'Should show Unblock button when isBlocked is true.',
    ).toBeInTheDocument();
  });

  it('shows confirmation modal before blocking', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} />);

    await user.click(screen.getByRole('button', { name: /block/i }));

    const dialog = screen.getByRole('dialog');
    expect(
      dialog,
      'Should show a confirmation modal when Block is clicked.',
    ).toBeInTheDocument();
    expect(within(dialog).getByText(/are you sure/i)).toBeInTheDocument();
  });

  it('calls onBlock when block is confirmed', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} />);

    await user.click(screen.getByRole('button', { name: /block/i }));
    await user.click(screen.getByRole('button', { name: /confirm/i }));

    expect(
      defaultProps.onBlock,
      'onBlock should be called after confirming the block modal.',
    ).toHaveBeenCalled();
  });

  it('does not call onBlock when block modal is cancelled', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} />);

    await user.click(screen.getByRole('button', { name: /block/i }));
    await user.click(screen.getByRole('button', { name: /cancel/i }));

    expect(
      defaultProps.onBlock,
      'onBlock should NOT be called when cancel is clicked.',
    ).not.toHaveBeenCalled();
  });

  it('renders Report button', () => {
    render(<ActionButtons {...defaultProps} />);

    expect(
      screen.getByRole('button', { name: /report/i }),
      'Should show Report button.',
    ).toBeInTheDocument();
  });

  it('shows confirmation modal before reporting', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} />);

    await user.click(screen.getByRole('button', { name: /report/i }));

    const dialog = screen.getByRole('dialog');
    expect(
      dialog,
      'Should show a confirmation modal when Report is clicked.',
    ).toBeInTheDocument();
  });

  it('calls onReport when report is confirmed', async () => {
    const user = userEvent.setup();
    render(<ActionButtons {...defaultProps} />);

    await user.click(screen.getByRole('button', { name: /report/i }));
    await user.click(screen.getByRole('button', { name: /confirm/i }));

    expect(
      defaultProps.onReport,
      'onReport should be called after confirming the report modal.',
    ).toHaveBeenCalled();
  });

  it('disables buttons when isLoading is true', () => {
    render(<ActionButtons {...defaultProps} isLoading={true} />);

    const buttons = screen.getAllByRole('button');
    buttons.forEach((btn) => {
      expect(btn).toBeDisabled();
    });
  });
});
