import { create } from "zustand";
import { WebSocketService } from "./websoket-service";
import envConfig from "@/config";
import { WebSocketMessage } from "@/schemaValidations/interface/type_websocker";
import { logWithLevel } from "@/lib/log";

const LOG_PATH = "quananqr1/zusstand/web-socket/websocketStore.ts";

interface WebSocketState {
  socket: WebSocketService | null;
  isConnected: boolean;
  wsToken: string | null;
  wsTokenExpiry: string | null;
  connect: (params: {
    userId: string;
    isGuest: boolean;
    userToken: string;
    tableToken: string;
    role: string;
    email: string;
  }) => Promise<void>;
  disconnect: () => void;
  sendMessage: (message: WebSocketMessage) => void;
  addMessageHandler: (
    handler: (message: WebSocketMessage) => void
  ) => () => void;
  messageHandlers: Array<(message: WebSocketMessage) => void>;
  fetchWsToken: (params: {
    userId: number;
    email: string;
    role: string;
  }) => Promise<WsAuthResponse>;
}

interface WsAuthResponse {
  token: string;
  expiresAt: string;
  role: string;
  userId: number;
  email: string;
}

export const useWebSocketStore = create<WebSocketState>((set, get) => ({
  socket: null,
  isConnected: false,
  messageHandlers: [],
  wsToken: null,
  wsTokenExpiry: null,

  fetchWsToken: async ({ userId, email, role }) => {
    const serverEndpoint = envConfig.NEXT_PUBLIC_API_ENDPOINT;

    logWithLevel(
      {
        action: "fetchWsToken",
        userId,
        email,
        role
      },
      LOG_PATH,
      "debug",
      1
    );

    try {
      const response = await fetch(`${serverEndpoint}${envConfig.wsAuth}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          userId,
          email,
          role
        })
      });

      if (!response.ok) {
        logWithLevel(
          {
            error: "Failed to fetch WS token",
            status: response.status,
            statusText: response.statusText
          },
          LOG_PATH,
          "error",
          3
        );
        throw new Error("Failed to fetch WS token");
      }

      const data: WsAuthResponse = await response.json();

      logWithLevel(
        {
          action: "tokenReceived",
          expiresAt: data.expiresAt,
          role: data.role
        },
        LOG_PATH,
        "info",
        2
      );

      set({
        wsToken: data.token,
        wsTokenExpiry: data.expiresAt
      });

      return data;
    } catch (error) {
      logWithLevel(
        {
          error: error instanceof Error ? error.message : "Unknown error",
          userId,
          email
        },
        LOG_PATH,
        "error",
        3
      );
      throw error;
    }
  },

  connect: async ({ userId, isGuest, userToken, tableToken, role, email }) => {
    logWithLevel(
      {
        action: "connectAttempt",
        userId,
        isGuest,
        role,
        email
      },
      LOG_PATH,
      "info",
      4
    );

    const currentTime = new Date();
    const wsTokenExpiry = get().wsTokenExpiry;
    const tokenExpiry = wsTokenExpiry ? new Date(wsTokenExpiry) : null;

    const isTokenExpired =
      !get().wsToken ||
      !wsTokenExpiry ||
      currentTime >= new Date(wsTokenExpiry);

    if (isTokenExpired) {
      logWithLevel(
        {
          action: "tokenRefreshRequired",
          currentTime: currentTime.toISOString(),
          expiry: wsTokenExpiry
        },
        LOG_PATH,
        "debug",
        5
      );

      try {
        await get().fetchWsToken({
          userId: parseInt(userId),
          email: userToken,
          role
        });
      } catch (error) {
        logWithLevel(
          {
            error: error instanceof Error ? error.message : "Unknown error",
            userId,
            role
          },
          LOG_PATH,
          "error",
          3
        );
        return;
      }
    }

    const socket = new WebSocketService(
      userId,
      role,
      userToken,
      tableToken,
      email
    );

    socket.onMessage((message: WebSocketMessage) => {
      logWithLevel(
        {
          action: "messageReceived",
          messageType: message.type
        },
        LOG_PATH,
        "debug",
        6
      );

      const handlers = get().messageHandlers;
      handlers.forEach((handler) => handler(message));
    });

    socket.onConnect(() => {
      logWithLevel(
        {
          action: "connected",
          userId,
          role
        },
        LOG_PATH,
        "info",
        7
      );
      set({ isConnected: true });
    });

    socket.onDisconnect(() => {
      logWithLevel(
        {
          action: "disconnected",
          userId,
          role
        },
        LOG_PATH,
        "info",
        8
      );
      set({ isConnected: false });
    });

    set({ socket });
  },

  disconnect: () => {
    const { socket } = get();
    if (socket) {
      logWithLevel(
        {
          action: "manualDisconnect"
        },
        LOG_PATH,
        "info",
        8
      );

      socket.disconnect();
      set({
        socket: null,
        isConnected: false,
        wsToken: null,
        wsTokenExpiry: null
      });
    }
  },

  sendMessage: (message: WebSocketMessage) => {
    const { socket } = get();
    if (socket) {
      logWithLevel(
        {
          action: "sendMessage",
          messageType: message.type
        },
        LOG_PATH,
        "debug",
        9
      );
      socket.sendMessage(message);
    }
  },

  addMessageHandler: (handler) => {
    logWithLevel(
      {
        action: "addMessageHandler",
        handlersCount: get().messageHandlers.length + 1
      },
      LOG_PATH,
      "debug",
      10
    );

    set((state) => ({
      messageHandlers: [...state.messageHandlers, handler]
    }));

    return () => {
      set((state) => ({
        messageHandlers: state.messageHandlers.filter((h) => h !== handler)
      }));
    };
  }
}));
