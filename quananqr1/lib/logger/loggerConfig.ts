export const loggerPaths: LogPath[] = [
  {
    path: "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
    enabled: false,
    description: "Restaurant Summary Component Logs",
    enabledLogIds: [1, 2, 3],
    disabledLogIds: [],
    logIds: [1, 2, 3],
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
    disabledLogIds: [],
    logIds: [1, 2, 3],
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
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8],
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
    path: "quananqr1/zusstand/order/order_zustand.ts",
    enabled: false,
    description: "Order and Delivery Store State Management Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
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
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
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
    enabledLogIds: [6],
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8],
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
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6],
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
    enabled: true,
    description: "Order Creation Component Logs",
    enabledLogIds: [3],
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8, 9],
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
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8],
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
  },

  {
    path: "quananqr1/zusstand/new_auth/new_auth_controller.ts",
    enabled: false,
    description: "Authentication Controller Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12],
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12],
    logDescriptions: {
      1: {
        description: "User registration attempt",
        location: "register function - initialization",
        status: "enabled"
      },
      2: {
        description: "User registration success",
        location: "register function - success path",
        status: "enabled"
      },
      3: {
        description: "User login attempt",
        location: "login function - initialization",
        status: "enabled"
      },
      4: {
        description: "User login success",
        location: "login function - success path",
        status: "enabled"
      },
      5: {
        description: "Guest login attempt",
        location: "guestLogin function - initialization",
        status: "enabled"
      },
      6: {
        description: "Guest login success",
        location: "guestLogin function - success path",
        status: "enabled"
      },
      7: {
        description: "Logout process",
        location: "logout function - initialization",
        status: "enabled"
      },
      8: {
        description: "Guest logout process",
        location: "guestLogout function - initialization",
        status: "enabled"
      },
      9: {
        description: "Token refresh attempt",
        location: "refreshAccessToken function - initialization",
        status: "enabled"
      },
      10: {
        description: "Auth state synchronization",
        location: "syncAuthState function",
        status: "enabled"
      },
      11: {
        description: "Cookie-based auth initialization",
        location: "initializeAuthFromCookies function",
        status: "enabled"
      },
      12: {
        description: "Error state management",
        location: "error handling across all functions",
        status: "enabled"
      }
    }
  },
  {
    path: "quananqr1/components/set/set_card.tsx",
    enabled: false,
    description: "Set Card Component Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6, 7, 8],
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6, 7, 8],
    logDescriptions: {
      1: {
        description: "Component initialization and props",
        location: "SetCard component - initialization",
        status: "enabled"
      },
      2: {
        description: "Current order state changes",
        location: "SetCard component - order state updates",
        status: "enabled"
      },
      3: {
        description: "Set quantity modifications",
        location: "handleIncrease/handleDecrease functions",
        status: "enabled"
      },
      4: {
        description: "Dish quantity modifications",
        location: "handleDishIncrease/handleDishDecrease functions",
        status: "enabled"
      },
      5: {
        description: "Price calculations",
        location: "totalPrice calculation",
        status: "enabled"
      },
      6: {
        description: "Set visibility toggle",
        location: "toggleList function",
        status: "enabled"
      },
      7: {
        description: "Error boundary triggers",
        location: "SetCard component - error handling",
        status: "enabled"
      },
      8: {
        description: "Order store interactions",
        location: "useOrderStore hook interactions",
        status: "enabled"
      }
    }
  },

  {
    path: "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
    enabled: false,
    description: "Restaurant Summary Component Logs",
    enabledLogIds: [1, 2, 3],
    disabledLogIds: [],
    logIds: [1, 2, 3],
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
    path: "quananqr1/app/(client)/table/[number]/component/total-dishes-detail.tsx",
    enabled: false,
    description: "Total Dishes Detail Component Logs",
    enabledLogIds: [1],
    disabledLogIds: [1, 2, 3, 4, 5],
    logIds: [1, 2, 3, 4, 5],
    logDescriptions: {
      1: {
        description: "Component initialization state",
        location: "OrderDetails component - initialization",
        status: "enabled"
      },
      2: {
        description: "Set expansion and price calculations",
        location: "toggleSetExpansion and calculateSetPrice functions",
        status: "enabled"
      },
      3: {
        description: "Dish totals calculation",
        location: "calculateDishTotals function",
        status: "enabled"
      },
      4: {
        description: "Price updates and state changes",
        location: "Price calculation and state updates",
        status: "enabled"
      },
      5: {
        description: "Error handling",
        location: "Error boundary and validation checks",
        status: "enabled"
      }
    }
  },

  {
    path: "quananqr1/app/(client)/table/[number]/component/order-list/list-order.tsx",
    enabled: false,
    description: "Orders List Component Logs",
    enabledLogIds: [1, 2, 3, 4, 5, 6],
    disabledLogIds: [],
    logIds: [1, 2, 3, 4, 5, 6],
    logDescriptions: {
      1: {
        description: "Order summary generation started",
        location: "getOrderSummaryForOrder function - initialization",
        status: "enabled"
      },
      2: {
        description: "Dish items transformation",
        location: "getOrderSummaryForOrder function - dish transformation",
        status: "enabled"
      },
      3: {
        description: "Set items transformation",
        location: "getOrderSummaryForOrder function - set transformation",
        status: "enabled"
      },
      4: {
        description: "Total calculations completed",
        location: "getOrderSummaryForOrder function - totals calculation",
        status: "enabled"
      },
      5: {
        description: "Order list rendering state",
        location: "OrdersList component - render phase",
        status: "enabled"
      },
      6: {
        description: "Empty orders list state",
        location: "OrdersList component - empty state handling",
        status: "enabled"
      }
    }
  }
];

export interface LogPath {
  path: string;
  enabled: boolean;
  description: string;
  enabledLogIds: number[];
  disabledLogIds: number[];
  logIds: number[];
  logDescriptions: {
    [key: number]: LogDescription;
  };
}

export interface LogDescription {
  description: string;
  location: string;
  status: "enabled" | "disabled";
}
