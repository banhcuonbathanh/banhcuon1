// log_zustand.ts
import { create } from "zustand";
import { useCallback } from "react";

// Types
export type LogLevel = "debug" | "info" | "warn" | "error";

export interface LoggerConfig {
  enabled: boolean;
  level: LogLevel;
  isDevelopment: boolean;
}
interface LogEntry {
  level: LogLevel;
  message: string;
  module: string;
  timestamp: string;
  data?: any;
}

interface LoggerStore {
  config: LoggerConfig;
  logs: LogEntry[];
  log: (level: LogLevel, module: string, message: string, data?: any) => void;
  clearLogs: () => void;
  updateConfig: (config: Partial<LoggerConfig>) => void;
}

// Config
const createLoggerConfig = (): LoggerConfig => {
  const isDevelopment = process.env.NODE_ENV === "development";

  return {
    enabled: isDevelopment,
    level: isDevelopment ? "debug" : ("error" as LogLevel), // Explicit type assertion
    isDevelopment
  };
};

// No-op logger for production
const createNoOpLogger = () => ({
  log: () => {},
  debug: () => {},
  info: () => {},
  warn: () => {},
  error: () => {},
  clearLogs: () => {},
  updateConfig: () => {}
});

// Store
export const useLoggerStore = create<LoggerStore>((set, get) => ({
  config: createLoggerConfig(),
  logs: [],

  log: (level, module, message, data) => {
    const { config } = get();

    // Only log if enabled and in development
    if (!config.enabled || !config.isDevelopment) {
      return;
    }

    set((state) => ({
      logs: [
        ...state.logs,
        {
          level,
          module,
          message,
          data,
          timestamp: new Date().toISOString()
        }
      ]
    }));

    // Console logging only in development
    if (config.isDevelopment) {
      const consoleMethod = console[level] || console.log;
      consoleMethod(`[${module}] ${message}`, data);
    }
  },

  clearLogs: () => set({ logs: [] }),

  updateConfig: (newConfig) =>
    set((state) => ({
      config: { ...state.config, ...newConfig }
    }))
}));

// Hook
export const useLogger = (module: string) => {
  const { config, log } = useLoggerStore();

  // If not in development, return no-op logger
  if (!config.isDevelopment) {
    return createNoOpLogger();
  }

  return {
    debug: useCallback(
      (message: string, data?: any) => {
        log("debug", module, message, data);
      },
      [module]
    ),

    info: useCallback(
      (message: string, data?: any) => {
        log("info", module, message, data);
      },
      [module]
    ),

    warn: useCallback(
      (message: string, data?: any) => {
        log("warn", module, message, data);
      },
      [module]
    ),

    error: useCallback(
      (message: string, data?: any) => {
        log("error", module, message, data);
      },
      [module]
    )
  };
};
