import React from "react";
import { Button } from "@/components/ui/button";

import useOrderStore from "@/zusstand/order/order_zustand";
import { useOrderCreationStore } from "./logic";
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
  const { tableNumber, getOrderSummary } = useOrderStore();

  const orderSummary = getOrderSummary();
  const isDisabled = isLoading || !tableNumber || orderSummary.totalItems === 0;

  const handleCreateOrder = () => {
    createOrder(bowlChili, bowlNoChili, table_token);
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
