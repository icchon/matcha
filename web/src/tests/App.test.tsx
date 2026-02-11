import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import App from '@/App';

describe('App', () => {
  it('renders the Matcha heading on the home page', () => {
    render(<App />);

    expect(
      screen.getByText('Matcha'),
      'The home page should display a heading with text "Matcha". Check that HomePage renders an h1 with "Matcha".',
    ).toBeInTheDocument();
  });
});
