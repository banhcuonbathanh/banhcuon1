"use client";

import React from "react";
import TitleButton from "./title-button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import useCartStore from "@/zusstand/new-order/new-order-zustand";

interface DishTotalItem {
  key: string;
  quantity: number;
  totalPrice: number;
  dishId: number;
}

const ItemsBreakdown = () => {
  const { dishTotal, deliveryData, remainingData } = useCartStore();

  // Create an array for dish totals
  const dishTotalsArray: DishTotalItem[] = dishTotal.map((dish) => ({
    key: `${dish.name}-${dish.dish_id}`,
    quantity: dish.quantity,
    totalPrice: dish.price * dish.quantity,
    dishId: dish.dish_id
  }));

  // Convert delivery data to required format
  // Now we need to extract from the map using dish IDs
  const deliveryMap: Record<string, number> = {};
  dishTotalsArray.forEach((dish) => {
    const dishDelivery = deliveryData[dish.dishId];
    if (dishDelivery) {
      // Use dish name as key for backwards compatibility with TitleButton
      const dishName = dish.key.split("-")[0];
      deliveryMap[dishName] = dishDelivery.quantity;
    }
  });

  // Convert remaining data to required format
  const remainingMap: Record<string, number> = {};
  if (remainingData && remainingData.name) {
    remainingMap[remainingData.name] = remainingData.quantity;
  }

  return (
    <div className="container mx-auto px-4 py-5 space-y-5">
      <Card>
        <CardHeader>
          <CardTitle className="flex justify-between items-center">
            total item break down
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid grid-cols-4 gap-4">
            {/* Header */}
            <div className="p-3 text-left font-medium text-gray-600">Title</div>
            <div className="p-3 text-right font-medium text-gray-600">
              Total Qty
            </div>
            <div className="p-3 text-right font-medium text-gray-600">
              Delivered
            </div>
            <div className="p-3 text-right font-medium text-gray-600">
              Remaining
            </div>

            {/* Grid Items */}
            {dishTotalsArray.map((details) => (
              <TitleButton
                key={`grid-${details.dishId}-${details.key}`}
                dishKey={details.key}
                details={details}
                deliveryData={deliveryMap}
                remainingData={remainingMap}
              />
            ))}

            {/* Footer Totals */}
            <div className="p-3 font-medium text-gray-500 border-t">Total</div>
            <div className="p-3 text-right font-medium text-gray-500 border-t">
              {dishTotalsArray.reduce((sum, item) => sum + item.quantity, 0)}
              <br />
              <span className="text-primary">
                {dishTotalsArray.reduce(
                  (sum, item) => sum + item.totalPrice,
                  0
                )}{" "}
              </span>
            </div>
            <div className="p-3 text-right font-medium text-gray-500 border-t">
              {Object.values(deliveryMap).reduce((sum, val) => sum + val, 0)}
              <br />
              <span className="text-primary">
                {dishTotalsArray.reduce((sum, details) => {
                  const title = details.key.split("-")[0];
                  const delivered = deliveryMap[title] || 0;
                  return (
                    sum + (delivered * details.totalPrice) / details.quantity
                  );
                }, 0)}
              </span>
            </div>
            <div className="p-3 text-right font-medium text-gray-500 border-t">
              {Object.values(remainingMap).reduce((sum, val) => sum + val, 0)}
              <br />
              <span className="text-gray-500">
                {dishTotalsArray.reduce((sum, details) => {
                  const title = details.key.split("-")[0];
                  const remaining =
                    remainingMap[title] ||
                    details.quantity - (deliveryMap[title] || 0);
                  return (
                    sum + (remaining * details.totalPrice) / details.quantity
                  );
                }, 0)}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default ItemsBreakdown;
