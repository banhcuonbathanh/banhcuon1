import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Package } from "lucide-react";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import { OrderSummary } from "./ordersummary";
import OrderSummaryPage from "./ordersummarypage";

const OrderListPage = () => {
  const { new_order, isLoading, tableToken, tableNumber } = useCartStore();
  const [expandedOrders, setExpandedOrders] = useState<Set<number>>(new Set());

  const toggleOrderDetails = (orderId: number) => {
    const newExpanded = new Set(expandedOrders);
    if (newExpanded.has(orderId)) {
      newExpanded.delete(orderId);
    } else {
      newExpanded.add(orderId);
    }
    setExpandedOrders(newExpanded);
  };

  const toggleAllOrders = () => {
    if (expandedOrders.size === new_order.length) {
      setExpandedOrders(new Set());
    } else {
      setExpandedOrders(new Set(new_order.map((order) => order.id)));
    }
  };

  if (isLoading) {
    return (
      <div className="p-4">
        <div className="text-center">Loading orders...</div>
      </div>
    );
  }

  return (
    <div className="p-4">
      {tableToken && tableNumber && (
        <div className="mb-4">
          <h2 className="text-2xl font-bold">Table #{tableNumber}</h2>
          <p className="text-sm text-gray-500">Token: {tableToken}</p>
        </div>
      )}

      <div className="flex gap-4 mb-6">
        <Button
          onClick={toggleAllOrders}
          className="flex items-center gap-2"
          variant="outline"
        >
          <Package className="h-4 w-4" />
          {expandedOrders.size === new_order.length
            ? "Hide All Details"
            : "Show All Details"}
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <OrderSummaryPage />

        {/* {new_order.map((order) => (
          <OrderSummary
            key={order.id}
            order={order}
            showDetails={expandedOrders.has(order.id)}
            onToggleDetails={() => toggleOrderDetails(order.id)}
          />
        ))} */}
      </div>
    </div>
  );
};

export default OrderListPage;
