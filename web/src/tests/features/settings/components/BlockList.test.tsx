import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BlockList } from '@/features/settings/components/BlockList';
import type { Block } from '@/types';

const mockFetchBlockList = vi.fn();
const mockUnblock = vi.fn();

const blocks: Block[] = [
  { blockerId: 'me', blockedId: 'user-2' },
  { blockerId: 'me', blockedId: 'user-3' },
];

let mockBlocks: Block[] = [];
let mockIsLoading = false;
let mockError: string | null = null;
let mockUnblockingId: string | null = null;

vi.mock('@/features/settings/hooks/useBlockList', () => ({
  useBlockList: () => ({
    blocks: mockBlocks,
    isLoading: mockIsLoading,
    error: mockError,
    unblockingId: mockUnblockingId,
    fetchBlockList: mockFetchBlockList,
    unblock: mockUnblock,
  }),
}));

beforeEach(() => {
  vi.resetAllMocks();
  mockBlocks = [];
  mockIsLoading = false;
  mockError = null;
  mockUnblockingId = null;
});

describe('BlockList', () => {
  it('calls fetchBlockList on mount', () => {
    render(<BlockList />);

    expect(
      mockFetchBlockList,
      'fetchBlockList should be called when BlockList mounts.',
    ).toHaveBeenCalled();
  });

  it('shows loading spinner when isLoading is true', () => {
    mockIsLoading = true;
    render(<BlockList />);

    expect(
      screen.getByRole('status'),
      'Should show a loading spinner when isLoading is true.',
    ).toBeInTheDocument();
  });

  it('shows empty message when no blocked users and no error', () => {
    mockBlocks = [];
    render(<BlockList />);

    expect(
      screen.getByText(/no blocked users/i),
      'Should show "No blocked users" when list is empty and no error.',
    ).toBeInTheDocument();
  });

  it('does not show empty message when error exists', () => {
    mockBlocks = [];
    mockError = 'Failed to load';
    render(<BlockList />);

    expect(
      screen.queryByText(/no blocked users/i),
      'Should not show "No blocked users" when there is an error.',
    ).not.toBeInTheDocument();
    expect(screen.getByText('Failed to load')).toBeInTheDocument();
  });

  it('renders list of blocked users with unblock buttons with aria-labels', () => {
    mockBlocks = blocks;
    render(<BlockList />);

    expect(
      screen.getByText('user-2'),
      'Should display blocked user IDs.',
    ).toBeInTheDocument();
    expect(screen.getByText('user-3')).toBeInTheDocument();

    const unblockButtons = screen.getAllByRole('button', { name: /unblock/i });
    expect(
      unblockButtons,
      'Should render an Unblock button for each blocked user.',
    ).toHaveLength(2);
    expect(
      screen.getByRole('button', { name: 'Unblock user user-2' }),
      'Unblock button should have aria-label identifying which user will be unblocked.',
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Unblock user user-3' }),
      'Unblock button should have aria-label identifying which user will be unblocked.',
    ).toBeInTheDocument();
  });

  it('calls unblock when unblock button is clicked', async () => {
    mockBlocks = blocks;
    mockUnblock.mockResolvedValue(undefined);
    const user = userEvent.setup();
    render(<BlockList />);

    const unblockButtons = screen.getAllByRole('button', { name: /unblock/i });
    await user.click(unblockButtons[0]);

    await waitFor(() => {
      expect(
        mockUnblock,
        'unblock should be called with the blockedId of the first user.',
      ).toHaveBeenCalledWith('user-2');
    });
  });

  it('shows error message when error exists', () => {
    mockError = 'Failed to load';
    render(<BlockList />);

    expect(
      screen.getByText('Failed to load'),
      'Should display the error message.',
    ).toBeInTheDocument();
  });

  it('disables unblock button for the row being unblocked', () => {
    mockBlocks = blocks;
    mockUnblockingId = 'user-2';
    render(<BlockList />);

    const unblockButtons = screen.getAllByRole('button', { name: /unblock/i });
    expect(
      unblockButtons[0],
      'Unblock button for the row being unblocked should be disabled.',
    ).toBeDisabled();
    expect(
      unblockButtons[1],
      'Unblock button for other rows should remain enabled.',
    ).not.toBeDisabled();
  });
});
