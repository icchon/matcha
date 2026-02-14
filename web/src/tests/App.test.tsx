import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import App from '@/App';

vi.mock('@/stores/authStore', () => ({
  useAuthStore: vi.fn(),
}));

import { useAuthStore } from '@/stores/authStore';

const mockedUseAuthStore = vi.mocked(useAuthStore);

const mockInitialize = vi.fn();

function setupAuthMock(overrides: { isAuthenticated?: boolean } = {}) {
  mockedUseAuthStore.mockReturnValue({
    isAuthenticated: false,
    initialize: mockInitialize,
    ...overrides,
  } as unknown as ReturnType<typeof useAuthStore>);
}

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset window location to root
    window.history.pushState({}, '', '/');
  });

  it('renders home page content when authenticated', () => {
    setupAuthMock({ isAuthenticated: true });

    render(<App />);

    expect(
      screen.getByText('HomePage'),
      'Authenticated users visiting "/" should see the HomePage placeholder. Check that the home route is inside ProtectedRoute and renders HomePage.',
    ).toBeInTheDocument();
  });

  it('redirects unauthenticated users from "/" to login', () => {
    setupAuthMock({ isAuthenticated: false });

    render(<App />);

    expect(
      screen.getByText('LoginPage'),
      'Unauthenticated users visiting "/" should be redirected to /login. Check that ProtectedRoute redirects to /login when isAuthenticated is false.',
    ).toBeInTheDocument();
  });

  it('renders login page at /login', () => {
    setupAuthMock({ isAuthenticated: false });
    window.history.pushState({}, '', '/login');

    render(<App />);

    expect(
      screen.getByText('LoginPage'),
      'The /login route should render the LoginPage component. Check that the public route for /login is configured.',
    ).toBeInTheDocument();
  });

  it('renders 404 page for unknown routes', () => {
    setupAuthMock({ isAuthenticated: false });
    window.history.pushState({}, '', '/some-nonexistent-route');

    render(<App />);

    expect(
      screen.getByText('NotFoundPage'),
      'Unknown routes should render the NotFoundPage component. Check that the catch-all "*" route is configured.',
    ).toBeInTheDocument();
  });

  it('calls initialize on mount to hydrate auth state', () => {
    setupAuthMock({ isAuthenticated: false });

    render(<App />);

    expect(
      mockInitialize,
      'App should call useAuthStore().initialize() on mount to hydrate auth state from localStorage. Check useEffect in App component.',
    ).toHaveBeenCalledOnce();
  });
});
