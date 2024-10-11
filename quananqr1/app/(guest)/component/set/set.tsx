import React, { useState } from "react";
import {
  Card,
  CardHeader,
  CardContent,
  CardFooter,
  CardTitle
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import { SetListResType, Dish } from "@/schemaValidations/dish.schema";

interface DishSelectionProps {
  sets: SetListResType;
}

export function DishSelection({ sets }: DishSelectionProps) {
  const [order, setOrder] = useState<Dish[]>([]);

  const addToOrder = (dish: Dish) => {
    setOrder([...order, dish]);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Menu</h1>
      {sets.map((set) => (
        <div key={set.id} className="mb-8">
          <h2 className="text-2xl font-bold mb-4">{set.name}</h2>
          <p className="text-gray-600 mb-4">{set.description}</p>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {set.dishes.map((dish) => (
              <DishCard key={dish.id} dish={dish} onAddToOrder={addToOrder} />
            ))}
          </div>
        </div>
      ))}
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
