import { create } from "zustand";
import { WebSocketMessage, WebSocketService } from "./websoket-service";

import { User } from "@/schemaValidations/user.schema";
import { GuestInfo } from "@/schemaValidations/interface/type_guest";

interface WebSocketState {
  socket: WebSocketService | null;
  isConnected: boolean;
  connect: (user: User | GuestInfo | null, isGuest: boolean) => void;
  disconnect: () => void;
  sendMessage: (message: WebSocketMessage) => void;
  addMessageHandler: (handler: (message: WebSocketMessage) => void) => () => void;
  messageHandlers: Array<(message: WebSocketMessage) => void>;
}


export const useWebSocketStore = create<WebSocketState>((set, get) => ({
  socket: null,
  isConnected: false,
  messageHandlers: [],

  connect: (user: User | GuestInfo | null, isGuest: boolean) => {
    console.log("quananqr1/zusstand/web-socket/websocketStore.ts");
    // const userId = user.id.toString();
    // const userName = user.name;
    const socket = new WebSocketService(
      "9",
      "dung_2024_11_08_12_43_15_0ed49e95-07c3-489f-a6f3-f6a8dcef835a",
      true
    );
    socket.onMessage((message: WebSocketMessage) => {
      const handlers = get().messageHandlers;
      handlers.forEach((handler) => handler(message)); // Call each registered handler
    });

    socket.onConnect(() => set({ isConnected: true }));
    socket.onDisconnect(() => set({ isConnected: false }));
    socket.onMessage((message: WebSocketMessage) => {
      console.log(
        "quananqr1/zusstand/web-socket/websocketStore.ts message",
        message
      ); // Log message for debugging
    });
    set({ socket });
    // if (!get().socket && user) {
    //   const userId = user.id.toString();
    //   const userName = user.name;
    //   const socket = new WebSocketService(userId, userName, isGuest);
    //   socket.onConnect(() => set({ isConnected: true }));
    //   socket.onDisconnect(() => set({ isConnected: false }));
    //   set({ socket });
    // }
  },
  //

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
  },
  addMessageHandler: (handler) => {
    set((state) => ({
      messageHandlers: [...state.messageHandlers, handler],
    }));

    // Return a function to unsubscribe this handler
    return () => {
      set((state) => ({
        messageHandlers: state.messageHandlers.filter((h) => h !== handler),
      }));
    };
  },

}));

// how to use

// // 1. Basic Component with WebSocket Connection
// import React, { useEffect } from 'react';
// import { useWebSocketStore } from './websocket-store';

// const WebSocketComponent: React.FC = () => {
//   const { connect, disconnect, isConnected, socket, sendMessage } = useWebSocketStore();

//   // Connect on component mount
//   useEffect(() => {
//     connect();

//     // Cleanup on unmount
//     return () => {
//       disconnect();
//     };
//   }, [connect, disconnect]);

//   // Example message handler
//   useEffect(() => {
//     if (!socket) return;

//     const cleanup = socket.onMessage((message) => {
//       try {
//         console.log('Received message:', message);
//         switch (message.type) {
//           case 'chat':
//             handleChatMessage(message);
//             break;
//           case 'notification':
//             handleNotification(message);
//             break;
//           default:
//             console.log('Unknown message type:', message.type);
//         }
//       } catch (error) {
//         console.error('Error handling message:', error);
//       }
//     });

//     return () => cleanup();
//   }, [socket]);

//   // Example of handling connection events
//   useEffect(() => {
//     if (!socket) return;

//     const onConnect = () => {
//       console.log('Connected to WebSocket');
//       // Send initial messages or perform setup
//       sendMessage({ type: 'init', data: { userId: 'user123' } });
//     };

//     const onDisconnect = () => {
//       console.log('Disconnected from WebSocket');
//       // Handle cleanup or show reconnection UI
//     };

//     socket.onConnect(onConnect);
//     socket.onDisconnect(onDisconnect);
//   }, [socket, sendMessage]);

//   return (
//     <div>
//       <div className="status">
//         Connection Status: {isConnected ? 'Connected' : 'Disconnected'}
//       </div>

//       {/* Example UI */}
//       <button
//         onClick={() => sendMessage({
//           type: 'chat',
//           data: { message: 'Hello!' }
//         })}
//         disabled={!isConnected}
//       >
//         Send Message
//       </button>
//     </div>
//   );
// };

// // 2. Custom Hook for WebSocket Functionality
// const useWebSocket = () => {
//   const { connect, disconnect, isConnected, socket, sendMessage } = useWebSocketStore();

//   const initialize = React.useCallback(() => {
//     if (!socket) {
//       connect();
//     }
//   }, [connect, socket]);

//   const handleMessage = React.useCallback((handler: (message: WebSocketMessage) => void) => {
//     if (!socket) return () => {};
//     return socket.onMessage(handler);
//   }, [socket]);

//   return {
//     initialize,
//     disconnect,
//     isConnected,
//     sendMessage,
//     handleMessage,
//   };
// };

// // 3. Example Chat Component Using Custom Hook
// const ChatComponent: React.FC = () => {
//   const { initialize, isConnected, sendMessage, handleMessage } = useWebSocket();
//   const [messages, setMessages] = React.useState<string[]>([]);

//   useEffect(() => {
//     initialize();

//     return () => {
//       // Cleanup
//     };
//   }, [initialize]);

//   useEffect(() => {
//     const cleanup = handleMessage((message) => {
//       if (message.type === 'chat') {
//         setMessages(prev => [...prev, message.data.content]);
//       }
//     });

//     return () => cleanup();
//   }, [handleMessage]);

//   const sendChatMessage = (content: string) => {
//     if (isConnected) {
//       sendMessage({
//         type: 'chat',
//         data: { content }
//       });
//     }
//   };

//   // Example error handling
//   const handleError = (error: Error) => {
//     console.error('WebSocket error:', error);
//     // Show error UI or retry connection
//   };

//   return (
//     <div>
//       {/* Chat UI implementation */}
//     </div>
//   );
// };

// // 4. Example of handling different message types
// type MessageHandlers = {
//   [K in WebSocketMessage['type']]: (data: any) => void;
// };

// const messageHandlers: MessageHandlers = {
//   chat: (data) => {
//     console.log('Chat message:', data);
//   },
//   notification: (data) => {
//     console.log('Notification:', data);
//   },
//   error: (data) => {
//     console.error('Error message:', data);
//   }
// };

// // 5. Example of reconnection handling
// const ReconnectionHandler: React.FC = () => {
//   const { isConnected, connect } = useWebSocketStore();
//   const [retryCount, setRetryCount] = useState(0);
//   const maxRetries = 5;

//   useEffect(() => {
//     if (!isConnected && retryCount < maxRetries) {
//       const timeout = setTimeout(() => {
//         console.log(`Attempting to reconnect (${retryCount + 1}/${maxRetries})`);
//         connect();
//         setRetryCount(prev => prev + 1);
//       }, 3000);

//       return () => clearTimeout(timeout);
//     }
//   }, [isConnected, retryCount, connect]);

//   return (
//     <div>
//       {!isConnected && retryCount >= maxRetries && (
//         <div className="error-message">
//           Unable to establish connection. Please try again later.
//           <button onClick={() => {
//             setRetryCount(0);
//             connect();
//           }}>
//             Retry Connection
//           </button>
//         </div>
//       )}
//     </div>
//   );
// };
