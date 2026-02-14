import { Link } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';

export function Header() {
  const { isAuthenticated, logout } = useAuthStore();

  return (
    <header className="border-b">
      <nav className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
        <Link to="/" className="text-xl font-bold">
          Matcha
        </Link>
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <>
              <Link to="/browse">Browse</Link>
              <Link to="/chat">Chat</Link>
              <Link to="/profile/edit">Profile</Link>
              <button type="button" onClick={logout}>
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login">Login</Link>
              <Link to="/signup">Sign Up</Link>
            </>
          )}
        </div>
      </nav>
    </header>
  );
}
