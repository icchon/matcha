import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from '@/components/ui';

describe('Button', () => {
  it('renders children text', () => {
    render(<Button>Click me</Button>);
    expect(
      screen.getByRole('button', { name: 'Click me' }),
    ).toBeInTheDocument();
  });

  it('calls onClick handler when clicked', async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click</Button>);

    await user.click(screen.getByRole('button'));
    expect(
      handleClick,
      'onClick should be called once after a single click.',
    ).toHaveBeenCalledTimes(1);
  });

  it('does not call onClick when disabled', async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(<Button disabled onClick={handleClick}>Disabled</Button>);

    await user.click(screen.getByRole('button'));
    expect(
      handleClick,
      'onClick should not fire when button is disabled.',
    ).not.toHaveBeenCalled();
  });

  it('renders as disabled when disabled prop is true', () => {
    render(<Button disabled>Disabled</Button>);
    expect(
      screen.getByRole('button'),
    ).toBeDisabled();
  });

  // Loading state
  it('shows spinner and disables button when loading', () => {
    render(<Button loading>Submit</Button>);
    const button = screen.getByRole('button');
    expect(
      button,
      'Button should be disabled while loading.',
    ).toBeDisabled();
    expect(
      button.querySelector('[role="status"]') ?? button.querySelector('.animate-spin'),
      'Button should display a spinner element when loading is true.',
    ).toBeTruthy();
  });

  it('still shows button text when loading', () => {
    render(<Button loading>Submit</Button>);
    expect(
      screen.getByRole('button', { name: /submit/i }),
      'Button text should remain visible during loading state.',
    ).toBeInTheDocument();
  });

  // Variants
  it('applies primary variant styles by default', () => {
    render(<Button>Primary</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Default variant should be primary with blue background.',
    ).toMatch(/bg-blue/);
  });

  it('applies secondary variant styles', () => {
    render(<Button variant="secondary">Secondary</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Secondary variant should use gray background.',
    ).toMatch(/bg-gray/);
  });

  it('applies danger variant styles', () => {
    render(<Button variant="danger">Delete</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Danger variant should use red background.',
    ).toMatch(/bg-red/);
  });

  it('applies ghost variant styles', () => {
    render(<Button variant="ghost">Ghost</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Ghost variant should use bg-transparent.',
    ).toMatch(/bg-transparent/);
  });

  // Sizes
  it('applies medium size by default', () => {
    render(<Button>Medium</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Default size should be md with px-4 py-2.',
    ).toMatch(/px-4/);
  });

  it('applies small size', () => {
    render(<Button size="sm">Small</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Small size should use smaller padding like px-3.',
    ).toMatch(/px-3/);
  });

  it('applies large size', () => {
    render(<Button size="lg">Large</Button>);
    const button = screen.getByRole('button');
    expect(
      button.className,
      'Large size should use larger padding like px-6.',
    ).toMatch(/px-6/);
  });

  it('forwards additional HTML button attributes', () => {
    render(<Button type="submit" aria-label="Submit form">Go</Button>);
    const button = screen.getByRole('button');
    expect(button).toHaveAttribute('type', 'submit');
    expect(button).toHaveAttribute('aria-label', 'Submit form');
  });
});
