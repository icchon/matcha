import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Footer } from '@/components/layout/Footer';

describe('Footer', () => {
  it('renders copyright text', () => {
    render(<Footer />);

    expect(
      screen.getByText('© 2025 Matcha. All rights reserved.'),
      'Footer should display copyright text "© 2025 Matcha. All rights reserved.". Check that Footer renders a <footer> with the correct text.',
    ).toBeInTheDocument();
  });

  it('renders as a footer element', () => {
    render(<Footer />);

    const footer = screen.getByRole('contentinfo');
    expect(
      footer,
      'Footer should render a semantic <footer> element (role="contentinfo"). Check the root element tag.',
    ).toBeInTheDocument();
  });
});
