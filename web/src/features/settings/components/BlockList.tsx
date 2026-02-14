import { useEffect, type FC } from 'react';
import { useBlockList } from '@/features/settings/hooks/useSettings';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';

const BlockList: FC = () => {
  const { blocks, isLoading, error, fetchBlockList, unblock } = useBlockList();

  useEffect(() => {
    fetchBlockList();
  }, [fetchBlockList]);

  if (isLoading) {
    return (
      <div className="flex justify-center py-4">
        <Spinner />
      </div>
    );
  }

  return (
    <div>
      {error ? (
        <p className="text-sm text-red-600">{error}</p>
      ) : null}

      {blocks.length === 0 ? (
        <p className="text-sm text-gray-500">No blocked users.</p>
      ) : (
        <ul className="divide-y divide-gray-200">
          {blocks.map((block) => (
            <li
              key={block.blockedId}
              className="flex items-center justify-between py-3"
            >
              <span className="text-sm text-gray-900">{block.blockedId}</span>
              <Button
                variant="secondary"
                size="sm"
                onClick={() => unblock(block.blockedId)}
              >
                Unblock
              </Button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export { BlockList };
