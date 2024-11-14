import envConfig from "@/config";
import { CreateOrderRequest } from "@/schemaValidations/interface/type_order";

export interface OrderPayload {
  orderId: number;

  orderData: CreateOrderRequest;
}

export type WebSocketMessage =
  | {
      type: "NEW_ORDER";
      data: OrderPayload;
    }
  | {
      type: "ORDER_STATUS_UPDATE";
      data: {
        orderId: number;
        status: string;
        timestamp: string;
      };
    };
// Add other message types as needed
export class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 3000;
  private messageHandlers: ((message: WebSocketMessage) => void)[] = [];
  private connectHandlers: (() => void)[] = [];
  private disconnectHandlers: (() => void)[] = [];
  private userId: string;
  private role: string;

  private user_Token: string;
  private table_Token: string;

  constructor(
    userId: string,
    role: string,

    user_Token: string,
    table_Token: string
  ) {
    this.userId = userId;
    this.role = role;

    this.user_Token = user_Token;
    this.table_Token = table_Token;
    this.connect();
  }

  private connect() {
    const link = `${envConfig.wslink}/${this.role}/${this.userId}?token=${this.user_Token}&tableToken=${this.table_Token}`;
 "   ws://localhost:8888/ws/admin/1?token=abc123&tableToken=table455"
    console.log("Connecting to WebSocket:", link);
    console.log("quananqr1/zusstand/web-socket/websoket-service.ts connect");
    try {
      this.ws = new WebSocket(link);

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

      this.ws.onclose = (event) => {
        console.log("WebSocket disconnected:", event.code, event.reason);
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

