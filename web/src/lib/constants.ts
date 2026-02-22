export const API_PATHS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    SIGNUP: '/auth/signup',
    VERIFY_EMAIL: (token: string) => `/auth/verify/${token}`,
    SEND_VERIFICATION: '/auth/verify/mail',
    PASSWORD_FORGOT: '/auth/password/forgot',
    PASSWORD_RESET: '/auth/password/reset',
    // TODO(BE-XX): Move to /auth/password/change once backend supports a dedicated authenticated endpoint.
    // Currently shares endpoint with PASSWORD_RESET; backend differentiates by presence of currentPassword vs token.
    // Backend MUST validate Bearer token when currentPassword is present in the request body.
    CHANGE_PASSWORD: '/auth/password/reset',
    OAUTH_GOOGLE: '/auth/oauth/google/login',
    OAUTH_GITHUB: '/auth/oauth/github/login',
  },
  USERS: {
    LIKE: (userId: string) => `/users/${userId}/like`,
    UNLIKE: (userId: string) => `/users/${userId}/like`,
    MY_LIKES: '/me/likes',
    MY_VIEWS: '/me/views',
    DELETE_ME: '/me/',
    DELETE_ACCOUNT: '/me/delete',
    MY_BLOCKS: '/me/blocks',
    BLOCK: (userId: string) => `/users/${userId}/block`,
    UNBLOCK: (userId: string) => `/users/${userId}/block`,
    MY_DATA: '/me/data/',
    MY_TAGS: '/me/tags/',
    DELETE_TAG: (tagId: number) => `/me/tags/${tagId}`,
  },
  TAGS: '/tags/',
  PROFILE: {
    CREATE: '/me/profile/',
    UPDATE: '/me/profile/',
    PICTURES: '/me/profile/pictures',
    DELETE_PICTURE: (pictureId: number) => `/me/profile/pictures/${pictureId}`,
    WHO_LIKED_ME: '/me/profile/likes',
    WHO_VIEWED_ME: '/me/profile/views',
    LIST: '/profiles/',
    GET: (userId: string) => `/users/${userId}/profile`,
    RECOMMEND: '/profiles/recommends',
  },
  CHATS: {
    MY_CHATS: '/me/chats',
    MESSAGES: (userId: string) => `/chats/${userId}/messages/`,
  },
  NOTIFICATIONS: {
    LIST: '/me/notifications',
  },
} as const;
