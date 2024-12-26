"use client";

import React, { useState, useEffect } from "react";
import useOrderStore from "@/zusstand/order/order_zustand";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ChevronDown, ChevronUp } from "lucide-react";
import { Order } from "@/schemaValidations/interface/type_order";

// Stateless OrderItem component remains the same
const OrderItem = React.memo(
  ({
    order,
    formattedDate,
    itemCount
  }: {
    order: Order;
    formattedDate: string;
    itemCount: number;
  }) => (
    <Card className="mb-4">
      <CardContent className="p-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <h3 className="font-semibold">Order #{order.id}</h3>
            <p>Table: {order.table_number}</p>
            <p>Status: {order.status}</p>
            <p>Total: ${order.total_price.toFixed(2)}</p>
          </div>
          <div className="text-right">
            <p>{formattedDate}</p>

            <p>Items: {itemCount}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
);

OrderItem.displayName = "OrderItem";

// Loading skeleton component remains the same
const OrderItemSkeleton = () => (
  <Card className="mb-4">
    <CardContent className="p-4">
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <div className="h-6 bg-gray-200 rounded w-1/2 animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-1/3 animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-2/3 animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-1/4 animate-pulse"></div>
        </div>
        <div className="text-right space-y-2">
          <div className="h-4 bg-gray-200 rounded w-1/2 ml-auto animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-1/3 ml-auto animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-1/4 ml-auto animate-pulse"></div>
        </div>
      </div>
    </CardContent>
  </Card>
);

const OrdersList = () => {
  const zustandOrders = useOrderStore((state) => state?.listOfOrders) || [];
  const [processedOrders, setProcessedOrders] = useState<
    Array<{
      order: Order;
      formattedDate: string;
      itemCount: number;
    }>
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isOrderListVisible, setIsOrderListVisible] = useState(false);

  useEffect(() => {
    try {
      if (Array.isArray(zustandOrders) && zustandOrders.length > 0) {
        // Modified filtering logic to handle null dish_items or set_items
        const processed = zustandOrders
          .filter((order) => order && order.created_at)
          .map((order) => ({
            order,
            formattedDate: new Date(order.created_at).toLocaleDateString(),
            itemCount:
              (order.dish_items?.length || 0) +
                (order.set_items?.length || 0) || 1 // Default to 1 if both are null/empty
          }));

        setProcessedOrders(processed);
      }
    } catch (error) {
      console.error("Error processing orders:", error);
    } finally {
      setIsLoading(false);
    }
  }, [zustandOrders]);

  const toggleOrderListVisibility = () => {
    setIsOrderListVisible((prev) => !prev);
  };

  return (
    <div className="w-full max-w-2xl mx-auto p-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Orders List</CardTitle>
          <Button
            variant="ghost"
            size="sm"
            onClick={toggleOrderListVisibility}
            className="ml-2"
          >
            {isOrderListVisible ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
          </Button>
        </CardHeader>

        {isOrderListVisible && (
          <CardContent>
            {isLoading ? (
              <div className="space-y-4">
                {[...Array(3)].map((_, index) => (
                  <OrderItemSkeleton key={index} />
                ))}
              </div>
            ) : processedOrders.length === 0 ? (
              <p className="text-center text-gray-500 py-4">
                No orders available
              </p>
            ) : (
              <div>
                <p className="mb-4">Total Orders: {processedOrders.length}</p>
                {processedOrders.map(({ order, formattedDate, itemCount }) => (
                  <OrderItem
                    key={`order-${order.id}`}
                    order={order}
                    formattedDate={formattedDate}
                    itemCount={itemCount}
                  />
                ))}
              </div>
            )}
          </CardContent>
        )}
      </Card>
    </div>
  );
};

export default OrdersList;
