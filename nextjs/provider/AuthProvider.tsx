"use client";

import { SessionProvider } from "next-auth/react";

export default function AuthProvider({
  children
}: {
  children: React.ReactNode;
}) {
  return <SessionProvider>{children}</SessionProvider>;
}

// import { NextComponentType, NextPageContext } from 'next'
// import { Session } from 'next-auth' // or wherever Session is imported from in your project

// interface IProps {
//   Component: NextComponentType<NextPageContext, any, {}>,
//   pageProps: { session: Session, [key: string]: any }
// }

// export default function AuthProvider({ Component, pageProps: { session, ...pageProps } }: IProps) {
//   return (
//     <SessionProvider session={session}>
//       <Component {...pageProps} />
//     </SessionProvider>
//   )
// }
