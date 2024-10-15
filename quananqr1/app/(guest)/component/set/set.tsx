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
import { SetInterface } from "@/schemaValidations/interface/types_set";
import { DishInterface } from "@/schemaValidations/interface/type_dish";

interface SetSelectionProps {
  set: SetInterface;
}

interface DishCardProps {
  dish: DishInterface;
  quantity: number;
  onIncrease: () => void;
  onDecrease: () => void;
}

const SetDishCard: React.FC<DishCardProps> = ({
  dish,
  quantity,
  onIncrease,
  onDecrease
}) => (
  <Card className="w-full">
    <CardHeader>
      <CardTitle className="text-lg">{dish.name}</CardTitle>
    </CardHeader>
    <CardContent>
      <img
        src={dish.image || "/api/placeholder/150/100"}
        alt={dish.name}
        className="w-full h-24 object-cover rounded-md mb-2"
      />
      <p className="text-sm">{dish.description}</p>
      <p className="font-semibold mt-2">Price: ${dish.price.toFixed(2)}</p>
      <div className="flex items-center justify-between mt-2">
        <span>Quantity:</span>
        <div className="flex items-center space-x-2">
          <Button variant="outline" size="sm" onClick={onDecrease}>
            <Minus className="h-3 w-3" />
          </Button>
          <span className="w-8 text-center">{quantity}</span>
          <Button variant="outline" size="sm" onClick={onIncrease}>
            <Plus className="h-3 w-3" />
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
);

export function SetCard({ set }: SetSelectionProps) {
  const [quantity, setQuantity] = useState(1);
  const [dishQuantities, setDishQuantities] = useState<Record<number, number>>(
    set.dishes.reduce(
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

  const handleIncrease = () => setQuantity((prev) => prev + 1);
  const handleDecrease = () => setQuantity((prev) => Math.max(1, prev - 1));

  const handleDishIncrease = (dishId: number) => {
    setDishQuantities((prev) => ({
      ...prev,
      [dishId]: (prev[dishId] || 0) + 1
    }));
  };

  const handleDishDecrease = (dishId: number) => {
    setDishQuantities((prev) => ({
      ...prev,
      [dishId]: Math.max(0, (prev[dishId] || 0) - 1)
    }));
  };

  return (
    <Card className="w-full max-w-sm mx-auto">
      <CardHeader>
        <CardTitle className="text-xl font-bold">{set.name}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <img
          src={set.image || "/api/placeholder/300/200"}
          alt={set.name}
          className="w-full h-48 object-cover rounded-md"
        />
        <p className="text-sm text-gray-600">{set.description}</p>
        <div className="flex justify-between items-center">
          <span className="font-semibold">
            Total Price: ${(totalPrice * quantity).toFixed(2)}
          </span>
          <span className="text-sm">({totalDishes} dishes)</span>
        </div>
        <div>
          <h3 className="font-semibold mb-2">Dishes:</h3>
          <ul className="list-disc pl-5">
            {set.dishes.map((dish) => (
              <li key={`${set.id}-${dish.dishId}`}>
                {dish.dish.name} - ${dish.dish.price.toFixed(2)} x{" "}
                {dishQuantities[dish.dishId] || 0}
              </li>
            ))}
          </ul>
        </div>
        <div className="flex items-center justify-between">
          <span className="font-semibold">Set Quantity:</span>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="icon" onClick={handleDecrease}>
              <Minus className="h-4 w-4" />
            </Button>
            <span className="w-8 text-center">{quantity}</span>
            <Button variant="outline" size="icon" onClick={handleIncrease}>
              <Plus className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Dialog>
          <DialogTrigger asChild>
            <Button className="w-full">
              <Info className="mr-2 h-4 w-4" /> View Details
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-4xl">
            <DialogHeader>
              <DialogTitle>{set.name} - Detailed View</DialogTitle>
            </DialogHeader>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {set.dishes.map((dish) => (
                <SetDishCard
                  key={`${set.id}-${dish.dishId}`}
                  dish={dish.dish}
                  quantity={dishQuantities[dish.dishId] || 0}
                  onIncrease={() => handleDishIncrease(dish.dishId)}
                  onDecrease={() => handleDishDecrease(dish.dishId)}
                />
              ))}
            </div>
          </DialogContent>
        </Dialog>
      </CardFooter>
    </Card>
  );
}

export default SetCard;
