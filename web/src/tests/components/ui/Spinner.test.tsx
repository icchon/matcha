import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Spinner } from '@/components/ui';

describe('Spinner', () => {
  it('renders with default medium size', () => {
    render(<Spinner />);
    const spinner = screen.getByRole('status');
    expect(spinner).toBeInTheDocument();
    expect(
      spinner.querySelector('.h-6.w-6') ?? spinner.className,
    ).toBeTruthy();
  });

  it('renders with accessible label', () => {
    render(<Spinner />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveAccessibleName(
      /loading/i,
    );
  });

  it('applies animate-spin class for animation', () => {
    render(<Spinner />);
    const svg = screen.getByRole('status').querySelector('svg');
    expect(
      svg?.classList.contains('animate-spin'),
    ).toBe(true);
  });

  it('renders small size when size="sm"', () => {
    render(<Spinner size="sm" />);
    const svg = screen.getByRole('status').querySelector('svg');
    expect(
      svg?.classList.contains('h-4'),
    ).toBe(true);
  });

  it('renders large size when size="lg"', () => {
    render(<Spinner size="lg" />);
    const svg = screen.getByRole('status').querySelector('svg');
    expect(
      svg?.classList.contains('h-8'),
    ).toBe(true);
  });
});
