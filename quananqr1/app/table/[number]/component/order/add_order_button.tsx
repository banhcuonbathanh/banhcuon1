"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
import { useOrderCreationStore } from "./logic";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";

interface OrderCreationComponentProps {
  bowlChili: number;
  bowlNoChili: number;
  table_token: string;
}



const OrderCreationComponent: React.FC<OrderCreationComponentProps> = ({
  bowlChili,
  bowlNoChili,
  table_token
}) => {
  const { isLoading, createOrder } = useOrderCreationStore();
  const { tableNumber, getOrderSummary, clearOrder } = useOrderStore();
  const { http } = useApiStore();
  const { guest, user, isGuest, openLoginDialog, userId, isLogin } =
    useAuthStore();
  const {
    connect,
    disconnect,
    isConnected,
    sendMessage,
    wsToken,
    fetchWsToken
  } = useWebSocketStore();

  // const MAX_RETRY_ATTEMPTS = 3;

  const orderSummary = getOrderSummary();
  const isDisabled = isLoading || !tableNumber || orderSummary.totalItems === 0;

  const getEmailIdentifier = () => {
    if (isGuest && guest) {
      return guest.name;
    }
    return user?.email;
  };

  useEffect(() => {
    if (isLogin && userId) {
      initializeWebSocket();
    }
  }, [isLogin, userId, user?.email, guest?.name]);

  const initializeWebSocket = async () => {
    console.log(
      "quananqr1/app/table/[number]/component/order/add_order_button.tsx  12121212 !initializeWebSocket"
    );
    const emailIdentifier = getEmailIdentifier();

    if (!isLogin || !userId || !emailIdentifier) {
      console.log(
        "quananqr1/app/table/[number]/component/order/add_order_button.tsx emailIdentifier",
        emailIdentifier
      );
      return;
    }

    try {
      const wstoken1 = await fetchWsToken({
        userId: Number(userId),
        email: emailIdentifier,
        role: isGuest ? "Guest" : "User"
      });
      console.log(
        "quananqr1/app/table/[number]/component/order/add_order_button.tsx wstoken1",
        wstoken1
      );
      // Verify token after fetching
      if (!wstoken1) {
        throw new Error("Failed to obtain WebSocket token");
      }

      // Step 2: Establish WebSocket Connection

      await connect({
        userId: userId.toString(),
        isGuest,
        userToken: wstoken1.token, // Now TypeScript knows this is non-null
        tableToken: table_token,
        role: isGuest ? "Guest" : "User"
      });

      // setConnectionAttempts(0);
    } catch (error) {
      console.error("[OrderCreation] Connection error:", error);

 
    }
  };

  const handleCreateOrder = async () => {
    if (!isLogin) {
      console.log("[OrderCreation] User not logged in, showing login dialog");
      openLoginDialog();
      return;
    }

    if (!isConnected) {
      console.log(
        "quananqr1/app/table/[number]/component/order/add_order_button.tsx 12121212 !isConnected"
      );
      await initializeWebSocket();
      if (!isConnected) {
        console.log(
          "[OrderCreation] Failed to establish connection, aborting order creation"
        );
        return;
      }
    }

    if (orderSummary.totalItems === 0) {
      console.log("[OrderCreation] No items in order, aborting");
      return;
    }

    console.log("[OrderCreation] Creating order with summary:", orderSummary);
    createOrder({
      bowlChili,
      bowlNoChili,
      Table_token: table_token,
      http,
      auth: { guest, user, isGuest },
      orderStore: { tableNumber, getOrderSummary, clearOrder },
      websocket: { disconnect, isConnected, sendMessage },
      openLoginDialog
    });
  };

  const getButtonText = () => {
    if (!isLogin) {
      return "Login to Order";
    }

    if (orderSummary.totalItems === 0) {
      return "Add Items to Order";
    }
  };

  const isButtonDisabled = () => {
    if (!isLogin) return false;
    if (orderSummary.totalItems === 0) return true;

    return isLoading;

    // (connectionStatus === "error" &&
    //   connectionAttempts >= MAX_RETRY_ATTEMPTS) ||
    // isDisabled
  };

  useEffect(() => {
    return () => {
      console.log(
        "[OrderCreation] Component unmounting, cleaning up connection"
      );
      disconnect();
    };
  }, []);

  return (
    <div className="mt-4">
      <Button
        className="w-full"
        onClick={handleCreateOrder}
        disabled={isButtonDisabled()}
      >
        {getButtonText()}
      </Button>

    </div>
  );
};

export default OrderCreationComponent;
