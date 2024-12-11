'use client';

import { useLoggerStore } from '@/zusstand/log/log_zustand';
import { useEffect } from 'react';


export function LoggerProvider({ children }: { children: React.ReactNode }) {
  const { updateConfig } = useLoggerStore();

  useEffect(() => {
    const isDevelopment = process.env.NODE_ENV === 'development';
    
    updateConfig({
      enabled: isDevelopment,
      level: isDevelopment ? 'debug' : 'error',
      isDevelopment
    });

    if (isDevelopment) {
      console.log('[Logger] Initialized in development mode');
    }
  }, [updateConfig]);

  return <>{children}</>;
}