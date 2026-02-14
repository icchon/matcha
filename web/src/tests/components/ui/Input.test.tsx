import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createRef } from 'react';
import { Input } from '@/components/ui';

describe('Input', () => {
  it('renders an input element', () => {
    render(<Input />);
    expect(
      screen.getByRole('textbox'),
    ).toBeInTheDocument();
  });

  it('renders a label when label prop is provided', () => {
    render(<Input label="Email" />);
    expect(
      screen.getByLabelText('Email'),
      'Input should be associated with label via htmlFor/id.',
    ).toBeInTheDocument();
  });

  it('renders error message when error prop is provided', () => {
    render(<Input label="Email" error="Email is required" />);
    expect(
      screen.getByText('Email is required'),
      'Error message should be visible when error prop is set.',
    ).toBeInTheDocument();
  });

  it('does not render error message when error is not provided', () => {
    render(<Input label="Email" />);
    expect(
      screen.queryByRole('alert'),
    ).not.toBeInTheDocument();
  });

  it('applies error styling to input when error is present', () => {
    render(<Input label="Email" error="Required" />);
    const input = screen.getByRole('textbox');
    expect(
      input.className,
      'Input should have red border when error is present.',
    ).toMatch(/border-red/);
  });

  it('accepts user input', async () => {
    const user = userEvent.setup();
    render(<Input label="Name" />);
    const input = screen.getByRole('textbox');

    await user.type(input, 'John');
    expect(input).toHaveValue('John');
  });

  it('forwards ref for react-hook-form compatibility', () => {
    const ref = createRef<HTMLInputElement>();
    render(<Input ref={ref} label="Test" />);
    expect(
      ref.current,
      'Input should forward ref to the underlying input element.',
    ).toBeInstanceOf(HTMLInputElement);
  });

  it('forwards standard input props like placeholder', () => {
    render(<Input placeholder="Enter email" />);
    expect(
      screen.getByPlaceholderText('Enter email'),
    ).toBeInTheDocument();
  });

  it('associates error message with input via aria-describedby', () => {
    render(<Input label="Email" error="Invalid email" />);
    const input = screen.getByRole('textbox');
    expect(
      input,
      'Input should reference error message via aria-describedby for accessibility.',
    ).toHaveAttribute('aria-describedby');
    const describedById = input.getAttribute('aria-describedby');
    const errorEl = document.getElementById(describedById!);
    expect(errorEl?.textContent).toBe('Invalid email');
  });

  it('marks input as aria-invalid when error is present', () => {
    render(<Input label="Email" error="Required" />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveAttribute('aria-invalid', 'true');
  });
});
