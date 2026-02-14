import { useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Toaster } from 'sonner';
import { Layout } from '@/components/layout';
import { ProtectedRoute } from '@/components/common/ProtectedRoute';
import { useAuthStore } from '@/stores/authStore';

// Placeholder page components - will be replaced in later phases
function LoginPage() {
  return <div>LoginPage</div>;
}

function SignupPage() {
  return <div>SignupPage</div>;
}

function VerifyEmailPage() {
  return <div>VerifyEmailPage</div>;
}

function ForgotPasswordPage() {
  return <div>ForgotPasswordPage</div>;
}

function ResetPasswordPage() {
  return <div>ResetPasswordPage</div>;
}

function HomePage() {
  return <div>HomePage</div>;
}

function BrowsePage() {
  return <div>BrowsePage</div>;
}

function SearchPage() {
  return <div>SearchPage</div>;
}

function ProfileCreatePage() {
  return <div>ProfileCreatePage</div>;
}

function EditProfilePage() {
  return <div>EditProfilePage</div>;
}

function UserProfilePage() {
  return <div>UserProfilePage</div>;
}

function ChatPage() {
  return <div>ChatPage</div>;
}

function LikesPage() {
  return <div>LikesPage</div>;
}

function ViewsPage() {
  return <div>ViewsPage</div>;
}

function NotificationsPage() {
  return <div>NotificationsPage</div>;
}

function SettingsPage() {
  return <div>SettingsPage</div>;
}

function NotFoundPage() {
  return <div>NotFoundPage</div>;
}

function App() {
  const { initialize } = useAuthStore();

  useEffect(() => {
    initialize();
  }, [initialize]);

  return (
    <BrowserRouter>
      <Toaster position="top-right" richColors />
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/verify/:token" element={<VerifyEmailPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        <Route path="/reset-password" element={<ResetPasswordPage />} />

        {/* Protected routes with Layout */}
        <Route element={<ProtectedRoute />}>
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/browse" element={<BrowsePage />} />
            <Route path="/search" element={<SearchPage />} />
            <Route path="/profile/create" element={<ProfileCreatePage />} />
            <Route path="/profile/edit" element={<EditProfilePage />} />
            <Route path="/users/:userId" element={<UserProfilePage />} />
            <Route path="/chat" element={<ChatPage />} />
            <Route path="/likes" element={<LikesPage />} />
            <Route path="/views" element={<ViewsPage />} />
            <Route path="/notifications" element={<NotificationsPage />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Route>
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
