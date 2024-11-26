"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
import { useOrderCreationStore } from "./logic";
import { useApiStore } from "@/zusstand/api/api-controller";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { useWebSocketStore } from "@/zusstand/web-socket/websocketStore";

interface OrderCreationComponentProps {
  table_token: string;
}

const OrderCreationComponent: React.FC<OrderCreationComponentProps> = ({
  table_token
}) => {
  const { isLoading, createOrder } = useOrderCreationStore();
  const {
    tableNumber,
    getOrderSummary,
    clearOrder,
    canhKhongRau,
    canhCoRau,
    smallBowl,
    wantChili,
    selectedFilling
  } = useOrderStore();
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

  // State to track authentication check status
  const [authChecked, setAuthChecked] = useState(false);

  let topping = `canhKhongRau ${canhKhongRau} - canhCoRau ${canhCoRau} - bat be ${smallBowl} - ot tuoi ${wantChili} - nhan ${selectedFilling} -`;
  const orderSummary = getOrderSummary();

  // Initialize auth state when component mounts
  useEffect(() => {
    const initializeAuth = async () => {
      // Sync auth state from cookies
      useAuthStore.getState().syncAuthState();
      setAuthChecked(true);
    };

    initializeAuth();
  }, []);

  // Effect for WebSocket initialization after auth check
  useEffect(() => {
    if (authChecked && isLogin && userId) {
      console.log("Initializing WebSocket connection for user:", userId);
      initializeWebSocket();
    }
  }, [authChecked, isLogin, userId, user?.email, guest?.name]);

  const getEmailIdentifier = () => {
    if (isGuest && guest) {
      return guest.name;
    }
    return user?.email;
  };

  const initializeWebSocket = async () => {
    const emailIdentifier = getEmailIdentifier();

    if (!isLogin || !userId || !emailIdentifier) {
      console.log(
        "WebSocket initialization failed: Missing required credentials"
      );
      return;
    }

    try {
      const wstoken1 = await fetchWsToken({
        userId: Number(userId),
        email: emailIdentifier,
        role: isGuest ? "Guest" : "User"
      });

      if (!wstoken1) {
        throw new Error("Failed to obtain WebSocket token");
      }

      await connect({
        userId: userId.toString(),
        isGuest,
        userToken: wstoken1.token,
        tableToken: table_token,
        role: isGuest ? "Guest" : "User"
      });
    } catch (error) {
      console.error("[OrderCreation] Connection error:", error);
    }
  };

  const handleCreateOrder = async () => {
    // Ensure auth state is synced before proceeding
    useAuthStore.getState().syncAuthState();
    const currentAuthState = useAuthStore.getState();

    if (!currentAuthState.isLogin) {
      console.log("[OrderCreation] User not logged in, showing login dialog");
      openLoginDialog();
      return;
    }

    if (!isConnected) {
      console.log("Attempting to establish WebSocket connection");
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
    if (tableNumber === null) {
      console.log("[OrderCreation] No items in order, aborting");
      return;
    }

    console.log("[OrderCreation] Creating order with summary:", orderSummary);
    createOrder({
      topping,
      Table_token: table_token,
      http,
      auth: { guest, user, isGuest },
      orderStore: { tableNumber, getOrderSummary, clearOrder },
      websocket: { disconnect, isConnected, sendMessage },
      openLoginDialog
    });
  };

  const getButtonText = () => {
    if (!authChecked) {
      return "Loading...";
    }
    if (!isLogin) {
      return "Login to Order";
    }
    if (orderSummary.totalItems === 0) {
      return "Add Items to Order";
    }
    return "Create Order";
  };

  const isButtonDisabled = () => {
    if (!authChecked) return true;
    if (!isLogin) return false;
    if (orderSummary.totalItems === 0) return true;
    return isLoading;
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
