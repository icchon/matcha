import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Card } from '@/components/ui';

describe('Card', () => {
  it('renders children content', () => {
    render(<Card><p>Card content</p></Card>);
    expect(
      screen.getByText('Card content'),
    ).toBeInTheDocument();
  });

  it('applies padding for content spacing', () => {
    render(<Card><p>Padded</p></Card>);
    const card = screen.getByText('Padded').closest('div');
    expect(
      card?.className,
      'Card should have padding via Tailwind p-* class.',
    ).toMatch(/p-/);
  });

  it('applies shadow for elevation', () => {
    render(<Card><p>Shadowed</p></Card>);
    const card = screen.getByText('Shadowed').closest('div');
    expect(
      card?.className,
      'Card should have shadow via Tailwind shadow class.',
    ).toMatch(/shadow/);
  });

  it('applies rounded corners', () => {
    render(<Card><p>Rounded</p></Card>);
    const card = screen.getByText('Rounded').closest('div');
    expect(
      card?.className,
      'Card should have rounded corners via Tailwind rounded-* class.',
    ).toMatch(/rounded/);
  });

  it('merges custom className with default classes', () => {
    render(<Card className="mt-4"><p>Custom</p></Card>);
    const card = screen.getByText('Custom').closest('div');
    expect(
      card?.classList.contains('mt-4'),
      'Custom className should be merged into Card element.',
    ).toBe(true);
    expect(
      card?.className,
      'Default shadow class should still be present when custom className is added.',
    ).toMatch(/shadow/);
  });
});
