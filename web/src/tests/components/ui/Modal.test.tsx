import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Modal } from '@/components/ui';

describe('Modal', () => {
  it('renders nothing when isOpen is false', () => {
    render(
      <Modal isOpen={false} onClose={vi.fn()} title="Test">
        <p>Content</p>
      </Modal>,
    );
    expect(
      screen.queryByText('Content'),
      'Modal content should not be in the DOM when isOpen is false.',
    ).not.toBeInTheDocument();
  });

  it('renders children when isOpen is true', () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test">
        <p>Modal content</p>
      </Modal>,
    );
    expect(
      screen.getByText('Modal content'),
    ).toBeInTheDocument();
  });

  it('renders the title', () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Confirm Action">
        <p>Are you sure?</p>
      </Modal>,
    );
    expect(
      screen.getByText('Confirm Action'),
    ).toBeInTheDocument();
  });

  it('calls onClose when backdrop is clicked', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();
    render(
      <Modal isOpen={true} onClose={handleClose} title="Test">
        <p>Content</p>
      </Modal>,
    );

    // Click the backdrop (the overlay behind the modal panel)
    const backdrop = screen.getByTestId('modal-backdrop');
    await user.click(backdrop);
    expect(
      handleClose,
      'onClose should be called when clicking the backdrop.',
    ).toHaveBeenCalledTimes(1);
  });

  it('does not call onClose when modal content is clicked', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();
    render(
      <Modal isOpen={true} onClose={handleClose} title="Test">
        <p>Content</p>
      </Modal>,
    );

    await user.click(screen.getByText('Content'));
    expect(
      handleClose,
      'Clicking inside the modal panel should not trigger onClose.',
    ).not.toHaveBeenCalled();
  });

  it('calls onClose when Escape key is pressed', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();
    render(
      <Modal isOpen={true} onClose={handleClose} title="Test">
        <p>Content</p>
      </Modal>,
    );

    await user.keyboard('{Escape}');
    expect(
      handleClose,
      'onClose should be called when Escape key is pressed.',
    ).toHaveBeenCalledTimes(1);
  });

  it('renders in a portal (outside the parent DOM hierarchy)', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={vi.fn()} title="Portal Test">
        <p>Portal content</p>
      </Modal>,
    );
    // The modal content should NOT be a child of the render container
    expect(
      container.querySelector('[data-testid="modal-backdrop"]'),
      'Modal should render in a portal, not inside the render container.',
    ).toBeNull();
    // But it should exist in the document
    expect(
      screen.getByText('Portal content'),
    ).toBeInTheDocument();
  });

  it('has appropriate dialog role for accessibility', () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Accessible Modal">
        <p>Content</p>
      </Modal>,
    );
    expect(
      screen.getByRole('dialog'),
      'Modal should have role="dialog" for accessibility.',
    ).toBeInTheDocument();
  });
});
