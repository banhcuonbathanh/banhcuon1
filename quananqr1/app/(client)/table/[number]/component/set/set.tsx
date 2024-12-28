"use client";

import Image from "next/image";
import React, { useEffect, useCallback, useMemo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Minus, ChevronDown, ChevronUp } from "lucide-react";
import DishList from "./set_dish";
import { logWithLevel } from "@/lib/logger/log";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import { SetOrderItem } from "@/schemaValidations/interface/type_order";
import { SetInterface } from "@/schemaValidations/interface/types_set";

const LOG_PATH = "quananqr1/components/set/set_card.tsx";

interface SetSelectionProps {
  set: SetInterface;
}

const SetCard: React.FC<SetSelectionProps> = ({ set }) => {
  const [isListVisible, setIsListVisible] = React.useState(false);
  const [localQuantity, setLocalQuantity] = React.useState(0);

  const { current_order, addSetToCart, updateSetQuantity, removeSetFromCart } =
    useCartStore();

  // Debug logging for initial render

  // Find current set in order with useMemo to optimize performance
  const setOrderItem = useMemo(() => {
    return current_order?.set_items.find((item) => item.set_id === set.id);
  }, [current_order?.set_items, set.id]);

  // Update local quantity when setOrderItem changes
  useEffect(() => {
    setLocalQuantity(setOrderItem?.quantity || 0);
  }, [setOrderItem?.quantity]);

  const handleIncrease = useCallback(() => {
    console.log("Increase clicked for set:", set.id);

    if (setOrderItem) {
      updateSetQuantity("increment", set.id);
    } else {
      const newSet: SetOrderItem = {
        set_id: set.id,
        quantity: 1,
        userId: set.userId,
        name: set.name,
        description: set.description || "",
        price: set.price,
        image: set.image || "",
        status: "active",
        created_at: set.created_at,
        updated_at: set.updated_at,
        is_favourite: set.is_favourite || false,
        like_by: set.like_by || [],
        is_public: set.is_public || false,
        dishes: set.dishes.map((dish) => ({
          dish_id: dish.dish_id,
          name: dish.name,
          price: dish.price,
          quantity: dish.quantity,
          description: dish.description || "",
          image: dish.image || "",
          status: dish.status || "active",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          is_favourite: false,
          like_by: []
        }))
      };
      addSetToCart(newSet);
    }
  }, [set, setOrderItem, addSetToCart, updateSetQuantity]);

  const handleDecrease = useCallback(() => {
    console.log("Decrease clicked for set:", set.id);

    if (!setOrderItem) return;

    if (setOrderItem.quantity <= 1) {
      removeSetFromCart(set.id);
    } else {
      updateSetQuantity("decrement", set.id);
    }
  }, [set.id, setOrderItem, updateSetQuantity, removeSetFromCart]);

  const toggleList = () => {
    setIsListVisible(!isListVisible);
  };

  // Safety check
  if (!set) return null;

  return (
    <Card className="w-full">
      <CardContent className="p-4 flex md:flex-row">
        <div className="w-full md:w-1/4 md:pr-4 mb-4 md:mb-0">
          <Image
            src={set.image || "/api/placeholder/300/400"}
            alt={set.name}
            className="w-full h-48 md:h-full object-cover rounded-md"
            width={300}
            height={400}
            priority
          />
        </div>

        <div className="w-full md:w-2/3 flex flex-col justify-between">
          <div className="space-y-2">
            <div className="flex flex-row justify-between items-center">
              <button
                onClick={toggleList}
                className="flex items-center gap-2 hover:opacity-75 transition-opacity"
              >
                <h2 className="text-xl md:text-2xl font-bold">{set.name}</h2>
                {isListVisible ? (
                  <ChevronUp className="h-4 w-4 mt-1" />
                ) : (
                  <ChevronDown className="h-4 w-4 mt-1" />
                )}
              </button>
              <span className="font-semibold text-lg">
                {typeof set.price === "number" ? `${set.price}k` : ""}
              </span>
            </div>

            <p className="text-sm text-gray-600">{set.description}</p>
          </div>

          {isListVisible && (
            <div className="mt-4">
              <DishList
                dishes={set.dishes}
                dishQuantities={set.dishes.reduce(
                  (acc, dish) => ({
                    ...acc,
                    [dish.dish_id]: dish.quantity
                  }),
                  {}
                )}
                onIncrease={(dishId) => {
                  console.log("Dish increase:", dishId);
                }}
                onDecrease={(dishId) => {
                  console.log("Dish decrease:", dishId);
                }}
              />
            </div>
          )}

          <div className="flex items-center justify-end mt-4">
            <div className="flex items-center space-x-4">
              <Button
                variant="outline"
                onClick={handleDecrease}
                disabled={!setOrderItem}
                className="h-8 w-8 p-0"
              >
                <Minus className="h-4 w-4" />
              </Button>

              <span className="font-semibold w-8 text-center">
                {localQuantity}
              </span>

              <Button
                variant="default"
                onClick={handleIncrease}
                className="h-8 w-8 p-0"
              >
                <Plus className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default SetCard;
