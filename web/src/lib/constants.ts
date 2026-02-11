export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'matcha_access_token',
  REFRESH_TOKEN: 'matcha_refresh_token',
} as const;

export const API_PATHS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    SIGNUP: '/auth/signup',
    VERIFY_EMAIL: (token: string) => `/auth/verify/${token}`,
    SEND_VERIFICATION: '/auth/verify',
    PASSWORD_FORGOT: '/auth/password/forgot',
    PASSWORD_RESET: '/auth/password/reset',
    OAUTH_GOOGLE: '/auth/oauth/google/login',
    OAUTH_GITHUB: '/auth/oauth/github/login',
  },
  USERS: {
    LIKE: (userId: string) => `/users/${userId}/like`,
    UNLIKE: (userId: string) => `/users/${userId}/like`,
    MY_LIKES: '/users/me/likes',
    MY_VIEWS: '/users/me/views',
    DELETE_ME: '/users/me',
    MY_BLOCKS: '/users/me/blocks',
    BLOCK: (userId: string) => `/users/${userId}/block`,
    UNBLOCK: (userId: string) => `/users/${userId}/block`,
    MY_DATA: '/users/me/data',
    MY_TAGS: '/users/me/tags',
    DELETE_TAG: (tagId: number) => `/users/me/tags/${tagId}`,
  },
  TAGS: '/tags',
  PROFILE: {
    CREATE: '/profile',
    UPDATE: '/profile',
    PICTURES: '/profile/pictures',
    DELETE_PICTURE: (pictureId: number) => `/profile/pictures/${pictureId}`,
    WHO_LIKED_ME: '/profile/likes',
    WHO_VIEWED_ME: '/profile/views',
    LIST: '/profiles',
    GET: (userId: string) => `/users/${userId}/profile`,
    RECOMMEND: '/profiles/recommends',
  },
} as const;
