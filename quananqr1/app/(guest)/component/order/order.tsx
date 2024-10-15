"use client";

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";

const OrderSummary = () => {
  const { items, updateQuantity, removeItem, getOrderSummary } =
    useOrderStore();
  const { dishes, sets, totalItems, totalPrice } = getOrderSummary();

  const handleQuantityChange = (
    type: "dish" | "set",
    id: number,
    change: number
  ) => {
    const item = items.find((i) => i.type === type && i.id === id);
    if (item) {
      const newQuantity = item.quantity + change;
      if (newQuantity > 0) {
        updateQuantity(type, id, newQuantity);
      } else {
        removeItem(type, id);
      }
    }
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Order Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Total Items: {totalItems}</p>
          <p>Total Price: ${totalPrice.toFixed(2)}</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Dishes</CardTitle>
        </CardHeader>
        <CardContent>
          {dishes.map((dish) => (
            <div
              key={dish.id}
              className="flex items-center justify-between mb-2"
            >
              <span>
                {dish.name} - ${dish.price.toFixed(2)} x {dish.quantity}
              </span>
              <div>
                <Button
                  onClick={() => handleQuantityChange("dish", dish.id, -1)}
                >
                  -
                </Button>
                <span className="mx-2">{dish.quantity}</span>
                <Button
                  onClick={() => handleQuantityChange("dish", dish.id, 1)}
                >
                  +
                </Button>
              </div>
            </div>
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Sets</CardTitle>
        </CardHeader>
        <CardContent>
          {sets.map((set) => (
            <div key={set.id} className="mb-4">
              <div className="flex items-center justify-between mb-2">
                <span>
                  {set.name} - ${set.price.toFixed(2)} x {set.quantity}
                </span>
                <div>
                  <Button
                    onClick={() => handleQuantityChange("set", set.id, -1)}
                  >
                    -
                  </Button>
                  <span className="mx-2">{set.quantity}</span>
                  <Button
                    onClick={() => handleQuantityChange("set", set.id, 1)}
                  >
                    +
                  </Button>
                </div>
              </div>
              <div className="ml-4">
                {set.dishes && Array.isArray(set.dishes) ? (
                  set.dishes.map((dish) => (
                    <div key={dish.dishId} className="text-sm">
                      {dish.dish.name} x {dish.quantity}
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
        </CardContent>
      </Card>
    </div>
  );
};

export default OrderSummary;
