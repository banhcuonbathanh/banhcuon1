"use client";

import React, { useState } from "react";
import Image from "next/image";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Dish, DishListResType } from "@/zusstand/dished/domain/dish.schema";

interface DishCardProps {
  dish: Dish;
  onAddToOrder: (dish: Dish) => void;
}

const DishCard: React.FC<DishCardProps> = ({ dish, onAddToOrder }) => (
  <Card className="w-full max-w-sm">
    <CardHeader>
      <CardTitle>{dish.name}</CardTitle>
    </CardHeader>
    <CardContent>
      <div className="aspect-square relative mb-2">
        <Image
          src={dish.image}
          alt={dish.name}
          layout="fill"
          objectFit="cover"
          className="rounded-md"
        />
      </div>
      <p className="text-sm text-gray-600 mb-2">{dish.description}</p>
      <p className="font-bold">${dish.price.toFixed(2)}</p>
    </CardContent>
    <CardFooter>
      <Button onClick={() => onAddToOrder(dish)} className="w-full">
        Add to Order
      </Button>
    </CardFooter>
  </Card>
);

interface DishSelectionProps {
  dishesData: DishListResType;
}

export function DishSelection({ dishesData }: DishSelectionProps) {
  const [order, setOrder] = useState<Dish[]>([]);

  const addToOrder = (dish: Dish) => {
    setOrder([...order, dish]);
  };

  const dishes = dishesData.data;

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Menu</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {dishes.map((dish) => (
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
