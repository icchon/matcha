import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useAuth } from './AuthContext';

interface WebSocketContextType {
    ws: WebSocket | null;
    isConnected: boolean;
    lastMessage: any; // Keep it simple for now
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export const useWebSocket = () => {
    return useContext(WebSocketContext);
};

interface WebSocketProviderProps {
    children: ReactNode;
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children }) => {
    const auth = useAuth();
    const [ws, setWs] = useState<WebSocket | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    const [lastMessage, setLastMessage] = useState<any>(null);

    useEffect(() => {
        if (!auth?.token) {
            if (ws) {
                ws.close();
            }
            return;
        }

        // Construct the URL with the token as a query parameter
        const url = `ws://localhost/ws?token=${auth.token}`;
        const newWs = new WebSocket(url);

        newWs.onopen = () => {
            console.log('Native WebSocket Connected');
            setIsConnected(true);
        };

        newWs.onclose = () => {
            console.log('Native WebSocket Disconnected');
            setIsConnected(false);
            setWs(null);
        };

        newWs.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        newWs.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                setLastMessage(message);
            } catch (e) {
                console.error('Failed to parse WebSocket message:', event.data);
            }
        };

        setWs(newWs);

        return () => {
            console.log('Closing Native WebSocket connection');
            newWs.close();
        };
    }, [auth?.token]);

    const value = {
        ws,
        isConnected,
        lastMessage,
    };

    return (
        <WebSocketContext.Provider value={value}>
            {children}
        </WebSocketContext.Provider>
    );
};