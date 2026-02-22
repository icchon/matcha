import { create } from 'zustand';
import { connect, disconnect, send } from '@/features/websocket/connection';
import { registerHandler, unregisterHandler, clearAllHandlers } from '@/features/websocket/messageRouter';
import type { ConnectionStatus, WsMessage, MessageHandler } from '@/features/websocket/types';

export type { WsMessage, ConnectionStatus, MessageHandler };

interface WsState {
  readonly connectionStatus: ConnectionStatus;
  readonly error: string | null;
}

interface WsActions {
  readonly connect: (url: string) => void;
  readonly disconnect: () => void;
  readonly send: (message: WsMessage) => void;
  readonly registerHandler: (type: string, handler: MessageHandler) => void;
  readonly unregisterHandler: (type: string, handler: MessageHandler) => void;
  readonly clearAllHandlers: () => void;
}

type WsStore = WsState & WsActions;

export const useWsStore = create<WsStore>()((set) => ({
  connectionStatus: 'disconnected',
  error: null,

  connect: (url: string) => {
    connect(url, set);
  },

  disconnect: () => {
    disconnect(set);
  },

  send: (message: WsMessage) => {
    send(message);
  },

  registerHandler: (type: string, handler: MessageHandler) => {
    registerHandler(type, handler);
  },

  unregisterHandler: (type: string, handler: MessageHandler) => {
    unregisterHandler(type, handler);
  },

  clearAllHandlers: () => {
    clearAllHandlers();
  },
}));
