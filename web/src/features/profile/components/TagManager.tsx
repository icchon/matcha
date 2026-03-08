import { useState, useMemo, useCallback, type FC } from 'react';
import { Badge } from '@/components/ui/Badge';
import type { Tag } from '@/types';

interface TagManagerProps {
  readonly tags: readonly Tag[];
  readonly allTags: readonly Tag[];
  readonly onAdd: (tagId: number) => void;
  readonly onRemove: (tagId: number) => void;
  readonly isLoading: boolean;
}

const TagManager: FC<TagManagerProps> = ({ tags, allTags, onAdd, onRemove, isLoading }) => {
  const [searchQuery, setSearchQuery] = useState('');

  const selectedTagIds = useMemo(
    () => new Set(tags.map((t) => t.id)),
    [tags],
  );

  const filteredSuggestions = useMemo(() => {
    if (searchQuery.length === 0) return [];
    return allTags.filter(
      (tag) =>
        !selectedTagIds.has(tag.id) &&
        tag.name.toLowerCase().includes(searchQuery.toLowerCase()),
    );
  }, [allTags, selectedTagIds, searchQuery]);

  const handleAdd = useCallback(
    (tagId: number) => {
      onAdd(tagId);
      setSearchQuery('');
    },
    [onAdd],
  );

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium">Interest Tags</h3>

      <div className="flex flex-wrap gap-2">
        {tags.map((tag) => (
          <Badge key={tag.id}>
            <span className="flex items-center gap-1">
              {tag.name}
              <button
                type="button"
                aria-label={`Remove ${tag.name}`}
                className="ml-1 text-gray-500 hover:text-gray-700"
                onClick={() => onRemove(tag.id)}
                disabled={isLoading}
              >
                &times;
              </button>
            </span>
          </Badge>
        ))}
      </div>

      <div className="relative">
        <input
          type="text"
          placeholder="Search tags..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          disabled={isLoading}
        />
        {filteredSuggestions.length > 0 ? (
          <ul className="absolute z-10 mt-1 max-h-48 w-full overflow-y-auto rounded-md border border-gray-200 bg-white shadow-lg">
            {filteredSuggestions.map((tag) => (
              <li key={tag.id}>
                <button
                  type="button"
                  data-testid="tag-suggestion"
                  className="w-full px-3 py-2 text-left text-sm hover:bg-gray-100"
                  onClick={() => handleAdd(tag.id)}
                  disabled={isLoading}
                >
                  {tag.name}
                </button>
              </li>
            ))}
          </ul>
        ) : null}
      </div>
    </div>
  );
};

export { TagManager };
