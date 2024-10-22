"use client";

import React, { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp } from "lucide-react";
import useOrderStore from "@/zusstand/order/order_zustand";

const OrderSummary = () => {
  const {
    getOrderSummary,
    updateDishQuantity,
    updateSetQuantity,
    removeDishItem,
    removeSetItem
  } = useOrderStore();

  const [showSets, setShowSets] = useState(true);
  const [showDishes, setShowDishes] = useState(true);
  const [bowlChili, setBowlChili] = useState(1);
  const [bowlNoChili, setBowlNoChili] = useState(2);

  const { dishes, sets, totalItems, totalPrice } = getOrderSummary();

  const handleDishQuantityChange = (id: number, change: number) => {
    const dish = dishes.find((d) => d.id === id);
    if (dish) {
      const newQuantity = dish.quantity + change;
      if (newQuantity > 0) {
        updateDishQuantity(id, newQuantity);
      } else {
        removeDishItem(id);
      }
    }
  };

  const handleSetQuantityChange = (id: number, change: number) => {
    const set = sets.find((s) => s.id === id);
    if (set) {
      const newQuantity = set.quantity + change;
      if (newQuantity > 0) {
        updateSetQuantity(id, newQuantity);
      } else {
        removeSetItem(id);
      }
    }
  };

  const handleBowlChange = (type: "chili" | "noChili", change: number) => {
    if (type === "chili") {
      const newValue = bowlChili + change;
      if (newValue >= 0) setBowlChili(newValue);
    } else {
      const newValue = bowlNoChili + change;
      if (newValue >= 0) setBowlNoChili(newValue);
    }
  };

  return (
    <div className="container mx-auto px-4 py-5 space-y-5">
      <Card>
        <CardHeader>
          <CardTitle>canh banh cuon</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span>canh khong rau</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", -1)}
                  >
                    -
                  </Button>
                  <span className="mx-2">{bowlChili}</span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span>canh rau </span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", -1)}
                  >
                    -
                  </Button>
                  <span className="mx-2">{bowlNoChili}</span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Order Summary</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div>
            <p>Total Items: {totalItems}</p>
            <p>Total Price: ${totalPrice.toFixed(2)}</p>
          </div>

          <div className="border-t border-gray-200 my-4" />

          <div>
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold text-lg">Sets</h3>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowSets(!showSets)}
              >
                {showSets ? <ChevronUp /> : <ChevronDown />}
              </Button>
            </div>
            {showSets &&
              sets.map((set) => (
                <div key={set.id} className="mb-4">
                  <div className="flex items-center justify-between mb-2">
                    <span>
                      {set.name} - ${set.price.toFixed(2)} x {set.quantity}
                    </span>
                    <div>
                      <Button
                        onClick={() => handleSetQuantityChange(set.id, -1)}
                      >
                        -
                      </Button>
                      <span className="mx-2">{set.quantity}</span>
                      <Button
                        onClick={() => handleSetQuantityChange(set.id, 1)}
                      >
                        +
                      </Button>
                    </div>
                  </div>
                  <div className="ml-4">
                    {set.dishes && Array.isArray(set.dishes) ? (
                      set.dishes.map((dish) => (
                        <div key={dish.dish_id} className="text-sm">
                          {dish.name} x {dish.quantity}
                        </div>
                      ))
                    ) : (
                      <div className="text-sm text-gray-500">
                        No dishes in this set
                      </div>
                    )}
                  </div>
                </div>
              ))}
          </div>

          <div className="border-t border-gray-200 my-4" />

          <div>
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold text-lg">Dishes</h3>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowDishes(!showDishes)}
              >
                {showDishes ? <ChevronUp /> : <ChevronDown />}
              </Button>
            </div>
            {showDishes &&
              dishes.map((dish) => (
                <div
                  key={dish.id}
                  className="flex items-center justify-between mb-2"
                >
                  <span>
                    {dish.name} - ${dish.price.toFixed(2)} x {dish.quantity}
                  </span>
                  <div>
                    <Button
                      onClick={() => handleDishQuantityChange(dish.id, -1)}
                    >
                      -
                    </Button>
                    <span className="mx-2">{dish.quantity}</span>
                    <Button
                      onClick={() => handleDishQuantityChange(dish.id, 1)}
                    >
                      +
                    </Button>
                  </div>
                </div>
              ))}
          </div>
        </CardContent>
      </Card>
      <div className="mt-4">
        <Button className="w-full" onClick={() => {}}>
          Add Order
        </Button>
      </div>
    </div>
  );
};

export default OrderSummary;
