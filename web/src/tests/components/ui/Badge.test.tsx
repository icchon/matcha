import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Badge } from '@/components/ui';

describe('Badge', () => {
  it('renders children text', () => {
    render(<Badge>Active</Badge>);
    expect(
      screen.getByText('Active'),
    ).toBeInTheDocument();
  });

  it('renders with default variant styling', () => {
    render(<Badge>Default</Badge>);
    const badge = screen.getByText('Default');
    expect(
      badge.className,
    ).toMatch(/bg-gray/);
  });

  it('renders success variant with green styling', () => {
    render(<Badge variant="success">Online</Badge>);
    const badge = screen.getByText('Online');
    expect(
      badge.className,
      'Success badge should use green background. Check variant class mapping.',
    ).toMatch(/bg-green/);
  });

  it('renders warning variant with yellow styling', () => {
    render(<Badge variant="warning">Pending</Badge>);
    const badge = screen.getByText('Pending');
    expect(
      badge.className,
      'Warning badge should use yellow background. Check variant class mapping.',
    ).toMatch(/bg-yellow/);
  });

  it('renders error variant with red styling', () => {
    render(<Badge variant="error">Offline</Badge>);
    const badge = screen.getByText('Offline');
    expect(
      badge.className,
      'Error badge should use red background. Check variant class mapping.',
    ).toMatch(/bg-red/);
  });

  it('has pill shape with rounded-full class', () => {
    render(<Badge>Status</Badge>);
    const badge = screen.getByText('Status');
    expect(
      badge.classList.contains('rounded-full'),
      'Badge should be pill-shaped using rounded-full class.',
    ).toBe(true);
  });

  it('renders as inline element with appropriate sizing', () => {
    render(<Badge>Small</Badge>);
    const badge = screen.getByText('Small');
    expect(
      badge.classList.contains('inline-flex'),
      'Badge should be inline-flex for inline display.',
    ).toBe(true);
  });
});
