import { useState, type FC } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { deleteAccountSchema, type DeleteAccountFormData } from '@/features/settings/validators';
import { useDeleteAccount } from '@/features/settings/hooks/useDeleteAccount';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';

const DeleteAccountSection: FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { deleteAccount, isLoading } = useDeleteAccount();
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<DeleteAccountFormData>({
    resolver: zodResolver(deleteAccountSchema),
  });

  const openModal = () => setIsModalOpen(true);
  const closeModal = () => {
    setIsModalOpen(false);
    reset();
  };

  const onSubmit = async () => {
    await deleteAccount();
    // Navigation to /login happens inside the hook — no need to close modal
    // as the component will unmount. Calling closeModal() here would trigger
    // a React state update on an unmounted component.
  };

  return (
    <div>
      <h3 className="text-lg font-semibold text-red-600">Danger Zone</h3>
      <p className="mt-1 text-sm text-gray-600">
        Once you delete your account, there is no going back.
      </p>
      <Button variant="danger" onClick={openModal} className="mt-3">
        Delete Account
      </Button>

      <Modal isOpen={isModalOpen} onClose={closeModal} title="Delete Account">
        <p className="mb-4 text-sm text-gray-600">
          This action is irreversible. Type DELETE to confirm.
        </p>
        <form onSubmit={handleSubmit(onSubmit)} noValidate className="flex flex-col gap-4">
          <Input
            label="Confirm"
            placeholder="DELETE"
            error={errors.confirmText?.message}
            {...register('confirmText')}
          />
          <div className="flex gap-3 justify-end">
            <Button type="button" variant="secondary" onClick={closeModal}>
              Cancel
            </Button>
            <Button type="submit" variant="danger" loading={isLoading}>
              Confirm Delete
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export { DeleteAccountSection };
