import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Profile from './pages/Profile';
import Profiles from './pages/Profiles';
import UserProfile from './pages/UserProfile';
import Chats from './pages/Chats';
import Chat from './pages/Chat';
import ForgotPassword from './pages/ForgotPassword';
import ResetPassword from './pages/ResetPassword';
import Notifications from './pages/Notifications';
import RecommendedProfiles from './pages/RecommendedProfiles';
import Tags from './pages/Tags';
import { AuthProvider } from './contexts/AuthContext';
import { WebSocketProvider, useWebSocket } from './contexts/WebSocketContext';
import './App.css';

const WebSocketStatus: React.FC = () => {
    const ws = useWebSocket();
    // Don't render anything if the context is not yet available
    if (!ws) return null;
    return (
        <span style={{ marginLeft: '20px', color: ws.isConnected ? 'green' : 'red', fontWeight: 'bold' }}>
            WS: {ws.isConnected ? 'Connected' : 'Disconnected'}
        </span>
    );
};

function App() {
  return (
    <AuthProvider>
      <WebSocketProvider>
        <Router>
          <div>
            <nav>
              <ul>
                <li><Link to="/">Home</Link></li>
                <li><Link to="/login">Login</Link></li>
                <li><Link to="/signup">Signup</Link></li>
                <li><Link to="/profile">My Profile</Link></li>
                <li><Link to="/profiles">Browse Profiles</Link></li>
                <li><Link to="/chats">Chats</Link></li>
                <li><Link to="/forgot-password">Forgot Password</Link></li>
                <li><Link to="/notifications">Notifications</Link></li>
                <li><Link to="/recommended-profiles">Recommended</Link></li>
                <li><Link to="/tags">Tags</Link></li>
              </ul>
              <WebSocketStatus />
            </nav>
            <hr />
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/login" element={<Login />} />
              <Route path="/signup" element={<Signup />} />
              <Route path="/profile" element={<Profile />} />
              <Route path="/profiles" element={<Profiles />} />
              <Route path="/users/:userId/profile" element={<UserProfile />} />
              <Route path="/chats" element={<Chats />} />
              <Route path="/chats/:userId" element={<Chat />} />
              <Route path="/forgot-password" element={<ForgotPassword />} />
              <Route path="/reset-password" element={<ResetPassword />} />
              <Route path="/notifications" element={<Notifications />} />
              <Route path="/recommended-profiles" element={<RecommendedProfiles />} />
              <Route path="/tags" element={<Tags />} />
            </Routes>
          </div>
        </Router>
      </WebSocketProvider>
    </AuthProvider>
  );
}

export default App;
