import type { FC } from 'react';
import { ChangePasswordForm } from '@/features/settings/components/ChangePasswordForm';
import { BlockList } from '@/features/settings/components/BlockList';
import { DeleteAccountSection } from '@/features/settings/components/DeleteAccountSection';

const SettingsPage: FC = () => {
  return (
    <div className="mx-auto max-w-2xl space-y-8 px-4 py-8">
      <h1 className="text-2xl font-bold text-gray-900">Settings</h1>

      <section>
        <h2 className="mb-4 text-lg font-semibold text-gray-900">Change Password</h2>
        <ChangePasswordForm />
      </section>

      <section>
        <h2 className="mb-4 text-lg font-semibold text-gray-900">Blocked Users</h2>
        <BlockList />
      </section>

      <section className="border-t border-gray-200 pt-6">
        <DeleteAccountSection />
      </section>
    </div>
  );
};

export { SettingsPage };
