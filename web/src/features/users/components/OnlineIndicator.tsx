import type { FC } from 'react';

interface OnlineIndicatorProps {
  readonly isOnline: boolean;
  readonly lastConnection: string | null;
}

function formatLastSeen(dateString: string): string {
  const date = new Date(dateString);
  return `Last seen ${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
}

const OnlineIndicator: FC<OnlineIndicatorProps> = ({ isOnline, lastConnection }) => {
  const dotClass = isOnline ? 'bg-green-500' : 'bg-gray-400';
  const label = isOnline
    ? 'Online'
    : lastConnection
      ? formatLastSeen(lastConnection)
      : 'Offline';

  return (
    <span data-testid="online-indicator" className="inline-flex items-center gap-1.5 text-sm">
      <span
        data-testid="status-dot"
        className={`inline-block h-2.5 w-2.5 rounded-full ${dotClass}`}
      />
      <span className="text-gray-600">{label}</span>
    </span>
  );
};

export { OnlineIndicator };
export type { OnlineIndicatorProps };
