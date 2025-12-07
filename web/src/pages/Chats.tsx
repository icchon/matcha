import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getChatsForUser, ChatOverview } from '../services/chatService';

const Chats: React.FC = () => {
  const [chats, setChats] = useState<ChatOverview[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchChats = async () => {
      try {
        const response = await getChatsForUser();
        setChats(response.data);
      } catch (error) {
        console.error('Error fetching chats', error);
      } finally {
        setLoading(false);
      }
    };

    fetchChats();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Chats</h1>
      <ul>
        {chats.map((chat) => (
          <li key={chat.other_user.user_id}>
            <Link to={`/chats/${chat.other_user.user_id}`}>
              <div>
                <strong>{chat.other_user.username}</strong>
                <p>{chat.last_message?.content}</p>
              </div>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Chats;
