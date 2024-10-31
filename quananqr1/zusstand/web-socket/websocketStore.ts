// import { create } from 'zustand';
// import { WebSocketService } from './websoket-service';


// interface WebSocketState {
//   socket: WebSocketService | null;
//   isConnected: boolean;
//   connect: () => void;
//   disconnect: () => void;
//   sendMessage: (message: any) => void;
// }

// export const useWebSocketStore = create<WebSocketState>((set, get) => ({
//   socket: null,
//   isConnected: false,
  
//   connect: () => {
//     if (!get().socket) {
//       const socket = new WebSocketService();
//       socket.onConnect(() => set({ isConnected: true }));
//       socket.onDisconnect(() => set({ isConnected: false }));
//       set({ socket });
//     }
//   },
  
//   disconnect: () => {
//     const { socket } = get();
//     if (socket) {
//       socket.disconnect();
//       set({ socket: null, isConnected: false });
//     }
//   },
  
//   sendMessage: (message: any) => {
//     const { socket } = get();
//     if (socket) {
//       socket.sendMessage(message);
//     }
//   }
// }));


import { create } from 'zustand';
import { WebSocketService } from './websoket-service';
import { WebSocketMessage } from '@/schemaValidations/interface/type_websocker';


interface WebSocketState {
  socket: WebSocketService | null;
  isConnected: boolean;
  connect: () => void;
  disconnect: () => void;
  sendMessage: (message: WebSocketMessage) => void;
}

export const useWebSocketStore = create<WebSocketState>((set, get) => ({
  socket: null,
  isConnected: false,
  
  connect: () => {
    if (!get().socket) {
      const socket = new WebSocketService();
      socket.onConnect(() => set({ isConnected: true }));
      socket.onDisconnect(() => set({ isConnected: false }));
      set({ socket });
    }
  },
  
  disconnect: () => {
    const { socket } = get();
    if (socket) {
      socket.disconnect();
      set({ socket: null, isConnected: false });
    }
  },
  
  sendMessage: (message: WebSocketMessage) => {
    const { socket } = get();
    if (socket) {
      socket.sendMessage(message);
    }
  }
}));