import React from 'react';

import Image from 'next/image';
import { Dish } from '@/zusstand/dished/domain/dish.schema';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card';


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

export default DishCard;
