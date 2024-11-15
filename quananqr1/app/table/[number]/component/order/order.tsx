"use client";

import React, { useEffect, useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
import OrderDetails from "../total-dishes-detail";
import OrderCreationComponent from "./add_order_button";

interface OrderProps {
  number: string;
  token: string;
}

export default function OrderSummary({ number, token }: OrderProps) {
  const { addTableNumber, addTableToken, getOrderSummary } = useOrderStore();

  const [canhkhongrauCount, setCanhKhongrau] = useState(0);
  const [canhCoRau, setCanhCoRau] = useState(0);
  const [smallBowl, setSmallBowl] = useState(0);
  const [toppingTotal, setToppingTotal] = useState("");

  useEffect(() => {
    if (token) {
      addTableToken(token);
    }
    if (number) {
      const tableNumber = addTableNumberconvert(number);
      addTableNumber(tableNumber);
    }
  }, [token, addTableToken, number, addTableNumber]);

  // New useEffect to calculate total
  useEffect(() => {
    const total = `Canh không rau: ${canhkhongrauCount} - Canh rau: ${canhCoRau} - Bát bé: ${smallBowl}`;
    setToppingTotal(total);
  }, [canhkhongrauCount, canhCoRau, smallBowl]);

  const orderSummary = getOrderSummary();

  const handleBowlChange = (
    type: "chili" | "noChili" | "small",
    change: number
  ) => {
    switch (type) {
      case "chili":
        const newToppingValue = canhkhongrauCount + change;
        if (newToppingValue >= 0) setCanhKhongrau(newToppingValue);
        break;
      case "noChili":
        const newNoChiliValue = canhCoRau + change;
        if (newNoChiliValue >= 0) setCanhCoRau(newNoChiliValue);
        break;
      case "small":
        const newSmallBowlValue = smallBowl + change;
        if (newSmallBowlValue >= 0) setSmallBowl(newSmallBowlValue);
        break;
    }
  };

  return (
    <div className="container mx-auto px-4 py-5 space-y-5">
      <Card>
        <CardHeader>
          <CardTitle>Canh Banh Cuon</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="space-y-3">
              {/* Bowl without vegetables */}
              <div className="flex items-center justify-between">
                <span>Canh không rau</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", -1)}
                    disabled={canhkhongrauCount === 0}
                  >
                    -
                  </Button>
                  <span className="mx-2 min-w-[2rem] text-center">
                    {canhkhongrauCount}
                  </span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>

              {/* Bowl with vegetables */}
              <div className="flex items-center justify-between">
                <span>Canh rau</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", -1)}
                    disabled={canhCoRau === 0}
                  >
                    -
                  </Button>
                  <span className="mx-2 min-w-[2rem] text-center">
                    {canhCoRau}
                  </span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>

              {/* Small bowl */}
              <div className="flex items-center justify-between">
                <span>Bát bé</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("small", -1)}
                    disabled={smallBowl === 0}
                  >
                    -
                  </Button>
                  <span className="mx-2 min-w-[2rem] text-center">
                    {smallBowl}
                  </span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("small", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>

              {/* Total summary */}
              <div className="mt-4 p-3 bg-gray-50 rounded-md">
                <span className="font-medium">Total Orders: </span>
                <span>{toppingTotal}</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <OrderDetails
        dishes={orderSummary.dishes}
        sets={orderSummary.sets}
        totalPrice={orderSummary.totalPrice}
        totalItems={orderSummary.totalItems}
      />

      <OrderCreationComponent
        topping={toppingTotal}
        bowlNoChili={canhCoRau}
        table_token={token}
      />
    </div>
  );
}

function addTableNumberconvert(value: string): number {
  let tableNumber: number;

  if (typeof value === "string") {
    if (/^\d+$/.test(value)) {
      tableNumber = parseInt(value, 10);
    } else {
      throw new Error("Invalid input: expected a string of digits.");
    }
  } else if (typeof value === "number") {
    tableNumber = value;
  } else {
    throw new Error("Invalid input: expected a string or number.");
  }

  return tableNumber;
}
