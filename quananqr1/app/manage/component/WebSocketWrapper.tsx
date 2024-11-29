"use client";

import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";
import { useEffect } from "react";

interface WebSocketWrapperProps {
  userId: string;
  role: string;
}

export default function WebSocketWrapper({
  userId,
  role
}: WebSocketWrapperProps) {
  console.log("quananqr1/app/manage/component/WebSocketWrapper.tsx");
  const { connect } = useWebSocketStore();

  useEffect(() => {
    // You'll need to get these tokens from your authentication context or state
    const userToken = "...";
    const tableToken = "...";

    connect({
      userId,
      isGuest: false,
      userToken,
      tableToken,
      role
    });

    // Optional: Disconnect on unmount
    return () => {
      useWebSocketStore.getState().disconnect();
    };
  }, [userId, role]);

  return null; // This component doesn't render anything
}
