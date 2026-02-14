# FE-03: Authentication — Detailed Implementation Plan

> Issue: #9 | Depends: FE-02 (#8) ✅ | Blocks: FE-04, FE-07, FE-08, FE-11

## Pre-Implementation Findings

| ID | Finding | Impact |
|----|---------|--------|
| F-1 | BE returns snake_case (`user_id`, `access_token`), FE types use camelCase | `mapLoginResponse()` mapper required in `api/auth.ts` |
| F-2 | Token refresh HTTP route not registered in BE (`IssueAccessToken` service exists but no handler route) | 401 interceptor → logout (`[MOCK]`). Session = tab lifetime only |
| F-3 | `SendVerificationEmail` needs `{ user_id, email }` but signup returns `{ message }` only | Resend form needs separate email input |

## Decisions

| # | Decision |
|---|----------|
| 1 | `STORAGE_KEYS` removal: atomic in Step 1. PR must build |
| 2 | OAuth: `[MOCK]` — buttons render but show "OAuth coming soon" toast. Provider registration needed for real flow |
| 3 | `RawLoginResponse`: separate `types/raw.ts` file |
| 4 | Token refresh: `[MOCK]` (401 → logout). Future BE route addition enables swap |
| 5 | Token storage: in-memory (closure), not localStorage (H-1 fix) |

## Implementation Steps

### Step 1: `apiClient` Refactoring + In-Memory Tokens (~20 tests)

Addresses H-1 (XSS) and H-2 (401 interceptor).

| Action | File | Changes |
|--------|------|---------|
| Modify | `src/types/api.ts` | Add `TokenRefreshRequest`, `SendVerificationRequest` |
| Create | `src/types/raw.ts` | `RawLoginResponse` (snake_case fields) |
| Modify | `src/lib/constants.ts` | Add `TOKEN_REFRESH` to `API_PATHS.AUTH`, remove `STORAGE_KEYS` |
| Modify | `src/api/client.ts` | Replace localStorage with in-memory closure, extract `request()`, add 401 interceptor |
| Modify | `src/stores/authStore.ts` | Remove localStorage read/write, `initialize()` becomes no-op |
| Test | `src/tests/api/client.test.ts` | Rewrite: in-memory token tests, 401 interceptor flow, no localStorage refs |
| Test | `src/tests/stores/authStore.test.ts` | Rewrite: remove localStorage assertions |

Key design:
```
// Module-scoped token storage (closure)
let accessToken: string | null = null;
let refreshToken: string | null = null;

export const setTokens = (access: string, refresh: string) => { ... };
export const clearTokens = () => { ... };
export const getAccessToken = () => accessToken;
```

401 interceptor:
```
// [MOCK] POST /auth/token/refresh not yet routed in BE
// On 401: clear tokens + redirect to login
// Future: attempt refresh before logout
```

### Step 2: Zod Validation Schemas (~25 tests)

| Action | File | Content |
|--------|------|---------|
| Create | `src/lib/validators.ts` | `loginSchema`, `signupSchema`, `forgotPasswordSchema`, `resetPasswordSchema`, `sendVerificationSchema` |
| Test | `src/tests/lib/validators.test.ts` | Valid/invalid cases per schema |

Rules (matching BE `util.go`):
- Email: `z.string().email()`
- Password: `z.string().min(8)` (BE: `MIN_PASSWORD_LENGTH = 8`)
- Username: `z.string().min(3).max(30)`
- `signupSchema` / `resetPasswordSchema`: `.refine()` for password confirmation match

### Step 3: Auth API Functions (~18 tests)

| Action | File | Content |
|--------|------|---------|
| Create | `src/api/auth.ts` | `login()`, `signup()`, `logout()`, `verifyEmail()`, `sendVerificationEmail()`, `forgotPassword()`, `resetPassword()`, `oauthLogin()` `[MOCK]`, `refreshToken()` `[MOCK]` |
| Test | `src/tests/api/auth.test.ts` | Mock `apiClient`, test path/method/body/response-mapping per function |

Includes `mapLoginResponse(raw: RawLoginResponse): LoginResponse` for snake_case → camelCase.

BE endpoint mapping:
| FE Function | Method | Path | BE Handler |
|---|---|---|---|
| `login()` | POST | `/auth/login` | `LoginHandler` |
| `signup()` | POST | `/auth/register` | `RegisterHandler` |
| `logout()` | POST | `/auth/logout` | `LogoutHandler` |
| `verifyEmail()` | POST | `/auth/verify` | `VerifyEmailHandler` |
| `sendVerificationEmail()` | POST | `/auth/send-verification` | `SendVerificationEmailHandler` |
| `forgotPassword()` | POST | `/auth/forgot-password` | `ForgotPasswordHandler` |
| `resetPassword()` | POST | `/auth/reset-password` | `ResetPasswordHandler` |
| `oauthLogin()` | GET | `/auth/oauth/{provider}` | `OAuthHandler` — `[MOCK]` |
| `refreshToken()` | — | — | No route — `[MOCK]` |

### Step 4: Auth Hook (~15 tests)

| Action | File | Content |
|--------|------|---------|
| Create | `src/features/auth/hooks/useAuth.ts` | Wraps API + store + navigate + toast |
| Test | `src/tests/features/auth/hooks/useAuth.test.ts` | Mock API + store, test each action |

Exposes: `login()`, `signup()`, `logout()`, `verifyEmail()`, `sendVerificationEmail()`, `forgotPassword()`, `resetPassword()` with `isLoading`, `error` state.

### Step 5: Auth Components (~24 tests)

| Action | File | Tests |
|--------|------|-------|
| Create | `src/features/auth/components/LoginForm.tsx` | ~10 |
| Create | `src/features/auth/components/SignupForm.tsx` | ~10 |
| Create | `src/features/auth/components/OAuthButtons.tsx` | ~4 |
| Create | `src/features/auth/components/index.ts` | — |
| Test | `src/tests/features/auth/components/LoginForm.test.tsx` | |
| Test | `src/tests/features/auth/components/SignupForm.test.tsx` | |
| Test | `src/tests/features/auth/components/OAuthButtons.test.tsx` | |

Components use `react-hook-form` + `@hookform/resolvers/zod` + existing `Button`/`Input` from `src/components/ui/`.

### Step 6: Auth Pages + App.tsx Update (~22 tests)

| Action | File | Tests |
|--------|------|-------|
| Create | `src/features/auth/pages/LoginPage.tsx` | ~4 |
| Create | `src/features/auth/pages/SignupPage.tsx` | ~4 |
| Create | `src/features/auth/pages/VerifyEmailPage.tsx` | ~5 |
| Create | `src/features/auth/pages/ForgotPasswordPage.tsx` | ~4 |
| Create | `src/features/auth/pages/ResetPasswordPage.tsx` | ~5 |
| Create | `src/features/auth/pages/index.ts` | — |
| Modify | `src/App.tsx` | Replace placeholder pages with real imports |
| Test | `src/tests/features/auth/pages/*.test.tsx` | |

Route structure (existing in App.tsx):
- `/login` → `LoginPage`
- `/signup` → `SignupPage`
- `/verify/:token` → `VerifyEmailPage`
- `/forgot-password` → `ForgotPasswordPage`
- `/reset-password` → `ResetPasswordPage` (token via query param `?token=xxx`)

## Total: ~124 tests across 6 steps
