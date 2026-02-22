import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { OnlineIndicator } from '@/features/users/components/OnlineIndicator';

describe('OnlineIndicator', () => {
  it('shows green dot and "Online" text when user is online', () => {
    render(<OnlineIndicator isOnline={true} lastConnection={null} />);

    const indicator = screen.getByTestId('online-indicator');
    expect(
      indicator.textContent,
      'Should display "Online" when isOnline is true.',
    ).toContain('Online');
    expect(
      indicator.querySelector('[data-testid="status-dot"]')?.classList.toString(),
      'Should have green background class when online.',
    ).toContain('bg-green');
  });

  it('shows grey dot and "Offline" when user is offline with no last connection', () => {
    render(<OnlineIndicator isOnline={false} lastConnection={null} />);

    const indicator = screen.getByTestId('online-indicator');
    expect(
      indicator.textContent,
      'Should display "Offline" when isOnline is false and no lastConnection.',
    ).toContain('Offline');
    expect(
      indicator.querySelector('[data-testid="status-dot"]')?.classList.toString(),
      'Should have gray background class when offline.',
    ).toContain('bg-gray');
  });

  it('shows last seen time when offline with lastConnection', () => {
    const lastSeen = '2024-01-15T10:30:00Z';
    render(<OnlineIndicator isOnline={false} lastConnection={lastSeen} />);

    const indicator = screen.getByTestId('online-indicator');
    expect(
      indicator.textContent,
      'Should display "Last seen" with formatted date when offline with lastConnection.',
    ).toContain('Last seen');
  });
});
