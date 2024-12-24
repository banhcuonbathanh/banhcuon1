"use client";

import Image from "next/image";
import React, { useEffect, useCallback, useMemo } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Minus, ChevronDown, ChevronUp } from "lucide-react";
import { SetInterface } from "@/schemaValidations/interface/types_set";
import useOrderStore from "@/zusstand/order/order_zustand";
import DishList from "./set_dish";
import { logWithLevel } from "@/lib/logger/log";

const LOG_PATH = "quananqr1/components/set/set_card.tsx";

interface SetSelectionProps {
  set: SetInterface;
}

const SetCard: React.FC<SetSelectionProps> = ({ set }) => {
  // Component initialization logging
  useEffect(() => {
    logWithLevel(
      {
        setId: set.id,
        setName: set.name,
        initialProps: set
      },
      LOG_PATH,
      "info",
      1
    );
  }, [set]);

  const [isListVisible, setIsListVisible] = React.useState(false);

  const {
    currentOrder,
    addSetToCurrentOrder,
    updateSetQuantityInCurrentOrder,
    removeSetFromCurrentOrder,
    clearCurrentOrder
  } = useOrderStore();

  // Initialize currentOrder if it's null
  useEffect(() => {
    if (!currentOrder) {
      logWithLevel({ action: "initializeCurrentOrder" }, LOG_PATH, "info", 2);
      clearCurrentOrder();
    }
  }, [currentOrder, clearCurrentOrder]);

  // Debug logging for store updates
  useEffect(() => {
    const unsubscribe = useOrderStore.subscribe((state) => {
      logWithLevel(
        {
          currentOrder: state.currentOrder,
          timestamp: new Date().toISOString()
        },
        LOG_PATH,
        "debug",
        2
      );
    });

    return () => unsubscribe();
  }, []);

  const setOrderItem = currentOrder?.set_items.find(
    (item) => item.set_id === set.id
  );

  const [dishQuantities, setDishQuantities] = React.useState<
    Record<number, number>
  >({});

  useEffect(() => {
    setDishQuantities(
      set.dishes.reduce(
        (acc, dish) => ({ ...acc, [dish.dish_id]: dish.quantity }),
        {}
      )
    );
  }, [set.dishes]);

  const totalPrice = useMemo(() => {
    const price = set.dishes.reduce(
      (sum, dish) => sum + dish.price * (dishQuantities[dish.dish_id] || 0),
      0
    );
    logWithLevel(
      {
        setId: set.id,
        totalPrice: price,
        dishQuantities
      },
      LOG_PATH,
      "debug",
      5
    );
    return price;
  }, [set.dishes, dishQuantities, set.id]);

  const handleIncrease = useCallback(() => {
    logWithLevel(
      {
        action: "increase",
        setId: set.id,
        currentQuantity: setOrderItem?.quantity || 0
      },
      LOG_PATH,
      "info",
      3
    );

    if (setOrderItem) {
      updateSetQuantityInCurrentOrder(set.id, (setOrderItem.quantity || 0) + 1);
    } else {
      addSetToCurrentOrder({
        set_id: set.id,
        quantity: 1
      });
    }
  }, [
    set.id,
    setOrderItem,
    addSetToCurrentOrder,
    updateSetQuantityInCurrentOrder
  ]);

  const handleDecrease = useCallback(() => {
    logWithLevel(
      {
        action: "decrease",
        setId: set.id,
        currentQuantity: setOrderItem?.quantity || 0
      },
      LOG_PATH,
      "info",
      3
    );

    if (setOrderItem) {
      if ((setOrderItem.quantity || 0) > 1) {
        updateSetQuantityInCurrentOrder(
          set.id,
          (setOrderItem.quantity || 0) - 1
        );
      } else {
        removeSetFromCurrentOrder(set.id);
      }
    }
  }, [
    set.id,
    setOrderItem,
    updateSetQuantityInCurrentOrder,
    removeSetFromCurrentOrder
  ]);

  const handleDishIncrease = useCallback(
    (dishId: number) => {
      logWithLevel(
        {
          action: "increaseDish",
          dishId,
          currentQuantity: dishQuantities[dishId] || 0
        },
        LOG_PATH,
        "info",
        4
      );

      setDishQuantities((prev) => ({
        ...prev,
        [dishId]: (prev[dishId] || 0) + 1
      }));
    },
    [dishQuantities]
  );

  const handleDishDecrease = useCallback(
    (dishId: number) => {
      logWithLevel(
        {
          action: "decreaseDish",
          dishId,
          currentQuantity: dishQuantities[dishId] || 0
        },
        LOG_PATH,
        "info",
        4
      );

      setDishQuantities((prev) => ({
        ...prev,
        [dishId]: Math.max(0, (prev[dishId] || 0) - 1)
      }));
    },
    [dishQuantities]
  );

  const toggleList = () => {
    logWithLevel(
      {
        action: "toggleList",
        setId: set.id,
        newVisibility: !isListVisible
      },
      LOG_PATH,
      "debug",
      6
    );
    setIsListVisible(!isListVisible);
  };

  // Error boundary for rendering
  if (!set) {
    logWithLevel(
      {
        error: "Set prop is undefined",
        component: "SetCard"
      },
      LOG_PATH,
      "error",
      7
    );
    return null;
  }

  return (
    <Card className="w-full">
      <CardContent className="p-4 flex">
        <div className="w-1/3 pr-4">
          <Image
            src={set.image || "/api/placeholder/300/400"}
            alt={set.name}
            className="w-full h-full object-cover rounded-md"
            width={200}
            height={200}
            priority
          />
        </div>
        <div className="w-2/3 flex flex-col justify-between">
          <div className="space-y-2">
            <div className="flex flex-row justify-between items-center">
              <button
                onClick={toggleList}
                className="flex items-center gap-2 hover:opacity-75 transition-opacity"
              >
                <h2 className="text-2xl font-bold">{set.name}</h2>
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
                dishQuantities={dishQuantities}
                onIncrease={handleDishIncrease}
                onDecrease={handleDishDecrease}
              />
            </div>
          )}

          <div className="flex items-center justify-end mt-4">
            <div className="flex items-center space-x-4">
              <Button
                variant="outline"
                onClick={handleDecrease}
                disabled={!setOrderItem}
              >
                <Minus className="h-4 w-4" />
              </Button>
              <span className="font-semibold w-8 text-center">
                {setOrderItem?.quantity || 0}
              </span>
              <Button variant="default" onClick={handleIncrease} className=" ">
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
