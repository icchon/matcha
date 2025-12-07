import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getChatMessages, Message } from '../services/chatService';
import { useWebSocket } from '../contexts/WebSocketContext';

const Chat: React.FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const wsContext = useWebSocket();

  // Effect for fetching initial messages
  useEffect(() => {
    const fetchMessages = async () => {
      if (!userId) return;
      try {
        const response = await getChatMessages(userId);
        setMessages(response.data);
      } catch (error) {
        console.error('Error fetching messages', error);
      }
    };
    fetchMessages();
  }, [userId]);

  // Effect for handling incoming WebSocket messages
  useEffect(() => {
    if (wsContext?.lastMessage) {
      const message = wsContext.lastMessage;
      // Check if the message is for the current chat
      if (message.type === 'chat_message' && (message.sender_id === userId || message.recipient_id === userId)) {
        setMessages((prevMessages) => [...prevMessages, message]);
      }
    }
  }, [wsContext?.lastMessage, userId]);

  const handleSendMessage = () => {
    if (wsContext?.ws && newMessage.trim() && userId) {
      const message = {
        type: 'chat_message',
        recipient_id: userId,
        content: newMessage,
      };
      // Native WebSocket API sends strings
      wsContext.ws.send(JSON.stringify(message));
      
      // Optimistically add sent message to the UI
      const sentMessage: Message = {
        id: Date.now(), // temporary ID
        sender_id: 'me', // temporary sender
        recipient_id: userId,
        content: newMessage,
        sent_at: new Date().toISOString(),
      };
      setMessages(prev => [...prev, sentMessage]);
      setNewMessage('');
    }
  };

  return (
    <div>
      <h1>Chat with User</h1>
      <div>
        {messages.map((msg, index) => (
          <div key={msg.id || index}>
            <strong>{msg.sender_id === 'me' ? 'You' : msg.sender_id}: </strong>
            <span>{msg.content}</span>
          </div>
        ))}
      </div>
      <input
        type="text"
        value={newMessage}
        onChange={(e) => setNewMessage(e.target.value)}
        onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
      />
      <button onClick={handleSendMessage} disabled={!wsContext?.isConnected}>Send</button>
    </div>
  );
};

export default Chat;