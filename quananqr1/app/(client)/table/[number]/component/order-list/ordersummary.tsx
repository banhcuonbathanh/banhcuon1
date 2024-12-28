// components/OrderSummary.tsx
import React from 'react';
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, CircleDollarSign } from "lucide-react";
import { Order } from "@/schemaValidations/interface/type_order";
import { OrderSummaryHeader } from './ordersummaryheader';
import { OrderDetailsList } from './orderdetaillist';


interface OrderSummaryProps {
  order: Order;
  showDetails: boolean;
  onToggleDetails: () => void;
}

export const OrderSummary: React.FC<OrderSummaryProps> = ({
  order,
  showDetails,
  onToggleDetails
}) => {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD"
    }).format(price);
  };

  const totalIndividualDishes = order.dish_items.reduce(
    (sum, item) => sum + item.quantity,
    0
  );
  const totalSetDishes = order.set_items.reduce(
    (sum, item) => sum + item.quantity,
    0
  );

  return (
    <Card className="w-full">
      <OrderSummaryHeader 
        order={order}
        totalItems={totalIndividualDishes + totalSetDishes}
        formattedPrice={formatPrice(order.total_price)}
      />

      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="font-medium">Individual Dishes:</span>
            <span>{totalIndividualDishes} items</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="font-medium">Set Menus:</span>
            <span>{totalSetDishes} sets</span>
          </div>
        </div>

        <Button onClick={onToggleDetails} variant="outline" className="w-full">
          {showDetails ? (
            <div className="flex items-center">
              Hide Details <ChevronUp className="ml-2 h-4 w-4" />
            </div>
          ) : (
            <div className="flex items-center">
              Show Details <ChevronDown className="ml-2 h-4 w-4" />
            </div>
          )}
        </Button>

        {showDetails && (
          <>
            <OrderDetailsList order={order} formatPrice={formatPrice} />
            {order.topping && (
              <div className="mt-4 text-sm text-gray-600">
                <span className="font-medium">Special Instructions:</span>{" "}
                {order.topping}
              </div>
            )}
          </>
        )}

        <div className="flex items-center justify-between border-t pt-4">
          <div className="flex items-center gap-2">
            <CircleDollarSign className="h-5 w-5" />
            <span className="font-semibold">Total:</span>
          </div>
          <span className="text-lg font-bold">
            {formatPrice(order.total_price)}
          </span>
        </div>
      </CardContent>
    </Card>
  );
};
