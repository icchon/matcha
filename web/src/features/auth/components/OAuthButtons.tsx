import type { FC } from 'react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/Button';

// [MOCK] OAuth requires provider registration (client IDs, redirect URIs).
// Buttons render but show "coming soon" toast until OAuth is configured.

const OAuthButtons: FC = () => {
  const handleOAuth = (provider: string) => {
    toast.info(`${provider} login coming soon`);
  };

  return (
    <div className="flex flex-col gap-2">
      <Button variant="secondary" type="button" onClick={() => handleOAuth('Google')}>
        Continue with Google
      </Button>
      <Button variant="secondary" type="button" onClick={() => handleOAuth('GitHub')}>
        Continue with GitHub
      </Button>
    </div>
  );
};

export { OAuthButtons };
