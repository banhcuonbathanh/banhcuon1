// order.tsx
import React, { memo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import type { Order } from "@/schemaValidations/interface/type_order";

const OrderItem: React.FC<{ order: Order }> = memo(({ order }) => {
  // Memoize computed values
  const formattedDate = React.useMemo(
    () => new Date(order.created_at).toLocaleDateString(),
    [order.created_at]
  );

  const itemCount = React.useMemo(
    () => order.dish_items.length + order.set_items.length,
    [order.dish_items.length, order.set_items.length]
  );

  return (
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
            <p>{order.takeAway ? "Takeaway" : "Dine-in"}</p>
            <p>Items: {itemCount}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
});

OrderItem.displayName = "OrderItem";

export default OrderItem;
