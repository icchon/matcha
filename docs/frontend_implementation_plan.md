# Frontend Implementation Plan

## Overview

Backend (Go) has basic features implemented. Frontend (React + Vite + TS) scaffold is complete (FE-01 merged). This plan covers FE-02 through FE-11 implementation with parallel agent development.

### Principles

- **TDD**: RED -> GREEN -> IMPROVE for all phases
- **Backend-independent**: All API calls are mockable. Unimplemented endpoints use stubs
- **Parallel maximization**: Phase 3 and 4 run 3 agents concurrently

---

## Dependency Graph

```
FE-01 (Complete)
  |
  v
FE-02 (#8) Layout / UI Foundation
  |
  v
FE-03 (#9) Authentication
  |
  +------------------------------+------------------+
  v                              v                  v
FE-04 (#10)                FE-08 (#14)         FE-07 (#13) + FE-11 (#17)
Profile Settings           WebSocket Foundation Profile View + Settings
  |                              |
  v                         +----+----+
FE-05 (#11)                 v         v
Browse/Discovery       FE-09 (#15) FE-10 (#16)
  |                    Chat        Notifications
  v
FE-06 (#12)
Search
```

---

## Known Issues

### 1. `constants.ts` API Path Mismatches (Fix in Phase 1)

| constants.ts value | Actual route | Status |
|---|---|---|
| `USERS.MY_LIKES: '/users/me/likes'` | `/me/likes` | Mismatch |
| `USERS.MY_VIEWS: '/users/me/views'` | `/me/views` | Mismatch |
| `USERS.DELETE_ME: '/users/me'` | `/me/` | Mismatch |
| `USERS.MY_BLOCKS: '/users/me/blocks'` | `/me/blocks` | Mismatch |
| `USERS.MY_DATA: '/users/me/data'` | `/me/data/` | Mismatch |
| `USERS.MY_TAGS: '/users/me/tags'` | `/me/tags/` | Mismatch |
| `PROFILE.CREATE: '/profile'` | `/me/profile/` | Mismatch |
| `PROFILE.UPDATE: '/profile'` | `/me/profile/` | Mismatch |
| `PROFILE.PICTURES: '/profile/pictures'` | `/me/profile/pictures` | Mismatch |
| `PROFILE.WHO_LIKED_ME: '/profile/likes'` | `/me/profile/likes` | Mismatch |
| `PROFILE.WHO_VIEWED_ME: '/profile/views'` | `/me/profile/views` | Mismatch |

### 2. WebSocket Authentication

- WS Gateway uses `Authorization: Bearer` header for auth
- Browser `new WebSocket()` API does not support custom headers
- **Solution**: nginx query param (`/ws?token=XXX`) to `Authorization` header conversion
- **[MOCK]** WS auth requires nginx config change. Frontend tests mock WebSocket connections until BE-03 (#20) is ready

### 3. Unimplemented Backend Endpoints (tracked in BE-08 #25)

| Endpoint | Purpose | FE issue | BE issue |
|---|---|---|---|
| `DELETE /users/{userID}/block` | Unblock user | FE-07 (#13), FE-11 (#17) | BE-08 (#25) |
| `POST /users/{userID}/report` | Report user | FE-07 (#13) | BE-08 (#25) |
| `PUT /me/notifications/{id}/read` | Mark notification read | FE-10 (#16) | BE-08 (#25) |

---

## TDD Workflow (All Phases)

1. **RED**: Write test files first. Describe expected behavior with assertions
   - Components: `@testing-library/react` for render + user interaction
   - Stores: Zustand store state transition tests
   - API functions: `vi.mock` for fetch mocking, success/error cases
   - Hooks: `renderHook` for lifecycle tests
2. **GREEN**: Minimal implementation to pass tests
3. **IMPROVE**: Refactor while keeping tests green

### Mock Strategy

- API calls: `vi.mock('@/api/client')` to mock `apiClient`
- WebSocket: `vi.mock` for WebSocket class, simulate events
- Router: `MemoryRouter` for navigation tests
- Unimplemented BE endpoints: Mock functions returning success. Test comments include `// [MOCK] BE-08 #25: endpoint not yet implemented`

---

## Phase 1: Foundation (FE-02) - 1 agent

**Issue**: #8 | **Blocks**: FE-03 and all subsequent

1. Fix `constants.ts` API path mismatches + update tests
2. Add `upload()` method to `apiClient` (FormData support) + tests
3. UI components (`components/ui/`) with tests first:
   - `Button`, `Input`, `Modal`, `Card`, `Badge`, `Spinner`
4. Layout components (`components/layout/`) + tests:
   - `Layout`, `Header`, `Footer`
5. `ProtectedRoute` (`components/common/`) + test (unauthenticated redirect)
6. `App.tsx` route structure definition
7. Sonner Toaster in `App.tsx`

## Phase 2: Authentication (FE-03) - 1 agent

**Issue**: #9 | **Depends**: FE-02 | **Blocks**: FE-04, FE-07, FE-08, FE-11

1. Zod validation schemas (`lib/validators.ts`) + tests (valid/invalid cases)
2. API functions (`api/auth.ts`) + tests (apiClient mock)
3. Auth hook (`features/auth/hooks/useAuth.ts`) + tests
4. Auth components + tests: `LoginForm`, `SignupForm`, `OAuthButtons`
5. Auth pages + tests: `LoginPage`, `SignupPage`, `VerifyEmailPage`, `ForgotPasswordPage`, `ResetPasswordPage`

## Phase 3: Parallel Development - 3 agents

After FE-03 completion, 3 independent tracks run concurrently.

| Agent | Issue | Tasks | Branch |
|-------|-------|-------|--------|
| A | FE-04 (#10) | Profile Settings | `feat/fe-04-profile` |
| B | FE-07 (#13) -> FE-11 (#17) | Profile View -> Settings (sequential) | `feat/fe-07-profile-view` |
| C | FE-08 (#14) | WebSocket Foundation | `feat/fe-08-websocket` |

### Agent A (FE-04)
- API function tests -> implementation (`api/profile.ts`): createProfile, updateProfile, uploadPicture (FormData), deletePicture, tag CRUD, userData CRUD
- `profileStore.ts` tests -> implementation
- Component tests -> implementation: `ProfileForm`, `PhotoUploader` (D&D, max 5 pics), `TagManager` (autocomplete)
- Page tests -> implementation: `ProfileCreatePage` (first login), `EditProfilePage`
- `useGeolocation` hook tests -> implementation: Geolocation API + manual input fallback

### Agent B (FE-07 -> FE-11)
- FE-07:
  - API function tests -> implementation (`api/users.ts`): getUserProfile, like/unlike, block/unblock, getLists
  - **[MOCK]** `unblockUser()`, `reportUser()` are mock implementations. Need connection when BE-08 (#25) implements `DELETE /users/{id}/block`, `POST /users/{id}/report`
  - `ProfileCard` tests -> implementation (shared component: reused in FE-05, FE-06)
  - `OnlineIndicator` tests -> implementation
  - Page tests -> implementation: `UserProfilePage`, `LikesPage`, `ViewsPage`
- FE-11:
  - `SettingsPage` tests -> implementation (account deletion + block list management)

### Agent C (FE-08)
- `wsStore.ts` tests -> implementation: WebSocket connection management, auto-reconnect (exponential backoff)
  - **[MOCK]** WebSocket connections fully mocked. FE tests use `MockWebSocket` class to simulate events
- Message routing tests -> implementation: event type -> store handler dispatch
- `useWebSocket` hook tests -> implementation: auto connect/disconnect on auth state change
- `chatStore.ts` / `notificationStore.ts` placeholders (interfaces only) + tests

## Phase 4: Parallel Development - 3 agents

Second parallel phase after Phase 3 completion.

| Agent | Issue | Tasks | Branch |
|-------|-------|-------|--------|
| A | FE-05 (#11) | Browse/Discovery | `feat/fe-05-browse` |
| B | FE-09 (#15) | Chat | `feat/fe-09-chat` |
| C | FE-10 (#16) | Notifications | `feat/fe-10-notifications` |

### Agent A (FE-05)
- API function tests -> implementation (`api/profiles.ts`): getRecommendedProfiles, getProfiles (with filters)
- Component tests -> implementation: `ProfileList`, `FilterPanel` (sliders), `SortControls`
- Page tests -> implementation: `BrowsePage` (`ProfileCard` grid + filter + sort)
- Pagination / infinite scroll

### Agent B (FE-09)
- `chatStore.ts` tests -> full implementation: conversations Map, unreadCount, onMessage/onAck/onRead
- API function tests -> implementation (`api/chat.ts`): getChats, getMessages (pagination)
- Component tests -> implementation: `ChatList`, `ChatWindow`, `MessageBubble`, `MessageInput`
- Page tests -> implementation: `ChatPage` (conversation list + chat window split)

### Agent C (FE-10)
- `notificationStore.ts` tests -> full implementation: notifications, unreadCount, onNotification + toast
- API function tests -> implementation: getNotifications
  - **[MOCK]** `markAsRead()` is mock implementation. Needs connection when BE-08 (#25) implements `PUT /me/notifications/{id}/read`
- Component tests -> implementation: `NotificationBell` (Header integration), `NotificationList`, `NotificationItem`
- Sonner toast integration tests -> implementation

## Phase 5: Search (FE-06) - 1 agent

**Issue**: #12 | **Depends**: FE-05

- Search form tests -> implementation: age/fame rating range sliders, location, multi-tag selection
- `SearchPage` tests -> implementation: reuse FE-05 components (`ProfileList`, `SortControls`)
- Reuse FE-04 `TagManager`

---

## Shared Component Dependency Map

| Component | Created In | Used By |
|---|---|---|
| `Button`, `Input`, `Modal`, `Card`, `Badge`, `Spinner` | Phase 1 (FE-02) | All features |
| `Layout`, `Header`, `Footer` | Phase 1 (FE-02) | All pages |
| `ProtectedRoute` | Phase 1 (FE-02) | All authenticated pages |
| `ProfileCard` | Phase 3 (FE-07) | FE-05, FE-06, FE-11 |
| `OnlineIndicator` | Phase 3 (FE-07) | FE-09 (Chat) |
| `TagManager` | Phase 3 (FE-04) | FE-06 (Search) |
| `FilterPanel`, `SortControls` | Phase 4 (FE-05) | FE-06 (Search) |
| `NotificationBell` | Phase 4 (FE-10) | Header |

---

## Backend Integration Mock Summary

| FE Location | Mock Content | Target BE Issue | Timeline |
|---|---|---|---|
| `api/users.ts: unblockUser()` | Success response stub | BE-08 (#25) | After BE impl |
| `api/users.ts: reportUser()` | Success response stub | BE-08 (#25) | After BE impl |
| `api/notifications.ts: markAsRead()` | Success response stub | BE-08 (#25) | After BE impl |
| `wsStore.ts: connect()` | MockWebSocket for tests | BE-03 (#20) + nginx | After nginx config |
| `constants.ts` | Paths fixed, integration test pending | - | docker compose up |
