import { useState, useCallback, type FC } from 'react';
import { Button } from '@/components/ui/Button';
import { Modal } from '@/components/ui/Modal';

interface ActionButtonsProps {
  readonly userId: string;
  readonly isLiked: boolean;
  readonly isBlocked: boolean;
  readonly isLoading?: boolean;
  readonly onLike: () => void;
  readonly onUnlike: () => void;
  readonly onBlock: () => void;
  readonly onUnblock: () => void;
  readonly onReport: () => void;
}

type ConfirmAction = 'block' | 'report' | null;

const ActionButtons: FC<ActionButtonsProps> = ({
  isLiked,
  isBlocked,
  isLoading = false,
  onLike,
  onUnlike,
  onBlock,
  onUnblock,
  onReport,
}) => {
  const [confirmAction, setConfirmAction] = useState<ConfirmAction>(null);

  const handleConfirm = useCallback(() => {
    if (confirmAction === 'block') {
      onBlock();
    } else if (confirmAction === 'report') {
      onReport();
    }
    setConfirmAction(null);
  }, [confirmAction, onBlock, onReport]);

  const handleCancel = useCallback(() => {
    setConfirmAction(null);
  }, []);

  const modalTitle = confirmAction === 'block' ? 'Block User' : 'Report User';
  const modalMessage =
    confirmAction === 'block'
      ? 'Are you sure you want to block this user? They will no longer be able to see your profile.'
      : 'Are you sure you want to report this user?';

  return (
    <>
      <div className="flex flex-wrap gap-2">
        {isLiked ? (
          <Button type="button" onClick={onUnlike} disabled={isLoading}>
            Unlike
          </Button>
        ) : (
          <Button type="button" onClick={onLike} disabled={isLoading}>
            Like
          </Button>
        )}

        {isBlocked ? (
          <Button type="button" onClick={onUnblock} disabled={isLoading}>
            Unblock
          </Button>
        ) : (
          <Button
            type="button"
            onClick={() => setConfirmAction('block')}
            disabled={isLoading}
          >
            Block
          </Button>
        )}

        <Button
          type="button"
          onClick={() => setConfirmAction('report')}
          disabled={isLoading}
        >
          Report
        </Button>
      </div>

      <Modal
        isOpen={confirmAction !== null}
        onClose={handleCancel}
        title={modalTitle}
      >
        <p className="mb-4 text-gray-600">{modalMessage}</p>
        <div className="flex justify-end gap-2">
          <Button type="button" onClick={handleCancel}>
            Cancel
          </Button>
          <Button type="button" onClick={handleConfirm}>
            Confirm
          </Button>
        </div>
      </Modal>
    </>
  );
};

export { ActionButtons };
export type { ActionButtonsProps };
