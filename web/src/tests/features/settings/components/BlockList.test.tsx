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

vi.mock('@/features/settings/hooks/useSettings', () => ({
  useBlockList: () => ({
    blocks: mockBlocks,
    isLoading: mockIsLoading,
    error: mockError,
    fetchBlockList: mockFetchBlockList,
    unblock: mockUnblock,
  }),
}));

beforeEach(() => {
  vi.resetAllMocks();
  mockBlocks = [];
  mockIsLoading = false;
  mockError = null;
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

  it('shows empty message when no blocked users', () => {
    mockBlocks = [];
    render(<BlockList />);

    expect(
      screen.getByText(/no blocked users/i),
      'Should show "No blocked users" when list is empty.',
    ).toBeInTheDocument();
  });

  it('renders list of blocked users with unblock buttons', () => {
    mockBlocks = blocks;
    render(<BlockList />);

    expect(
      screen.getByText('user-2'),
      'Should display blocked user IDs.',
    ).toBeInTheDocument();
    expect(screen.getByText('user-3')).toBeInTheDocument();
    expect(
      screen.getAllByRole('button', { name: /unblock/i }),
      'Should render an Unblock button for each blocked user.',
    ).toHaveLength(2);
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
});
