import React from 'react';
import { logout } from '../services/authService';
import { useAuth } from '../contexts/AuthContext';

const Home: React.FC = () => {
  const auth = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
      auth?.setToken(null);
      alert('Logout successful');
    } catch (error) {
      alert('Logout failed');
    }
  };

  return (
    <div>
      <h1>Home</h1>
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
};

export default Home;