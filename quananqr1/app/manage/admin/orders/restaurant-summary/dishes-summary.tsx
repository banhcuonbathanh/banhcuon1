import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface OrderDetailedDish {
  dish_id: number;
  quantity: number;
  name: string;
  price: number;
  description: string;
  iamge: string;
  status: string;
}

interface AggregatedDish extends OrderDetailedDish {}

interface NumPadProps {
  dishName: string;
  onSubmit: (value: number) => void;
  onClose: () => void;
}

const NumPad: React.FC<NumPadProps> = ({ dishName, onSubmit, onClose }) => {
  const [value, setValue] = useState<string>("");

  const handleNumClick = (num: number) => {
    setValue((prev) => prev + num.toString());
  };

  const handleClear = () => setValue("");

  const handleSubmit = () => {
    onSubmit(parseInt(value || "0", 10));
    setValue("");
    onClose();
  };

  return (
    <div className="p-4 bg-gray-700">
      <input
        type="text"
        value={value}
        className="w-full p-2 mb-4 text-right text-xl  rounded"
        readOnly
      />
      <div className="grid grid-cols-3 gap-2">
        {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((num) => (
          <button
            key={num}
            onClick={() => handleNumClick(num)}
            className="p-4 text-xl  rounded "
          >
            {num}
          </button>
        ))}
        <button onClick={handleClear} className="p-4 text-xl  rounded ">
          C
        </button>
        <button
          onClick={() => handleNumClick(0)}
          className="p-4 text-xl  rounded "
        >
          0
        </button>
        <button
          onClick={handleSubmit}
          className="p-4 text-xl  rounded bg-blue-500 text-white hover:bg-blue-600"
        >
          ✓
        </button>
      </div>
    </div>
  );
};

const DishSummary: React.FC<{ dish: AggregatedDish }> = ({ dish }) => {
  const [showDetails, setShowDetails] = useState(false);
  const [showNumPad, setShowNumPad] = useState(false);

  const handleDeliverySubmit = async (quantity: number) => {
    try {
      const response = await fetch("/api/delivery", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          dishId: dish.dish_id,
          quantity: quantity
        })
      });

      if (!response.ok) {
        throw new Error("Delivery request failed");
      }

      setShowNumPad(false);
    } catch (error) {
      console.error("Error submitting delivery:", error);
    }
  };

  const handleClose = () => {
    setShowNumPad(false);
  };

  return (
    <div className="p-2 mb-2 rounded ">
      <div className="flex items-center justify-between">
        <div
          className="flex-1 cursor-pointer"
          onClick={() => setShowDetails(!showDetails)}
        >
          <span className="font-bold">
            {dish.name} :{dish.quantity}
          </span>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setShowNumPad(true)}
          className="ml-2"
        >
          Delivery
        </Button>
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

      <Dialog open={showNumPad} onOpenChange={setShowNumPad}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Enter Delivery Quantity for {dish.name}</DialogTitle>
          </DialogHeader>
          <NumPad
            dishName={dish.name}
            onSubmit={handleDeliverySubmit}
            onClose={handleClose}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default DishSummary;
