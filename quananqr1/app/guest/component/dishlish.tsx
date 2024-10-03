"use client"; // This marks the component as a client component

import React, { useState } from "react";
import { Dish } from "@/zusstand/dished/domain/dish.schema";

// Component to represent a single dish
const DishItem = ({ dish, handleAddToOrder }: { dish: Dish; handleAddToOrder: (dish: Dish) => void }) => {
  return (
    <div className="dish-item">
      <img src={dish.image} alt={dish.name} width={100} height={100} />
      <div>
        <h3>{dish.name}</h3>
        <p>{dish.description}</p>
        <p>Price: ${dish.price}</p>
        <button onClick={() => handleAddToOrder(dish)}>Add to Order</button>
      </div>
    </div>
  );
};

// This is a client component to manage dishes and the order
export default function DishList({ dishes = [] }: { dishes: Dish[] | undefined }) { // Ensure dishes is an array
  const [order, setOrder] = useState<Dish[]>([]);

  const handleAddToOrder = (dish: Dish) => {
    setOrder([...order, dish]);
  };

  return (
    <div>
      <div className="dish-list">
        {dishes.map((dish) => (
          <DishItem key={dish.id} dish={dish} handleAddToOrder={handleAddToOrder} />
        ))}
      </div>

      <h2>Your Order</h2>
      {order.length > 0 ? (
        <ul>
          {order.map((dish, index) => (
            <li key={index}>
              {dish.name} - ${dish.price}
            </li>
          ))}
        </ul>
      ) : (
        <p>No items in your order.</p>
      )}
    </div>
  );
}
