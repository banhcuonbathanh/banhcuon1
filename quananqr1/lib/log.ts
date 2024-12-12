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
    enabled: true,
    description: "Restaurant Summary Component Logs",
    enabledLogIds: [], // Only logs #1 and #2 are enabled
    logDescriptions: {
      1: {
        description: "Initial dishMap state in aggregateDishes function",
        location: "aggregateDishes function - before processing orders",
        status: "enabled"
      },
      2: {
        description: "Aggregated dishes result for each group",
        location: "RestaurantSummary component - inside groupedOrders.map",
        status: "enabled"
      }
    }
  },
  {
    path: "/manage/admin/table",
    enabled: true,
    description: "Admin Table Management Logs",
    enabledLogIds: [1, 2, 3],
    logDescriptions: {
      1: {
        description: "Table initialization",
        location: "Table component initialization",
        status: "enabled"
      },
      2: {
        description: "Table data updates",
        location: "Table data manipulation functions",
        status: "enabled"
      },
      3: {
        description: "Table error states",
        location: "Error handling in table components",
        status: "enabled"
      }
    }
  },

  {
    path: "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx",
    enabled: true,
    description: "Restaurant Summary Component Logs",
    enabledLogIds: [1, 2, 3], // Only logs #1 and #2 are enabled
    logDescriptions: {
      1: {
        description: "Initial dishMap state in aggregateDishes function",
        location: "aggregateDishes function - before processing orders",
        status: "enabled"
      },
      2: {
        description: "Initial dishMap state in aggregateDishes function",
        location: "aggregateDishes function - before processing orders",
        status: "enabled"
      },

      3: {
        description: "handleDeliverySubmitn createDelivery",
        location: "aggregateDishes function - before processing orders",
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

// Example usage in restaurant-summary.tsx:
/*
// Get log status for a specific file
getLogStatus("quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx");

// Log with ID 1
logWithLevel(
  { dishMap },
  "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
  "info",
  1
);

// Log with ID 2
logWithLevel(
  { aggregatedDishes },
  "quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary.tsx",
  "info",
  2
);
*/
