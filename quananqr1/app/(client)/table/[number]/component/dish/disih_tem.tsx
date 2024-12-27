"use client";

import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Minus } from "lucide-react";
import { DishInterface } from "@/schemaValidations/interface/type_dish";

import { DishOrderItem } from "@/schemaValidations/interface/type_order";
import useCartStore from "@/zusstand/new-order/new-order-zustand";

interface DishCardProps {
  dish: DishInterface;
}

export const DishCard: React.FC<DishCardProps> = ({ dish }) => {
  const { current_order, addDishToCart, updateDishQuantity, removeDishFromCart } = useCartStore();

  // Find the dish in the current order
  const currentDish = current_order?.dish_items.find(
    (item) => item.dish_id === dish.id
  );
  const quantity = currentDish ? currentDish.quantity : 0;

  const handleIncrease = () => {
    if (currentDish) {
      updateDishQuantity("increment", dish.id);
    } else {
      const newDishItem: DishOrderItem = {
        dish_id: dish.id,
        quantity: 1,
        price: dish.price,
        name: dish.name,
        description: dish.description,
        image: dish.image,
        status: dish.status,
        created_at: "",
        updated_at: "",
        is_favourite: false,
        like_by: []
      };
      addDishToCart(newDishItem);
    }
  };

  const handleDecrease = () => {
    if (currentDish) {
      if (currentDish.quantity > 1) {
        updateDishQuantity("decrement", dish.id);
      } else {
        removeDishFromCart(dish.id);
      }
    }
  };

  return (
    <Card className="w-full">
      <CardContent className="p-4 flex">
        <div className="w-1/3 pr-4">
          <img
            src={dish.image || "/api/placeholder/150/150"}
            alt={dish.name}
            className="w-full h-full object-cover rounded-md"
          />
        </div>
        <div className="w-2/3 flex flex-col justify-between">
          <div>
            <h3 className="text-lg font-semibold mb-2">{dish.name}</h3>
            <p className="text-sm mb-2">{dish.description}</p>
            <p className="font-semibold">Price: ${dish.price.toFixed(2)}</p>
          </div>
          <div className="flex items-center justify-end mt-2">
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleDecrease}
                disabled={!currentDish}
              >
                <Minus className="h-3 w-3" />
              </Button>
              <span className="w-8 text-center">{quantity}</span>
              <Button variant="default" size="sm" onClick={handleIncrease}>
                <Plus className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default DishCard;