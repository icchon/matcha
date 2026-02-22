import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  deleteAccount,
  changePassword,
  getBlockList,
  unblockUser,
} from '@/api/settings';
import { apiClient } from '@/api/client';
import { API_PATHS } from '@/lib/constants';
import type { Block, MessageResponse } from '@/types';

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
  ApiClientError: class extends Error {
    readonly status: number;
    constructor(status: number, body: { error: string }) {
      super(body.error);
      this.status = status;
    }
  },
}));

const mockGet = vi.mocked(apiClient.get);
const mockPost = vi.mocked(apiClient.post);
const mockDelete = vi.mocked(apiClient.delete);

beforeEach(() => {
  vi.resetAllMocks();
});

describe('deleteAccount', () => {
  it('calls POST /me/delete with currentPassword to delete the current user account', async () => {
    mockPost.mockResolvedValue(undefined);

    await deleteAccount({ currentPassword: 'mypass123' });

    expect(
      mockPost,
      'deleteAccount should call POST on USERS.DELETE_ACCOUNT path with currentPassword.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.DELETE_ACCOUNT, { currentPassword: 'mypass123' });
  });
});

describe('changePassword', () => {
  it('calls POST AUTH.CHANGE_PASSWORD with currentPassword and newPassword', async () => {
    const response: MessageResponse = { message: 'Password changed' };
    mockPost.mockResolvedValue(response);

    const result = await changePassword({
      currentPassword: 'oldpass123',
      newPassword: 'newpass123',
    });

    expect(
      mockPost,
      'changePassword should POST to AUTH.CHANGE_PASSWORD with current and new passwords.',
    ).toHaveBeenCalledWith(API_PATHS.AUTH.CHANGE_PASSWORD, {
      currentPassword: 'oldpass123',
      newPassword: 'newpass123',
    });
    expect(
      result.message,
      'changePassword should return the server message. Check return type.',
    ).toBe('Password changed');
  });
});

describe('getBlockList', () => {
  it('calls GET /me/blocks and returns list of blocked users', async () => {
    const blocks: Block[] = [
      { blockerId: 'user-1', blockedId: 'user-2' },
      { blockerId: 'user-1', blockedId: 'user-3' },
    ];
    mockGet.mockResolvedValue(blocks);

    const result = await getBlockList();

    expect(
      mockGet,
      'getBlockList should call GET on USERS.MY_BLOCKS path.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.MY_BLOCKS);
    expect(
      result,
      'getBlockList should return the array of Block objects from the API.',
    ).toEqual(blocks);
  });
});

describe('unblockUser', () => {
  it('[MOCK] calls DELETE /users/:userId/block and returns success stub', async () => {
    mockDelete.mockResolvedValue(undefined);

    await unblockUser('user-2');

    expect(
      mockDelete,
      'unblockUser should call DELETE on USERS.UNBLOCK(userId) path.',
    ).toHaveBeenCalledWith(API_PATHS.USERS.UNBLOCK('user-2'));
  });
});
