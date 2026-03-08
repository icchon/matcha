import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { TagManager } from '@/features/profile/components/TagManager';
import type { Tag } from '@/types';

const mockOnAdd = vi.fn();
const mockOnRemove = vi.fn();

const allTags: Tag[] = [
  { id: 1, name: 'hiking' },
  { id: 2, name: 'cooking' },
  { id: 3, name: 'music' },
  { id: 4, name: 'photography' },
];

const userTags: Tag[] = [
  { id: 1, name: 'hiking' },
];

beforeEach(() => {
  vi.resetAllMocks();
});

describe('TagManager', () => {
  it('renders user tags as badges', () => {
    render(
      <TagManager
        tags={userTags}
        allTags={allTags}
        onAdd={mockOnAdd}
        onRemove={mockOnRemove}
        isLoading={false}
      />,
    );

    expect(
      screen.getByText('hiking'),
      'TagManager should display user tags as visible badges.',
    ).toBeInTheDocument();
  });

  it('renders search input for tags', () => {
    render(
      <TagManager
        tags={userTags}
        allTags={allTags}
        onAdd={mockOnAdd}
        onRemove={mockOnRemove}
        isLoading={false}
      />,
    );

    expect(
      screen.getByPlaceholderText(/search.*tag/i),
      'TagManager should have a search input for filtering tags.',
    ).toBeInTheDocument();
  });

  it('filters available tags by search query', async () => {
    const user = userEvent.setup();
    render(
      <TagManager
        tags={userTags}
        allTags={allTags}
        onAdd={mockOnAdd}
        onRemove={mockOnRemove}
        isLoading={false}
      />,
    );

    await user.type(screen.getByPlaceholderText(/search.*tag/i), 'cook');

    await waitFor(() => {
      expect(
        screen.getByText('cooking'),
        'Autocomplete should show tags matching the search query.',
      ).toBeInTheDocument();
    });
    expect(
      screen.queryByText('music'),
      'Tags not matching the query should be hidden.',
    ).not.toBeInTheDocument();
  });

  it('calls onAdd when a tag suggestion is clicked', async () => {
    const user = userEvent.setup();
    render(
      <TagManager
        tags={userTags}
        allTags={allTags}
        onAdd={mockOnAdd}
        onRemove={mockOnRemove}
        isLoading={false}
      />,
    );

    await user.type(screen.getByPlaceholderText(/search.*tag/i), 'cook');

    await waitFor(() => {
      expect(screen.getByText('cooking')).toBeInTheDocument();
    });

    await user.click(screen.getByText('cooking'));

    expect(
      mockOnAdd,
      'onAdd should be called with the tag ID when a suggestion is clicked.',
    ).toHaveBeenCalledWith(2);
  });

  it('calls onRemove when a tag remove button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <TagManager
        tags={userTags}
        allTags={allTags}
        onAdd={mockOnAdd}
        onRemove={mockOnRemove}
        isLoading={false}
      />,
    );

    const removeButton = screen.getByRole('button', { name: /remove hiking/i });
    await user.click(removeButton);

    expect(
      mockOnRemove,
      'onRemove should be called with the tag ID when remove button is clicked.',
    ).toHaveBeenCalledWith(1);
  });

  it('excludes already-selected tags from suggestions', async () => {
    const user = userEvent.setup();
    render(
      <TagManager
        tags={userTags}
        allTags={allTags}
        onAdd={mockOnAdd}
        onRemove={mockOnRemove}
        isLoading={false}
      />,
    );

    await user.type(screen.getByPlaceholderText(/search.*tag/i), 'hi');

    await waitFor(() => {
      // "hiking" is already selected, should not appear in suggestions
      const suggestions = screen.queryAllByTestId('tag-suggestion');
      const hikingSuggestion = suggestions.find((s) => s.textContent === 'hiking');
      expect(
        hikingSuggestion,
        'Already-selected tags should not appear in the autocomplete suggestions.',
      ).toBeUndefined();
    });
  });
});
