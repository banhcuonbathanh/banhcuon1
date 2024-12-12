// config/loggerConfig.ts
export interface LogPath {
    path: string;
    enabled: boolean;
    description?: string;
  }
  
  export const loggerPaths: LogPath[] = [
    {
      path: '/manage/admin/table',
      enabled: true,
      description: 'Admin table management page logs'
    },
    {
      path: '/manage/admin/set',
      enabled: false,
      description: 'Admin settings page logs'
    },
    {
      path: '/(client)/table',
      enabled: true,
      description: 'Client table viewing page logs'
    }
  ];
  
 
  
  const isDevelopment = process.env.NODE_ENV !== 'production';
  
  export function log(message: any, path: string) {
    // Early return if in production
    if (!isDevelopment) return;
  
    const pathConfig = loggerPaths.find(p => path.includes(p.path));
    
    if (!pathConfig?.enabled) return;
    
    const description = pathConfig?.description || 'No description';
    const isServer = typeof window === 'undefined';
    const context = isServer ? '[Server]' : '[Client]';
    
    console.log(`${context}[${path}] (${description}):`, message);
  }
  
  // Alternative version with selective production logging
  export function logWithLevel(message: any, path: string, level: 'debug' | 'info' | 'warn' | 'error' = 'debug') {
    const pathConfig = loggerPaths.find(p => path.includes(p.path));
    
    // Only continue if path is enabled
    if (!pathConfig?.enabled) return;
    
    // In production, only show warnings and errors
    if (!isDevelopment && level !== 'warn' && level !== 'error') return;
    
    const description = pathConfig?.description || 'No description';
    const isServer = typeof window === 'undefined';
    const context = isServer ? '[Server]' : '[Client]';
    
    switch (level) {
      case 'debug':
        isDevelopment && console.log(`${context}[${path}] (${description}):`, message);
        break;
      case 'info':
        isDevelopment && console.info(`${context}[${path}] (${description}):`, message);
        break;
      case 'warn':
        console.warn(`${context}[${path}] (${description}):`, message);
        break;
      case 'error':
        console.error(`${context}[${path}] (${description}):`, message);
        break;
    }
  }


//   'use client';

  
//   export default function AdminTablePage() {
//     const pathname = usePathname();



  // First, let's set up our logger paths
// config/loggerConfig.ts
// export interface LogPath {
//     path: string;
//     enabled: boolean;
//     description?: string;
//   }
  
//   export const loggerPaths: LogPath[] = [
//     {
//       path: '/manage/admin/table',
//       enabled: true,
//       description: 'Admin table management page logs'
//     },
//     {
//       path: '/manage/admin/set',
//       enabled: false,
//       description: 'Admin settings page logs'
//     },
//     {
//       path: '/(client)/table',
//       enabled: true,
//       description: 'Client table viewing page logs'
//     }
//   ];
  
//   // Example 1: Page with logging enabled
//   // app/manage/admin/table/page.tsx
//   'use client';
//   import { log } from '@/utils/logger';
//   import { usePathname } from 'next/navigation';
  
//   export default function AdminTablePage() {
//     const pathname = usePathname();
    
//     // This WILL show up because enabled: true for this path
//     log('Admin table page loaded', pathname);
    
//     return (
//       <div>
//         <button onClick={() => {
//           // This WILL also show up
//           log('Button clicked in admin table', pathname);
//         }}>
//           Click Me
//         </button>
//       </div>
//     );
//   }
  
//   // Example 2: Page with logging disabled
//   // app/manage/admin/set/page.tsx
//   'use client';
//   import { log } from '@/utils/logger';
//   import { usePathname } from 'next/navigation';
  
//   export default function AdminSettingsPage() {
//     const pathname = usePathname();
    
//     // This will NOT show up because enabled: false for this path
//     log('Admin settings page loaded', pathname);
    
//     return (
//       <div>
//         <button onClick={() => {
//           // This will NOT show up either
//           log('Button clicked in settings', pathname);
//         }}>
//           Click Me
//         </button>
//       </div>
//     );
//   }
  
//   // Example 3: Server Component with logging enabled
//   // app/(client)/table/[number]/page.tsx
//   import { log } from '@/utils/logger';
  
//   export default async function ClientTablePage() {
//     // This WILL show up because enabled: true for this path
//     log('Client table page loaded', '/(client)/table');
    
//     const data = await fetchData();
//     // This WILL also show up
//     log('Data fetched', '/(client)/table');
    
//     return <div>Table Content</div>;
//   }