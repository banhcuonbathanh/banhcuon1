"use client";

import React, { useState } from "react";
import { Dialog, DialogContent } from "@/components/ui/dialog";

import { AggregatedDish } from "./aggregateDishes";
import NumericKeypad from "./num-pad";

interface DishSummaryProps {
  dish: AggregatedDish;
}

export const DishSummary: React.FC<DishSummaryProps> = ({ dish }) => {
  const [showDetails, setShowDetails] = useState(false);
  const [showKeypad, setShowKeypad] = useState(false);
  const [deliveryNumber, setDeliveryNumber] = useState(0);

  const handleDeliveryClick = () => {
    setShowKeypad(true);
  };

  const handleKeypadSubmit = () => {
    // Here you can handle what happens after the number is submitted
    console.log(`Delivery number set to: ${deliveryNumber}`);
    setShowKeypad(false);
  };

  return (
    <div className="p-2 mb-2">
      <div className="flex items-center justify-between">
        <div
          className="flex-1 cursor-pointer"
          onClick={() => setShowDetails(!showDetails)}
        >
          <span className="font-bold">
            {dish.name} :{dish.quantity} -
          </span>
          <span
            className="text-blue-600 cursor-pointer hover:text-blue-800"
            onClick={(e) => {
              e.stopPropagation();
              handleDeliveryClick();
            }}
          >
            {deliveryNumber > 0 ? `delivery (${deliveryNumber})` : "delivery"}
          </span>
        </div>
      </div>

      {showDetails && (
        <div className="mt-2 pl-4 text-gray-600">
          <div className="grid grid-cols-2 gap-1">
            <div className="font-medium">Price per Unit:</div>
            <div>${dish.price.toFixed(2)}</div>
            <div className="font-medium">Total Price:</div>
            <div>${(dish.price * dish.quantity).toFixed(2)}</div>
          </div>
        </div>
      )}

      <Dialog open={showKeypad} onOpenChange={setShowKeypad}>
        <DialogContent className="sm:max-w-md">
          <div className="py-4">
            <h2 className="text-lg font-semibold mb-4 text-center">
              Enter Delivery Number for {dish.name}
            </h2>
            <NumericKeypad
              value={deliveryNumber}
              onChange={setDeliveryNumber}
              onSubmit={handleKeypadSubmit}
              min={0}
              max={dish.quantity}
              className="w-full"
            />
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default DishSummary;
