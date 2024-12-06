"use client";

import React, { useState } from "react";
import {
  OrderDetailedResponse,
  OrderSetDetailed,
  OrderDetailedDish
} from "../component/new-order-column";
import { ChevronDown, ChevronRight } from "lucide-react";
import { TotalPriceSummary } from "./totalpricesummary";

interface RestaurantSummaryProps {
  restaurantLayoutProps: OrderDetailedResponse[];
}

interface AggregatedDish {
  dish_id: number;
  name: string;
  quantity: number;
  price: number;
}

// CollapsibleSection component
const CollapsibleSection: React.FC<{
  title: string;
  children: React.ReactNode;
}> = ({ title, children }) => {
  // Set initial state based on section title
  const [isOpen, setIsOpen] = useState(
    title === "Topping Details" || title === "Aggregated Dishes"
  );

  return (
    <div className="mt-2">
      <div
        className="flex items-center cursor-pointer select-none p-2 rounded"
        onClick={() => setIsOpen(!isOpen)}
      >
        <h3 className="text-md font-semibold">{title}</h3>
        {isOpen ? (
          <ChevronDown className="ml-2 h-4 w-4" />
        ) : (
          <ChevronRight className="ml-2 h-4 w-4" />
        )}
      </div>
      {isOpen && (
        <div className="p-2 border-l border-r border-b rounded-b">
          {children}
        </div>
      )}
    </div>
  );
};

export const RestaurantSummary: React.FC<RestaurantSummaryProps> = ({
  restaurantLayoutProps
}) => {
  // Function to extract the first part of the order name
  const extractOrderPrefix = (orderName: string) => {
    return orderName.split("-")[0];
  };

  // Function to determine status based on takeAway
  const determineStatus = (order: OrderDetailedResponse) => {
    return order.takeAway ? "Take Away" : order.status;
  };

  // Function to get status color
  const getStatusColor = (order: OrderDetailedResponse) => {
    return order.takeAway ? "text-red-600 font-bold" : "";
  };

  // Function to aggregate dishes from both individual dishes and set dishes
  const aggregateDishes = (order: OrderDetailedResponse): AggregatedDish[] => {
    const dishMap = new Map<number, AggregatedDish>();

    // Add individual dishes
    order.data_dish.forEach((dish) => {
      const existingDish = dishMap.get(dish.dish_id);
      if (existingDish) {
        existingDish.quantity += dish.quantity;
      } else {
        dishMap.set(dish.dish_id, { ...dish, quantity: dish.quantity });
      }
    });

    // Add dishes from sets
    order.data_set.forEach((set) => {
      set.dishes.forEach((setDish) => {
        const existingDish = dishMap.get(setDish.dish_id);
        if (existingDish) {
          existingDish.quantity += setDish.quantity * set.quantity;
        } else {
          dishMap.set(setDish.dish_id, {
            ...setDish,
            quantity: setDish.quantity * set.quantity
          });
        }
      });
    });

    return Array.from(dishMap.values());
  };

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Restaurant Order Summary</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {restaurantLayoutProps.map((order) => {
          const aggregatedDishes = aggregateDishes(order);

          return (
            <div key={order.id} className="shadow-md rounded-lg p-4">
              <div className="font-semibold">
                Table Number: {order.table_number}
              </div>

              {/* Order Status Section */}
              <CollapsibleSection title="Order Status">
                <div className="grid grid-cols-2 gap-2">
                  <div className="font-semibold">Status:</div>
                  <div className={getStatusColor(order)}>
                    {determineStatus(order)}
                  </div>
                  <div className="font-semibold">Total Price:</div>
                  <div>${order.total_price.toFixed(2)}</div>
                  <div className="font-semibold">Tracking Order:</div>
                  <div>{order.tracking_order}</div>
                  <div className="font-semibold">Chili Number:</div>
                  <div>{order.chiliNumber}</div>
                  <div className="font-semibold">Order Name:</div>
                  <div>{extractOrderPrefix(order.order_name)}</div>
                </div>
              </CollapsibleSection>

              {/* Topping Section */}
              <CollapsibleSection title="Topping Details">
                <div className="grid grid-cols-2 gap-2">
                  <div className="font-semibold">Topping:</div>
                  <div>{order.topping}</div>
                </div>
              </CollapsibleSection>

              {/* Aggregated Dishes Section */}

              {/* Individual Dishes Section */}
              <CollapsibleSection title="Individual Dishes">
                {order.data_dish.map((dish, index) => (
                  <div
                    key={`${dish.dish_id}-${index}`}
                    className="p-2 rounded mb-2"
                  >
                    <div className="grid grid-cols-2 gap-1">
                      <div className="font-medium">Dish Name:</div>
                      <div>{dish.name}</div>
                      <div className="font-medium">Quantity:</div>
                      <div>{dish.quantity}</div>
                      <div className="font-medium">Price:</div>
                      <div>${dish.price.toFixed(2)}</div>
                    </div>
                  </div>
                ))}
              </CollapsibleSection>

              {/* Order Sets Section */}
              <CollapsibleSection title="Order Sets">
                {order.data_set.map((set, index) => (
                  <div key={`${set.id}-${index}`} className="p-2 rounded mb-2">
                    <div className="grid grid-cols-2 gap-1">
                      <div className="font-medium">Set Name:</div>
                      <div>{set.name}</div>
                      <div className="font-medium">Price:</div>
                      <div>${set.price.toFixed(2)}</div>
                      <div className="font-medium">Quantity:</div>
                      <div>{set.quantity}</div>
                      <div className="font-medium">Dishes in Set:</div>
                      <div>
                        {set.dishes
                          .map((d) => `${d.name} (${d.quantity})`)
                          .join(", ")}
                      </div>
                    </div>
                  </div>
                ))}
              </CollapsibleSection>
              <CollapsibleSection title="Aggregated Dishes">
                {aggregatedDishes.map((dish, index) => (
                  <div
                    key={`${dish.dish_id}-${index}`}
                    className="p-2 rounded mb-2 border-b last:border-b-0"
                  >
                    <div className="grid grid-cols-2 gap-1">
                      <div className="font-medium">Dish Name:</div>
                      <div>{dish.name}</div>
                      <div className="font-medium">Total Quantity:</div>
                      <div>{dish.quantity}</div>
                    </div>
                    <DishPriceDetails
                      price={dish.price}
                      quantity={dish.quantity}
                    />
                  </div>
                ))}
              </CollapsibleSection>
              <div className="mt-6">
                <TotalPriceSummary TotalPriceProps={restaurantLayoutProps} />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default RestaurantSummary;

const DishPriceDetails: React.FC<{
  price: number;
  quantity: number;
}> = ({ price, quantity }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="mt-1">
      <div
        className="flex items-center cursor-pointer select-none text-sm text-blue-600 hover:text-blue-800"
        onClick={() => setIsOpen(!isOpen)}
      >
        {isOpen ? (
          <ChevronDown className="h-3 w-3 mr-1" />
        ) : (
          <ChevronRight className="h-3 w-3 mr-1" />
        )}
        <span>Price Details</span>
      </div>
      {isOpen && (
        <div className="pl-4 py-2 text-sm">
          <div className="grid grid-cols-2 gap-1">
            <div className="font-medium">Price per Dish:</div>
            <div>${price.toFixed(2)}</div>
            <div className="font-medium">Total Dish Price:</div>
            <div>${(quantity * price).toFixed(2)}</div>
          </div>
        </div>
      )}
    </div>
  );
};
