"use client";

import React, { useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import useOrderStore from "@/zusstand/order/order_zustand";
import OrderDetails from "../total-dishes-detail/total-dishes-detail";
import OrderCreationComponent from "./add_order_button";
import { decodeTableToken } from "@/lib/utils";
import { logWithLevel } from "@/lib/logger/log";

import ChoosingTopping from "../topping/canh-banh-cuon";
import useCartStore from "@/zusstand/new-order/new-order-zustand";

interface OrderProps {
  number: string;
  token: string;
}

const LOG_PATH =
  "quananqr1/app/(client)/table/[number]/component/order/order.tsx";

export default function OrderSummary({ number, token }: OrderProps) {
  const decoded = decodeTableToken(token);
  logWithLevel({ decoded, token }, LOG_PATH, "debug", 1);
  const { addTableToken, addTableNumber } = useCartStore();
  // const { addTableNumber, addTableToken, getOrderSummary } = useOrderStore();

  useEffect(() => {
    try {
      if (token) {
        addTableToken(token);
      }
      if (number) {
        const tableNumber = addTableNumberconvert(number);
        addTableNumber(tableNumber);
      }

      logWithLevel(
        {
          token,
          number,
          tableNumber: number ? addTableNumberconvert(number) : null
        },
        LOG_PATH,
        "info",
        2
      );
    } catch (error) {
      logWithLevel({ error, token, number }, LOG_PATH, "error", 2);
    }
  }, [token, addTableToken, number, addTableNumber]);

  return (
    <div className="container mx-auto px-4 py-5 space-y-5">
      <ChoosingTopping />
      <OrderDetails />
      {/* <OrdersDetails /> */}
      {/* <div>
        <h1>Orders</h1>
        <OrderListPage />
      </div> */}

      <OrderCreationComponent />
    </div>
  );
}

function addTableNumberconvert(value: string): number {
  let tableNumber: number;

  try {
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

    logWithLevel(
      {
        inputValue: value,
        convertedNumber: tableNumber
      },
      LOG_PATH,
      "debug",
      7
    );

    return tableNumber;
  } catch (error) {
    logWithLevel(
      {
        error,
        inputValue: value
      },
      LOG_PATH,
      "error",
      7
    );
    throw error;
  }
}
