import { logWithLevel } from "@/lib/log";
import { WebSocketMessage } from "@/schemaValidations/interface/type_websocker";

const LOG_PATH = "quananqr1/zusstand/web-socket/websoket-service.ts";

export class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 3000;
  private messageHandlers: ((message: WebSocketMessage) => void)[] = [];
  private connectHandlers: (() => void)[] = [];
  private disconnectHandlers: (() => void)[] = [];
  private userName: string;
  private role: string;
  private userToken: string;
  private tableToken: string;
  private email: string;

  constructor(
    userName: string,
    role: string,
    userToken: string,
    tableToken: string,
    email: string
  ) {
    logWithLevel(
      {
        event: "service_initialization",
        userName,
        role,
        email,
        hasToken: !!userToken,
        hasTableToken: !!tableToken
      },
      LOG_PATH,
      "info",
      1
    );

    this.email = email;
    this.userName = userName;
    this.role = role;
    this.userToken = userToken;
    this.tableToken = tableToken;
    this.connect();
  }

  public connect() {
    try {
      const wsUrl = `ws://localhost:8888/ws/${this.role.toLowerCase()}/${
        this.userName
      }?token=${this.userToken}&tableToken=${this.tableToken}&email=${
        this.email
      }`;

      logWithLevel(
        {
          event: "connection_attempt",
          wsUrl: wsUrl.replace(this.userToken, "***"),
          role: this.role,
          userName: this.userName
        },
        LOG_PATH,
        "debug",
        1
      );

      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        logWithLevel(
          {
            event: "connection_established",
            userName: this.userName,
            role: this.role
          },
          LOG_PATH,
          "info",
          2
        );
        this.reconnectAttempts = 0;
        this.connectHandlers.forEach((handler) => handler());
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage;
          logWithLevel(
            {
              event: "message_received",
              messageType: message.type,
              action: message.action
            },
            LOG_PATH,
            "debug",
            3
          );
          this.messageHandlers.forEach((handler) => handler(message));
        } catch (error) {
          logWithLevel(
            {
              event: "message_parse_error",
              error: "error.message",
              rawData: event.data.slice(0, 100) // Log only first 100 chars for safety
            },
            LOG_PATH,
            "error",
            7
          );
        }
      };

      this.ws.onclose = (event) => {
        logWithLevel(
          {
            event: "connection_closed",
            code: event.code,
            reason: event.reason,
            wasClean: event.wasClean
          },
          LOG_PATH,
          "info",
          2
        );
        this.disconnectHandlers.forEach((handler) => handler());
        if (event.code !== 1000) {
          this.attemptReconnect();
        }
      };

      this.ws.onerror = (error) => {
        logWithLevel(
          {
            event: "websocket_error",
            error: "WebSocket error occurred"
          },
          LOG_PATH,
          "error",
          7
        );
      };
    } catch (error) {
      logWithLevel(
        {
          event: "connection_creation_error",
          error: "error.message"
        },
        LOG_PATH,
        "error",
        7
      );
      this.attemptReconnect();
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const backoffTime =
        this.reconnectTimeout * Math.pow(2, this.reconnectAttempts - 1);

      logWithLevel(
        {
          event: "reconnection_attempt",
          attempt: this.reconnectAttempts,
          maxAttempts: this.maxReconnectAttempts,
          backoffTime
        },
        LOG_PATH,
        "info",
        4
      );

      setTimeout(() => this.connect(), backoffTime);
    } else {
      logWithLevel(
        {
          event: "max_reconnection_attempts_reached",
          attempts: this.reconnectAttempts
        },
        LOG_PATH,
        "error",
        4
      );
    }
  }

  public sendMessage(message: WebSocketMessage) {
    if (!this.ws) {
      logWithLevel(
        {
          event: "send_message_failed",
          reason: "websocket_not_initialized"
        },
        LOG_PATH,
        "error",
        5
      );
      return;
    }

    if (this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
      logWithLevel(
        {
          event: "message_sent",
          messageType: message.type,
          action: message.action
        },
        LOG_PATH,
        "debug",
        5
      );
    } else {
      logWithLevel(
        {
          event: "send_message_failed",
          reason: "invalid_socket_state",
          state: this.ws.readyState
        },
        LOG_PATH,
        "error",
        5
      );
    }
  }

  public onMessage(handler: (message: WebSocketMessage) => void) {
    logWithLevel(
      {
        event: "message_handler_added",
        totalHandlers: this.messageHandlers.length + 1
      },
      LOG_PATH,
      "debug",
      6
    );
    this.messageHandlers.push(handler);
    return () => {
      this.messageHandlers = this.messageHandlers.filter((h) => h !== handler);
      logWithLevel(
        {
          event: "message_handler_removed",
          totalHandlers: this.messageHandlers.length
        },
        LOG_PATH,
        "debug",
        6
      );
    };
  }

  public onConnect(handler: () => void) {
    logWithLevel(
      {
        event: "connect_handler_added",
        totalHandlers: this.connectHandlers.length + 1
      },
      LOG_PATH,
      "debug",
      6
    );
    this.connectHandlers.push(handler);
    return () => {
      this.connectHandlers = this.connectHandlers.filter((h) => h !== handler);
      logWithLevel(
        {
          event: "connect_handler_removed",
          totalHandlers: this.connectHandlers.length
        },
        LOG_PATH,
        "debug",
        6
      );
    };
  }

  public onDisconnect(handler: () => void) {
    logWithLevel(
      {
        event: "disconnect_handler_added",
        totalHandlers: this.disconnectHandlers.length + 1
      },
      LOG_PATH,
      "debug",
      6
    );
    this.disconnectHandlers.push(handler);
    return () => {
      this.disconnectHandlers = this.disconnectHandlers.filter(
        (h) => h !== handler
      );
      logWithLevel(
        {
          event: "disconnect_handler_removed",
          totalHandlers: this.disconnectHandlers.length
        },
        LOG_PATH,
        "debug",
        6
      );
    };
  }

  public disconnect() {
    logWithLevel(
      {
        event: "manual_disconnect_initiated",
        userName: this.userName
      },
      LOG_PATH,
      "info",
      8
    );

    if (this.ws) {
      this.ws.close(1000, "Normal closure");
      this.ws = null;
      logWithLevel(
        {
          event: "disconnect_completed",
          userName: this.userName
        },
        LOG_PATH,
        "info",
        8
      );
    }
  }
}
