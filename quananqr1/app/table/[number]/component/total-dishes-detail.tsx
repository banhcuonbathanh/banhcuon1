import React, { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ChevronDown, ChevronUp } from "lucide-react";
import useOrderStore, {
  DishOrderItemustand,
  SetOrderItemustand
} from "@/zusstand/order/order_zustand";

interface OrderDetailsProps {
  dishes: DishOrderItemustand[];
  sets: SetOrderItemustand[];
  totalPrice: number;
  totalItems: number;
}

export default function OrderDetails({
  dishes,
  sets,
  totalPrice,
  totalItems
}: OrderDetailsProps) {
  const { tableNumber } = useOrderStore();
  const [expandedSets, setExpandedSets] = useState<Record<string, boolean>>({});

  const toggleSetExpansion = (setId: number) => {
    setExpandedSets((prev) => ({
      ...prev,
      [setId]: !prev[setId]
    }));
  };

  // Calculate dish totals from both sets and individual dishes
  const calculateDishTotals = () => {
    const dishTotals = new Map<
      string,
      { quantity: number; totalPrice: number }
    >();

    // Calculate totals from sets
    sets.forEach((set) => {
      set.dishes.forEach((dish) => {
        const totalQuantity = set.quantity * dish.quantity;
        const totalPrice = totalQuantity * dish.price;
        const current = dishTotals.get(dish.name) || {
          quantity: 0,
          totalPrice: 0
        };
        dishTotals.set(dish.name, {
          quantity: current.quantity + totalQuantity,
          totalPrice: current.totalPrice + totalPrice
        });
      });
    });

    // Add individual dishes to totals
    dishes.forEach((dish) => {
      const current = dishTotals.get(dish.name) || {
        quantity: 0,
        totalPrice: 0
      };
      dishTotals.set(dish.name, {
        quantity: current.quantity + dish.quantity,
        totalPrice: current.totalPrice + dish.quantity * dish.price
      });
    });

    return dishTotals;
  };

  const dishTotals = calculateDishTotals();

  // Calculate prices for display
  const calculateSetPrice = (set: SetOrderItemustand) => {
    return set.dishes.reduce(
      (acc, dish) => acc + dish.price * dish.quantity,
      0
    );
  };

  const setsTotalPrice = sets.reduce(
    (acc, set) => acc + calculateSetPrice(set) * set.quantity,
    0
  );
  const dishesTotalPrice = dishes.reduce(
    (acc, dish) => acc + dish.price * dish.quantity,
    0
  );

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          <span>Order Summary - Table {tableNumber}</span>
          <span className="text-base font-bold">${totalPrice.toFixed(2)}</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Sets Section */}
        {sets.length > 0 && (
          <div>
            <h3 className="font-semibold mb-2 flex justify-between">
              <span>Sets</span>
              <span className="text-primary">${setsTotalPrice.toFixed(2)}</span>
            </h3>
            <div className="space-y-2">
              {sets.map((set) => (
                <div key={set.id} className="border rounded-lg p-3">
                  <div
                    className="flex justify-between items-center cursor-pointer"
                    onClick={() => toggleSetExpansion(set.id)}
                  >
                    <span className="text-gray-400 font-medium">
                      {set.name}
                    </span>
                    <div className="flex items-center space-x-4">
                      <span className="text-sm text-gray-400">
                        {set.quantity} x ${set.price.toFixed(2)} = $
                        {(set.quantity * set.price).toFixed(2)}
                      </span>
                      {expandedSets[set.id] ? (
                        <ChevronUp className="h-4 w-4" />
                      ) : (
                        <ChevronDown className="h-4 w-4" />
                      )}
                    </div>
                  </div>
                  {expandedSets[set.id] && (
                    <div className="mt-2 ml-4 text-sm text-gray-400 space-y-1">
                      {set.dishes.map((dish, idx) => (
                        <div key={idx} className="flex justify-between">
                          <span>{dish.name}</span>
                          <span>
                            {dish.quantity} x ${dish.price.toFixed(2)} = $
                            {(dish.quantity * dish.price).toFixed(2)}
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

        {/* Individual Dishes Section */}
        {dishes.length > 0 && (
          <div>
            <h3 className="font-semibold mb-2 flex justify-between">
              <span>Individual Dishes</span>
              <span className="text-primary">
                ${dishesTotalPrice.toFixed(2)}
              </span>
            </h3>
            <div className="space-y-2">
              {dishes.map((dish) => (
                <div key={dish.id} className="border rounded-lg p-3">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">{dish.name}</span>
                    <span className="text-sm text-gray-400">
                      {dish.quantity} x ${dish.price.toFixed(2)} = $
                      {(dish.quantity * dish.price).toFixed(2)}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Total Items Breakdown */}
        <div className="border rounded-lg p-4">
          <h3 className="font-semibold mb-3">Total Items Breakdown</h3>
          <div className="space-y-2">
            {Array.from(dishTotals.entries()).map(([dishName, details]) => (
              <div
                key={dishName}
                className="flex justify-between items-center text-sm border-b pb-2"
              >
                <div className="flex items-center space-x-2">
                  <span className="text-gray-400 font-medium">{dishName}</span>
                  <span className="text-gray-400">x {details.quantity}</span>
                </div>
                <span className="font-medium">
                  ${details.totalPrice.toFixed(2)}
                </span>
              </div>
            ))}
            <div className="pt-3 mt-2">
              <div className="flex justify-between font-bold text-lg">
                <span>Total Items</span>
                <span>{totalItems}</span>
              </div>
              <div className="flex justify-between font-bold text-lg text-primary mt-2">
                <span>Grand Total</span>
                <span>${totalPrice.toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
