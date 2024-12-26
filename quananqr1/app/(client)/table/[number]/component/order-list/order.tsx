"use client";

import React, { memo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import type { Order } from "@/schemaValidations/interface/type_order";
import { logWithLevel } from "@/lib/logger/log";

const OrderItem: React.FC<{ order: Order }> = memo(({ order }) => {
  const logPath =
    "quananqr1/app/(client)/table/[number]/component/order-list/order.tsx";

  // Basic item count calculation
  const itemCount = order.dish_items.length + order.set_items.length;

  // Logging
  logWithLevel(
    { orderId: order.id, tableNumber: order.table_number },
    logPath,
    "info",
    1
  );

  console.log("Order data:", {
    id: order.id,
    dishes: order.dish_items,
    sets: order.set_items,
    total: order.total_price
  });

  return (
    <Card className="mb-4">
      <CardContent className="p-4">
        <div className="grid grid-cols-2 gap-4">
          {/* Left Column - Basic Info */}
          <div>
            <h3 className="font-semibold">Order #{order.id}</h3>
            <p>Table: {order.table_number}</p>
            <p>Status: {order.status}</p>
            <p>Total: ${order.total_price.toFixed(2)}</p>
          </div>

          {/* Right Column - Additional Info */}
          <div className="text-right">
            <p>{new Date(order.created_at).toLocaleDateString()}</p>
            <p>Items: {itemCount}</p>
          </div>

          {/* Dishes Summary - Initially Collapsed */}
          <div className="col-span-2 mt-4">
            <details>
              <summary className="cursor-pointer">Show Order Details</summary>

              {/* Individual Dishes */}
              {order.dish_items.length > 0 && (
                <div className="mt-2">
                  <h4 className="font-medium">Individual Dishes:</h4>
                  {order.dish_items.map((dish, index) => (
                    <div
                      key={index}
                      className="flex justify-between text-sm py-1"
                    >
                      <span>
                        {dish.name} × {dish.quantity}
                      </span>
                      <span>${(dish.price * dish.quantity).toFixed(2)}</span>
                    </div>
                  ))}
                </div>
              )}

              {/* Set Meals */}
              {order.set_items.length > 0 && (
                <div className="mt-2">
                  <h4 className="font-medium">Set Meals:</h4>
                  {order.set_items.map((set, index) => (
                    <div key={index} className="mb-2">
                      <div className="flex justify-between text-sm">
                        <span>
                          {set.name} × {set.quantity}
                        </span>
                        <span>${(set.price * set.quantity).toFixed(2)}</span>
                      </div>
                      <div className="pl-4 text-xs text-gray-600">
                        {set.dishes.map((dish, dishIndex) => (
                          <div key={dishIndex}>
                            - {dish.name} × {dish.quantity}
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </details>
          </div>
        </div>
      </CardContent>
    </Card>
  );
});

OrderItem.displayName = "OrderItem";

export default OrderItem;
