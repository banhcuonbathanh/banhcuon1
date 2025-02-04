"use client";

import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ChevronDown, ChevronUp } from "lucide-react";
import { Order, SetOrderItem } from "@/schemaValidations/interface/type_order";
import useCartStore from "@/zusstand/new-order/new-order-zustand";
import ItemsBreakdown from "./total-items-break-down";
const deliveryData = {
  "banh cuon": 30,
  trung: 10,
  gio: 12
};

const remainingData = {
  "banh cuon": 15,
  trung: 7,
  gio: 7
};

export default function OrderDetails() {
  const [isMounted, setIsMounted] = useState(false);
  const [expandedSets, setExpandedSets] = useState<Record<string, boolean>>({});

  const cartStore = useCartStore();
  const current_order = cartStore.current_order;
  const tableNumber = cartStore.tableNumber;

  useEffect(() => {
    setIsMounted(true);
  }, []);

  const toggleSetExpansion = (setId: number) => {
    setExpandedSets((prev) => ({
      ...prev,
      [setId]: !prev[setId]
    }));
  };

  const sets = current_order?.set_items || [];
  const dishes = current_order?.dish_items || [];

  const calculateDishTotals = () => {
    const dishTotals = new Map<
      string,
      { quantity: number; totalPrice: number; dishId: number }
    >();

    // Calculate totals from sets
    sets.forEach((set) => {
      set.dishes.forEach((dish) => {
        const totalQuantity = set.quantity * dish.quantity;
        const totalPrice = totalQuantity * dish.price;
        const dishKey = `${dish.name}-${dish.dish_id}`;
        const current = dishTotals.get(dishKey) || {
          quantity: 0,
          totalPrice: 0,
          dishId: dish.dish_id
        };
        dishTotals.set(dishKey, {
          quantity: current.quantity + totalQuantity,
          totalPrice: current.totalPrice + totalPrice,
          dishId: dish.dish_id
        });
      });
    });

    // Add individual dishes to totals
    dishes.forEach((dish) => {
      const dishKey = `${dish.name}-${dish.dish_id}`;
      const current = dishTotals.get(dishKey) || {
        quantity: 0,
        totalPrice: 0,
        dishId: dish.dish_id
      };
      dishTotals.set(dishKey, {
        quantity: current.quantity + dish.quantity,
        totalPrice: current.totalPrice + dish.quantity * dish.price,
        dishId: dish.dish_id
      });
    });

    return dishTotals;
  };

  const calculateSetPrice = (set: SetOrderItem) => {
    return set.dishes.reduce(
      (acc, dish) => acc + dish.price * dish.quantity,
      0
    );
  };

  const dishTotals = calculateDishTotals();

  const setsTotalPrice = sets.reduce(
    (acc, set) => acc + calculateSetPrice(set) * set.quantity,
    0
  );
  const dishesTotalPrice = dishes.reduce(
    (acc, dish) => acc + dish.price * dish.quantity,
    0
  );

  const totalPrice = setsTotalPrice + dishesTotalPrice;

  if (!isMounted) {
    return null;
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          <span>Order Summary - Table {tableNumber}</span>
          <span className="text-base font-bold">{totalPrice} K</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {sets.length > 0 && (
          <div>
            <h3 className="font-semibold mb-2 flex justify-between">
              <span>Sets</span>
              <span className="text-primary">{setsTotalPrice} K</span>
            </h3>
            <div className="space-y-2">
              {sets.map((set) => (
                <div
                  key={`set-${set.set_id}`}
                  className="border rounded-lg p-3"
                >
                  <div
                    className="flex justify-between items-center cursor-pointer"
                    onClick={() => toggleSetExpansion(set.set_id)}
                  >
                    <span className="text-gray-400 font-medium">
                      {set.name}
                    </span>
                    <div className="flex items-center space-x-4">
                      <span className="text-sm text-primary">
                        {set.quantity} x {set.price} K =
                        {set.quantity * set.price} K
                      </span>
                      {expandedSets[set.set_id] ? (
                        <ChevronUp className="h-4 w-4" />
                      ) : (
                        <ChevronDown className="h-4 w-4" />
                      )}
                    </div>
                  </div>
                  {expandedSets[set.set_id] && (
                    <div className="mt-2 ml-4 text-sm text-gray-400 space-y-1">
                      {set.dishes.map((dish) => (
                        <div
                          key={`set-${set.set_id}-dish-${dish.dish_id}`}
                          className="flex justify-between text-primary"
                        >
                          <span>{dish.name}</span>
                          <span>
                            {dish.quantity} x {dish.price} K =
                            {dish.quantity * dish.price} K
                          </span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {dishes.length > 0 && (
          <div>
            <h3 className="font-semibold mb-2 flex justify-between">
              <span>Individual Dishes</span>
              <span className="text-primary">{dishesTotalPrice} K</span>
            </h3>
            <div className="space-y-2">
              {dishes.map((dish) => (
                <div
                  key={`dish-${dish.dish_id}`}
                  className="border rounded-lg p-3"
                >
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">{dish.name}</span>
                    <span className="text-sm text-primary">
                      {dish.quantity} x {dish.price} K =
                      {dish.quantity * dish.price} K
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <h3 className="font-semibold mb-2 flex justify-between">
          <span>Total Items Breakdown</span>
        </h3>
        <ItemsBreakdown
          dishTotals={dishTotals}
          deliveryData={deliveryData}
          remainingData={remainingData}
        />
      </CardContent>
    </Card>
  );
}
