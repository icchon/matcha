import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { ProtectedRoute } from '@/components/common/ProtectedRoute';

vi.mock('@/stores/authStore', () => ({
  useAuthStore: vi.fn(),
}));

import { useAuthStore } from '@/stores/authStore';

const mockedUseAuthStore = vi.mocked(useAuthStore);

function ChildPage() {
  return <div>Protected Content</div>;
}

function LoginPage() {
  return <div>Login Page</div>;
}

function renderWithRouter(initialPath: string) {
  return render(
    <MemoryRouter initialEntries={[initialPath]}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<ChildPage />} />
        </Route>
      </Routes>
    </MemoryRouter>,
  );
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders child route when user is authenticated', () => {
    mockedUseAuthStore.mockReturnValue({
      isAuthenticated: true,
    } as ReturnType<typeof useAuthStore>);

    renderWithRouter('/dashboard');

    expect(
      screen.getByText('Protected Content'),
      'Authenticated users should see the child route content. Check that ProtectedRoute renders <Outlet /> when isAuthenticated is true.',
    ).toBeInTheDocument();
  });

  it('redirects to /login when user is not authenticated', () => {
    mockedUseAuthStore.mockReturnValue({
      isAuthenticated: false,
    } as ReturnType<typeof useAuthStore>);

    renderWithRouter('/dashboard');

    expect(
      screen.getByText('Login Page'),
      'Unauthenticated users should be redirected to /login. Check that ProtectedRoute renders <Navigate to="/login" replace /> when isAuthenticated is false.',
    ).toBeInTheDocument();

    expect(
      screen.queryByText('Protected Content'),
      'Protected content should not be visible to unauthenticated users.',
    ).not.toBeInTheDocument();
  });
});
