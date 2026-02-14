import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import { Layout } from '@/components/layout/Layout';

vi.mock('@/stores/authStore', () => ({
  useAuthStore: vi.fn(() => ({
    isAuthenticated: false,
    userId: null,
    isVerified: false,
    authMethod: null,
    login: vi.fn(),
    logout: vi.fn(),
    initialize: vi.fn(),
  })),
}));

function renderLayoutWithRoute(childContent: string) {
  const router = createMemoryRouter(
    [
      {
        path: '/',
        element: <Layout />,
        children: [
          {
            index: true,
            element: <div>{childContent}</div>,
          },
        ],
      },
    ],
    { initialEntries: ['/'] },
  );

  return render(<RouterProvider router={router} />);
}

describe('Layout', () => {
  it('renders the Header', () => {
    renderLayoutWithRoute('Test Content');

    expect(
      screen.getByRole('link', { name: 'Matcha' }),
      'Layout should render the Header component which contains a "Matcha" link. Check that Layout includes <Header />.',
    ).toBeInTheDocument();
  });

  it('renders the Footer', () => {
    renderLayoutWithRoute('Test Content');

    expect(
      screen.getByRole('contentinfo'),
      'Layout should render the Footer component (a <footer> element). Check that Layout includes <Footer />.',
    ).toBeInTheDocument();
  });

  it('renders child route content via Outlet', () => {
    renderLayoutWithRoute('Child Route Content');

    expect(
      screen.getByText('Child Route Content'),
      'Layout should render nested route content via <Outlet />. Check that Layout includes <Outlet /> inside a <main> element.',
    ).toBeInTheDocument();
  });

  it('renders a main content area', () => {
    renderLayoutWithRoute('Test Content');

    const main = screen.getByRole('main');
    expect(
      main,
      'Layout should wrap the Outlet in a <main> element for semantic structure. Check the element wrapping <Outlet />.',
    ).toBeInTheDocument();
  });
});
