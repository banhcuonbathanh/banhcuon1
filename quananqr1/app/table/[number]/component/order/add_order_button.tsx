"use client";

import React, { useEffect } from "react";
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
  const { guest, user, isGuest, openLoginDialog } = useAuthStore();
  const { connect, disconnect, isConnected, sendMessage } = useWebSocketStore();

  const orderSummary = getOrderSummary();
  const isDisabled = isLoading || !tableNumber || orderSummary.totalItems === 0;

  const handleCreateOrder = () => {
    createOrder({
      bowlChili,
      bowlNoChili,
      Table_token: table_token,
      http,
      auth: { guest, user, isGuest },
      orderStore: { tableNumber, getOrderSummary, clearOrder },
      websocket: { connect, disconnect, isConnected, sendMessage },
      openLoginDialog
    });
  };

  return (
    <div className="mt-4">
      <Button
        className="w-full"
        onClick={handleCreateOrder}
        disabled={isDisabled}
      >
        {isLoading ? "Creating Order..." : "Place Order"}
      </Button>
    </div>
  );
};

export default OrderCreationComponent;
