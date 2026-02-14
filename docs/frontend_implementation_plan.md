# Frontend Implementation Plan

> Last updated: 2026-02-14

## Overview

Backend (Go) has basic features implemented. Frontend (React + Vite + TS) scaffold is complete (FE-01 merged). This plan covers FE-02 through FE-11 implementation with parallel agent development.

### Principles

- **TDD**: RED -> GREEN -> IMPROVE for all phases
- **Backend-independent**: All API calls are mockable. Unimplemented endpoints use stubs
- **Parallel maximization**: Phase 3 runs 4 agents concurrently, Phase 4 runs 3 agents concurrently

---

## Progress

| Phase | Issue | Status | PR | Tests |
|-------|-------|--------|-----|-------|
| 1 | FE-02 (#8) | **Complete** | #33 | 104 |
| 2 | FE-03 (#9) | Pending | — | — |
| 3 | FE-04 (#10), FE-07 (#13), FE-08 (#14), FE-11 (#17) — 4 agents parallel | Pending | — | — |
| 4 | FE-05 (#11), FE-09 (#15), FE-10 (#16) | Pending | — | — |
| 5 | FE-06 (#12) | Pending | — | — |

---

## Dependency Graph

```
FE-01 ✅
  |
  v
FE-02 (#8) ✅ Layout / UI Foundation
  |
  v
FE-03 (#9) Authentication
  |
  +------------------------------+------------------+
  v                              v                  v
FE-04 (#10)        FE-08 (#14)         FE-07 (#13)      FE-11 (#17)
Profile Settings   WebSocket Foundation Profile View     Settings
  |                      |
  v                 +----+----+
FE-05 (#11)         v         v
Browse/Discovery FE-09 (#15) FE-10 (#16)
  |              Chat        Notifications
  v
FE-06 (#12)
Search
```

---

## Known Issues

### 1. `constants.ts` API Path Mismatches — RESOLVED in Phase 1

11 箇所のパス不一致を修正済み。`CHATS`, `NOTIFICATIONS` グループも追加。

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

### 4. Security Review Findings (Phase 1)

Phase 1 の security-reviewer で検出された HIGH 項目。Phase 2 以降で対応:

| ID | Issue | Resolution |
|----|-------|-----------|
| H-1 | localStorage にトークン保存 (XSS リスク) | FE-03 でインメモリ保存 (Zustand store + closure) に移行。BE 変更不要。トレードオフ: ページリロード時に再ログインが必要 |
| H-2 | トークンリフレッシュ機構なし | FE-03 で 401 インターセプタ + リフレッシュ実装 |

### 5. apiClient 構造改善 (Phase 2 で対応)

Phase 1 レビューで `apiClient` の各メソッド (`get`/`post`/`put`/`delete`) の fetch パターン重複が指摘された。FE-03 でトークンリフレッシュの 401 インターセプタを追加するタイミングで、共通 `request()` ベースメソッドへのリファクタリングを行う。

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

## Phase 1: Foundation (FE-02) — COMPLETE

**Issue**: #8 | **PR**: #33 | **Branch**: `feat/fe-02-layout` | **Tests**: 104

### Deliverables

| Category | Files |
|----------|-------|
| API path fixes | `constants.ts` — 11 mismatches fixed, `CHATS`/`NOTIFICATIONS` groups added |
| apiClient | `client.ts` — `upload()` method (FormData), immutable header builders |
| UI components | `Button`, `Input`, `Modal`, `Card`, `Badge`, `Spinner` |
| Layout components | `Header` (auth-aware nav), `Footer`, `Layout` |
| Routing | `ProtectedRoute`, `App.tsx` (full route structure + Sonner Toaster) |
| Auth store fixes | Base64URL decoding, runtime type guards for JWT payload |

### Review Results

| Reviewer | Result |
|----------|--------|
| code-reviewer | Approve (WARNING 3 件修正済み) |
| security-reviewer | CRITICAL: 0, HIGH: 2 (FE-03/BE scope), M-2 修正済み |

---

## Phase 2: Authentication (FE-03) — 1 agent

**Issue**: #9 | **Depends**: FE-02 ✅ | **Blocks**: FE-04, FE-07, FE-08, FE-11

### Tasks

1. `apiClient` リファクタリング: 共通 `request()` メソッド + 401 インターセプタ (H-2 対応)
2. Zod validation schemas (`lib/validators.ts`) + tests (valid/invalid cases)
3. API functions (`api/auth.ts`) + tests (apiClient mock)
4. Auth hook (`features/auth/hooks/useAuth.ts`) + tests
5. Auth components + tests: `LoginForm`, `SignupForm`, `OAuthButtons`
6. Auth pages + tests: `LoginPage`, `SignupPage`, `VerifyEmailPage`, `ForgotPasswordPage`, `ResetPasswordPage`

---

## Phase 3: Parallel Development — 4 agents

After FE-03 completion, 4 independent tracks run concurrently. FE-11 depends only on FE-03 so it runs as a separate agent.

| Agent | Issue | Tasks | Branch |
|-------|-------|-------|--------|
| A | FE-04 (#10) | Profile Settings | `feat/fe-04-profile` |
| B | FE-07 (#13) | Profile View | `feat/fe-07-profile-view` |
| C | FE-08 (#14) | WebSocket Foundation | `feat/fe-08-websocket` |
| D | FE-11 (#17) | Settings / Account Management | `feat/fe-11-settings` |

### Agent A (FE-04)
- API function tests -> implementation (`api/profile.ts`): createProfile, updateProfile, uploadPicture (FormData), deletePicture, tag CRUD, userData CRUD
- `profileStore.ts` tests -> implementation
- Component tests -> implementation: `ProfileForm`, `PhotoUploader` (D&D, max 5 pics), `TagManager` (autocomplete)
- Page tests -> implementation: `ProfileCreatePage` (first login), `EditProfilePage`
- `useGeolocation` hook tests -> implementation: Geolocation API + manual input fallback

### Agent B (FE-07)
- API function tests -> implementation (`api/users.ts`): getUserProfile, like/unlike, block/unblock, getLists
- **[MOCK]** `unblockUser()`, `reportUser()` are mock implementations. Need connection when BE-08 (#25) implements `DELETE /users/{id}/block`, `POST /users/{id}/report`
- `ProfileCard` tests -> implementation (shared component: reused in FE-05, FE-06)
- `OnlineIndicator` tests -> implementation
- Page tests -> implementation: `UserProfilePage`, `LikesPage`, `ViewsPage`

### Agent C (FE-08)
- `wsStore.ts` tests -> implementation: WebSocket connection management, auto-reconnect (exponential backoff)
  - **[MOCK]** WebSocket connections fully mocked. FE tests use `MockWebSocket` class to simulate events
- Message routing tests -> implementation: event type -> store handler dispatch
- `useWebSocket` hook tests -> implementation: auto connect/disconnect on auth state change
- `chatStore.ts` / `notificationStore.ts` placeholders (interfaces only) + tests

### Agent D (FE-11)
- `SettingsPage` tests -> implementation: account deletion, password change, block list management
- API function tests -> implementation: deleteAccount, changePassword, getBlockList, unblockUser
- **[MOCK]** `unblockUser()` is mock implementation (same as FE-07). Needs BE-08 (#25)

---

## Phase 4: Parallel Development — 3 agents

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

---

## Phase 5: Search (FE-06) — 1 agent

**Issue**: #12 | **Depends**: FE-05

> **並列化見送り理由**: FE-06 は FE-05 の `ProfileList`, `FilterPanel`, `SortControls` を直接再利用する設計。これらを事前に共通コンポーネントとして切り出す方法もあるが、ブラウジング UI と密結合であり、先に FE-05 でコンポーネント設計を確定してから FE-06 で再利用する方が手戻りリスクが低い。FE-06 自体のスコープも小さいため、sequential でもボトルネックにならない。

- Search form tests -> implementation: age/fame rating range sliders, location, multi-tag selection
- `SearchPage` tests -> implementation: reuse FE-05 components (`ProfileList`, `SortControls`)
- Reuse FE-04 `TagManager`

---

## Shared Component Dependency Map

| Component | Created In | Status | Used By |
|---|---|---|---|
| `Button`, `Input`, `Modal`, `Card`, `Badge`, `Spinner` | Phase 1 (FE-02) | **Done** | All features |
| `Layout`, `Header`, `Footer` | Phase 1 (FE-02) | **Done** | All pages |
| `ProtectedRoute` | Phase 1 (FE-02) | **Done** | All authenticated pages |
| `ProfileCard` | Phase 3 (FE-07) | Pending | FE-05, FE-06, FE-11 |
| `OnlineIndicator` | Phase 3 (FE-07) | Pending | FE-09 (Chat) |
| `TagManager` | Phase 3 (FE-04) | Pending | FE-06 (Search) |
| `FilterPanel`, `SortControls` | Phase 4 (FE-05) | Pending | FE-06 (Search) |
| `NotificationBell` | Phase 4 (FE-10) | Pending | Header |

---

## Backend Integration Mock Summary

| FE Location | Mock Content | Target BE Issue | Timeline |
|---|---|---|---|
| `api/users.ts: unblockUser()` | Success response stub | BE-08 (#25) | After BE impl |
| `api/users.ts: reportUser()` | Success response stub | BE-08 (#25) | After BE impl |
| `api/notifications.ts: markAsRead()` | Success response stub | BE-08 (#25) | After BE impl |
| `wsStore.ts: connect()` | MockWebSocket for tests | BE-03 (#20) + nginx | After nginx config |
