// components/OrderSummaryHeader.tsx
import React from 'react';
import { Badge } from "@/components/ui/badge";
import { Flame } from "lucide-react";
import { CardTitle, CardDescription, CardHeader } from "@/components/ui/card";
import { Order } from "@/schemaValidations/interface/type_order";

interface OrderSummaryHeaderProps {
  order: Order;
  totalItems: number;
  formattedPrice: string;
}

export const OrderSummaryHeader: React.FC<OrderSummaryHeaderProps> = ({
  order,
  totalItems,
  formattedPrice
}) => (
  <CardHeader>
    <div className="flex items-center justify-between">
      <div>
        <CardTitle className="text-xl">Order #{order.id}</CardTitle>
        <CardDescription>
          {totalItems} items - {formattedPrice}
        </CardDescription>
      </div>
      {order.takeAway && (
        <Badge className="bg-purple-500 text-white">Takeaway</Badge>
      )}
    </div>
    {order.chiliNumber > 0 && (
      <div className="flex items-center gap-2 mt-2">
        <Flame className="h-4 w-4 text-red-500" />
        <span>Spice Level: {order.chiliNumber}</span>
      </div>
    )}
  </CardHeader>
);
