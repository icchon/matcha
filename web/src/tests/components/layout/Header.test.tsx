import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { Header } from '@/components/layout/Header';

const mockLogout = vi.fn();

vi.mock('@/stores/authStore', () => ({
  useAuthStore: vi.fn(),
}));

import { useAuthStore } from '@/stores/authStore';

const mockUseAuthStore = vi.mocked(useAuthStore);

function renderHeader() {
  return render(
    <MemoryRouter>
      <Header />
    </MemoryRouter>,
  );
}

describe('Header', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('when user is NOT authenticated', () => {
    beforeEach(() => {
      mockUseAuthStore.mockReturnValue({
        isAuthenticated: false,
        userId: null,
        isVerified: false,
        authMethod: null,
        login: vi.fn(),
        logout: mockLogout,
        initialize: vi.fn(),
      });
    });

    it('renders app name "Matcha" as a link to home', () => {
      renderHeader();

      const link = screen.getByRole('link', { name: 'Matcha' });
      expect(
        link,
        'Header should display "Matcha" as a link. Check that an <a> element with text "Matcha" is rendered.',
      ).toBeInTheDocument();
      expect(
        link.getAttribute('href'),
        'The "Matcha" link should navigate to "/". Check the "to" prop on the Link component.',
      ).toBe('/');
    });

    it('shows Login link', () => {
      renderHeader();

      const loginLink = screen.getByRole('link', { name: 'Login' });
      expect(
        loginLink,
        'When not authenticated, Header should show a "Login" link. Check conditional rendering based on isAuthenticated.',
      ).toBeInTheDocument();
      expect(
        loginLink.getAttribute('href'),
        'Login link should navigate to "/login". Check the "to" prop.',
      ).toBe('/login');
    });

    it('shows Signup link', () => {
      renderHeader();

      const signupLink = screen.getByRole('link', { name: 'Sign Up' });
      expect(
        signupLink,
        'When not authenticated, Header should show a "Sign Up" link. Check conditional rendering based on isAuthenticated.',
      ).toBeInTheDocument();
      expect(
        signupLink.getAttribute('href'),
        'Sign Up link should navigate to "/signup". Check the "to" prop.',
      ).toBe('/signup');
    });

    it('does NOT show authenticated navigation links', () => {
      renderHeader();

      expect(
        screen.queryByRole('link', { name: 'Browse' }),
        'When not authenticated, "Browse" link should not be visible. Check conditional rendering.',
      ).not.toBeInTheDocument();
      expect(
        screen.queryByRole('link', { name: 'Chat' }),
        'When not authenticated, "Chat" link should not be visible. Check conditional rendering.',
      ).not.toBeInTheDocument();
      expect(
        screen.queryByRole('link', { name: 'Profile' }),
        'When not authenticated, "Profile" link should not be visible. Check conditional rendering.',
      ).not.toBeInTheDocument();
    });

    it('does NOT show logout button', () => {
      renderHeader();

      expect(
        screen.queryByRole('button', { name: 'Logout' }),
        'When not authenticated, "Logout" button should not be visible. Check conditional rendering.',
      ).not.toBeInTheDocument();
    });
  });

  describe('when user IS authenticated', () => {
    beforeEach(() => {
      mockUseAuthStore.mockReturnValue({
        isAuthenticated: true,
        userId: 'user-1',
        isVerified: true,
        authMethod: 'local' as const,
        login: vi.fn(),
        logout: mockLogout,
        initialize: vi.fn(),
      });
    });

    it('shows Browse, Chat, and Profile navigation links', () => {
      renderHeader();

      const browseLink = screen.getByRole('link', { name: 'Browse' });
      expect(
        browseLink,
        'When authenticated, Header should show a "Browse" link. Check conditional rendering based on isAuthenticated.',
      ).toBeInTheDocument();
      expect(browseLink.getAttribute('href'), 'Browse link should navigate to "/browse".').toBe(
        '/browse',
      );

      const chatLink = screen.getByRole('link', { name: 'Chat' });
      expect(
        chatLink,
        'When authenticated, Header should show a "Chat" link.',
      ).toBeInTheDocument();
      expect(chatLink.getAttribute('href'), 'Chat link should navigate to "/chat".').toBe('/chat');

      const profileLink = screen.getByRole('link', { name: 'Profile' });
      expect(
        profileLink,
        'When authenticated, Header should show a "Profile" link.',
      ).toBeInTheDocument();
      expect(
        profileLink.getAttribute('href'),
        'Profile link should navigate to "/profile/edit".',
      ).toBe('/profile/edit');
    });

    it('does NOT show Login or Sign Up links', () => {
      renderHeader();

      expect(
        screen.queryByRole('link', { name: 'Login' }),
        'When authenticated, "Login" link should not be visible.',
      ).not.toBeInTheDocument();
      expect(
        screen.queryByRole('link', { name: 'Sign Up' }),
        'When authenticated, "Sign Up" link should not be visible.',
      ).not.toBeInTheDocument();
    });

    it('shows logout button that calls logout on click', async () => {
      const user = userEvent.setup();
      renderHeader();

      const logoutButton = screen.getByRole('button', { name: 'Logout' });
      expect(
        logoutButton,
        'When authenticated, Header should show a "Logout" button.',
      ).toBeInTheDocument();

      await user.click(logoutButton);

      expect(
        mockLogout,
        'Clicking the Logout button should call the logout function from useAuthStore. Check the onClick handler.',
      ).toHaveBeenCalledOnce();
    });
  });
});
