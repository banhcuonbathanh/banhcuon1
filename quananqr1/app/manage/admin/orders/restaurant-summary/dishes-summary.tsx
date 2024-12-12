"use client";

import React, { useState, useCallback } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import useDeliveryStore from "@/zusstand/delivery/delivery_zustand";
import { toast } from "@/components/ui/use-toast";
import { useAuthStore } from "@/zusstand/new_auth/new_auth_controller";
import { logWithLevel } from "@/lib/log";

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

  const handleClear = () => {
    setValue("");
  };

  const handleSubmit = () => {
    console.log(
      "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleSubmit"
    );
    const parsedValue = parseInt(value || "0", 10);
    onSubmit(parsedValue);
    setValue("");
  };

  return (
    <div className="p-4 bg-gray-700">
      <input
        type="text"
        value={value}
        className="w-full p-2 mb-4 text-right text-xl rounded"
        readOnly
      />
      <div className="grid grid-cols-3 gap-2">
        {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((num) => (
          <button
            key={num}
            onClick={() => handleNumClick(num)}
            className="p-4 text-xl rounded"
          >
            {num}
          </button>
        ))}
        <button onClick={handleClear} className="p-4 text-xl rounded">
          C
        </button>
        <button
          onClick={() => handleNumClick(0)}
          className="p-4 text-xl rounded"
        >
          0
        </button>
        <button
          onClick={handleSubmit}
          className="p-4 text-xl rounded bg-blue-500 text-white hover:bg-blue-600"
        >
          ✓
        </button>
      </div>
    </div>
  );
};

const DishSummary: React.FC<{
  dish: AggregatedDish;
  http: any;
  orderStore: {
    tableNumber: number;
    getOrderSummary: () => any;
    clearOrder: () => void;
  };
}> = ({ dish, http, orderStore }) => {
  const [showDetails, setShowDetails] = useState(false);
  const [showNumPad, setShowNumPad] = useState(false);

  const guest = useAuthStore((state) => state.guest);
  const user = useAuthStore((state) => state.user);
  const isGuest = useAuthStore((state) => state.isGuest);
  const createDelivery = useDeliveryStore((state) => state.createDelivery);

  const handleDeliverySubmit = useCallback(
    async (quantity: number) => {
      try {
        logWithLevel(
          {
            message: "Starting delivery submission",
            quantity,
            dish: dish.name
          },
          "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleDeliverySubmit start",
          "info",
          3
        );

        const deliveryDetails = {
          deliveryAddress: "Default Address",
          deliveryContact: "Default Contact",
          deliveryNotes: "",
          scheduledTime: new Date().toISOString(),
          deliveryFee: 0
        };

        logWithLevel(
          {
            http,
            guest,
            user,
            isGuest,
            orderStore,
            deliveryDetails
          },
          "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleDeliverySubmit 12121",
          "info",
          3
        );

        try {
          // const orderSummary = orderStore.getOrderSummary();
          // logWithLevel(
          //   {
          //     message: "Order summary retrieved",
          //     orderSummary
          //   },
          //   "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleDeliverySubmit before createDelivery",
          //   "info",
          //   3
          // );

          const result = await createDelivery({
            http,
            guest,
            user,
            isGuest,
            orderStore,
            deliveryDetails
          });

          logWithLevel(
            {
              result,
              http,
              guest,
              user,
              isGuest,
              orderStore,
              deliveryDetails
            },
            "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleDeliverySubmit 131313131",
            "info",
            3
          );

          toast({
            title: "Success",
            description: `Delivery created for ${quantity} ${dish.name}`
          });

          setShowNumPad(false);
        } catch (createDeliveryError) {
          logWithLevel(
            {
              error: createDeliveryError,
              http,
              guest,
              user,
              isGuest,
              orderStore,
              deliveryDetails
            },
            "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleDeliverySubmit createDelivery error",
            "error",
            3
          );
          throw createDeliveryError;
        }
      } catch (error) {
        logWithLevel(
          {
            error,
            http,
            guest,
            user,
            isGuest,
            orderStore
          },
          "quananqr1/app/manage/admin/orders/restaurant-summary/dishes-summary.tsx handleDeliverySubmit outer error",
          "error",
          3
        );

        toast({
          variant: "destructive",
          title: "Error",
          description:
            error instanceof Error ? error.message : "Failed to create delivery"
        });
      }
    },
    [createDelivery, dish.name, guest, http, isGuest, orderStore, user]
  );

  const handleClose = useCallback(() => {
    setShowNumPad(false);
  }, []);

  const toggleDetails = useCallback(() => {
    setShowDetails((prev) => !prev);
  }, []);

  const handleShowNumPad = useCallback(() => {
    setShowNumPad(true);
  }, []);

  return (
    <div className="p-2 mb-2 rounded">
      <div className="flex items-center justify-between">
        <div className="flex-1 cursor-pointer" onClick={toggleDetails}>
          <span className="font-bold">
            {dish.name} : {dish.quantity}
          </span>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleShowNumPad}
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

      {showNumPad && (
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
      )}
    </div>
  );
};

export default DishSummary;
