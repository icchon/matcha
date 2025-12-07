import React, { useState, useEffect } from 'react';
import { getAllTags, getMyTags, addMyTag, deleteMyTag, Tag } from '../services/commonService';

const Tags: React.FC = () => {
  const [tags, setTags] = useState<Tag[]>([]);
  const [myTags, setMyTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchTags = async () => {
    try {
      const [tagsResponse, myTagsResponse] = await Promise.all([
        getAllTags(),
        getMyTags(),
      ]);
      setTags(tagsResponse.data);
      setMyTags(myTagsResponse.data);
    } catch (error) {
      console.error('Error fetching tags', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTags();
  }, []);

  const handleAddTag = async (tagId: number) => {
    try {
      await addMyTag(tagId);
      fetchTags(); // Refresh tags
    } catch (error) {
      alert('Failed to add tag');
    }
  };

  const handleRemoveTag = async (tagId: number) => {
    try {
      await deleteMyTag(tagId);
      fetchTags(); // Refresh tags
    } catch (error) {
      alert('Failed to remove tag');
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  const myTagIds = new Set(myTags.map((tag) => tag.id));

  return (
    <div>
      <h1>Tags</h1>
      <h2>My Tags</h2>
      <ul>
        {myTags.map((tag) => (
          <li key={tag.id}>
            {tag.name}{' '}
            <button onClick={() => handleRemoveTag(tag.id)}>Remove</button>
          </li>
        ))}
      </ul>

      <h2>All Tags</h2>
      <ul>
        {tags.map((tag) => (
          <li key={tag.id}>
            {tag.name}{' '}
            {!myTagIds.has(tag.id) && (
              <button onClick={() => handleAddTag(tag.id)}>Add</button>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Tags;
