import type { FC, ReactNode } from 'react';

interface CardProps {
  readonly children: ReactNode;
  readonly className?: string;
}

const Card: FC<CardProps> = ({ children, className = '' }) => {
  return (
    <div className={`rounded-lg bg-white p-6 shadow ${className}`.trim()}>
      {children}
    </div>
  );
};

export { Card };
export type { CardProps };
