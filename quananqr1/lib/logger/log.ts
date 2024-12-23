// config/loggerConfig.ts

import { loggerPaths, LogPath } from "./loggerConfig";

// Types


type LogLevel = "debug" | "info" | "warn" | "error";

// Validation Functions
function validateLogConfiguration(logPath: LogPath): {
  isValid: boolean;
  errors: string[];
  warnings: string[];
} {
  const result = {
    isValid: true,
    errors: [] as string[],
    warnings: [] as string[]
  };

  const definedLogIds = new Set(
    Object.keys(logPath.logDescriptions).map(Number)
  );
  const enabledLogIds = new Set<number>(logPath.enabledLogIds);

  // Check if all enabled IDs exist in logDescriptions
  logPath.enabledLogIds.forEach((enabledId) => {
    if (!definedLogIds.has(enabledId)) {
      result.errors.push(
        `Enabled log ID ${enabledId} does not exist in logDescriptions`
      );
      result.isValid = false;
    }
  });

  // Check status consistency
  for (const [idStr, description] of Object.entries(logPath.logDescriptions)) {
    const id = Number(idStr);
    const isEnabled = enabledLogIds.has(id);

    if (description.status === "enabled" && !isEnabled) {
      result.errors.push(
        `Log ID ${id} has status "enabled" but is not in enabledLogIds`
      );
      result.isValid = false;
    }

    if (description.status === "disabled" && isEnabled) {
      result.errors.push(
        `Log ID ${id} has status "disabled" but is in enabledLogIds`
      );
      result.isValid = false;
    }
  }

  // Verify sequential IDs
  const sortedIds = Array.from(definedLogIds).sort((a, b) => a - b);
  if (sortedIds[0] !== 1) {
    result.errors.push("Log IDs should start from 1");
    result.isValid = false;
  }

  for (let i = 1; i < sortedIds.length; i++) {
    if (sortedIds[i] !== sortedIds[i - 1] + 1) {
      result.errors.push(
        `Missing sequential log ID: expected ${sortedIds[i - 1] + 1}, found ${
          sortedIds[i]
        }`
      );
      result.isValid = false;
    }
  }

  return result;
}

function updateLogConfiguration(logPath: LogPath, enable: boolean): LogPath {
  const updatedConfig = { ...logPath };

  // Update main enabled flag
  updatedConfig.enabled = enable;

  // Update all log descriptions
  for (const [id, desc] of Object.entries(updatedConfig.logDescriptions)) {
    updatedConfig.logDescriptions[Number(id)] = {
      ...desc,
      status: enable ? "enabled" : "disabled"
    };
  }

  // Update enabledLogIds
  updatedConfig.enabledLogIds = enable
    ? Object.keys(updatedConfig.logDescriptions).map(Number)
    : [];

  return updatedConfig;
}

// Environment configuration
const isDevelopment = process.env.NODE_ENV !== "production";

// Utility Functions
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

// Status Display Function
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

// Main Logging Function
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

  // Check if logging is enabled
  if (!pathConfig.enabled || !pathConfig.enabledLogIds.includes(logId)) {
    return;
  }

  // Get log description
  const logInfo = pathConfig.logDescriptions[logId];

  // Format log prefix
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

// Validation example usage
export function validateAllConfigurations(): void {
  for (const config of loggerPaths) {
    const result = validateLogConfiguration(config);
    if (!result.isValid) {
      console.error(`Validation failed for path: ${config.path}`);
      console.error("Errors:", result.errors);
      console.error("Warnings:", result.warnings);
    }
  }
}
