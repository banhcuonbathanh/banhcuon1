// config/loggerConfig.ts

export interface LogPath {
  path: string;
  enabled: boolean;
  description: string;
  enabledLogIds: number[];
  logDescriptions: {
    [key: number]: {
      description: string;
      location: string;
      status: "enabled" | "disabled";
    };
  };
}

export const loggerPaths: LogPath[] = [
  {
    path: "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
    enabled: false,
    description: "Restaurant Summary Component Logs",
    enabledLogIds: [1, 2, 3],
    logDescriptions: {
      1: {
        description: "Log initial dish aggregation state",
        location: "aggregateDishes function - initialization",
        status: "enabled"
      },
      2: {
        description: "Log aggregated dishes for order groups",
        location: "RestaurantSummary component - groupedOrders processing",
        status: "enabled"
      },
      3: {
        description: "Log aggregation completion state",
        location: "RestaurantSummary component - final state",
        status: "enabled"
      }
    }
  },
  {
    path: "/manage/admin/table",
    enabled: false,
    description: "Admin Table Management Logs",
    enabledLogIds: [1, 2, 3],
    logDescriptions: {
      1: {
        description: "Table component initialization state",
        location: "Table component - mount phase",
        status: "enabled"
      },
      2: {
        description: "Table data state changes",
        location: "Table component - data updates",
        status: "enabled"
      },
      3: {
        description: "Table component error handling",
        location: "Table component - error boundaries",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx",
    enabled: false,
    description: "Dishes Summary Component Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8],
    logDescriptions: {
      1: {
        description: "Delivery button click tracking",
        location: "handleDeliveryClick function",
        status: "enabled"
      },
      2: {
        description: "Keypad submission attempt",
        location: "handleKeypadSubmit function - initial validation",
        status: "enabled"
      },
      3: {
        description: "Delivery quantity validation error",
        location: "handleKeypadSubmit function - quantity validation",
        status: "enabled"
      },
      4: {
        description: "Delivery creation success",
        location: "handleKeypadSubmit function - success path",
        status: "enabled"
      },
      5: {
        description: "Delivery creation error",
        location: "handleKeypadSubmit function - error handling",
        status: "enabled"
      },
      6: {
        description: "Reference order validation",
        location: "handleKeypadSubmit function - order validation",
        status: "enabled"
      },
      7: {
        description: "Guest order processing",
        location: "handleKeypadSubmit function - guest handling",
        status: "enabled"
      },
      8: {
        description: "Delivery details compilation",
        location: "handleKeypadSubmit function - delivery preparation",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/app/(client)/table/[number]/component/order/logic.ts",
    enabled: false,
    description: "Order and Delivery Store State Management Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
    logDescriptions: {
      1: {
        description: "Initial order/delivery data preparation",
        location: "createOrder/createDelivery function - initialization",
        status: "enabled"
      },
      2: {
        description: "API request monitoring",
        location: "createOrder/createDelivery function - API communication",
        status: "enabled"
      },
      3: {
        description: "Item state changes",
        location: "Item manipulation functions",
        status: "enabled"
      },
      4: {
        description: "Price calculation events",
        location: "Price calculation functions",
        status: "enabled"
      },
      5: {
        description: "Status change events",
        location: "Status management functions",
        status: "enabled"
      },
      6: {
        description: "State persistence events",
        location: "Persistence middleware",
        status: "enabled"
      },
      7: {
        description: "Error handling",
        location: "Error handling middleware",
        status: "enabled"
      },
      8: {
        description: "General state updates",
        location: "State update functions",
        status: "enabled"
      },
      9: {
        description: "Order summary state changes",
        location: "Order summary processing",
        status: "enabled"
      },
      10: {
        description: "Authentication state validation",
        location: "Auth validation middleware",
        status: "enabled"
      },
      11: {
        description: "Data transformation",
        location: "createOrder/createDelivery function - data transformation",
        status: "enabled"
      },
      12: {
        description: "Request validation",
        location: "createOrder/createDelivery function - request validation",
        status: "enabled"
      },
      13: {
        description: "WebSocket communication",
        location: "WebSocket message handling",
        status: "enabled"
      },
      14: {
        description: "User/Guest validation",
        location: "Authentication validation",
        status: "enabled"
      },
      15: {
        description: "Order cleanup",
        location: "Order completion and cleanup",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/zusstand/web-socket/websocketStore.ts",
    enabled: false,
    description: "WebSocket Store Management Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    logDescriptions: {
      1: {
        description: "WebSocket token fetch initialization",
        location: "fetchWsToken function - start",
        status: "enabled"
      },
      2: {
        description: "WebSocket token received successfully",
        location: "fetchWsToken function - success",
        status: "enabled"
      },
      3: {
        description: "WebSocket token fetch error",
        location: "fetchWsToken function - error handling",
        status: "enabled"
      },
      4: {
        description: "WebSocket connection attempt",
        location: "connect function - initialization",
        status: "enabled"
      },
      5: {
        description: "WebSocket token refresh check",
        location: "connect function - token validation",
        status: "enabled"
      },
      6: {
        description: "WebSocket message received",
        location: "socket.onMessage handler",
        status: "enabled"
      },
      7: {
        description: "WebSocket connection established",
        location: "socket.onConnect handler",
        status: "enabled"
      },
      8: {
        description: "WebSocket disconnection",
        location: "socket.onDisconnect/disconnect function",
        status: "enabled"
      },
      9: {
        description: "WebSocket message sent",
        location: "sendMessage function",
        status: "enabled"
      },
      10: {
        description: "Message handler management",
        location: "addMessageHandler function",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/app/(client)/table/[number]/component/order/order.tsx",
    enabled: false,
    description: "Client Order Component Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8],
    logDescriptions: {
      1: {
        description: "Table token decode result",
        location: "OrderSummary component - initialization",
        status: "enabled"
      },
      2: {
        description: "Table number and token state updates",
        location: "OrderSummary useEffect hook",
        status: "enabled"
      },
      3: {
        description: "Bowl quantity changes",
        location: "handleBowlChange function",
        status: "enabled"
      },
      4: {
        description: "Chili preference updates",
        location: "Chili option click handler",
        status: "enabled"
      },
      5: {
        description: "Filling selection changes",
        location: "Filling selection handlers",
        status: "enabled"
      },
      6: {
        description: "Order summary calculation",
        location: "getOrderSummary function call",
        status: "enabled"
      },
      7: {
        description: "Table number conversion",
        location: "addTableNumberconvert function",
        status: "enabled"
      },
      8: {
        description: "Component state changes",
        location: "OrderSummary state updates",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/components/form/login-dialog.tsx",
    enabled: false,
    description: "Login Dialog Component Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6],
    logDescriptions: {
      1: {
        description: "Component initialization and pathname changes",
        location: "LoginDialog1 component - initialization and useEffect",
        status: "enabled"
      },
      2: {
        description: "Form submission attempts",
        location: "onSubmit function",
        status: "enabled"
      },
      3: {
        description: "Login API response handling",
        location: "onSubmit function - API response",
        status: "enabled"
      },
      4: {
        description: "Dialog state changes",
        location: "Dialog state handlers",
        status: "enabled"
      },
      5: {
        description: "Form validation errors",
        location: "Form validation and error handling",
        status: "enabled"
      },
      6: {
        description: "Navigation events",
        location: "handleLoginRedirect and navigation handlers",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/app/(client)/table/[number]/component/order/add_order_button.tsx",
    enabled: false,
    description: "Order Creation Component Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8, 9],
    logDescriptions: {
      1: {
        description: "Component initialization and authentication state",
        location: "OrderCreationComponent - initialization",
        status: "enabled"
      },
      2: {
        description: "WebSocket connection management",
        location: "initializeWebSocket function",
        status: "enabled"
      },
      3: {
        description: "Order creation process",
        location: "handleCreateOrder function",
        status: "enabled"
      },
      4: {
        description: "WebSocket message handling",
        location: "sendMessage1 function",
        status: "enabled"
      },
      5: {
        description: "Authentication state changes",
        location: "Authentication state handlers",
        status: "enabled"
      },
      6: {
        description: "Order validation",
        location: "Order validation checks",
        status: "enabled"
      },
      7: {
        description: "WebSocket token management",
        location: "WebSocket token handling",
        status: "enabled"
      },
      8: {
        description: "Component cleanup",
        location: "Cleanup useEffect",
        status: "enabled"
      },
      9: {
        description: "User identification process",
        location: "getEmailIdentifier function",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/zusstand/web-socket/websoket-service.ts",
    enabled: false,
    description: "WebSocket Service Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8],
    logDescriptions: {
      1: {
        description: "WebSocket initialization and connection attempts",
        location: "WebSocketService constructor and connect method",
        status: "enabled"
      },
      2: {
        description: "WebSocket connection state changes",
        location: "WebSocket event handlers (onopen, onclose)",
        status: "enabled"
      },
      3: {
        description: "Message handling and processing",
        location: "onmessage handler and message processing",
        status: "enabled"
      },
      4: {
        description: "Reconnection attempts and backoff",
        location: "attemptReconnect method",
        status: "enabled"
      },
      5: {
        description: "Message sending operations",
        location: "sendMessage method",
        status: "enabled"
      },
      6: {
        description: "Event handler management",
        location: "Handler registration methods",
        status: "enabled"
      },
      7: {
        description: "Error handling and validation",
        location: "Error handlers and validation checks",
        status: "enabled"
      },
      8: {
        description: "Service cleanup and disconnection",
        location: "disconnect method",
        status: "enabled"
      }
    }
  }
];

const isDevelopment = process.env.NODE_ENV !== "production";

type LogLevel = "debug" | "info" | "warn" | "error";

function validateLogId(pathConfig: LogPath, logId: number): void {
  if (!pathConfig.logDescriptions[logId]) {
    throw new Error(`
      Invalid log ID: ${logId}
      Available log IDs for ${pathConfig.path}:
      ${Object.keys(pathConfig.logDescriptions)
        .map(
          (id) =>
            `\n- Log #${id}: ${
              pathConfig.logDescriptions[Number(id)].description
            }`
        )
        .join("")}
    `);
  }
}

function validatePath(path: string): LogPath {
  const pathConfig = loggerPaths.find((p) => path.startsWith(p.path));
  if (!pathConfig) {
    throw new Error(`
      Invalid path: ${path}
      Available paths:
      ${loggerPaths.map((p) => `\n- ${p.path}`).join("")}
    `);
  }
  return pathConfig;
}

export function getLogStatus(path: string): void {
  const pathConfig = validatePath(path);
  console.log(`
Log Status for: ${pathConfig.path}
Description: ${pathConfig.description}
Enabled: ${pathConfig.enabled}

Available Logs:
${Object.entries(pathConfig.logDescriptions)
  .map(
    ([id, log]) => `
Log #${id}:
- Description: ${log.description}
- Location: ${log.location}
- Status: ${log.status}
- Enabled: ${pathConfig.enabledLogIds.includes(Number(id)) ? "Yes" : "No"}

`
  )
  .join("")}
  `);
}

export function logWithLevel(
  message: Record<string, unknown>,
  path: string,
  level: LogLevel,
  logId: number
): void {
  // Validate message
  if (!message || typeof message !== "object") {
    throw new Error(`
      message: Must be an object (e.g., { dishMap }, { users })
      Received: ${JSON.stringify(message)}
    `);
  }

  // Validate path and get config
  const pathConfig = validatePath(path);

  // Validate log ID
  validateLogId(pathConfig, logId);

  // Check if logging is enabled for this path and log ID
  if (!pathConfig.enabled || !pathConfig.enabledLogIds.includes(logId)) {
    return;
  }

  // Get log description
  const logInfo = pathConfig.logDescriptions[logId];

  // Format the log prefix
  const logIdPrefix = `[Log #${logId}]`;
  const locationPrefix = `[${logInfo.location}]`;

  // Log based on level
  switch (level) {
    case "debug":
      isDevelopment &&
        console.log(
          `🔍 DEBUG ${logIdPrefix} ${locationPrefix} [${path}]:`,
          message
        );
      break;
    case "info":
      isDevelopment &&
        console.log(
          `ℹ️ INFO ${logIdPrefix} ${locationPrefix} [${path}]:`,
          message
        );
      break;
    case "warn":
      console.log(
        `⚠️ WARN ${logIdPrefix} ${locationPrefix} [${path}]:`,
        message
      );
      break;
    case "error":
      console.log(
        `❌ ERROR ${logIdPrefix} ${locationPrefix} [${path}]:`,
        message
      );
      break;
  }
}


// summary 
// # Logger Configuration Summary

// ## System Overview 🔧
// - **Global Status**: All paths currently disabled
// - **Log Levels**: `debug`, `info`, `warn`, `error`
// - **Environment Control**: Development logs controlled by `isDevelopment` flag

// ## Components Breakdown 📑

// ### 1. Restaurant Summary Component 🍽️
// - **Path**: `quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx`
// - **Status**: Disabled
// - **Active Log IDs**: 1-3
// - **Purpose**: Tracks dish aggregation and order group processing

// ### 2. Admin Table Management 📊
// - **Path**: `/manage/admin/table`
// - **Status**: Disabled
// - **Active Log IDs**: 1-3
// - **Purpose**: Monitors table component lifecycle and data updates

// ### 3. Dishes Summary Component 🍜
// - **Path**: `quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx`
// - **Status**: Disabled
// - **Active Log IDs**: 1-8
// - **Purpose**: Handles delivery tracking and order processing

// ### 4. Order/Delivery Store Management 📦
// - **Path**: `quananqr1/app/(client)/table/[number]/component/order/logic.ts`
// - **Status**: Disabled
// - **Active Log IDs**: 1-15
// - **Purpose**: Complete order lifecycle management and state tracking

// ### 5. WebSocket Store 🌐
// - **Path**: `quananqr1/zusstand/web-socket/websocketStore.ts`
// - **Status**: Disabled
// - **Active Log IDs**: 1-10
// - **Purpose**: WebSocket communication and token management

// ### 6. Client Order Component 🛒
// - **Path**: `quananqr1/app/(client)/table/[number]/component/order/order.tsx`
// - **Status**: Disabled
// - **Active Log IDs**: 1-8
// - **Purpose**: Customer order management and UI state tracking

// ### 7. Login Dialog 🔐
// - **Path**: `quananqr1/components/form/login-dialog.tsx`
// - **Status**: Disabled
// - **Active Log IDs**: 1-6
// - **Purpose**: Authentication flow and form submission tracking

// ### 8. Order Creation Button 📝
// - **Path**: `quananqr1/app/(client)/table/[number]/component/order/add_order_button.tsx`
// - **Status**: Disabled
// - **Active Log IDs**: 1-9
// - **Purpose**: Order creation process and authentication state management

// ### 9. WebSocket Service ⚡
// - **Path**: `quananqr1/zusstand/web-socket/websoket-service.ts`
// - **Status**: Disabled
// - **Active Log IDs**: 1-8
// - **Purpose**: Core WebSocket service functionality and connection management

// ## Quick Reference Table 📋

// | Component               | Log IDs | Status  | Path Type    |
// |------------------------|---------|----------|-------------|
// | Restaurant Summary     | 1-3     | Disabled | Admin       |
// | Admin Table            | 1-3     | Disabled | Admin       |
// | Dishes Summary         | 1-8     | Disabled | Admin       |
// | Order Store            | 1-15    | Disabled | Client      |
// | WebSocket Store        | 1-10    | Disabled | Core        |
// | Client Order           | 1-8     | Disabled | Client      |
// | Login Dialog           | 1-6     | Disabled | Auth        |
// | Order Creation         | 1-9     | Disabled | Client      |
// | WebSocket Service      | 1-8     | Disabled | Core        |

// ## Notes 📌
// 1. All components have properly configured log IDs
// 2. Each path includes comprehensive error handling
// 3. Logging levels are consistently implemented across components
// 4. Development-mode specific logging is properly segregated
// 5. All paths implement proper validation checks