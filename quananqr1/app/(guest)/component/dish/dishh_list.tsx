"use client";
import React, { useState } from "react";

import { Dish, DishListResType } from "@/schemaValidations/dish.schema";
import { DishCard } from "./disih_tem";
import { DishInterface } from "@/schemaValidations/interface/type_dish";

interface DishSelectionProps {
  dishes: DishInterface[];
}

export function DishSelection({ dishes }: DishSelectionProps) {
  const [order, setOrder] = useState<DishInterface[]>([]);

  const addToOrder = (dish: DishInterface) => {
    setOrder([...order, dish]);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Menu</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {dishes.map((dish: DishInterface) => (
          <DishCard key={dish.id} dish={dish} onAddToOrder={addToOrder} />
        ))}
      </div>
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">Your Order</h2>
        {order.length === 0 ? (
          <p>No items in your order yet.</p>
        ) : (
          <ul>
            {order.map((item, index) => (
              <li key={index} className="mb-2">
                {item.name} - ${item.price.toFixed(2)}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
