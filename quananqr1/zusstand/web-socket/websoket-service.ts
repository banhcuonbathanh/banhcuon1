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
  private userName: string;
  private isGuest: boolean;

  constructor(userId: string, userName: string, isGuest: boolean) {
    this.userId = userId;
    this.userName = userName;
    this.isGuest = isGuest;
    this.connect();
  }
  private createWebSocketUrl(): string {
    const endpoint = envConfig.NEXT_PUBLIC_API_ENDPOINT.replace(
      /^(http|https):\/\//,
      ""
    ).replace(/\/$/, "");
    const wsUrl = `ws://${endpoint}/ws?userId=${
      this.userId
    }&userName=${encodeURIComponent(this.userName)}&isGuest=${this.isGuest}`;
    return wsUrl;
  }

  private connect() {
    const wsUrl = this.createWebSocketUrl();

    const linktest =
      "ws://localhost:8888/ws?userId=9&userName=dung_2024_11_08_12_43_15_0ed49e95-07c3-489f-a6f3-f6a8dcef835a&isGuest=true";
    console.log("Connecting to WebSocket:", wsUrl);
    console.log("quananqr1/zusstand/web-socket/websoket-service.ts connect");
    try {
      this.ws = new WebSocket(linktest);

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

const createWebSocketUrl = (
  userId: string | number,
  userName: string,
  isGuest: boolean
) => {
  // Remove 'http://' or 'https://' from the API endpoint
  const cleanEndpoint = envConfig.NEXT_PUBLIC_API_ENDPOINT.replace(
    /^(http|https):\/\//,
    ""
  );

  const wsUrl = `ws://${cleanEndpoint}/ws?userId=${userId}&userName=${encodeURIComponent(
    userName
  )}&isGuest=${isGuest}`;

  console.log("Connecting to WebSocket:", wsUrl);
  return wsUrl;
};
