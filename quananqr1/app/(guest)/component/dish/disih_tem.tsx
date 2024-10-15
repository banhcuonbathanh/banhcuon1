"use client";

import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Minus } from "lucide-react";
import { DishInterface } from "@/schemaValidations/interface/type_dish";
import useOrderStore from "@/zusstand/order/order_zustand";

interface DishCardProps {
  dish: DishInterface;
}

export const DishCard: React.FC<DishCardProps> = ({ dish }) => {
  const { addItem, removeItem, findOrderItem } = useOrderStore();

  // Use the useOrderStore hook to get the current quantity
  const orderItem = useOrderStore((state) =>
    state.findOrderItem("dish", dish.id)
  );
  const quantity = orderItem ? orderItem.quantity : 0;

  const handleIncrease = () => {
    addItem(dish, 1);
  };

  const handleDecrease = () => {
    if (quantity > 0) {
      if (quantity === 1) {
        removeItem("dish", dish.id);
      } else {
        addItem(dish, -1);
      }
    }
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="text-lg">{dish.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <img
          src={dish.image || "/api/placeholder/150/100"}
          alt={dish.name}
          className="w-full h-24 object-cover rounded-md mb-2"
        />
        <p className="text-sm">{dish.description}</p>
        <p className="font-semibold mt-2">Price: ${dish.price.toFixed(2)}</p>
        <div className="flex items-center justify-between mt-2">
          <span>Quantity:</span>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={handleDecrease}>
              <Minus className="h-3 w-3" />
            </Button>
            <span className="w-8 text-center">{quantity}</span>
            <Button variant="outline" size="sm" onClick={handleIncrease}>
              <Plus className="h-3 w-3" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
