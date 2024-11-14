"use client";

import React, { useEffect, useState } from "react";

import useOrderStore from "@/zusstand/order/order_zustand";
import OrderCreationComponent from "./add_order_button";
import OrderDetails from "../total-dishes-detail";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

interface OrderProps {
  number: string;
  token: string;
}
export default function OrderSummary({ number, token }: OrderProps) {
  const {
    addTableNumber,
    addTableToken,

    dishItems,
    setItems
  } = useOrderStore();
  useEffect(() => {
    if (token) {
      addTableToken(token);
    }
    if (number) {
      const tablenumber = addTableNumberconvert(number);
      addTableNumber(tablenumber);
    }
  }, [token, addTableToken, number]);

  const [bowlChili, setBowlChili] = useState(0);
  const [bowlNoChili, setBowlNoChili] = useState(0);

  const { getOrderSummary } = useOrderStore();
  const orderSummary = getOrderSummary();

  // -------
  const handleBowlChange = (type: "chili" | "noChili", change: number) => {
    if (type === "chili") {
      const newValue = bowlChili + change;
      if (newValue >= 0) setBowlChili(newValue);
    } else {
      const newValue = bowlNoChili + change;
      if (newValue >= 0) setBowlNoChili(newValue);
    }
  };
  // ------
  return (
    <div className="container mx-auto px-4 py-5 space-y-5">
      {/*         canh banh cuon  */}

      <Card>
        <CardHeader>
          <CardTitle>canh banh cuon</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span>canh khong rau</span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", -1)}
                  >
                    -
                  </Button>
                  <span className="mx-2">{bowlChili}</span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("chili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span>canh rau </span>
                <div className="flex items-center gap-2">
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", -1)}
                  >
                    -
                  </Button>
                  <span className="mx-2">{bowlNoChili}</span>
                  <Button
                    size="sm"
                    onClick={() => handleBowlChange("noChili", 1)}
                  >
                    +
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
      {/*         canh banh cuon  */}

      <OrderDetails
        dishes={orderSummary.dishes}
        sets={orderSummary.sets}
        totalPrice={orderSummary.totalPrice}
        totalItems={orderSummary.totalItems}
      />


      <OrderCreationComponent
        bowlChili={bowlChili}
        bowlNoChili={bowlNoChili}
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
