"use client";
import React, { useState } from "react";
import {
  Card,
  CardHeader,
  CardContent,
  CardFooter,
  CardTitle
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import { Minus, Plus, Info } from "lucide-react";
import {
  SetInterface,
  SetProtoDish
} from "@/schemaValidations/interface/types_set";
import { DishInterface } from "@/schemaValidations/interface/type_dish";
import useOrderStore from "@/zusstand/order/order_zustand";

interface SetSelectionProps {
  set: SetInterface;
}

interface DishCardProps {
  dish: DishInterface;
  quantity: number;
  onIncrease: () => void;
  onDecrease: () => void;
}

// const SetDishCard: React.FC<DishCardProps> = ({
//   dish,
//   quantity,
//   onIncrease,
//   onDecrease
// }) => (
//   <Card className="w-full">
//     <CardHeader>
//       <CardTitle className="text-lg">{dish.name}</CardTitle>
//     </CardHeader>
//     <CardContent>
//       <img
//         src={dish.image || "/api/placeholder/150/100"}
//         alt={dish.name}
//         className="w-full h-24 object-cover rounded-md mb-2"
//       />
//       <p className="text-sm">{dish.description}</p>
//       <p className="font-semibold mt-2">Price: ${dish.price.toFixed(2)}</p>
//       <div className="flex items-center justify-center mt-2">
//         <div className="flex items-center space-x-2">
//           <Button variant="outline" size="sm" onClick={onDecrease}>
//             <Minus className="h-3 w-3" />
//           </Button>
//           <span className="w-8 text-center">{quantity}</span>
//           <Button variant="outline" size="sm" onClick={onIncrease}>
//             <Plus className="h-3 w-3" />
//           </Button>
//         </div>
//       </div>
//     </CardContent>
//   </Card>
// );

export function SetCard({ set }: SetSelectionProps) {
  const {
    addSetItem,
    updateSetDishes,
    findSetOrderItem,
    updateSetQuantity,
    removeSetItem
  } = useOrderStore();

  const setOrderItem = findSetOrderItem(set.id);

  const [dishQuantities, setDishQuantities] = React.useState<
    Record<number, number>
  >(
    setOrderItem
      ? setOrderItem.modifiedDishes.reduce(
          (acc, dish) => ({ ...acc, [dish.dishId]: dish.quantity }),
          {}
        )
      : set.dishes.reduce(
          (acc, dish) => ({ ...acc, [dish.dishId]: dish.quantity }),
          {}
        )
  );

  const totalPrice = set.dishes.reduce(
    (sum, dish) => sum + dish.dish.price * (dishQuantities[dish.dishId] || 0),
    0
  );
  const totalDishes = Object.values(dishQuantities).reduce(
    (sum, q) => sum + q,
    0
  );

  const handleIncrease = () => {
    if (setOrderItem) {
      updateSetQuantity(set.id, setOrderItem.quantity + 1);
    } else {
      const modifiedDishes: SetProtoDish[] = set.dishes.map((dish) => ({
        ...dish,
        quantity: dishQuantities[dish.dishId] || 0
      }));
      addSetItem(set, 1, modifiedDishes);
    }
  };

  const handleDecrease = () => {
    if (setOrderItem) {
      if (setOrderItem.quantity > 1) {
        updateSetQuantity(set.id, setOrderItem.quantity - 1);
      } else {
        removeSetItem(set.id);
      }
    }
  };

  const handleDishIncrease = (dishId: number) => {
    setDishQuantities((prev) => {
      const newQuantities = {
        ...prev,
        [dishId]: (prev[dishId] || 0) + 1
      };
      if (setOrderItem) {
        updateSetDishes(
          set.id,
          set.dishes.map((dish) => ({
            ...dish,
            quantity: newQuantities[dish.dishId] || 0
          }))
        );
      }
      return newQuantities;
    });
  };

  const handleDishDecrease = (dishId: number) => {
    setDishQuantities((prev) => {
      const newQuantities = {
        ...prev,
        [dishId]: Math.max(0, (prev[dishId] || 0) - 1)
      };
      if (setOrderItem) {
        updateSetDishes(
          set.id,
          set.dishes.map((dish) => ({
            ...dish,
            quantity: newQuantities[dish.dishId] || 0
          }))
        );
      }
      return newQuantities;
    });
  };

  return (
    <Card className="w-full">
      <CardContent className="p-4 flex">
        <div className="w-1/3 pr-4">
          <img
            src={set.image || "/api/placeholder/300/400"}
            alt={set.name}
            className="w-full h-full object-cover rounded-md"
          />
        </div>
        <div className="w-2/3 flex flex-col justify-between">
          <h2 className="text-2xl font-bold">{set.name}</h2>
          <p className="text-sm text-gray-600">{set.description}</p>
          <div className="flex justify-between items-center">
            <span className="font-semibold">
              Total Price: ${totalPrice.toFixed(2)}
            </span>
            <span className="text-sm">({totalDishes} dishes)</span>
          </div>
          <div>
            <h3 className="font-semibold mb-2">Dishes:</h3>
            <ul className="space-y-2">
              {set.dishes.map((dish) => (
                <li
                  key={`${set.id}-${dish.dishId}`}
                  className="flex items-center justify-between"
                >
                  <span>
                    {dish.dish.name} - ${dish.dish.price.toFixed(2)}
                  </span>
                  <div className="flex items-center space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDishDecrease(dish.dishId)}
                    >
                      <Minus className="h-3 w-3" />
                    </Button>
                    <span className="w-8 text-center">
                      {dishQuantities[dish.dishId] || 0}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDishIncrease(dish.dishId)}
                    >
                      <Plus className="h-3 w-3" />
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          </div>
          <div className="flex items-center justify-end mt-2">
            <div className="flex items-center space-x-4">
              <Button onClick={handleDecrease} disabled={!setOrderItem}>
                <Minus className="h-4 w-4" />
              </Button>
              <span className="font-semibold w-8 text-center">
                {setOrderItem ? setOrderItem.quantity : 0}
              </span>
              <Button onClick={handleIncrease}>
                <Plus className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default SetCard;
