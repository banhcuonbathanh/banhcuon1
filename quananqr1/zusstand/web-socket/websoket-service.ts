// import envConfig from '@/config';

// export class WebSocketService {
//   private ws: WebSocket | null = null;
//   private reconnectAttempts = 0;
//   private maxReconnectAttempts = 5;
//   private reconnectTimeout = 3000;
//   private messageHandlers: ((message: any) => void)[] = [];
//   private connectHandlers: (() => void)[] = [];
//   private disconnectHandlers: (() => void)[] = [];

//   constructor() {
//     this.connect();
//   }

//   private connect() {
//     try {
//       this.ws = new WebSocket(`ws://${envConfig.NEXT_PUBLIC_API_ENDPOINT}/ws`);

//       this.ws.onopen = () => {
//         console.log('WebSocket connected');
//         this.reconnectAttempts = 0;
//         this.connectHandlers.forEach(handler => handler());
//       };

//       this.ws.onmessage = (event) => {
//         try {
//           const message = JSON.parse(event.data);
//           this.messageHandlers.forEach(handler => handler(message));
//         } catch (error) {
//           console.error('Error parsing WebSocket message:', error);
//         }
//       };

//       this.ws.onclose = () => {
//         console.log('WebSocket disconnected');
//         this.disconnectHandlers.forEach(handler => handler());
//         this.attemptReconnect();
//       };

//       this.ws.onerror = (error) => {
//         console.error('WebSocket error:', error);
//       };
//     } catch (error) {
//       console.error('Error creating WebSocket connection:', error);
//       this.attemptReconnect();
//     }
//   }

//   private attemptReconnect() {
//     if (this.reconnectAttempts < this.maxReconnectAttempts) {
//       this.reconnectAttempts++;
//       console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
//       setTimeout(() => this.connect(), this.reconnectTimeout);
//     }
//   }

//    public sendMessage(message: any) {
//     if (this.ws && this.ws.readyState === WebSocket.OPEN) {
//       this.ws.send(JSON.stringify(message));
//     } else {
//       console.error('WebSocket is not connected');
//     }
//   }

//   public onMessage(handler: (message: any) => void) {
//     this.messageHandlers.push(handler);
//     return () => {
//       this.messageHandlers = this.messageHandlers.filter(h => h !== handler);
//     };
//   }

//   public onConnect(handler: () => void) {
//     this.connectHandlers.push(handler);
//     return () => {
//       this.connectHandlers = this.connectHandlers.filter(h => h !== handler);
//     };
//   }

//   public onDisconnect(handler: () => void) {
//     this.disconnectHandlers.push(handler);
//     return () => {
//       this.disconnectHandlers = this.disconnectHandlers.filter(h => h !== handler);
//     };
//   }

//   public disconnect() {
//     if (this.ws) {
//       this.ws.close();
//       this.ws = null;
//     }
//   }
// }

import envConfig from "@/config";
import { WebSocketMessage } from "@/schemaValidations/interface/type_websocker";

export class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 3000;
  private messageHandlers: ((message: WebSocketMessage) => void)[] = [];
  private connectHandlers: (() => void)[] = [];
  private disconnectHandlers: (() => void)[] = [];

  constructor() {
    this.connect();
  }

  private connect() {
    try {
      this.ws = new WebSocket(`ws://${envConfig.NEXT_PUBLIC_API_ENDPOINT}/ws`);

      this.ws.onopen = () => {
        console.log("WebSocket connected");
        this.reconnectAttempts = 0;
        this.connectHandlers.forEach((handler) => handler());
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage;
          this.messageHandlers.forEach((handler) => handler(message));
        } catch (error) {
          console.error("Error parsing WebSocket message:", error);
        }
      };

      this.ws.onclose = () => {
        console.log("WebSocket disconnected");
        this.disconnectHandlers.forEach((handler) => handler());
        this.attemptReconnect();
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    } catch (error) {
      console.error("Error creating WebSocket connection:", error);
      this.attemptReconnect();
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(
        `Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`
      );
      setTimeout(() => this.connect(), this.reconnectTimeout);
    }
  }

  public sendMessage(message: WebSocketMessage) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not connected");
    }
  }

  public onMessage(handler: (message: WebSocketMessage) => void) {
    this.messageHandlers.push(handler);
    return () => {
      this.messageHandlers = this.messageHandlers.filter((h) => h !== handler);
    };
  }

  public onConnect(handler: () => void) {
    this.connectHandlers.push(handler);
    return () => {
      this.connectHandlers = this.connectHandlers.filter((h) => h !== handler);
    };
  }

  public onDisconnect(handler: () => void) {
    this.disconnectHandlers.push(handler);
    return () => {
      this.disconnectHandlers = this.disconnectHandlers.filter(
        (h) => h !== handler
      );
    };
  }

  public disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
