import { describe, it, expect } from 'vitest';
import { API_PATHS } from '@/lib/constants';

describe('API_PATHS', () => {
  describe('AUTH paths', () => {
    it('has correct static auth paths', () => {
      expect(API_PATHS.AUTH.LOGIN, 'Login path should be /auth/login').toBe('/auth/login');
      expect(API_PATHS.AUTH.LOGOUT, 'Logout path should be /auth/logout').toBe('/auth/logout');
      expect(API_PATHS.AUTH.SIGNUP, 'Signup path should be /auth/signup').toBe('/auth/signup');
      expect(
        API_PATHS.AUTH.SEND_VERIFICATION,
        'Send verification path should be /auth/verify/mail to match BE route POST /auth/verify/mail',
      ).toBe('/auth/verify/mail');
      expect(API_PATHS.AUTH.PASSWORD_FORGOT).toBe('/auth/password/forgot');
      expect(API_PATHS.AUTH.PASSWORD_RESET).toBe('/auth/password/reset');
      expect(API_PATHS.AUTH.OAUTH_GOOGLE).toBe('/auth/oauth/google/login');
      expect(API_PATHS.AUTH.OAUTH_GITHUB).toBe('/auth/oauth/github/login');
    });

    it('has correct dynamic auth paths', () => {
      expect(
        API_PATHS.AUTH.VERIFY_EMAIL('abc-123'),
        'Verify email path should interpolate token into /auth/verify/{token}',
      ).toBe('/auth/verify/abc-123');
    });
  });

  describe('USERS paths (current user actions under /me)', () => {
    it('has correct "me" paths under /me prefix', () => {
      expect(
        API_PATHS.USERS.MY_LIKES,
        'MY_LIKES should be /me/likes (not /users/me/likes). Backend route: GET /me/likes',
      ).toBe('/me/likes');
      expect(
        API_PATHS.USERS.MY_VIEWS,
        'MY_VIEWS should be /me/views (not /users/me/views). Backend route: GET /me/views',
      ).toBe('/me/views');
      expect(
        API_PATHS.USERS.DELETE_ME,
        'DELETE_ME should be /me/ (not /users/me). Backend route: DELETE /me/',
      ).toBe('/me/');
      expect(
        API_PATHS.USERS.MY_BLOCKS,
        'MY_BLOCKS should be /me/blocks (not /users/me/blocks). Backend route: GET /me/blocks',
      ).toBe('/me/blocks');
      expect(
        API_PATHS.USERS.MY_DATA,
        'MY_DATA should be /me/data/ (not /users/me/data). Backend route: /me/data/',
      ).toBe('/me/data/');
      expect(
        API_PATHS.USERS.MY_TAGS,
        'MY_TAGS should be /me/tags/ (not /users/me/tags). Backend route: /me/tags/',
      ).toBe('/me/tags/');
    });

    it('has correct dynamic user paths', () => {
      expect(API_PATHS.USERS.LIKE('user-1')).toBe('/users/user-1/like');
      expect(API_PATHS.USERS.UNLIKE('user-1')).toBe('/users/user-1/like');
      expect(API_PATHS.USERS.BLOCK('user-1')).toBe('/users/user-1/block');
      expect(API_PATHS.USERS.UNBLOCK('user-1')).toBe('/users/user-1/block');
      expect(
        API_PATHS.USERS.DELETE_TAG(42),
        'DELETE_TAG should use /me/tags/{tagID} path. Backend route: DELETE /me/tags/{tagID}',
      ).toBe('/me/tags/42');
    });
  });

  describe('PROFILE paths (under /me/profile)', () => {
    it('has correct profile paths under /me/profile prefix', () => {
      expect(
        API_PATHS.PROFILE.CREATE,
        'PROFILE.CREATE should be /me/profile/ (not /profile). Backend route: POST /me/profile/',
      ).toBe('/me/profile/');
      expect(
        API_PATHS.PROFILE.UPDATE,
        'PROFILE.UPDATE should be /me/profile/ (not /profile). Backend route: PUT /me/profile/',
      ).toBe('/me/profile/');
      expect(
        API_PATHS.PROFILE.PICTURES,
        'PROFILE.PICTURES should be /me/profile/pictures (not /profile/pictures). Backend route: POST /me/profile/pictures',
      ).toBe('/me/profile/pictures');
      expect(
        API_PATHS.PROFILE.WHO_LIKED_ME,
        'WHO_LIKED_ME should be /me/profile/likes (not /profile/likes). Backend route: GET /me/profile/likes',
      ).toBe('/me/profile/likes');
      expect(
        API_PATHS.PROFILE.WHO_VIEWED_ME,
        'WHO_VIEWED_ME should be /me/profile/views (not /profile/views). Backend route: GET /me/profile/views',
      ).toBe('/me/profile/views');
    });

    it('has correct dynamic profile paths', () => {
      expect(
        API_PATHS.PROFILE.DELETE_PICTURE(5),
        'DELETE_PICTURE should use /me/profile/pictures/{pictureID} path',
      ).toBe('/me/profile/pictures/5');
      expect(API_PATHS.PROFILE.GET('user-1')).toBe('/users/user-1/profile');
    });

    it('has correct discovery paths', () => {
      expect(API_PATHS.PROFILE.LIST).toBe('/profiles/');
      expect(API_PATHS.PROFILE.RECOMMEND).toBe('/profiles/recommends');
    });
  });

  describe('TAGS path', () => {
    it('has correct tags path', () => {
      expect(API_PATHS.TAGS, 'TAGS should be /tags/ to match backend route GET /tags/').toBe(
        '/tags/',
      );
    });
  });

  describe('CHATS paths', () => {
    it('has correct chat paths', () => {
      expect(
        API_PATHS.CHATS.MY_CHATS,
        'MY_CHATS should be /me/chats. Backend route: GET /me/chats',
      ).toBe('/me/chats');
      expect(
        API_PATHS.CHATS.MESSAGES('user-1'),
        'MESSAGES should be /chats/{userID}/messages/. Backend route: GET /chats/{userID}/messages/',
      ).toBe('/chats/user-1/messages/');
    });
  });

  describe('NOTIFICATIONS path', () => {
    it('has correct notification paths', () => {
      expect(
        API_PATHS.NOTIFICATIONS.LIST,
        'NOTIFICATIONS.LIST should be /me/notifications. Backend route: GET /me/notifications',
      ).toBe('/me/notifications');
    });
  });
});
